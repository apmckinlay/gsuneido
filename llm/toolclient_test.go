// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestToolClientGetTools(t *testing.T) {
	resetToolSpecsForTests()
	defer resetToolSpecsForTests()
	_ = addTool(toolSpec{
		name:        "test_tool",
		description: "A test tool",
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			return "test result", nil
		},
	})

	client, err := NewToolClient()
	assert.T(t).This(err, nil)
	defer client.Close()

	tools := client.GetTools()
	assert.T(t).True(len(tools) == 1)
	assert.T(t).This(tools[0].Function.Name, "test_tool")
}

func TestToolClientCallTool(t *testing.T) {
	resetToolSpecsForTests()
	defer resetToolSpecsForTests()
	_ = addTool(toolSpec{
		name:        "echo",
		description: "Echo back the input",
		params: []stringParam{
			{name: "message", kind: paramString},
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			s, _ := args["message"].(string)
			return s, nil
		},
	})

	client, err := NewToolClient()
	assert.T(t).This(err, nil)
	defer client.Close()

	result, err := client.CallTool(context.Background(), "echo", map[string]any{"message": "hello"})
	assert.T(t).This(err, nil)
	assert.T(t).This(result, "hello")
}

func TestToolClientCallToolFromLLM(t *testing.T) {
	resetToolSpecsForTests()
	defer resetToolSpecsForTests()
	_ = addTool(toolSpec{
		name:        "add",
		description: "Add two numbers",
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			a := int(args["a"].(float64))
			b := int(args["b"].(float64))
			return intStr(a + b), nil
		},
	})

	client, err := NewToolClient()
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

func TestToolClientCallToolStructuredResult(t *testing.T) {
	resetToolSpecsForTests()
	defer resetToolSpecsForTests()
	_ = addTool(toolSpec{
		name: "obj",
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			return map[string]any{"ok": true}, nil
		},
	})
	client, err := NewToolClient()
	assert.T(t).This(err, nil)
	defer client.Close()

	result, err := client.CallTool(context.Background(), "obj", nil)
	assert.T(t).This(err, nil)

	var got map[string]any
	err = json.Unmarshal([]byte(result), &got)
	assert.T(t).This(err, nil)
	assert.T(t).This(got["ok"], true)
}

func resetToolSpecsForTests() {
	toolSpecs = nil
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
