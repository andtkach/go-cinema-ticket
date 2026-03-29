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

## Ports

| Service | Host port |
|---|---|
| Server | 17080 |
| Client | 17070 |
| Redis | 16379 |
| Postgres | 15432 |
| Redis Commander | 16378 |
