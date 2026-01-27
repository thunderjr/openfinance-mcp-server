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

// jsonResponse marshals data to compact JSON (empty/zero/null fields stripped)
// and returns it as a tool response.
func jsonResponse(data any) (*mcp.ToolResponse, error) {
	bytes, err := compactMarshal(data)
	if err != nil {
		return toolError("Error marshalling response", err)
	}
	return mcp.NewToolResponse(mcp.NewTextContent(string(bytes))), nil
}

// compactMarshal marshals data to JSON, stripping empty/zero/null fields recursively.
// Uses a two-pass approach: marshal → generic interface → strip → re-marshal.
func compactMarshal(data any) ([]byte, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var generic interface{}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return nil, err
	}

	cleaned := stripEmpty(generic)
	return json.Marshal(cleaned)
}

// stripEmpty recursively removes empty/zero/null values from a generic JSON structure.
func stripEmpty(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		cleaned := make(map[string]interface{})
		for k, child := range val {
			stripped := stripEmpty(child)
			if !isEmpty(stripped) {
				cleaned[k] = stripped
			}
		}
		if len(cleaned) == 0 {
			return nil
		}
		return cleaned
	case []interface{}:
		cleaned := make([]interface{}, 0, len(val))
		for _, item := range val {
			cleaned = append(cleaned, stripEmpty(item))
		}
		return cleaned
	default:
		return val
	}
}

// isEmpty returns true if a value should be stripped from compact JSON output.
func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case float64:
		return val == 0
	case bool:
		return !val
	case map[string]interface{}:
		return len(val) == 0
	case []interface{}:
		return len(val) == 0
	}
	return false
}
