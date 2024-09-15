migrate-up:
	@migrate -path internal/migration -database "$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

dev-migrate-up:
	@godotenv -f .env.dev -- make migrate-up

migrate-down:
	@migrate -path internal/migration -database "$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down

dev-migrate-down:
	@godotenv -f .env.dev -- make migrate-down

migrate-force:
	@read -p "Enter migration version: " version; \
	migrate -path internal/migration -database "$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose force $$version

dev-migrate-force:
	@godotenv -f .env.dev -- make migrate-force

create-migrate:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/migration -seq $$name

sqlc:
	@godotenv -f .env.dev DB_URI="$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" sqlc generate

test:
	@go test -v ./test

dev-test:
	@godotenv -f .env.test go test -v ./test

build:
	@go build -o ./bin/app 

start:
	@GIN_MODE="release" ./bin/app

dev:
	@GIN_MODE="debug" godotenv -f .env.dev air

install:
	@go get -u

.PHONY: migrate-up migrate-down migrate-force-prod sqlc sqlc test dev prod install
