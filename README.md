# tap-go

Go library for the [Transaction Authorization Protocol (TAP)](https://tap.rsvp) — typed message wrappers for all 20 TAP message types, built on [go-didcomm](https://github.com/Notabene-id/go-didcomm).

## Installation

```bash
go get github.com/TransactionAuthorizationProtocol/tap-go
```

## Quick Start

### Creating a Transfer message

```go
package main

import (
    "fmt"
    tap "github.com/TransactionAuthorizationProtocol/tap-go"
)

func main() {
    msg, err := tap.NewTransferMessage(
        "did:web:originator.vasp",
        []string{"did:web:beneficiary.vasp"},
        &tap.TransferBody{
            Asset:  "eip155:1/slip44:60",
            Amount: "1.23",
            Originator: &tap.Party{ID: "did:eg:bob"},
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
    fmt.Printf("Message ID: %s\n", msg.ID)
    fmt.Printf("Type: %s\n", msg.Type)
}
```

### Parsing a received message

```go
import (
    "encoding/json"
    tap "github.com/TransactionAuthorizationProtocol/tap-go"
    didcomm "github.com/Notabene-id/go-didcomm"
)

// Given a didcomm.Message (from Unpack or plain JSON)
body, err := tap.ParseBody(msg)
if err != nil {
    // handle error
}

switch b := body.(type) {
case *tap.TransferBody:
    fmt.Printf("Transfer of %s %s\n", b.Amount, b.Asset)
case *tap.PaymentBody:
    fmt.Printf("Payment of %s %s\n", b.Amount, b.Currency)
case *tap.AuthorizeBody:
    fmt.Println("Transaction authorized")
// ... handle all 20 message types
}
```

### Using the TAP Client

The `Client` wraps a `didcomm.Client` to handle DIDComm unpacking and TAP body parsing in one step:

```go
dc := didcomm.NewClient(resolver, secrets)
client := tap.NewClient(dc)

result, err := client.Receive(ctx, envelope)
if err != nil {
    // handle error
}

fmt.Printf("Message: %s\n", result.Message.Type)
fmt.Printf("Encrypted: %v, Signed: %v\n", result.Encrypted, result.Signed)

// Type-assert the body
if transfer, ok := result.Body.(*tap.TransferBody); ok {
    fmt.Printf("Transfer: %s %s\n", transfer.Amount, transfer.Asset)
}
```

## Message Types

### Transaction Messages

| Constructor | Body Struct | TAIP | Required Fields |
|------------|-------------|------|-----------------|
| `NewTransferMessage` | `TransferBody` | [TAIP-3](https://tap.rsvp/TAIPs/taip-3) | `asset`, `agents` |
| `NewPaymentMessage` | `PaymentBody` | [TAIP-14](https://tap.rsvp/TAIPs/taip-14) | `amount`, `merchant`, `agents`, `asset` or `currency` |
| `NewExchangeMessage` | `ExchangeBody` | [TAIP-18](https://tap.rsvp/TAIPs/taip-18) | `fromAssets`, `toAssets`, `requester`, `agents`, `fromAmount` or `toAmount` |
| `NewQuoteMessage` | `QuoteBody` | [TAIP-18](https://tap.rsvp/TAIPs/taip-18) | `fromAsset`, `toAsset`, `fromAmount`, `toAmount`, `provider`, `agents`, `expiresAt` |
| `NewEscrowMessage` | `EscrowBody` | [TAIP-17](https://tap.rsvp/TAIPs/taip-17) | `amount`, `originator`, `beneficiary`, `expiry`, `agents`, `asset` or `currency` |

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
