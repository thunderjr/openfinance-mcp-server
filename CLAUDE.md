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

### Configuration & Data Files
- `~/.config/finops/pluggy.json` - Legacy plaintext credentials. Read **once** on first run to seed the encrypted store, then no longer needed (the encrypted DB becomes authoritative). Can also be supplied via `PLUGGY_CLIENT_ID`/`PLUGGY_CLIENT_SECRET` env.
- `~/.config/finops/accounts.json` - **Legacy/optional.** No longer required. Connections are now auto-registered in the encrypted store (`connections:items` key) the first time any tool resolves an item successfully; `list_connections`/`refresh_connections` read from there and fetch names/accounts live from Pluggy. If this file exists it is imported into the store once on startup (never modified/removed).
- `~/.config/finops/openfinance-mcp.log` - Application logs
- `~/.openfinance/pluggy.db` - Encrypted SQLite store (mode 0700 dir). Holds the Pluggy credentials, the cached API key (2h TTL), and connect tokens (30m TTL). All secret values are sealed with AES-256-GCM before hitting disk.

### Encryption Key
The AES-256 master key is resolved in order: `$OPENFINANCE_DB_KEY` (base64) → OS keychain item `openfinance-mcp-server/db-encryption-key` (auto-generated on first run) → `~/.openfinance/master.key` (0600 fallback). The key never lives inside the DB file.

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

**`internal/provider/sqlite/`** - Encrypted key/value store (pure-Go `modernc.org/sqlite`, no CGO) backing auth credentials and cached tokens with optional TTL. Replaces the old Redis cache.

**`internal/provider/secret/`** - Master-key resolution (env → OS keychain → keyfile) and AES-256-GCM seal/open used by the sqlite store.

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
- SQLite (pure Go, no CGO): `modernc.org/sqlite`
- OS keychain: `github.com/zalando/go-keyring`
- No external services required (Redis removed)

### Calling the server with mcporter
Registered as `openfinance` in `~/.mcporter/mcporter.json` (stdio → built binary).
```bash
npx mcporter call openfinance.list_connections
npx mcporter list openfinance --schema   # show all tool signatures
```
