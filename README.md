# go-cinema-ticket

Tickets booking application.

## Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose

## Getting started

Copy the environment file and adjust credentials if needed:

```bash
cp .env.example .env
```

Generate the TLS certificate once (required for nginx):

```bash
make nginx-certs
```

Then start everything with a single command:

```bash
make run-all
```

This starts the infrastructure (Docker), builds and publishes the client, starts the server, and opens [https://localhost:17091/app/](https://localhost:17091/app/) in the browser.

## Infrastructure

| Command | Description |
|---|---|
| `make infra-up` | Start Redis and Postgres in the background |
| `make infra-down` | Stop and remove containers |
| `make infra-restart` | Restart all containers |
| `make infra-logs` | Tail logs from all containers |
| `make infra-ps` | Show container status |

## Server

| Command | Description |
|---|---|
| `make build-server` | Compile the server binary to `server/main` |
| `make run-server` | Start infrastructure and run the server |
| `make test-server` | Start infrastructure and run all tests |
| `make tidy-server` | Run `go mod tidy` to sync dependencies |
| `make clean-server` | Remove the compiled binary |

The server starts on [http://localhost:17080](http://localhost:17080).
Redis Commander (Redis UI) is available at [http://localhost:16378](http://localhost:16378).

## Client

A React (Vite) web app that talks to the server via a dev proxy.

| Command | Description |
|---|---|
| `make install-client` | Install npm dependencies |
| `make dev-client` | Start the dev server on [http://localhost:17070](http://localhost:17070) |
| `make build-client` | Build production assets to `client/dist/` |
| `make publish-client` | Build and copy assets to `server/static/` for the server to host |

## Nginx gateway

An nginx container acts as an HTTPS reverse proxy and gateway, routing traffic by path prefix to backend services.

**First-time setup** — generate the self-signed TLS certificate (required before starting the container):

```bash
make nginx-certs
```

This writes `config/nginx/certs/server.key` and `config/nginx/certs/server.crt` (valid 365 days, `CN=localhost`).

| Command | Description |
|---|---|
| `make nginx-certs` | Generate a self-signed TLS certificate |
| `make nginx-reload` | Reload nginx config without downtime |

**Routes** (all via `https://localhost:17091`):

| Path | Upstream |
|---|---|
| `/app/` | Cinema server → `localhost:17080` |
| `/idp/` | Identity provider → `localhost:17090` (future) |
| `/health` | nginx health check — returns `200 ok` |

> Browsers will warn about the self-signed certificate. Accept the exception or add `config/nginx/certs/server.crt` to your system trust store.

## Ports

| Service | Host port | Protocol |
|---|---|---|
| Nginx gateway | 17091 | HTTPS |
| Server | 17080 | HTTP |
| Client (dev) | 17070 | HTTP |
| Redis | 16379 | TCP |
| Postgres | 15432 | TCP |
| Redis Commander | 16378 | HTTP |
