package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mark3labs/mcp-go/server"

	"github.com/thunderjr/openfinance-mcp-server/internal/infra/logger"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/pluggy"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/redis"

	"github.com/thunderjr/openfinance-mcp-server/internal/mcp/resources"
	"github.com/thunderjr/openfinance-mcp-server/internal/mcp/tools"
)

func main() {
	if err := logger.Init("openfinance-mcp.log"); err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("Starting OpenFinance MCP Server")

	pluggyClient := pluggy.NewClient(pluggy.NewAuth(redis.Instance()))

	s := server.NewMCPServer(
		"OpenFinance MCP",
		"0.0.1",
		server.WithToolCapabilities(true),
	)

	// Register Resources
	resourceRegistry := mcp.NewResourceRegistry(
		resources.NewGetTransactionsResource(pluggyClient),
		resources.NewGetInvestmentsResource(pluggyClient),
		resources.NewGetAccountsResource(pluggyClient),
		resources.NewGetAccountResource(pluggyClient),
		resources.NewGetItemResource(pluggyClient),
		resources.NewGetBillsResource(pluggyClient),
		resources.NewGetBillResource(pluggyClient),
	)

	resourceRegistry.Register(s)

	// Register Tools
	toolRegistry := mcp.NewToolRegistry(
		tools.NewPluggyApiKeyTool(pluggyClient),
		tools.NewPluggyConnectTokenTool(pluggyClient),
		tools.NewPluggyWaitItemUpdatedTool(pluggyClient),
	)

	toolRegistry.Register(s)

	if err := server.ServeStdio(s); err != nil {
		logger.Errorf("Server error: %v", err)
	}
}
