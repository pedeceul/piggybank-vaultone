.PHONY: dev build docker-build up down logs test lint gen tidy buf

dev:
	- APP_ENV=dev go run ./cmd/api

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

docker-build:
	docker build -t vaultone/api:dev .

up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f api

test:
	go test ./...

lint:
	@echo "lint placeholder (golangci-lint to be added in Phase 5)"

gen:
	buf generate

buf:
	buf lint

e2e:
	bash scripts/demo.sh

