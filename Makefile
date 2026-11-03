.PHONY: dev build test vet migrate

dev:
	go run ./cmd/jobtrack

build:
	go build -o bin/jobtrack ./cmd/jobtrack

test:
	go test ./...

vet:
	go vet ./...

migrate:
	goose -dir migrations sqlite3 ./data/jobtrack.db up
