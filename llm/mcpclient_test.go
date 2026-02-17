// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPClientGetTools(t *testing.T) {
	// Create a simple test server
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, struct{}, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "test result"}},
		}, struct{}{}, nil
	})

	client, err := NewMCPClient(server)
	assert.T(t).This(err, nil)
	defer client.Close()

	tools := client.GetTools()
	assert.T(t).True(len(tools) == 1)
	assert.T(t).This(tools[0].Function.Name, "test_tool")
}

func TestMCPClientCallTool(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	type echoArgs struct {
		Message string `json:"message"`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "echo",
		Description: "Echo back the input",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]any{"type": "string"},
			},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args echoArgs) (*mcp.CallToolResult, struct{}, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: args.Message}},
		}, struct{}{}, nil
	})

	client, err := NewMCPClient(server)
	assert.T(t).This(err, nil)
	defer client.Close()

	result, err := client.CallTool(context.Background(), "echo", map[string]any{"message": "hello"})
	assert.T(t).This(err, nil)
	assert.T(t).This(result, "hello")
}

func TestMCPClientCallToolFromLLM(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	type addArgs struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "Add two numbers",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"a": map[string]any{"type": "number"},
				"b": map[string]any{"type": "number"},
			},
		},
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addArgs) (*mcp.CallToolResult, struct{}, error) {
		sum := args.A + args.B
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: intStr(sum)}},
		}, struct{}{}, nil
	})

	client, err := NewMCPClient(server)
	assert.T(t).This(err, nil)
	defer client.Close()

	tc := ToolCall{
		ID:   "call_123",
		Type: "function",
		Function: ToolCallFunction{
			Name:      "add",
			Arguments: `{"a": 2, "b": 3}`,
		},
	}

	result, err := client.CallToolFromLLM(context.Background(), tc)
	assert.T(t).This(err, nil)
	assert.T(t).This(result, "5")
}

func intStr(n int) string {
	if n == 0 {
		return "0"
	}
	var neg bool
	if n < 0 {
		neg = true
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
