include .env

migrate-up:
	@migrate -path internal/migration -database "$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

migrate-down:
	@migrate -path internal/migration -database "$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down

migrate-force:
	@read -p "Enter migration version: " version; \
	migrate -path internal/migration -database "$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose force $$version

create-migrate:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/migration -seq $$name

setup:
	@godotenv -f .env docker compose up -d
sqlc:
	@DB_URI="$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" sqlc generate
test:
	@go test -v -cover ./...
build:
	@go build -o ./bin/app 
dev:
	@GIN_MODE="debug" godotenv -f .env.development air
start:
	@GIN_MODE="release" godotenv -f .env ./bin/app
install:
	@go get -u

.PHONY: migrate-up migrate-down migrate-force sqlc create-db test dev prod install
