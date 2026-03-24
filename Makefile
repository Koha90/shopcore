APP_NAME := github.com/koha90/shopcore
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=$(DB_SSLMODE)

include .env
export

.PHONY: help up down restart logs ps psql migrate run-tui test test-manager test-botconfig-postgres fmt lint

help:
	@echo "Available targets:"
	@echo "  up                    - start postgres container"
	@echo "  down                  - stop postgres container"
	@echo "  restart               - restart postgres container"
	@echo "  logs                  - show postgres logs"
	@echo "  ps                    - show docker compose status"
	@echo "  psql                  - open psql shell"
	@echo "  migrate               - run app migrations"
	@echo "  run-tui               - run TUI app"
	@echo "  test                  - run all tests"
	@echo "  test-manager          - run manager tests"
	@echo "  test-botconfig-postgres - run postgres repo tests"
	@echo "  fmt                   - run gofmt"
	@echo "  lint                  - run go test build check"

up:
	docker compose up -d postgres

down:
	docker compose down

restart: down up

logs:
	docker compose logs -f postgres

ps:
	docker compose ps

psql:
	docker compose exec postgres psql -U $(DB_USER) -d $(DB_DATABASE)

migrate:
	go run ./cmd/migrate

run-tui:
	go run ./cmd/tui

test:
	go test ./...

test-manager:
	go test -cover ./internal/manager

test-botconfig-postgres:
	go test -cover ./internal/botconfig/postgres -v

fmt:
	gofmt -w $$(find . -type f -name '*.go')

lint:
	go test ./... >/dev/null

dev:
	docker compose up -d postgres
	go run ./cmd/migrate
	go run ./cmd/tui
