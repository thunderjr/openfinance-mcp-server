# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

OpenFinance MCP Server is a Go application that exposes Brazilian Open Finance data (via Pluggy.ai) to LLMs through the Model Context Protocol (MCP). It provides tools for querying bank accounts, transactions, investments, and credit card bills.

## Build Commands

```bash
make run      # Run the server directly
make build    # Build binary to ./bin/openfinance-mcp-server
make start    # Run the built binary
```

## Architecture

### Entry Point
- `main.go` - Initializes logger, Pluggy client, registers MCP tools, starts stdio server
- Embeds `.env` file at compile time via `//go:embed`

### Configuration Files (in `~/.config/finops/`)
- `pluggy.json` - Pluggy API credentials (`PLUGGY_CLIENT_ID`, `PLUGGY_CLIENT_SECRET`)
- `accounts.json` - Saved bank connections for `list_connections` tool
- `openfinance-mcp.log` - Application logs

### Core Packages

**`internal/mcp/tools/`** - MCP tool implementations
- Each tool implements `ToolProvider` interface: `Name()`, `Description()`, `Handle()`
- Tools use financial domain language (e.g., `list_accounts`, `connection_id` not `item_id`)
- Tool descriptions follow template: what it does, when to use it, what it returns

**`internal/provider/pluggy/`** - Pluggy.ai API client
- `factory.go` - Client initialization with auto-auth on startup
- `auth.go` - Credentials loading from config file, API key caching via Redis
- Domain files (`account.go`, `transaction.go`, etc.) - API endpoint implementations

**`internal/provider/mcp/`** - MCP framework abstractions
- `tool-registry.go` - `ToolRegistry` that registers all tools with the MCP server

**`internal/provider/redis/`** - Redis client for caching API keys and connect tokens

### Adding New Tools

1. Create tool struct implementing `ToolProvider` interface
2. Define args struct with `json` and `jsonschema` tags
3. Register in `main.go` via `mcp.NewToolRegistry(...)`

Example pattern from `accounts.go`:
```go
type ListAccountsTool struct {
    client *pluggy.Client
}

func (t *ListAccountsTool) Name() string { return "list_accounts" }
func (t *ListAccountsTool) Description() string { return `...` }
func (t *ListAccountsTool) Handle() internalMcp.ToolHandlerFunc {
    return t.handleListAccounts
}
```

### Dependencies
- MCP server: `github.com/metoro-io/mcp-golang`
- Redis: `github.com/redis/go-redis/v9`
- Requires running Redis instance for token caching
