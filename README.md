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

### Run with Docker Compose

```bash
docker compose up --build
```

Starts the app (port `8080`) and a Postgres container together; `DATABASE_URL` is preset in `docker-compose.yml` to point at the `db` service.

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
| DELETE | `/wallets/:id` | delete wallet |

Requests for each endpoint are in `requests.http` (VS Code REST Client / IntelliJ HTTP client).

## Test

```bash
go test ./...
```

## CI

`.github/workflow/main.yaml` runs `go vet` and `go test -race` on push/PR to `main`, and builds/pushes a Docker image to GHCR on version tags (`v*`).
