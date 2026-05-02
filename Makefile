.PHONY: help tidy build run test fmt vet lint clean migrate-up migrate-down docker-db

APP_NAME       := chiphouse
CMD_PATH       := ./cmd/rentapp
BIN_DIR        := bin
BIN            := $(BIN_DIR)/$(APP_NAME)

include .env
export

DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

help:
	@echo "ChipHouse — available targets:"
	@echo "  make tidy          — go mod tidy"
	@echo "  make build         — build binary into $(BIN)"
	@echo "  make run           — run the app"
	@echo "  make test          — run all tests"
	@echo "  make fmt           — format the code"
	@echo "  make vet           — go vet"
	@echo "  make lint          — golangci-lint (must be installed)"
	@echo "  make clean         — remove build artifacts"
	@echo "  make migrate-up    — apply ./migration/*.sql to \$$DB_NAME"
	@echo "  make migrate-down  — drop \$$DB_NAME (DESTRUCTIVE)"
	@echo "  make docker-db     — start a local Postgres in Docker"

tidy:
	go mod tidy

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./... -count=1

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BIN_DIR)

migrate-up:
	@for f in migration/*.sql; do \
		echo ">> applying $$f"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f $$f || exit 1; \
	done

migrate-down:
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
		-c "DROP DATABASE IF EXISTS $(DB_NAME);"
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
		-c "CREATE DATABASE $(DB_NAME);"

docker-db:
	docker run --rm -d --name chiphouse-pg \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		postgres:16-alpine
