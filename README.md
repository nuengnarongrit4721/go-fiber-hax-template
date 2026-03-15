# gofiber-hax

Hexagonal (Ports & Adapters) starter for GoFiber with manual DI, multi-DB (Mongo/MySQL), Cobra CLI, and flexible logging + auth.

## Highlights
- Hexagonal structure (core isolated from adapters)
- Manual DI (composition root in one place)
- Multi-DB support (Mongo + MySQL, optional MySQL replica)
- Versioned API routes (v1/v2 handlers separated)
- Pretty logs for dev (multi-line, array-friendly)
- Auth middleware: shared token or JWT/Google ID token

## Project Structure (simplified)
```
cmd/                    # cobra commands
internal/
  app/                  # composition root (DI wiring)
  core/
    domain/             # entities
    ports/              # in/out ports (interfaces)
    service/            # business logic
  adapters/
    http/               # fiber server, routes, handlers, middleware
    db/                 # mongo + mysql repositories
  infra/                # config, logging
  shared/               # common helpers
```

## Quick Start
```bash
# 1) install deps

go mod tidy

# 2) run

go run . serve
```

## Environment
`.env` is loaded automatically via `godotenv`.

Minimal example:
```env
APP_NAME=gofiber-hax
APP_ENV=dev
HTTP_ADDR=:8080

# DB: auto | mongo | mysql | both
DB_DRIVER=auto

# Mongo
MONGO_URI=mongodb://localhost:27017
MONGO_DB=example_db

# MySQL
MYSQL_DSN=root:password@tcp(127.0.0.1:3306)/app?parseTime=true
MYSQL_REPLICA_DSN=
MYSQL_AUTO_MIGRATE=false

# Logging
LOG_LEVEL=debug
LOG_FORMAT=pretty

# Auth (protected routes)
AUTH_ENABLED=false
AUTH_MODE=token   # token | jwt | google
AUTH_TOKEN=
AUTH_HEADER=Authorization
AUTH_SCHEME=Bearer

# JWT / Google OAuth (ID token)
GOOGLE_CLIENT_ID=
JWT_ISSUER=
JWT_AUDIENCE=
JWT_JWKS_URL=
JWT_JWKS_TTL_SEC=3600
JWT_CLOCK_SKEW_SEC=60
```

## Routes
Base prefix: `/api/{version}`

- Public:
  - `GET /api/v1/health`
  - `GET /api/v1/ready`
- Protected:
  - `GET /api/v1/users/:id`

Versioned handlers are separated in code (`V1`, `V2`) to allow different behavior per version.

## Auth Modes
### 1) Simple Token
```
AUTH_ENABLED=true
AUTH_MODE=token
AUTH_TOKEN=secret123
```
Request example:
```
Authorization: Bearer secret123
```

### 2) JWT (generic)
```
AUTH_ENABLED=true
AUTH_MODE=jwt
JWT_ISSUER=https://issuer.example
JWT_AUDIENCE=my-api
JWT_JWKS_URL=https://issuer.example/.well-known/jwks.json
```

### 3) Google ID Token
```
AUTH_ENABLED=true
AUTH_MODE=google
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
```

## Logging
- `LOG_FORMAT=pretty` for multi-line output (debug-friendly)
- See `internal/infra/logs/README.md` for usage and formats

## Notes
- MySQL uses GORM; auto-migration is optional via `MYSQL_AUTO_MIGRATE=true`.
- For real OAuth (access token introspection), add an introspection validator (not included yet).

## Adding New Feature (example flow)
```
handler -> service -> repo (port out) -> adapter
```
- Define domain
- Add repo interface in `core/ports/out`
- Implement repo in `adapters/db/*`
- Add service in `core/service`
- Wire in `internal/app/compose.go`
- Register routes in `adapters/http/routes`

## License
MIT (or update as needed)
# go-fiber-hax-template
