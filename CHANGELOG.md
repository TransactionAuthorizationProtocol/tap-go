# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.8.0] - 2026-09-11

### Changed

- **BREAKING** `Agent.For` (`for`) is enforced on the blockchain-address roles
  (`SourceAddress`, `SettlementAddress`). `NewTransferMessage`,
  `NewPaymentMessage`, `NewConnectMessage`, `NewQuoteMessage`, `NewRFQMessage`,
  `NewAddAgentsMessage`, `NewLockMessage`, `NewUpdateAgentMessage` and
  `NewReplaceAgentMessage` reject an address agent that omits it, and reject an
  empty DID inside `for` whatever the role.

  [TAIP-5] marks `for` REQUIRED on every agent, but provides no way to express
  that who owns an agent is not established yet — a state real flows pass
  through, since an address can be seen before anybody has resolved who
  custodies it. Enforcing the letter of the spec would mean inventing owners,
  and an invented `for` is worse than an absent one: a receiver stores it as
  fact, and it decides whether a wallet is treated as self-hosted (must prove
  ownership) or custodied (must not). So `for` is enforced where the sender
  cannot honestly be unsure — whoever supplies an address knows whose it is —
  and institutional agents may still travel without one.

  Validation is send-side only: `ParseBody` accepts inbound agents that omit
  `for` regardless of role, so peers that do not set it keep working.
- **BREAKING** `UpdatePartyBody.Role` is renamed to `PartyType` and retagged
  `partyType`, per [TAIP-6]. The previous `role` spelling matched no other
  implementation, so these messages were silently ignored by conformant peers.
  Inbound bodies still accept `role` as a fallback.

### Fixed

- An agent with no owner no longer serialises `"for": null`: `omitempty` has no
  effect on a struct field, so the key was always written. `Agent.For` is now
  tagged `omitzero` and the key is omitted. Unmarshalling `null` also stopped
  producing a single empty DID.

### Added

- `Agent.Validate()` and `ValidateAgents()` report TAIP-5 violations.

[TAIP-5]: https://github.com/TransactionAuthorizationProtocol/TAIPs/blob/main/TAIPs/taip-5.md
[TAIP-6]: https://github.com/TransactionAuthorizationProtocol/TAIPs/blob/main/TAIPs/taip-6.md

## [0.7.0] - 2026-07-23

### Added

- `AuthorizeBody` gains `MemoTag` (`memoTag`) — a Notabene extension (TAIP-4
  defines no memo-tag carriage) for the settlement memo/destination tag, as a
  fallback to embedding it as a `:tag` suffix on `settlementAddress`.

## [0.6.0] - 2026-07-08

### Added

- **TAIP-15 trust connections** (per the pending spec revision
  [TAIPs#53](https://github.com/TransactionAuthorizationProtocol/TAIPs/pull/53)):
  `ConnectBody` gains `ConnectionTypes` (`connectionTypes`) and `Action`
  (`action`) fields, plus `ConnectionType*` and `ConnectAction*` constants.
  `AuthorizeBody` gains `ApprovedTypes` (`approvedTypes`) for connection
  approvals.

### Changed

- `NewConnectMessage` validates `requester`/`principal`/`agents`/`constraints`
  only for transactional connections (`connectionTypes` absent or containing
  `"transaction"`). Trust connections (`ddq-access`, `mutual-trust`,
  `whitelist`) omit them, and the four fields are now `omitempty` in JSON.
- **BREAKING (TAIP-9):** Reshaped `ConfirmRelationshipBody` to match the TAIP-9
  spec, whose confirmation payload is an `Agent` payload. The body is now flat:
  `@context`, `@type` (set to `https://tap.rsvp/schema/1.0#Agent`), `@id` (the
  DID of the agent being confirmed, REQUIRED), `for` (the DID of the entity the
  agent acts for, REQUIRED), and `role` (OPTIONAL). The non-spec
  `relationship` (`Relationship *Relationship`), `status`, `validFrom`,
  `validUntil`, and `details` fields were removed. `NewConfirmRelationshipMessage`
  now validates `@id` and `for` instead of `relationship` and `status`. Consumers
  must read `body.ID` for the confirmed address and `body.For` for the owner
  instead of `body.Relationship.Parties` / `body.Status`. Added the `TypeAgent`
  constant for the body's JSON-LD `@type`.
- **BREAKING:** `Client.Receive` now authenticates the sender — plain and
  anonymously-encrypted messages are rejected. Use the new
  `Client.ReceiveUnverified` to accept them.
- **BREAKING:** `TAPResult` drops `Signed` and gains `SenderDID` (the
  cryptographically verified sender, empty when unverified); `Anonymous` is now
  derived from an empty `SenderDID`.

### Removed

- The `Relationship` struct (`types.go`), which was used only by the now-reshaped
  `ConfirmRelationshipBody`.

### Dependencies

- Bumped `go-didcomm` to v0.5.0 (module path lowercased to
  `github.com/notabene-id/go-didcomm`; pinned at
  `v0.4.1-0.20260708110526-ba45976d2288`). Its `Unpack` now returns
  `(*Message, *Metadata, error)` and exposes `UnpackUnverified`.

## [0.4.0] - 2026-05-05

### Changed

- Bumped `github.com/Notabene-id/go-didcomm` from v0.2.0 to v0.4.0

## [0.3.0] - 2026-05-04

### Added

- GitHub Actions CI with test, lint (golangci-lint), and vulncheck jobs
- `TransferBody.TransactionValue` (TAIP-3) — optional fiat-equivalent value
  (`amount`, `currency`) for Travel Rule threshold determination when an asset
  is not widely traded

### Changed

- Bumped Go from 1.25.0 to 1.26.2 to fix stdlib vulnerabilities
- Use published go-didcomm v0.1.0 instead of local replace directive
- **BREAKING (TAIP-17):** Renamed `Escrow` message type to `Lock`. Constants,
  types, constructors, files, and CLI subcommand all renamed:
  `TypeEscrow` → `TypeLock`, `EscrowBody` → `LockBody`,
  `NewEscrowMessage` → `NewLockMessage`, `escrow.go` → `lock.go`,
  `tap message escrow` → `tap message lock`. The `EscrowAgent` role name is
  preserved.
- **BREAKING (TAIP-18):** Renamed `Exchange` message type to `RFQ` (Request
  for Quote). `TypeExchange` → `TypeRFQ`, `ExchangeBody` → `RFQBody`,
  `NewExchangeMessage` → `NewRFQMessage`, `exchange.go` → `rfq.go`,
  `tap message exchange` → `tap message rfq`.
