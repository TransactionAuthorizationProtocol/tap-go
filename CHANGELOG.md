# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

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
