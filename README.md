# GofiberExpenseTracker

A personal learning project for practicing Go and the Fiber web framework. This REST API implements a personal expense tracking system with JWT authentication, PostgreSQL persistence, and clean architecture patterns.

---

## Overview

This project was built as a hands-on exercise to explore the Go ecosystem — specifically Fiber, JWT-based auth, layered architecture, SQL migrations, and containerized deployment. It is not intended for production use, but follows production-oriented practices throughout.

**Core capabilities:**
- User registration and login with bcrypt password hashing
- JWT-protected transaction management (income and expense records)
- Filtering, pagination, and financial summary aggregation
- Automatic SQL migration system
- Docker and Kubernetes deployment support

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22 |
| Web Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| Database | PostgreSQL |
| Driver | `lib/pq` |
| Authentication | `golang-jwt/jwt` v5 |
| Password Hashing | `golang.org/x/crypto` (bcrypt) |
| Testing | `testify` |
| Containerization | Docker, Docker Compose |
| Orchestration | Kubernetes |
| CI/CD | GitHub Actions |

---

## Project Structure

```
.
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── config/                    # Environment configuration
│   ├── database/                  # DB connection and auto-migration
│   ├── handlers/                  # HTTP request/response layer
│   ├── middleware/                 # JWT authentication middleware
│   ├── models/                    # Domain models and request structs
│   ├── repositories/              # Data access layer (SQL)
│   ├── routes/                    # Route registration
│   └── services/                  # Business logic and tests
├── migrations/                    # SQL migration files
├── k8s/                           # Kubernetes manifests
├── .github/workflows/             # CI/CD pipeline
├── Dockerfile
└── docker-compose.yml
```

The architecture follows a clean layered pattern:

```
HTTP Request -> Handler -> Service -> Repository -> PostgreSQL
```

Each layer has a single responsibility. Services contain all business logic. Repositories own all SQL. Handlers handle only request parsing and response formatting.

---

## API Endpoints

### Public

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive a JWT token |
| `GET` | `/health` | Health check |

### Protected (requires `Authorization: Bearer <token>`)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/transactions` | Create a transaction |
| `GET` | `/api/v1/transactions` | List transactions (with filters and pagination) |
| `GET` | `/api/v1/transactions/summary` | Get income, expense, and balance summary |
| `GET` | `/api/v1/transactions/:id` | Get transaction by ID |
| `PUT` | `/api/v1/transactions/:id` | Update a transaction |
| `DELETE` | `/api/v1/transactions/:id` | Delete a transaction |

### Query Parameters for `GET /api/v1/transactions`

| Parameter | Type | Description |
|---|---|---|
| `type` | string | Filter by `income` or `expense` |
| `category` | string | Filter by category name |
| `from` | date | Start date (`YYYY-MM-DD`) |
| `to` | date | End date (`YYYY-MM-DD`) |
| `page` | int | Page number (default: 1) |
| `limit` | int | Items per page (default: 20, max: 100) |

---

## Database Schema

```sql
-- Users
CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Transactions
CREATE TABLE transactions (
    id          BIGSERIAL PRIMARY KEY,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    amount      DECIMAL(12, 2) NOT NULL CHECK (amount > 0),
    category    VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    date        DATE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

Migrations are applied automatically on startup via a tracked migration system stored in the `schema_migrations` table.

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker and Docker Compose

### Run with Docker Compose

```bash
git clone <repository-url>
cd GofiberExpenseTracker

docker-compose up
```

The API will be available at `http://localhost:3000`.

### Run Locally (without Docker)

1. Ensure PostgreSQL is running and a database is created.
2. Copy and configure the environment file:

```bash
cp .env.example .env
```

3. Install dependencies and run:

```bash
go mod download
go run ./cmd/main.go
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3000` | Server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `expense_tracker` | Database name |
| `DB_SSL_MODE` | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | — | Secret key for signing JWT tokens |
| `JWT_EXPIRE_HOURS` | `24` | Token expiration in hours |
| `CORS_ORIGINS` | `*` | Allowed CORS origins |

---

## Running Tests

```bash
go test ./internal/services/...
```

Tests cover the service layer using `testify` for assertions.

---

## Deployment

### Docker

```bash
docker build -t gofiber-expense-tracker .
docker run -p 3000:3000 --env-file .env gofiber-expense-tracker
```

The Dockerfile uses a multi-stage build (Go builder on `golang:1.22-alpine`, runtime on `alpine:3.19`) with a non-root user and stripped binary for a minimal image.

### Kubernetes

Manifests are provided in the `k8s/` directory:

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

---

## Learning Goals

This project was used to practice and explore:

- Structuring a Go application with clean architecture
- Using Fiber as an alternative to standard `net/http`
- Implementing JWT authentication from scratch
- Writing raw SQL with parameterized queries (no ORM)
- Building an automatic migration runner
- Containerizing a Go service with Docker
- Writing unit tests for service-layer logic
- Setting up a CI/CD pipeline with GitHub Actions

---

## License

This is a personal practice project. No license is applied.
