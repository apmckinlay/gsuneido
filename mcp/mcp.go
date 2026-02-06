// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package mcp provides an MCP (Model Context Protocol) server for gSuneido.
package mcp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var httpServer *http.Server

// Start starts the MCP server in the background on the given port.
// It binds to localhost only (127.0.0.1) for security.
// Uses Streamable HTTP transport for streaming responses.
func Start(port string) {
	s := server.NewMCPServer(
		"gSuneido MCP Server",
		"0.0.1",
		server.WithToolCapabilities(true),
	)
	addTools(s)

	streamableServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
		server.WithStateLess(false),
	)

	mux := http.NewServeMux()
	mux.Handle("/mcp", streamableServer)

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
func addTools(s *server.MCPServer) {
	for _, spec := range toolSpecs {
		spec := spec
		t := makeTool(spec.name, spec.description, spec.outputSchema, spec.params...)
		s.AddTool(t, func(ctx context.Context, request mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("MCP tool panic: %s: %v\n%s", spec.name, r, debug.Stack())
					result = mcp.NewToolResultError("panic in tool: " + spec.name + ": " + fmt.Sprint(r))
					err = nil
				}
			}()
			out, herr := spec.handler(ctx, request.GetArguments())
			if herr != nil {
				return mcp.NewToolResultError(herr.Error()), nil
			}
			if s, ok := out.(string); ok {
				return mcp.NewToolResultText(s), nil
			}
			return mcp.NewToolResultJSON(out)
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
	outputSchema mcp.ToolOption
	handler      func(context.Context, map[string]any) (any, error)
}

var toolSpecs []toolSpec

// addTool adds a tool specification to the toolSpecs list.
func addTool(spec toolSpec) bool {
	toolSpecs = append(toolSpecs, spec)
	return true
}

func makeTool(name, desc string, outputSchema mcp.ToolOption, params ...stringParam) mcp.Tool {
	opts := []mcp.ToolOption{mcp.WithDescription(desc), outputSchema}
	for _, p := range params {
		prop := []mcp.PropertyOption{mcp.Description(p.description)}
		if p.required {
			prop = append(prop, mcp.Required())
		}
		switch p.kind {
		case paramNumber:
			opts = append(opts, mcp.WithNumber(p.name, prop...))
		case paramBool:
			opts = append(opts, mcp.WithBoolean(p.name, prop...))
		default:
			opts = append(opts, mcp.WithString(p.name, prop...))
		}
	}
	return mcp.NewTool(name, opts...)
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
