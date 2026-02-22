// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool represents a tool definition for LLM function calling.
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction describes a tool's function interface.
type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// ToolCall represents a tool call request from the LLM.
type ToolCall struct {
	Index    int              `json:"index,omitempty"`
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction contains the function name and arguments for a tool call.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolMessage represents a message containing tool call results.
type ToolMessage struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id"`
}

// MCPClient provides a local interface to MCP tools without HTTP transport.
type MCPClient struct {
	clientSession *mcp.ClientSession
	serverSession *mcp.ServerSession
	tools         []*mcp.Tool
}

// NewMCPClient creates a new local MCP client connected to the given server.
func NewMCPClient(server *mcp.Server) (*MCPClient, error) {
	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect server: %w", err)
	}

	client := mcp.NewClient(
		&mcp.Implementation{Name: "gsuneido-llm", Version: "0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		serverSession.Close()
		return nil, fmt.Errorf("connect client: %w", err)
	}

	// List available tools
	listResult, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		clientSession.Close()
		serverSession.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}

	return &MCPClient{
		clientSession: clientSession,
		serverSession: serverSession,
		tools:         listResult.Tools,
	}, nil
}

// Close closes the MCP client connections.
func (c *MCPClient) Close() error {
	if c.clientSession != nil {
		c.clientSession.Close()
	}
	return nil
}

// GetTools returns the available MCP tools in OpenAI function calling format.
func (c *MCPClient) GetTools() []Tool {
	tools := make([]Tool, len(c.tools))
	for i, t := range c.tools {
		// InputSchema may be *jsonschema.Schema or map[string]any depending on transport
		var params map[string]any
		switch schema := t.InputSchema.(type) {
		case map[string]any:
			params = schema
		default:
			// Marshal and unmarshal to convert any schema type to map[string]any
			b, err := json.Marshal(schema)
			if err != nil {
				panic(fmt.Sprintf("marshal input schema: %v", err))
			}
			if err := json.Unmarshal(b, &params); err != nil {
				panic(fmt.Sprintf("unmarshal input schema: %v", err))
			}
		}
		// Ensure properties field exists for compatibility with some providers
		if _, ok := params["properties"]; !ok {
			params["properties"] = map[string]any{}
		}
		tools[i] = Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		}
	}
	return tools
}

// CallTool invokes an MCP tool by name with the given arguments.
func (c *MCPClient) CallTool(ctx context.Context, name string, args map[string]any) (string, error) {
	var argsJSON json.RawMessage
	if args != nil {
		b, err := json.Marshal(args)
		if err != nil {
			return "", fmt.Errorf("marshal args: %w", err)
		}
		argsJSON = b
	} else {
		argsJSON = json.RawMessage("{}")
	}

	result, err := c.clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: argsJSON,
	})
	if err != nil {
		return "", fmt.Errorf("call tool: %w", err)
	}

	if result.IsError {
		if len(result.Content) > 0 {
			if tc, ok := result.Content[0].(*mcp.TextContent); ok {
				return "", fmt.Errorf("tool error: %s", tc.Text)
			}
		}
		return "", fmt.Errorf("tool error (unknown)")
	}

	// Extract text content from result
	if len(result.Content) > 0 {
		if tc, ok := result.Content[0].(*mcp.TextContent); ok {
			return tc.Text, nil
		}
		// Try structured content
		if result.StructuredContent != nil {
			if sc, ok := result.StructuredContent.(json.RawMessage); ok {
				return string(sc), nil
			}
			b, err := json.Marshal(result.StructuredContent)
			if err != nil {
				return "", fmt.Errorf("marshal structured content: %w", err)
			}
			return string(b), nil
		}
	}

	return "", nil
}

// CallToolFromLLM parses LLM tool call arguments and invokes the tool.
func (c *MCPClient) CallToolFromLLM(ctx context.Context, tc ToolCall) (string, error) {
	var args map[string]any
	argsStr := tc.Function.Arguments
	if argsStr != "" && argsStr != "null" {
		if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
			return "", fmt.Errorf("parse arguments: %w", err)
		}
	}

	result, err := c.CallTool(ctx, tc.Function.Name, args)
	if err != nil {
		return "", err
	}

	return result, nil
}
