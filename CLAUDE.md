# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All commands run from the `server/` directory.

```bash
# Start Redis (required before running server or tests)
docker-compose up -d

# Run the server (port 8080)
go run ./cmd

# Build binary
go build -o main ./cmd

# Run all tests (requires Redis on localhost:6379)
go test ./...

# Run tests with verbose output
go test -v ./internal/booking

# Run a single test
go test -v ./internal/booking -run TestConcurrentBooking
```

## Architecture

The module is `github.com/andtkach/cinema` and lives under `server/`.

**Request flow:** HTTP handler (`handler.go`) → Service (`service.go`) → BookingStore interface (`domain.go`) → RedisStore (`redis_store.go`)

### Booking domain (`server/internal/booking/`)

The core domain models a two-phase seat reservation:
1. **Hold** — a temporary reservation with a Redis TTL (2 minutes). Only one user wins via `SET NX` atomicity.
2. **Confirm** — makes the hold permanent by removing the TTL with `PERSIST`.
3. **Release** — deletes both keys to free the seat.

**Redis key design:**
- `seat:{movieID}:{seatID}` → JSON-encoded `Booking` (TTL = held; no TTL = confirmed)
- `session:{sessionID}` → reverse-lookup key pointing back to the seat key

`memory_store.go` and `concurrent_store.go` are alternative store implementations not used in production — `redis_store.go` is the active implementation instantiated in `cmd/main.go`.

### HTTP routes (`server/cmd/main.go`)

```
GET    /                                     → static/index.html
GET    /movies                               → hardcoded list (inception, dune)
GET    /movies/{movieID}/seats               → list seat statuses
POST   /movies/{movieID}/seats/{seatID}/hold → hold a seat, returns session_id
PUT    /sessions/{sessionID}/confirm         → confirm hold
DELETE /sessions/{sessionID}                 → release hold
```

Movies and seat grids are hardcoded in `main.go` (inception: 5×8, dune: 4×6).

### Infrastructure

Docker Compose provides two services:
- **redis** on port 6379
- **redis-commander** on port 8081 (web UI for inspecting Redis data)
