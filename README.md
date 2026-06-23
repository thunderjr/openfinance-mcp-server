# OpenFinance MCP Server

The OpenFinance MCP Server is a Go-based application that connects to the Brazilian Open Finance ecosystem through [Pluggy.ai](https://www.pluggy.ai/en), providing secure and structured access to financial data such as account balances, transactions, and investments. It is designed to expose these resources to Large Language Models (LLMs) and other AI agents via the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction), facilitating the development of intelligent financial applications.

*Note: Ensure to grant the necessary consents via [Pluggy Connect](https://meu.pluggy.ai/en) before using.*

### Installation

1. Clone the repository:
```bash
git clone https://github.com/thunderjr/openfinance-mcp-server.git
cd openfinance-mcp-server
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
make build
```

4. Use binary on your preferred MCP client (e.g. Claude Desktop):
```json
{
  "mcpServers": {
    "openfinance-mcp-server": {
      "command": "<project_path>/bin/openfinance-mcp-server",
      "args": [],
      "env": {
        "PLUGGY_CLIENT_ID": "${PLUGGY_CLIENT_ID}",
        "PLUGGY_CLIENT_SECRET": "${PLUGGY_CLIENT_SECRET}"
      }
    }
  }
}
```

### Credential storage

On first run the server reads `PLUGGY_CLIENT_ID`/`PLUGGY_CLIENT_SECRET` (env or
`~/.config/finops/pluggy.json`) **once** and seeds them into an encrypted SQLite
database at `~/.openfinance/pluggy.db`. From then on the database is the source
of truth, and all secrets (credentials, cached API key, connect tokens) are
sealed at rest with AES-256-GCM.

The AES master key is resolved as: `$OPENFINANCE_DB_KEY` (base64) → OS keychain
(auto-generated on first run) → `~/.openfinance/master.key` (0600 fallback). No
external services are required.

### Calling tools with mcporter

```bash
# Register once (stdio → built binary)
npx mcporter config add openfinance \
  --command "<project_path>/bin/openfinance-mcp-server" --scope home

npx mcporter list openfinance --schema      # show tool signatures
npx mcporter call openfinance.list_connections
```


## 🤝 Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.
