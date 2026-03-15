# Logging Guide

This package provides a small logging wrapper around `slog` + `zap` with a consistent format for application logs and errors. It is designed to be simple for day-to-day use and easy to read in local development.

**Quick Usage**
```go
import "gofiber-hax/internal/infra/logs"

logs.Debug("debug payload", "user_id", id)
logs.Info("user created", "user_id", id)
logs.Warn("cache miss", "key", key)
logs.Error(err)
```

**What You Get**
- `logs.Debug/Info/Warn`: single-line log entries with source location.
- `logs.Error`: when `LOG_FORMAT=pretty`, prints a clean, multi-line error block with `Message` and `Trace`.
- Optional HTTP access log (Common Log Format) via `HTTP_ACCESS_LOG=true`.

**Environment Config**
```
LOG_LEVEL=debug      # debug | info | warn | error
LOG_FORMAT=pretty    # pretty | text | json | logfmt | line

HTTP_ACCESS_LOG=false
HTTP_ACCESS_LOG_FORMAT=${ip} - - [${time}] "${method} ${url} ${protocol}" ${status} ${bytesSent}
HTTP_ACCESS_LOG_TIME_FORMAT=02/Jan/2006:15:04:05 -0700
```

**Formats**
- `pretty`: colored, dev-friendly output. `logs.Error` becomes multi-line.
- `text/logfmt/line`: one line per log entry, suitable for terminals or simple log viewers.
- `json`: machine-friendly format for ELK/Datadog/etc.

**Pretty Error Output (example)**
```
2026-03-15T17:20:10.557+0700	ERROR	handlers/auth.handler.go:33
Message: invalid value, should be pointer to struct or slice
Trace: authservice.register -> userservice.create -> mysql.userrepo.create
```

**How Trace Is Built**
`logs.Error` parses error strings using the separator `" error: "` to build a chain.
Example:
```
authservice.register error: userservice.create error: mysql.userrepo.create error: invalid value
```

**Notes**
- Source location is derived from the call site (no extra logger is required).
- Access log is a separate middleware; when enabled, it replaces the default request log.
- If you need custom fields, pass key/value pairs to any log call.
