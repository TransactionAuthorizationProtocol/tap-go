# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- GitHub Actions CI with test, lint (golangci-lint), and vulncheck jobs

### Changed

- Bumped Go from 1.25.0 to 1.25.3 to fix stdlib vulnerabilities
- Use published go-didcomm v0.1.0 instead of local replace directive
