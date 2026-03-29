.PHONY: infra-up infra-down infra-restart infra-logs infra-ps build-server run-server test-server tidy-server install-client dev-client build-client publish-client nginx-certs nginx-reload run-all

# All-in-one
run-all: infra-up publish-client
	cd server && env $$(grep -v '^#' ../.env | xargs) go run ./cmd &
	@sleep 2 && xdg-open https://localhost:17091/app/ &

# Infrastructure
infra-up:
	mkdir -p data/redis-app data/postgres-app data/postgres-idp data/server-idp data/worker-idp data/certs data/custom-templates blueprints
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
	cd server && env $$(grep -v '^#' ../.env | xargs) go run ./cmd

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
	cd client && npx vite build --base=/app/
	rm -rf server/static
	cp -r client/dist server/static

# Nginx
nginx-certs:
	openssl req -x509 -nodes -newkey rsa:2048 -days 365 \
		-keyout config/nginx/certs/server.key \
		-out config/nginx/certs/server.crt \
		-subj "/CN=localhost" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

nginx-reload:
	docker exec nginx-app nginx -s reload
