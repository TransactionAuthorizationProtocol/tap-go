# tap-go

Go library for the [Transaction Authorization Protocol (TAP)](https://tap.rsvp) — typed message wrappers for all 20 TAP message types, built on [go-didcomm](https://github.com/Notabene-id/go-didcomm).

## Installation

### Library

```bash
go get github.com/TransactionAuthorizationProtocol/tap-go
```

### CLI

```bash
go install github.com/TransactionAuthorizationProtocol/tap-go/cmd/tap@latest
```

## How It Works

`tap-go` sits on top of [go-didcomm](https://github.com/Notabene-id/go-didcomm), which handles DIDComm v2 message packing (signing, encryption) and unpacking. This library adds TAP-specific typed bodies, validation, and parsing.

```
┌─────────────────────────────────────────────────┐
│                  Your Application                │
├─────────────────────────────────────────────────┤
│  tap-go                                         │
│  • Typed body structs (TransferBody, etc.)      │
│  • Constructors with validation                 │
│  • ParseBody() dispatch                         │
│  • tap.Client.Receive() for typed parsing       │
├─────────────────────────────────────────────────┤
│  go-didcomm                                     │
│  • didcomm.Message (ID, Type, From, To, Body)   │
│  • Pack: Signed / Anoncrypt / Authcrypt → []byte│
│  • Unpack: []byte → UnpackResult                │
│  • DID resolution, key management               │
└─────────────────────────────────────────────────┘
```

Every `New*Message()` constructor returns a `*didcomm.Message` — the standard go-didcomm type — with the TAP body serialized into `msg.Body` as `json.RawMessage`. You can then use go-didcomm's `PackSigned`, `PackAnoncrypt`, or `PackAuthcrypt` to send it, and `Unpack` + `ParseBody` (or `tap.Client.Receive`) to receive it.

## Quick Start

### Sending: Create and pack a TAP message

```go
package main

import (
    "context"
    "fmt"

    tap "github.com/TransactionAuthorizationProtocol/tap-go"
    didcomm "github.com/Notabene-id/go-didcomm"
)

func main() {
    // 1. Create a TAP message (returns a *didcomm.Message)
    msg, err := tap.NewTransferMessage(
        "did:web:originator.vasp",                    // from
        []string{"did:web:beneficiary.vasp"},          // to
        &tap.TransferBody{
            Asset:  "eip155:1/slip44:60",
            Amount: "1.23",
            Originator:  &tap.Party{ID: "did:eg:bob"},
            Beneficiary: &tap.Party{ID: "did:eg:alice"},
            Agents: []tap.Agent{
                {ID: "did:web:originator.vasp", For: tap.NewForField("did:eg:bob")},
                {ID: "did:web:beneficiary.vasp", For: tap.NewForField("did:eg:alice")},
            },
        },
    )
    if err != nil {
        panic(err)
    }

    // 2. Pack with go-didcomm for transport (signed, encrypted, or both)
    dc := didcomm.NewClient(resolver, secrets)

    // Option A: Sign only (sender authenticated, not encrypted)
    envelope, err := dc.PackSigned(context.Background(), msg)

    // Option B: Encrypt only (anonymous, no sender identification)
    envelope, err = dc.PackAnoncrypt(context.Background(), msg)

    // Option C: Sign then encrypt (authenticated + encrypted)
    envelope, err = dc.PackAuthcrypt(context.Background(), msg)

    // 3. Send envelope ([]byte) over any transport (HTTP, WebSocket, etc.)
    fmt.Printf("Sending %d bytes\n", len(envelope))
}
```

### Receiving: Unpack and parse a TAP message

```go
// Option A: Two-step — unpack with go-didcomm, then parse with tap-go
dc := didcomm.NewClient(resolver, secrets)
result, err := dc.Unpack(ctx, envelope)  // → *didcomm.UnpackResult
if err != nil { /* handle */ }

// Check if it's a TAP message
if tap.IsTAPMessage(result.Message) {
    body, err := tap.ParseBody(result.Message)
    if err != nil { /* handle */ }

    switch b := body.(type) {
    case *tap.TransferBody:
        fmt.Printf("Transfer of %s %s\n", b.Amount, b.Asset)
    case *tap.PaymentBody:
        fmt.Printf("Payment of %s %s\n", b.Amount, b.Currency)
    case *tap.AuthorizeBody:
        fmt.Println("Transaction authorized")
    // ... handle all 20 types
    }
}

// Option B: One-step — tap.Client does both unpack + parse
client := tap.NewClient(dc)
tapResult, err := client.Receive(ctx, envelope)  // → *tap.TAPResult
if err != nil { /* handle */ }

fmt.Printf("Type: %s\n", tapResult.Message.Type)
fmt.Printf("Encrypted: %v, Signed: %v\n", tapResult.Encrypted, tapResult.Signed)

if transfer, ok := tapResult.Body.(*tap.TransferBody); ok {
    fmt.Printf("Transfer: %s %s\n", transfer.Amount, transfer.Asset)
}
```

### Replying to a message (thread-based messages)

Authorization flow and agent management messages reference an existing thread via `thid`:

```go
// Original transfer message starts a thread
transferMsg, _ := tap.NewTransferMessage(from, to, transferBody)

// Reply messages reference the original via thid
authorizeMsg, _ := tap.NewAuthorizeMessage(
    "did:web:beneficiary.vasp",                  // from (replier)
    []string{"did:web:originator.vasp"},          // to (original sender)
    transferMsg.ID,                               // thid — links to original
    &tap.AuthorizeBody{
        SettlementAddress: "eip155:1:0x1234a96D359eC26a11e2C2b3d8f8B8942d5Bfcdb",
    },
)

// Pack and send the reply the same way
envelope, _ := dc.PackAuthcrypt(ctx, authorizeMsg)
```

## Message Types

### Transaction Messages

| Constructor | Body Struct | TAIP | Required Fields |
|------------|-------------|------|-----------------|
| `NewTransferMessage` | `TransferBody` | [TAIP-3](https://tap.rsvp/TAIPs/taip-3) | `asset`, `agents` |
| `NewPaymentMessage` | `PaymentBody` | [TAIP-14](https://tap.rsvp/TAIPs/taip-14) | `amount`, `merchant`, `agents`, `asset` or `currency` |
| `NewRFQMessage` | `RFQBody` | [TAIP-18](https://tap.rsvp/TAIPs/taip-18) | `fromAssets`, `toAssets`, `requester`, `agents`, `fromAmount` or `toAmount` |
| `NewQuoteMessage` | `QuoteBody` | [TAIP-18](https://tap.rsvp/TAIPs/taip-18) | `fromAsset`, `toAsset`, `fromAmount`, `toAmount`, `provider`, `agents`, `expiresAt` |
| `NewLockMessage` | `LockBody` | [TAIP-17](https://tap.rsvp/TAIPs/taip-17) | `amount`, `originator`, `beneficiary`, `expiry`, `agents`, `asset` or `currency` |

### Authorization Flow Messages

| Constructor | Body Struct | TAIP | Required Fields |
|------------|-------------|------|-----------------|
| `NewAuthorizeMessage` | `AuthorizeBody` | [TAIP-4](https://tap.rsvp/TAIPs/taip-4) | _(all optional)_ |
| `NewAuthorizationRequiredMessage` | `AuthorizationRequiredBody` | [TAIP-4](https://tap.rsvp/TAIPs/taip-4) | `authorizationUrl`, `expires` |
| `NewSettleMessage` | `SettleBody` | [TAIP-4](https://tap.rsvp/TAIPs/taip-4) | `settlementAddress` |
| `NewRejectMessage` | `RejectBody` | [TAIP-4](https://tap.rsvp/TAIPs/taip-4) | _(all optional)_ |
| `NewCancelMessage` | `CancelBody` | [TAIP-4](https://tap.rsvp/TAIPs/taip-4) | `by` |
| `NewRevertMessage` | `RevertBody` | [TAIP-4](https://tap.rsvp/TAIPs/taip-4) | `settlementAddress`, `reason` |
| `NewCaptureMessage` | `CaptureBody` | [TAIP-17](https://tap.rsvp/TAIPs/taip-17) | _(all optional)_ |

### Agent Management Messages

| Constructor | Body Struct | TAIP | Required Fields |
|------------|-------------|------|-----------------|
| `NewUpdateAgentMessage` | `UpdateAgentBody` | [TAIP-5](https://tap.rsvp/TAIPs/taip-5) | `agent` |
| `NewUpdatePartyMessage` | `UpdatePartyBody` | [TAIP-6](https://tap.rsvp/TAIPs/taip-6) | `party`, `role` |
| `NewAddAgentsMessage` | `AddAgentsBody` | [TAIP-5](https://tap.rsvp/TAIPs/taip-5) | `agents` |
| `NewReplaceAgentMessage` | `ReplaceAgentBody` | [TAIP-5](https://tap.rsvp/TAIPs/taip-5) | `original`, `replacement` |
| `NewRemoveAgentMessage` | `RemoveAgentBody` | [TAIP-5](https://tap.rsvp/TAIPs/taip-5) | `agent` |

### Relationship Messages

| Constructor | Body Struct | TAIP | Required Fields |
|------------|-------------|------|-----------------|
| `NewConfirmRelationshipMessage` | `ConfirmRelationshipBody` | [TAIP-9](https://tap.rsvp/TAIPs/taip-9) | `relationship`, `status` |
| `NewUpdatePoliciesMessage` | `UpdatePoliciesBody` | [TAIP-7](https://tap.rsvp/TAIPs/taip-7) | `policies` |
| `NewConnectMessage` | `ConnectBody` | [TAIP-15](https://tap.rsvp/TAIPs/taip-15) | `requester`, `principal`, `agents`, `constraints` |

## Shared Types

### Party

Represents a real-world entity (legal or natural person):

```go
party := tap.Party{
    ID:   "did:eg:bob",       // required — DID or IRI
    Type: "Party",            // optional — JSON-LD type
    Name: "Bob's Exchange",   // optional
    MCC:  "5734",             // optional — Merchant Category Code
}
```

### Agent

Represents software acting on behalf of a participant:

```go
agent := tap.Agent{
    ID:   "did:web:originator.vasp",              // required
    For:  tap.NewForField("did:eg:bob"),           // DID(s) of represented party
    Role: "SettlementAddress",                     // optional
    Policies: []tap.Policy{                        // optional
        {Type: "RequireAuthorization", FromAgent: "beneficiary"},
    },
}
```

### ForField

The `Agent.For` field handles both single DIDs and arrays:

```go
// Single DID — marshals as "did:eg:alice"
single := tap.NewForField("did:eg:alice")

// Multiple DIDs — marshals as ["did:eg:alice", "did:eg:bob"]
multi := tap.NewForField("did:eg:alice", "did:eg:bob")

// Access values
single.String()   // "did:eg:alice"
multi.Values()    // ["did:eg:alice", "did:eg:bob"]
single.IsEmpty()  // false
```

### Policy

```go
policy := tap.Policy{
    Type:                   "RequirePresentation",
    FromAgent:              "beneficiary",
    AboutParty:             "beneficiary",
    Purpose:                "FATF Travel Rule Compliance",
    PresentationDefinition: "https://tap.rsvp/presentation-definitions/ivms-101/eu/tfr",
}
```

### TransactionConstraints

Used in `Connect` messages to define boundaries:

```go
constraints := tap.TransactionConstraints{
    Purposes:         []string{"BEXP", "SUPP"},
    CategoryPurposes: []string{"CASH"},
    Limits: &tap.Limits{
        PerTransaction: "10000.00",
        PerDay:         "50000.00",
        Currency:       "USD",
    },
    AllowedAssets: []string{"eip155:1/slip44:60"},
}
```

## Utility Functions

```go
// Check if a DIDComm message is a TAP message
tap.IsTAPMessage(msg) // bool

// Get all known TAP type URLs
tap.AllTypes() // []string (20 types)

// Parse body into typed struct
body, err := tap.ParseBody(msg) // (TAPBody, error)
body.TAPType() // e.g. "https://tap.rsvp/schema/1.0#Transfer"
```

## CLI

The `tap` CLI wraps all [go-didcomm](https://github.com/Notabene-id/go-didcomm) CLI commands and adds TAP-specific message creation and receiving.

### Commands

```
tap did generate-key [--output-dir <dir>]
tap did generate-web --domain <d> [--path <p>] [--service-endpoint <url>] [--output-dir <dir>]
tap pack signed    --key-file <f> [--send] [--did-doc <f>] [--message <m>]
tap pack anoncrypt [--send] [--did-doc <f>] [--message <m>]
tap pack authcrypt --key-file <f> [--send] [--did-doc <f>] [--message <m>]
tap unpack         --key-file <f> [--did-doc <f>] [--message <m>]
tap send           --to <url> [--message <m>]
tap message <type> --from <did> --to <did> [--thid <id>] [--body <json>]
tap receive        --key-file <f> [--did-doc <f>] [--message <m>]
```

### Create and pack a TAP message

```bash
# Generate identities
tap did generate-key --output-dir alice
tap did generate-key --output-dir bob

ALICE=$(jq -r .id alice/did-doc.json)
BOB=$(jq -r .id bob/did-doc.json)

# Create a TAP transfer message
tap message transfer --from $ALICE --to $BOB \
  --body '{"asset":"eip155:1/slip44:60","amount":"1.5","agents":[{"@id":"'$ALICE'","role":"OriginatingVASP"}]}'

# Pipe: create → pack → send
tap message transfer --from $ALICE --to $BOB --body @body.json | \
  tap pack authcrypt --key-file alice/keys.json
```

### Receive a TAP message

```bash
# Unpack a DIDComm envelope and parse the TAP body
echo '<packed-message>' | tap receive --key-file bob/keys.json
```

The `receive` command outputs JSON with the unpacked message, typed body, and envelope metadata (`encrypted`, `signed`, `anonymous`).

### TAP message types

**Initiating (no `--thid`):** `transfer`, `payment`, `rfq`, `lock`, `connect`

**Reply (require `--thid`):** `authorize`, `authorization-required`, `settle`, `reject`, `cancel`, `revert`, `capture`, `quote`, `add-agents`, `remove-agent`, `replace-agent`, `update-agent`, `update-party`, `update-policies`, `confirm-relationship`

### Body input (`--body` flag)

- `'{"json"}'` — inline JSON string
- `@file.json` — read from file
- `-` — read from stdin
- _(omitted)_ — defaults to `{}`

The body JSON should contain only message-specific fields (e.g. `asset`, `amount`, `agents`). The CLI automatically sets `@context` and `@type`.

## Error Handling

```go
import "errors"

msg, err := tap.NewTransferMessage(from, to, body)
if errors.Is(err, tap.ErrInvalidBody) {
    // missing required fields
}

body, err := tap.ParseBody(msg)
if errors.Is(err, tap.ErrUnknownMessageType) {
    // not a recognized TAP message type
}
if errors.Is(err, tap.ErrInvalidBody) {
    // body JSON could not be unmarshaled
}
```

## License

MIT
