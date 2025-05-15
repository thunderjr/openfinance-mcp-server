package mcp

import (
	mcp "github.com/metoro-io/mcp-golang"
)

type ToolHandlerFunc interface{}

type ToolProvider interface {
	Name() string
	Description() string
	Handle() ToolHandlerFunc
}

type ToolRegistry struct {
	handlers []ToolProvider
}

func NewToolRegistry(handlers ...ToolProvider) *ToolRegistry {
	return &ToolRegistry{handlers}
}

func (r *ToolRegistry) Register(s *mcp.Server) error {
	for _, provider := range r.handlers {
		err := s.RegisterTool(
			provider.Name(),
			provider.Description(),
			provider.Handle(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
