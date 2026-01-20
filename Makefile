run:
	go run ./main.go

build:
	go build -o ./bin/openfinance-mcp-server ./main.go

start:
	./bin/openfinance-mcp-server
