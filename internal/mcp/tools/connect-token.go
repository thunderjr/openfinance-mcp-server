package tools

import (
	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type CreateBankLinkTokenTool struct {
	client *pluggy.Client
}

func NewCreateBankLinkTokenTool(client *pluggy.Client) *CreateBankLinkTokenTool {
	return &CreateBankLinkTokenTool{client}
}

func (t *CreateBankLinkTokenTool) Name() string {
	return "create_bank_link_token"
}

func (t *CreateBankLinkTokenTool) Description() string {
	return `Creates a token for linking or updating a bank connection.

Use this tool when:
- The user wants to connect a new bank account
- A bank connection needs to be re-authenticated
- You need to generate a link for the user to update credentials

Returns: A connect token URL that can be used to link/update bank credentials.`
}

func (t *CreateBankLinkTokenTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleCreateBankLinkToken
}

type CreateBankLinkTokenArgs struct {
	ConnectionID string `json:"connection_id" jsonschema:"required,description=The connection ID to generate a link token for"`
}

func (t *CreateBankLinkTokenTool) handleCreateBankLinkToken(args CreateBankLinkTokenArgs) (*mcp.ToolResponse, error) {
	if args.ConnectionID == "" {
		return validationError("Connection ID is required")
	}

	logger.Info("Creating bank link token for connection:", args.ConnectionID)

	token, err := t.client.ConnectToken(args.ConnectionID)
	if err != nil {
		return toolError("Error creating bank link token", err)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(token)), nil
}
