# OpenFinance MCP Server

The OpenFinance MCP Server is a Go-based application that connects to the Brazilian Open Finance ecosystem through [Pluggy.ai](https://www.pluggy.ai/en), providing secure and structured access to financial data such as account balances, transactions, and investments. It is designed to expose these resources to Large Language Models (LLMs) and other AI agents via the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction), facilitating the development of intelligent financial applications.

*Note: Ensure to grant the necessary consents via [Pluggy Connect](https://meu.pluggy.ai/en) before accessing the endpoints.*

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
    "golang-mcp-server": {
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


## 🤝 Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.
