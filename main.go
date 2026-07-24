package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	// autoload loads a .env from the working directory at runtime (dev
	// convenience for seeding a fresh DB). Credentials are no longer embedded
	// into the binary at compile time.
	_ "github.com/joho/godotenv/autoload"
	server "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"

	"github.com/thunderjr/openfinance-mcp-server/internal/mcp/tools"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/connections"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/secret"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/sqlite"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	logDir := filepath.Join(homeDir, ".config", "finops")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Error creating log directory: %v", err)
	}

	handleErr("logger.Init", logger.Init(filepath.Join(logDir, "openfinance-mcp.log")))
	defer logger.Close()

	dataDir := filepath.Join(homeDir, ".openfinance")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	masterKey, err := secret.MasterKey(dataDir)
	handleErr("secret.MasterKey", err)

	store, err := sqlite.Open(filepath.Join(dataDir, "pluggy.db"), masterKey)
	handleErr("sqlite.Open", err)
	defer store.Close()

	pluggyClient := pluggy.NewClient(pluggy.NewAuth(store))

	connRegistry := connections.NewRegistry(store)
	// Back-compat: seed the store from a legacy accounts.json if one exists.
	if err := connRegistry.ImportLegacy(filepath.Join(logDir, "accounts.json")); err != nil {
		logger.Info("connections: legacy accounts.json import skipped:", err)
	}

	toolRegistry := mcp.NewToolRegistry(
		tools.NewListConnectionsTool(connRegistry, pluggyClient),
		tools.NewCreateBankLinkTokenTool(pluggyClient, connRegistry),
		tools.NewWaitForSyncTool(pluggyClient, connRegistry),
		tools.NewGetConnectionStatusTool(pluggyClient, connRegistry),
		tools.NewListAccountsTool(pluggyClient, connRegistry),
		tools.NewGetAccountTool(pluggyClient),
		tools.NewListTransactionsTool(pluggyClient),
		tools.NewListInvestmentsTool(pluggyClient, connRegistry),
		tools.NewListCreditCardBillsTool(pluggyClient),
		tools.NewGetCreditCardBillTool(pluggyClient),
		tools.NewGetCurrentBillTool(pluggyClient),
		tools.NewRefreshConnectionsTool(connRegistry, pluggyClient),
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
