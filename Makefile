.PHONY: infra-up infra-down infra-restart infra-logs infra-ps build-server run-server test-server tidy-server install-client dev-client build-client publish-client

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
build-server:
	cd server && go build -o main ./cmd

run-server: infra-up
	cd server && go run ./cmd

test-server: infra-up
	cd server && go test -v -race ./...

tidy-server:
	cd server && go mod tidy

clean-server:
	cd server && rm -f main

# Client
install-client:
	cd client && npm install

dev-client:
	cd client && npm run dev

build-client:
	cd client && npm run build

publish-client:
	cd client && npm run build
	rm -rf server/static
	cp -r client/dist server/static
