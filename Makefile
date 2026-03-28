.PHONY: infra-up infra-down infra-restart infra-logs infra-ps build run test tidy

# Infrastructure
infra-up:
	mkdir -p data/redis data/postgres
	docker compose up -d

infra-down:
	docker compose down

infra-restart:
	docker compose restart

infra-logs:
	docker compose logs -f

infra-ps:
	docker compose ps

# Server
build:
	cd server && go build -o main ./cmd

run: infra-up
	cd server && go run ./cmd

test: infra-up
	cd server && go test -v -race ./...

tidy:
	cd server && go mod tidy

clean:
	cd server && rm -f main
