# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run ./cmd            # run the server (reads .env, falls back to defaults)
go build ./...           # build
go test ./...            # run all tests
go test ./internal/core/service -run TestName -v   # run a single test
```

Requires a running Postgres reachable via `DATABASE_URL` (see `.env`); `PORT` defaults to `8080`.

## Architecture

Hexagonal (ports & adapters). Dependency direction: adapters → port interface ← core. The core package has zero imports of Gin or GORM.

- `internal/core/model` — domain types (`Wallet`)
- `internal/core/port` — interfaces the core depends on, e.g. `WalletRepository` (`internal/core/port/repository.go`)
- `internal/core/service` — business logic (`WalletService`); depends only on `port.WalletRepository`, never on a concrete adapter
- `internal/adapter/handler` — driving adapter: Gin HTTP handlers + request/response DTOs, translates HTTP into service calls
- `internal/adapter/repository` — driven adapter: `PostgresWalletRepository` (GORM), implements `port.WalletRepository`; keeps a separate `WalletGORM` struct with `toDomain`/`fromDomain` conversions so the domain model stays persistence-agnostic
- `internal/config` — env loading (`godotenv`)
- `cmd/main.go` — composition root: constructs the Postgres adapter, injects it into `WalletService` via the port interface, wraps the service in the HTTP handler, and registers Gin routes. This is the only place that knows all concrete types.

When adding a new capability to the core, add the method to the port interface first, implement it in the adapter, then expose it through the handler — don't let `core/service` reach for Gin or GORM directly.

### Known bugs (do not silently "fix" as unrelated cleanup — call out explicitly if touching this code)

- `cmd/main.go`: the `POST /wallets/:id/deposit` route is wired to `walletHdr.CreateWallet` instead of `walletHdr.Deposit`.
- `cmd/main.go`: wallet delete is registered as `POST /wallets/:id` instead of `DELETE /wallets/:id`.
