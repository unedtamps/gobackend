include .env

created-db:
	docker compose up
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
	@sqlc generate
test:
	@ENV="test" go test -v -cover ./...
dev:
	@ENV="dev" GIN_MODE="debug" air
prod:
	@go build -o ./bin/app 
	@ENV="prod" GIN_MODE="release" ./bin/app
install:
	@go get -u


.PHONY: migrate-up migrate-down migrate-force sqlc create-db test dev prod install
