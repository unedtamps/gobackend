gen:
	@sqlc generate

test:
	@go test -v ./test

build:
	@go build -o ./bin/app 

start:
	@GIN_MODE="release" ./bin/app

dev:
	@GIN_MODE="debug" air

install:
	@go get -u cmd/api && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: gen test build start dev install
