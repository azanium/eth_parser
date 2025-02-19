.PHONY: run test build

run:
	go run cmd/main.go

test:
	go test -race ./...

build:
	go build -o bin/parser cmd/main.go
