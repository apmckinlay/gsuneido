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

	"github.com/apmckinlay/gsuneido/core"
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
		t := makeTool(spec.name, spec.description, spec.params...)
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
	name        string
	description string
	params      []stringParam
	handler     func(context.Context, map[string]any) (any, error)
}

var toolSpecs = []toolSpec{
	{
		name:        "suneido_libraries",
		description: "Get a list of the libraries currently in use in Suneido",
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			libs := core.GetDbms().Libraries()
			return map[string]any{"libraries": libs}, nil
		},
	},
	{
		name:        "suneido_tables",
		description: "Get a list of database table names that start with the given prefix (limit of 100)",
		params:      []stringParam{{name: "prefix", description: "Only return tables whose names start with this prefix (empty string for all)", required: true}},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			prefix, err := requireString(args, "prefix")
			if err != nil {
				return nil, err
			}
			tables, err := tablesTool(prefix)
			if err != nil {
				return nil, err
			}
			return map[string]any{"tables": tables}, nil
		},
	},
	{
		name:        "suneido_schema",
		description: "Get the schema for a Suneido database table",
		params:      []stringParam{{name: "table", description: "Name of the table to get schema for", required: true}},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			table, err := requireString(args, "table")
			if err != nil {
				return nil, err
			}
			schema := core.GetDbms().Schema(table)
			return map[string]any{"schema": schema}, nil
		},
	},
	{
		name:        "suneido_query",
		description: "Execute a Suneido database query and return the results as Suneido-format text (Value.String) in a simple row/column array format (limit 100)",
		params: []stringParam{
			{name: "query", description: "Suneido query (e.g. 'tables sort table')", required: true},
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			qs, err := requireString(args, "query")
			if err != nil {
				return nil, err
			}
			return queryTool(qs)
		},
	},
	{
		name:        "suneido_read_code",
		description: "Get the source code from a library for a specific name",
		params: []stringParam{
			{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
			{name: "name", description: "Name of the definition (e.g. 'Alert')", required: true, kind: paramString},
			{name: "start_line", description: "1-based line number to start from (default 1)", required: false, kind: paramNumber},
			{name: "plain", description: "If true, don't add line numbers (default false)", required: false, kind: paramBool},
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			library, err := requireString(args, "library")
			if err != nil {
				return nil, err
			}
			name, err := requireString(args, "name")
			if err != nil {
				return nil, err
			}
			startLine, err := optionalInt(args, "start_line", 1)
			if err != nil {
				return nil, err
			}
			plain, err := optionalBool(args, "plain", false)
			if err != nil {
				return nil, err
			}
			return codeTool(library, name, startLine, plain)
		},
	},
	{
		name:        "suneido_read_book",
		description: "Read from a Suneido book (documentation) table. Returns a JSON object containing the page 'text' and a 'children' array of sub-topic names.",
		params: []stringParam{
			{name: "book", description: "Name of the book table (e.g. 'suneidoc')", required: true},
			{name: "path", description: "The path to the book page. If sub-topics are returned in 'children', append them to this path to dive deeper. (e.g. 'Database/Reference/Query'). Empty or omitted for root.", required: false},
		},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			book, err := requireString(args, "book")
			if err != nil {
				return nil, err
			}
			path := optionalString(args, "path")
			return bookTool(book, path)
		},
	},
	{
		name: "suneido_execute",
		description: "Executes Suneido code for its return value or side effects (e.g., database updates). Returns a text report containing:\n" +
			"- code: The code that was executed.\n" +
			"- results: Array of 0, 1, or multiple return values as Suneido-format strings (Value.String())\n" +
			"- warnings: Array of compiler warnings\n" +
			"Note: A single returned object will appear as the first result (e.g., [[1,2]]), while multiple return values appear as separate elements (e.g., [1,2])." +
			"Use this for calculations, data manipulation, or system commands.",
		params: []stringParam{{name: "code", description: "Suneido code to execute (as the body of a function)", required: true}},
		handler: func(ctx context.Context, args map[string]any) (any, error) {
			code, err := requireString(args, "code")
			if err != nil {
				return nil, err
			}
			return exectool(code)
		},
	},
}

func makeTool(name, desc string, params ...stringParam) mcp.Tool {
	opts := []mcp.ToolOption{mcp.WithDescription(desc)}
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
