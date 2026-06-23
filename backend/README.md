# EWallet Core Engine Backend

## Tech Stack

| Layer            | Technology                                                |
| ---------------- | --------------------------------------------------------- |
| Language         | Go 1.26                                                   |
| Router           | chi/v5                                                    |
| Database         | PostgreSQL 16                                             |
| SQL Code Gen     | sqlc                                                      |
| Auth             | JWT (golang-jwt/v5) + bcrypt                              |
| Email (dev)      | Mailpit                                                   |
| API Docs         | Swagger (swaggo)                                          |
| Testing          | testify                                                   |
| Containerization | Docker, Docker Compose                                    |
| CI/CD            | GitHub Actions                                            |

## Features

- User registration with password validation, login returning JWT, and JWT auth middleware
- Forgot/reset password flow with time-limited reset tokens sent via email
- Wallet creation in any currency (default VND)
- Background worker auto-creates a VND wallet for each new user
- ACID-compliant P2P money transfer with deadlock prevention, currency validation, and full rollback
- Wallet deletion (zero balance only), wallet details, and ownership enforcement
- Paginated transfer history and audit entry history per wallet

## Project Structure

```
backend
├── cmd
│   └── api
│       └── main.go
├── database
│   ├── migration
│   │   └── schema.sql
│   └── query
│       ├── users.sql
│       └── wallets.sql
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── user
│   │   ├── delivery
│   │   ├── domain
│   │   ├── infrastructure
│   │   │   └── sqlc
│   │   └── usecase
│   └── wallet
│       ├── delivery
│       ├── domain
│       ├── infrastructure
│       │   └── sqlc
│       └── usecase
├── middleware
├── response
├── vendor
├── .air.toml
├── .dockerignore
├── .env
├── .gitignore
├── docker-compose.dev.yaml
├── docker-compose.prod.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── sqlc.yaml
```

## Quick Start

Requirements: Go 1.26+, Docker and Docker Compose.

```bash
git clone <repo-url>
cd backend
cp .env.example .env    # Configure your environment variables
make dev-docker-up      # Start PostgreSQL and Mailpit
make dev                # Run the application with hot-reload
```

- Swagger UI: `http://localhost:{BE_PORT}/swagger/index.html`
- Mailpit UI: `http://localhost:{MAILPIT_UI_PORT}`

## Development

```bash
make sqlc        # Regenerate Go code from SQL queries
make tidy        # Tidy Go module dependencies
make vendor      # Vendor dependencies
make uu-test     # Run user usecase tests
make wu-test     # Run wallet usecase tests
```

Swagger documentation is auto-generated from annotations. After changing handler signatures or request models, run:

```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

## Docker

Two Docker Compose setups are provided:

- **dev** (`docker-compose.dev.yaml`) — includes PostgreSQL, Mailpit, and the app with Air hot-reload and live code mounting.
- **prod** (`docker-compose.prod.yaml`) — multi-stage Alpine build for production deployment.

```bash
make dev-docker-up      # Start dev environment
make dev-docker-down    # Stop dev environment
make prod-docker-up     # Start production environment
```

## CI/CD

- **CI** (`.github/workflows/ewallet_ci.yml`) — triggered on push to main and pull requests. Runs `go mod download`, `sqlc generate`, checks Swagger docs are up to date, verifies `gofmt`, runs `go vet`, executes all unit tests, and builds the Docker image.
- **CD** (`.github/workflows/ewallet_cd.yml`) — triggered on push to main. Builds and pushes the Docker image to GitHub Container Registry with both `latest` and commit SHA tags.

## Roadmap

### Backend (Near-term)
- Refresh token rotation for better security
- Rate limiting on auth endpoints
- Structured logging with zerolog or zap
- Two-factor authentication (TOTP)
- Admin management APIs
- Webhook notifications for transfers
- Formal database migrations with golang-migrate or goose

### Frontend (Future)
- React or Next.js web application
- React Native or Flutter mobile app
- Real-time updates via WebSocket
