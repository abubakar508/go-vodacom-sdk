# Changelog

All notable changes to `go-vodacom-sdk` will be documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project uses Go module semantic versioning.

## [Unreleased]

### Added

- Complete endpoint coverage for the Vodacom/Vodafone M-Pesa OpenAPI pages provided so far:
  - Generate SessionKey
  - C2B Single Stage
  - B2C Single Stage
  - B2B Single Stage
  - Reversal
  - Query Transaction Status
  - C2B Multi Stage
  - Update Transaction Status
  - Direct Debit Create
  - Direct Debit Payment
  - Query Beneficiary Name
  - Query Direct Debit
  - Cancel Direct Debit
- Configurable platform public key through `Config.PublicKey`, `WithPublicKey`, and `MPESA_PUBLIC_KEY`.
- Environment-based config through `ConfigFromEnv`.
- Built-in market support for Ghana, Tanzania, Lesotho, DRC, and Mozambique.
- Custom market support for future portal additions.
- Session helper flow: `GenerateSession`, `GenerateSessionAndWait`, and `*WithSession` endpoint methods.
- Async callback request/acknowledgement helpers.
- Request validation helpers.
- HTTP callback helper utilities.
- GitHub Actions CI workflow.
- Documentation site entrypoint for GitHub Pages.

## [0.1.0] - Planned

Initial public release.
