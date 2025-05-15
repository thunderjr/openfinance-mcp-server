package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	server "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"

	"github.com/thunderjr/openfinance-mcp-server/internal/mcp/tools"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/redis"
)

func main() {
	handleErr("logger.Init", logger.Init("openfinance-mcp.log"))
	defer logger.Close()

	pluggyClient := pluggy.NewClient(pluggy.NewAuth(redis.Instance()))

	toolRegistry := mcp.NewToolRegistry(
		tools.NewPluggyApiKeyTool(pluggyClient),
		tools.NewPluggyConnectTokenTool(pluggyClient),
		tools.NewPluggyWaitItemUpdatedTool(pluggyClient),
		tools.NewPluggyAccountsTool(pluggyClient),
		tools.NewPluggyAccountTool(pluggyClient),
		tools.NewPluggyTransactionsTool(pluggyClient),
		tools.NewPluggyInvestmentsTool(pluggyClient),
		tools.NewPluggyItemTool(pluggyClient),
		tools.NewPluggyBillsTool(pluggyClient),
		tools.NewPluggyBillTool(pluggyClient),
	)

	logger.Info("Starting OpenFinance MCP Server")

	server := server.NewServer(stdio.NewStdioServerTransport())
	handleErr("Tools Registration", toolRegistry.Register(server))
	handleErr("Server Startup", server.Serve())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}

func handleErr(prefix string, err error) {
	if err != nil {
		log.Fatalf("[%s] Err: %v", prefix, err)
	}
}
