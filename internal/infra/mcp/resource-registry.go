package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ResourceHandler interface {
	Resource() mcp.ResourceTemplate
	Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
}

type ResourceRegistry struct {
	handlers []ResourceHandler
}

func NewResourceRegistry(handlers ...ResourceHandler) *ResourceRegistry {
	return &ResourceRegistry{handlers}
}

func (r *ResourceRegistry) Register(s *server.MCPServer) {
	for _, handler := range r.handlers {
		s.AddResourceTemplate(handler.Resource(), handler.Handle)
	}
}
