// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
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

type localTool struct {
	Name        string
	Description string
	InputSchema map[string]any
	Summarize   func(args map[string]any) string
	Handler     func(context.Context, map[string]any) (any, error)
}

// ToolClient provides a direct local interface to tool handlers.
type ToolClient struct {
	tools       []localTool
	openAITools []Tool
}

// NewToolClient creates a direct local tool client.
func NewToolClient() (*ToolClient, error) {
	tools := make([]localTool, 0, len(toolSpecs))
	openAITools := make([]Tool, 0, len(toolSpecs))
	for _, spec := range toolSpecs {
		if spec.summarize == nil {
			panic("tool " + spec.name + " missing summarize function")
		}
		params := spec.inputSchema()
		tools = append(tools, localTool{
			Name:        spec.name,
			Description: spec.description,
			InputSchema: params,
			Summarize:   spec.summarize,
			Handler:     spec.handler,
		})
		openAITools = append(openAITools, Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        spec.name,
				Description: spec.description,
				Parameters:  params,
			},
		})
	}
	return &ToolClient{tools: tools, openAITools: openAITools}, nil
}

func (c *ToolClient) getTool(name string) (localTool, bool) {
	for _, t := range c.tools {
		if t.Name == name {
			return t, true
		}
	}
	return localTool{}, false
}

// Close closes the tool client.
func (c *ToolClient) Close() error {
	return nil
}

// GetTools returns available tools in OpenAI function calling format.
func (c *ToolClient) GetTools() []Tool {
	return c.openAITools
}

// CallTool invokes a local tool by name with the given arguments.
func (c *ToolClient) CallTool(ctx context.Context, name string, args map[string]any) (string, error) {
	t, ok := c.getTool(name)
	if !ok {
		return "", fmt.Errorf("call tool: unknown tool %q", name)
	}
	out, err := callToolHandler(ctx, t, args)
	if err != nil {
		return "", err
	}
	if out == nil {
		return "", nil
	}
	if s, ok := out.(string); ok {
		return s, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshal result: %w", err)
	}
	return string(b), nil
}

func callToolHandler(ctx context.Context, t localTool, args map[string]any) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			result = nil
			err = fmt.Errorf("panic in tool: %s: %v", t.Name, r)
		}
	}()
	out, err := t.Handler(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("tool error: %w", err)
	}
	return out, nil
}

// CallToolFromLLM parses LLM tool call arguments and invokes the tool.
func (c *ToolClient) CallToolFromLLM(ctx context.Context, tc ToolCall) (string, error) {
	args, err := parseToolArgs(tc.Function.Arguments)
	if err != nil {
		return "", err
	}

	result, err := c.CallTool(ctx, tc.Function.Name, args)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *ToolClient) FormatToolCallForDisplay(tc ToolCall) (string, error) {
	t, ok := c.getTool(tc.Function.Name)
	if !ok {
		return "", fmt.Errorf("format tool call: unknown tool %q", tc.Function.Name)
	}
	args, err := parseToolArgs(tc.Function.Arguments)
	if err != nil {
		return "", err
	}
	return t.Summarize(args), nil
}

func parseToolArgs(argsStr string) (map[string]any, error) {
	var args map[string]any
	argsStr = strings.TrimSpace(argsStr)
	if argsStr != "" && argsStr != "null" && !isEmptyJSONObject(argsStr) {
		if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
			return nil, fmt.Errorf("parse arguments: %w", err)
		}
	}
	return args, nil
}

func isEmptyJSONObject(s string) bool {
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return false
	}
	for i := 1; i < len(s)-1; i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
		default:
			return false
		}
	}
	return true
}
