# hexagonal_intro

Wallet service built with hexagonal (ports & adapters) architecture in Go.

## Layout

```
cmd/main.go                        composition root: wiring, HTTP routes
internal/core/model                domain types
internal/core/port                 interfaces the core depends on (e.g. WalletRepository)
internal/core/service               business logic, depends only on ports
internal/adapter/handler           driving adapter: Gin HTTP handlers + DTOs
internal/adapter/repository        driven adapter: Postgres/GORM implementation of the port
internal/config                    env loading
```

Core has no dependency on Gin or GORM; adapters implement/consume the core's interfaces.

## Run

```bash
cp .env.example .env   # set DATABASE_URL, PORT
go run ./cmd
```

## Env vars

| Var | Default |
|---|---|
| `DATABASE_URL` | `host=localhost user=postgres password=postgres dbname=walletdb port=5432 sslmode=disable TimeZone=Asia/Bangkok` |
| `PORT` | `8080` |

## API

| Method | Path | Description |
|---|---|---|
| GET | `/health` | health check |
| POST | `/wallets` | create wallet (`{"user_id": "..."}`) |
| GET | `/wallets/:id` | get wallet |
| POST | `/wallets/:id/deposit` | deposit (`{"amount": 10.5}`) |
| POST | `/wallets/:id/withDraw` | withdraw (`{"amount": 10.5}`) |
| POST | `/wallets/:id` | delete wallet |

> Known bug: in `cmd/main.go` the deposit route is wired to `CreateWallet` instead of `Deposit`, and delete uses `POST` instead of `DELETE`.

## Test

```bash
go test ./...
```
