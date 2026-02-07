// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package mcp provides an MCP (Model Context Protocol) server for gSuneido.
package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var httpServer *http.Server

// Start starts the MCP server in the background on the given port.
// It binds to localhost only (127.0.0.1) for security.
// Uses Streamable HTTP transport for streaming responses.
func Start(port string) {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "gSuneido MCP Server",
		Version: "0.0.1",
	}, nil)
	addTools(s)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return s
	}, nil))

	addr := "127.0.0.1:" + port
	httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("MCP server listening on http://%s/mcp", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("MCP server error: %v", err)
		}
	}()
}

// Stop gracefully shuts down the MCP server.
func Stop() {
	if httpServer != nil {
		exit.Progress("MCP server stopping")
		if err := httpServer.Close(); err != nil {
			log.Printf("ERROR: MCP server stop error: %v", err)
		}
		httpServer = nil
		exit.Progress("MCP server stopped")
	}
}

// addTools adds tools to the server.
func addTools(s *mcp.Server) {
	for _, spec := range toolSpecs {
		spec := spec
		t := makeTool(spec.name, spec.description, spec.outputSchema, spec.params...)
		s.AddTool(t, func(ctx context.Context, request *mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("MCP tool panic: %s: %v\n%s", spec.name, r, debug.Stack())
					result = toolError(fmt.Errorf("panic in tool: %s: %v", spec.name, r))
					err = nil
				}
			}()
			args, herr := requestArgs(request)
			if herr != nil {
				return toolError(herr), nil
			}
			out, herr := spec.handler(ctx, args)
			if herr != nil {
				return toolError(herr), nil
			}
			result, err = toolResult(out)
			if err != nil {
				return toolError(err), nil
			}
			return result, nil
		})
	}
}

type stringParam struct {
	name        string
	description string
	required    bool
	kind        paramKind
}

type paramKind int

const (
	paramString paramKind = iota
	paramNumber
	paramBool
)

type toolSpec struct {
	name         string
	description  string
	params       []stringParam
	outputSchema any
	handler      func(context.Context, map[string]any) (any, error)
}

var toolSpecs []toolSpec

// addTool adds a tool specification to the toolSpecs list.
func addTool(spec toolSpec) bool {
	toolSpecs = append(toolSpecs, spec)
	return true
}

func makeTool(name, desc string, outputSchema any, params ...stringParam) *mcp.Tool {
	return &mcp.Tool{
		Name:         name,
		Description:  desc,
		InputSchema:  inputSchema(params),
		OutputSchema: outputSchema,
	}
}

func inputSchema(params []stringParam) *jsonschema.Schema {
	schema := &jsonschema.Schema{Type: "object"}
	if len(params) == 0 {
		return schema
	}
	props := make(map[string]*jsonschema.Schema, len(params))
	required := []string{}
	for _, p := range params {
		prop := &jsonschema.Schema{Description: p.description}
		switch p.kind {
		case paramNumber:
			prop.Type = "integer"
		case paramBool:
			prop.Type = "boolean"
		default:
			prop.Type = "string"
		}
		props[p.name] = prop
		if p.required {
			required = append(required, p.name)
		}
	}
	schema.Properties = props
	if len(required) > 0 {
		schema.Required = required
	}
	return schema
}

func outputSchema[T any]() any {
	schema, err := jsonschema.For[T](nil)
	if err != nil {
		panic(fmt.Sprintf("tool output schema: %v", err))
	}
	return schema
}

func requestArgs(request *mcp.CallToolRequest) (map[string]any, error) {
	args := map[string]any{}
	if request == nil || request.Params == nil {
		return args, nil
	}
	if len(request.Params.Arguments) == 0 {
		return args, nil
	}
	if err := json.Unmarshal(request.Params.Arguments, &args); err != nil {
		return nil, err
	}
	return args, nil
}

func toolResult(out any) (*mcp.CallToolResult, error) {
	if out == nil {
		return &mcp.CallToolResult{Content: []mcp.Content{}}, nil
	}
	if s, ok := out.(string); ok {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: s}}}, nil
	}
	encoded, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content:           []mcp.Content{&mcp.TextContent{Text: string(encoded)}},
		StructuredContent: json.RawMessage(encoded),
	}, nil
}

func toolError(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		IsError: true,
	}
}

func requireString(args map[string]any, name string) (string, error) {
	s, ok := args[name].(string)
	if !ok {
		return "", errors.New(name + " must be a string")
	}
	return s, nil
}

func optionalString(args map[string]any, name string) string {
	s, _ := args[name].(string)
	return s
}

func optionalInt(args map[string]any, name string, def int) (int, error) {
	val, ok := args[name]
	if !ok || val == nil {
		return def, nil
	}
	switch n := val.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case float64:
		if n != float64(int(n)) {
			return 0, errors.New(name + " must be an integer")
		}
		return int(n), nil
	case float32:
		if n != float32(int(n)) {
			return 0, errors.New(name + " must be an integer")
		}
		return int(n), nil
	case string:
		if n == "" {
			return def, nil
		}
		parsed, err := strconv.Atoi(n)
		if err != nil {
			return 0, errors.New(name + " must be an integer")
		}
		return parsed, nil
	default:
		return 0, errors.New(name + " must be an integer")
	}
}

func optionalBool(args map[string]any, name string, def bool) (bool, error) {
	val, ok := args[name]
	if !ok || val == nil {
		return def, nil
	}
	switch b := val.(type) {
	case bool:
		return b, nil
	case string:
		if b == "" {
			return def, nil
		}
		switch b {
		case "true", "True", "TRUE", "1":
			return true, nil
		case "false", "False", "FALSE", "0":
			return false, nil
		default:
			return false, errors.New(name + " must be a boolean")
		}
	default:
		return false, errors.New(name + " must be a boolean")
	}
}
