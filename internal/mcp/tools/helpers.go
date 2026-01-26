package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"
)

// validationError creates a tool response for validation errors.
// Returns both a user-facing response and an error for logging.
func validationError(msg string) (*mcp.ToolResponse, error) {
	return mcp.NewToolResponse(mcp.NewTextContent(msg)), fmt.Errorf(msg)
}

// toolError creates a tool response for operational errors.
// Formats the error with a prefix for consistency.
func toolError(prefix string, err error) (*mcp.ToolResponse, error) {
	msg := fmt.Sprintf("%s: %v", prefix, err)
	return mcp.NewToolResponse(mcp.NewTextContent(msg)), fmt.Errorf(msg)
}

// jsonResponse marshals data to JSON and returns it as a tool response.
func jsonResponse(data any) (*mcp.ToolResponse, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return toolError("Error marshalling response", err)
	}
	return mcp.NewToolResponse(mcp.NewTextContent(string(bytes))), nil
}
