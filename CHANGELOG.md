# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

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

### Removed

- The `Relationship` struct (`types.go`), which was used only by the now-reshaped
  `ConfirmRelationshipBody`.

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
