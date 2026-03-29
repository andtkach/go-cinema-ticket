# go-cinema-ticket

Tickets booking application.

## Prerequisites

- [Go 1.21+](https://golang.org/dl/)
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
| `make build` | Compile the server binary to `server/main` |
| `make run` | Start infrastructure and run the server |
| `make test` | Start infrastructure and run all tests |
| `make tidy` | Run `go mod tidy` to sync dependencies |
| `make clean` | Remove the compiled binary |

The server starts on [http://localhost:17080](http://localhost:17080).
Redis Commander (Redis UI) is available at [http://localhost:16378](http://localhost:16378).

## Ports

| Service | Host port |
|---|---|
| Server | 17080 |
| Redis | 16379 |
| Postgres | 15432 |
| Redis Commander | 16378 |
