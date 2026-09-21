BINARY=bin/api
MAIN=./cmd/api

.PHONY: run build test migrate-up migrate-down tidy vet

run:
	go run $(MAIN)/main.go

build:
	go build -o $(BINARY) $(MAIN)/main.go

test:
	go test ./...

tidy:
	go mod tidy

vet:
	go vet ./...

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down
