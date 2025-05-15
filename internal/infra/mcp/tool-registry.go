package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ToolHandler interface {
	Tool() mcp.Tool
	Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

type ToolRegistry struct {
	handlers []ToolHandler
}

func NewToolRegistry(handlers ...ToolHandler) *ToolRegistry {
	return &ToolRegistry{handlers}
}

func (r *ToolRegistry) Register(s *server.MCPServer) {
	tools := make([]server.ServerTool, 0, len(r.handlers))
	for _, handler := range r.handlers {
		tools = append(tools, server.ServerTool{
			Tool:    handler.Tool(),
			Handler: handler.Handle,
		})
	}

	s.AddTools(tools...)
}
