# gofiber-hax

Hexagonal (Ports & Adapters) starter for GoFiber with manual DI, multi-DB (Mongo/MySQL), Cobra CLI, and flexible logging + auth.

## Highlights
- Hexagonal structure (core isolated from adapters)
- Manual DI (composition root in one place)
- Multi-DB support (Mongo + MySQL, optional MySQL replica)
- Versioned API routes (v1/v2 handlers separated)
- Request validation and centralized error handling
- Production hardening: request ID, recover, timeout, rate limit, access log
- Auth modes: static token, internal JWT, external JWKS, or Google ID token

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
go mod tidy
go run . migrate
go run . start
```

## Environment
`.env` is loaded automatically via `godotenv`.
- In `APP_ENV=prod|production` it **does not** override existing system env.
- In other environments it **does** override system env (dev-friendly).

Minimal example:
```env
APP_NAME=gofiber-hax
APP_ENV=dev

FIBER_ADDR=5001
FIBER_HOST=127.0.0.1
FIBER_ALLOW_ORIGINS=*
FIBER_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization
FIBER_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS

# DB: auto | mongo | mysql | both
DB_DRIVER=auto

# Mongo
MONGO_URI=mongodb://localhost:27017
MONGO_DB=example_db
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_USER=
MONGO_PASS=

# MySQL
MYSQL_DSN=root:password@tcp(127.0.0.1:3306)/app?parseTime=true
MYSQL_REPLICA_DSN=
MYSQL_AUTO_MIGRATE=false

# Logging
LOG_LEVEL=info
LOG_FORMAT=json   # json | pretty | text | logfmt | line

# HTTP access log (Common Log Format)
HTTP_ACCESS_LOG_FORMAT=${ip} - - [${time}] "${method} ${url} ${protocol}" ${status} ${bytesSent}
HTTP_ACCESS_LOG_TIME_FORMAT=02/Jan/2006:15:04:05 -0700

# Auth (protected routes)
AUTH_ENABLED=true
AUTH_MODE=jwt   # token | jwt | jwks | google
AUTH_HEADER=Authorization
AUTH_SCHEME=Bearer

# Static token mode
AUTH_TOKEN=

# Internal JWT mode
JWT_ALG=RS256
JWT_PRIVATE_KEY_PATH=keys/jwt_private.pem
JWT_KEY_ID=gofiber-hax-key-1
JWT_ISSUER=gofiber-hax
JWT_AUDIENCE=gofiber-hax
JWT_JWKS_TTL_SEC=3600
JWT_CLOCK_SKEW_SEC=60
JWT_ACCESS_TTL_SEC=900
JWT_REFRESH_TTL_SEC=604800

# External JWKS / Google
JWT_JWKS_URL=
GOOGLE_CLIENT_ID=

# Production recommendation
MYSQL_AUTO_MIGRATE=false
```

## Routes
Base prefix: `/api/{version}`
- `GET /api/v1/health`
- `GET /api/v1/ready`
- `POST /api/v1/auth/register` (creates user via bcrypt)
- `POST /api/v1/auth/login` (issues RS256 JWT)
- `GET /api/v1/users/:account_id` (protected)

Identity Provider (internal JWT mode only):
- `GET /api/.well-known/jwks.json`

Versioned handlers are separated in code (`V1`, `V2`) to allow different behavior per version.

## Auth Modes
Static token:
```
AUTH_ENABLED=true
AUTH_MODE=token
AUTH_TOKEN=secret123
```
Request example:
```
Authorization: Bearer secret123
```

Internal JWT (recommended default):
```env
AUTH_ENABLED=true
AUTH_MODE=jwt
JWT_ALG=RS256
JWT_PRIVATE_KEY_PATH=keys/jwt_private.pem
JWT_KEY_ID=gofiber-hax-key-1
JWT_ISSUER=gofiber-hax
JWT_AUDIENCE=gofiber-hax
```
Behavior:
- `POST /api/v1/auth/login` issues RS256 JWT
- protected middleware validates that JWT with the local public key
- `GET /api/.well-known/jwks.json` exposes the public key set for downstream consumers

Note:
- in `dev`, if `JWT_PRIVATE_KEY_PATH` does not exist, the key is generated automatically
- in `prod`, the key must already exist before startup
- this template uses a `single active key` model for simplicity
- the key is loaded once at startup; changing the key file requires an app restart
- if you replace the key and restart, previously issued JWTs may become invalid immediately

External JWKS:
```env
AUTH_ENABLED=true
AUTH_MODE=jwks
JWT_ISSUER=https://issuer.example.com
JWT_AUDIENCE=my-api
JWT_JWKS_URL=https://issuer.example.com/.well-known/jwks.json
```

Google ID Token:
```
AUTH_ENABLED=true
AUTH_MODE=google
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
```

## Logging
- `LOG_FORMAT=pretty` for multi-line output (debug-friendly)
- `LOG_FORMAT=text|logfmt|line` for single-line console
- `LOG_FORMAT=json` for log systems
- HTTP access log uses Common Log Format (configure `HTTP_ACCESS_LOG_FORMAT` and `HTTP_ACCESS_LOG_TIME_FORMAT`)
- See `internal/infra/logs/README.md` for usage and formats

## Production Key Management
- `AUTH_MODE=jwt` on production expects an existing private key file before the app starts.
- Do not rely on runtime key generation in production.
- Keep the private key outside the image, managed by DevOps as a secret.
- Mount or inject the key file into the runtime, then set `JWT_PRIVATE_KEY_PATH` to that path.
- If a deploy replaces the key file with a different key, previously issued JWTs may stop working after restart.
- For this template, key management intentionally stays at `single active key` to keep the auth flow simple.

## Notes
- Run `go run . migrate` to apply MySQL migrations and ensure Mongo indexes.
- `MYSQL_AUTO_MIGRATE=true` is allowed for dev only. In production it is rejected at config load.
- When `DB_DRIVER=both`, the User repo currently prefers MySQL.
- `AUTH_MODE=jwt` is the only mode that wires internal `/auth/register` and `/auth/login`.
- `AUTH_MODE=jwks` and `AUTH_MODE=google` are validation-only modes.
- `AUTH_MODE=jwt` intentionally stays on `single active key` for this template to keep the auth flow clean and maintainable.

## Adding New Feature (example flow)
```
handler -> service -> repo (port out) -> adapter
```
- Define domain
- Add repo interface in `core/ports/out`
- Implement repo in `adapters/db/*`
- Add service in `core/service`
- Wire in `internal/app/app.go`
- Register routes in `adapters/http/routes`

## License
MIT (or update as needed)
