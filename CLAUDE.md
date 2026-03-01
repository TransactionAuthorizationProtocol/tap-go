# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detector and coverage
go test -race -coverprofile=coverage.out ./...

# Run a single test
go test -run TestNewTransferMessage ./...

# Build
go build ./...

# Vet
go vet ./...
```

## Architecture

This is a Go library that provides typed wrappers for all 20 TAP (Transaction Authorization Protocol) message types, built on top of `go-didcomm`.

### Package structure

Single flat `tap` package — all types and helpers in the root. One file per message type.

### Key types

- **`TAPBody`** — interface implemented by all body structs (`TAPType() string`)
- **`Client`** — wraps `didcomm.Client`, adds `Receive()` for typed TAP parsing
- **`ParseBody(msg)`** — dispatches `didcomm.Message` to the correct typed body struct
- **`IsTAPMessage(msg)`** — checks if a DIDComm message is a known TAP type

### File layout

| File | Purpose |
|------|---------|
| `message_types.go` | Type URL constants (`TypeTransfer`, `TypePayment`, etc.) and `AllTypes()` |
| `types.go` | Shared types: `Party`, `Agent`, `ForField`, `Policy`, `TransactionConstraints`, `Limits`, `Invoice`, `LineItem`, `Relationship` |
| `errors.go` | `ErrInvalidBody`, `ErrUnknownMessageType` |
| `parse.go` | `TAPBody` interface, `ParseBody()` dispatch, `IsTAPMessage()` |
| `client.go` | `Client` wrapping `didcomm.Client`, `Receive()` returning `TAPResult` |
| `transfer.go` | `TransferBody` + `NewTransferMessage()` |
| `payment.go` | `PaymentBody` + `NewPaymentMessage()` |
| `exchange.go` | `ExchangeBody` + `NewExchangeMessage()` |
| `quote.go` | `QuoteBody` + `NewQuoteMessage()` |
| `escrow.go` | `EscrowBody` + `NewEscrowMessage()` |
| `authorize.go` | `AuthorizeBody` + `NewAuthorizeMessage()` |
| `authorization_required.go` | `AuthorizationRequiredBody` + `NewAuthorizationRequiredMessage()` |
| `settle.go` | `SettleBody` + `NewSettleMessage()` |
| `reject.go` | `RejectBody` + `NewRejectMessage()` |
| `cancel.go` | `CancelBody` + `NewCancelMessage()` |
| `revert.go` | `RevertBody` + `NewRevertMessage()` |
| `capture.go` | `CaptureBody` + `NewCaptureMessage()` |
| `update_agent.go` | `UpdateAgentBody` + `NewUpdateAgentMessage()` |
| `update_party.go` | `UpdatePartyBody` + `NewUpdatePartyMessage()` |
| `add_agents.go` | `AddAgentsBody` + `NewAddAgentsMessage()` |
| `replace_agent.go` | `ReplaceAgentBody` + `NewReplaceAgentMessage()` |
| `remove_agent.go` | `RemoveAgentBody` + `NewRemoveAgentMessage()` |
| `confirm_relationship.go` | `ConfirmRelationshipBody` + `NewConfirmRelationshipMessage()` |
| `update_policies.go` | `UpdatePoliciesBody` + `NewUpdatePoliciesMessage()` |
| `connect.go` | `ConnectBody` + `NewConnectMessage()` |

### Pattern for each message type

Each message file follows the same pattern:

1. Body struct with JSON tags (`@context`, `@type`, plus message-specific fields)
2. `TAPType()` method returning the type constant
3. `New*Message()` constructor that:
   - Validates required fields (returns `ErrInvalidBody` on failure)
   - Sets `@context` and `@type` automatically
   - Marshals body into `json.RawMessage`
   - Returns a `*didcomm.Message` with a UUID `ID`

### Dependencies

- `github.com/Notabene-id/go-didcomm` — DIDComm v2 (`Message`, `Client`, resolvers)
- `github.com/google/uuid` — message ID generation

### Test vectors

Test vectors from `TAIPs/test-vectors/` are loaded in tests. The `TAIPs` directory is a git submodule (symlinked to `../TAIPs` for local development). Tests skip gracefully if vectors are not available.

### Agent `for` field

The `Agent.For` field uses `ForField` type, which handles JSON marshaling of both single DID strings and arrays of DIDs. Use `NewForField("did:eg:alice")` or `NewForField("did:eg:alice", "did:eg:bob")`.

## Development guidelines

- Every new message type needs: body struct, `TAPType()`, `New*Message()` constructor, a case in `ParseBody()`, and a matching `_test.go`
- Constructors validate required fields and return `fmt.Errorf("%w: ...", ErrInvalidBody)`
- All body structs must have `Context` (`@context`) and `Type` (`@type`) fields, set automatically by constructors
- Thread-based messages (replies) take a `thid` parameter; initiating messages do not
- Test files include: JSON round-trip, constructor validation (required fields), `ParseBody` dispatch, and test vector loading where available

## Documentation requirements

- **CHANGELOG.md** — Maintain a `CHANGELOG.md` in the project root using [Keep a Changelog](https://keepachangelog.com/) format. Update it with every user-facing change (new features, bug fixes, breaking changes, dependency updates). Group entries under `Added`, `Changed`, `Fixed`, `Removed` sections within version headings.
- **README.md** — Update `README.md` whenever changes affect public API, usage examples, installation instructions, or project capabilities.
- **CLAUDE.md** — Update this file whenever changes affect architecture, file layout, commands, dependencies, or development guidelines (e.g., new message types added to the file layout table, new commands, changed patterns).
