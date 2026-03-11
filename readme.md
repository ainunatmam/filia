A production-ready Go REST API boilerplate using Clean Architecture pattern with Fiber framework.

## Features

- **Clean Architecture** - Handler → Service → Repository pattern
- **Fiber Framework** - Fast HTTP server with middleware support
- **MySQL + Goqu** - Database with type-safe query builder
- **JWT Authentication** - Secure token-based auth with algorithm validation
- **Rate Limiting** - Configurable request rate limiting
- **Graceful Shutdown** - Proper cleanup of resources on exit
- **Database Migrations** - Using Goose for versioned migrations
- **Structured Logging** - Zerolog with configurable output format
- **Input Sanitization** - XSS protection middleware
- **Health Checks** - Liveness and readiness endpoints

---

## Requirements

- Go 1.21+
- MySQL 8.0+

---

## Quick Start

### 1. Clone and Setup

```bash
# Copy environment file
cp .env.example .env

# Edit configuration
nano .env
```

### 2. Run Migrations

```bash
go run main.go migrate up
```

### 3. Start Server

```bash
go run main.go
```

Server runs at `http://localhost:3000` by default.

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| **Database** |
| `DATABASE_USERNAME` | `root` | MySQL username |
| `DATABASE_PASSWORD` | `` | MySQL password |
| `DATABASE_HOST` | `127.0.0.1` | MySQL host |
| `DATABASE_PORT` | `3306` | MySQL port |
| `DATABASE_NAME` | `wallet_api` | Database name |
| `DB_MAX_OPEN_CONNS` | `100` | Max open connections |
| `DB_MAX_IDLE_CONNS` | `25` | Max idle connections |
| `DB_CONN_MAX_LIFETIME` | `300` | Connection max lifetime (seconds) |
| `DB_CONN_MAX_IDLE_TIME` | `60` | Connection max idle time (seconds) |
| `DB_QUERY_TIMEOUT` | `30` | Query timeout (seconds) |
| **Application** |
| `APP_PORT` | `:3000` | Server port |
| `APP_ENV` | `development` | Environment (development/staging/production) |
| **JWT** |
| `JWT_SECRET` | - | **Required in production** |
| `JWT_EXPIRATION_HOURS` | `24` | Token expiration time |
| **Rate Limiting** |
| `RATE_LIMIT_MAX` | `100` | Max requests per window |
| `RATE_LIMIT_EXPIRATION` | `60` | Window duration (seconds) |
| **Logging** |
| `LOG_LEVEL` | `info` | Log level (debug/info/warn/error) |
| `LOG_FORMAT` | `console` | Output format (console/json) |

---

## Project Structure

```
wallet-api/
├── main.go                 # Application entry point
├── bootstrap/
│   ├── app.go              # Application bootstrap & middleware setup
│   └── database.go         # Database connection
├── config/
│   └── config.go           # Configuration loading & validation
├── app/
│   ├── ctx/
│   │   └── ctx.go          # Context utilities (timeout, request ID)
│   ├── entity/
│   │   └── example.go      # Domain entities
│   ├── errors/
│   │   └── errors.go       # Custom error types
│   ├── handlers/
│   │   └── example/        # HTTP handlers (controllers)
│   ├── services/
│   │   └── example/        # Business logic
│   ├── repositories/
│   │   ├── contract.go     # Repository interfaces
│   │   └── example.go      # Repository implementations
│   ├── presentation/
│   │   ├── request.go      # Request DTOs
│   │   └── response.go     # Response DTOs
│   ├── middleware/
│   │   ├── auth.go         # JWT authentication
│   │   ├── error_handler.go# Centralized error handling
│   │   ├── ratelimit.go    # Rate limiting
│   │   └── sanitizer.go    # Input sanitization
│   ├── libraries/
│   │   ├── goqu.go         # Goqu database wrapper
│   │   ├── logger.go       # Zerolog configuration
│   │   └── transaction_manager.go
│   ├── routes/
│   │   ├── route.go        # Router setup
│   │   └── api.go          # API route definitions
│   └── utilities/
│       └── *.go            # Helper functions
└── database/
    └── migration/          # SQL migration files
```

---

## Architecture

This boilerplate follows **Clean Architecture** principles:

```
┌─────────────────────────────────────────────────┐
│                   Handlers                       │  ← HTTP layer
│              (fiber.Ctx → Response)              │
├─────────────────────────────────────────────────┤
│                   Services                       │  ← Business logic
│           (context.Context → error)              │
├─────────────────────────────────────────────────┤
│                 Repositories                     │  ← Data access
│            (context.Context → entity)            │
├─────────────────────────────────────────────────┤
│                   Database                       │  ← Infrastructure
└─────────────────────────────────────────────────┘
```

**Key principles:**
- Services use `context.Context`, NOT `fiber.Ctx` (decoupled from HTTP)
- Repositories handle all database operations
- Errors are wrapped with context for debugging
- Internal errors are never exposed to clients

---

## API Endpoints

### Health Checks

```bash
# Liveness check
curl http://localhost:3000/health

# Readiness check (includes DB)
curl http://localhost:3000/ready
```

### Example Endpoint

```bash
# Create example
curl -X POST http://localhost:3000/api/example \
  -H "Content-Type: application/json" \
  -d '{"name": "Test"}'
```

---

## Database Migrations

```bash
# Run all pending migrations
go run main.go migrate up

# Rollback last migration
go run main.go migrate down

# Check migration status
go run main.go migrate status

# Create new migration
go run main.go migrate create <migration_name>
```

---

## Development

### Adding a New Feature

1. **Create Entity** in `app/entity/`
2. **Create Migration** with `go run main.go migrate create <name>`
3. **Create Repository** interface in `app/repositories/contract.go`
4. **Implement Repository** in `app/repositories/<name>.go`
5. **Create Service** in `app/services/<name>/`
6. **Create Handler** in `app/handlers/<name>/`
7. **Register Routes** in `app/routes/api.go`

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run specific package tests
go test ./app/services/example/... -v
```

---

## Security Features

| Feature | Implementation |
|---------|----------------|
| JWT Algorithm Validation | Prevents algorithm confusion attacks |
| Error Message Sanitization | Internal errors not exposed to clients |
| Input Sanitization | XSS protection on all inputs |
| Rate Limiting | Configurable per-IP limits |
| Request Body Limit | 4MB max payload size |
| Security Headers | Helmet middleware (CSP, etc.) |
| CORS | Configurable cross-origin policy |

---

## Production Checklist

- [ ] Change `JWT_SECRET` to a strong random value
- [ ] Set `APP_ENV=production`
- [ ] Set `LOG_FORMAT=json` for structured logging
- [ ] Configure proper database credentials
- [ ] Set up TLS/HTTPS (via reverse proxy)
- [ ] Configure rate limits appropriately
- [ ] Set up monitoring and alerting

---

## License

MIT
