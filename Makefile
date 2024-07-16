ENV ?=development

include .env.$(ENV)
export $(shell sed 's/=.*//' .env.$(ENV))

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

sqlc:
	@DB_URI="$(DB_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" sqlc generate
test:
	@go test -v -cover ./...
dev:
	@GIN_MODE="debug" godotenv -f .env.development air
prod:
	@go build -o ./bin/app 
	@GIN_MODE="release" godotenv -f .env.production ./bin/prod
install:
	@go get -u

.PHONY: migrate-up migrate-down migrate-force sqlc create-db test dev prod install
