#!make
include .env
export $(shell sed 's/=.*//' .env)

run:
	go run ./cmd/main.go

build:
	go build -o ./bin/openfinance-mcp-server ./cmd/main.go

start:
	./bin/openfinance-mcp-server
