// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package mcp provides an MCP (Model Context Protocol) server for gSuneido.
package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/mark3labs/mcp-go/mcp"
)

// libraries
var _ = addTool(toolSpec{
	name:         "suneido_libraries",
	description:  "Get a list of the libraries currently in use in Suneido",
	outputSchema: mcp.WithOutputSchema[librariesOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		libs := core.GetDbms().Libraries()
		return librariesOutput{Libraries: libs}, nil
	},
})

type librariesOutput struct {
	Libraries []string `json:"libraries" jsonschema:"description=List of libraries currently in use"`
}

// tables
var _ = addTool(toolSpec{
	name:         "suneido_tables",
	description:  "Get a list of database table names that start with the given prefix (limit of 100)",
	params:       []stringParam{{name: "prefix", description: "Only return tables whose names start with this prefix (empty string for all)", required: true}},
	outputSchema: mcp.WithOutputSchema[tablesOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		prefix, err := requireString(args, "prefix")
		if err != nil {
			return nil, err
		}
		tables, err := tablesTool(prefix)
		if err != nil {
			return nil, err
		}
		return tablesOutput{Tables: tables}, nil
	},
})

type tablesOutput struct {
	Tables []string `json:"tables" jsonschema:"description=Table names matching the requested prefix"`
}

// schema
var _ = addTool(toolSpec{
	name:         "suneido_schema",
	description:  "Get the schema for a Suneido database table",
	params:       []stringParam{{name: "table", description: "Name of the table to get schema for", required: true}},
	outputSchema: mcp.WithOutputSchema[schemaOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		table, err := requireString(args, "table")
		if err != nil {
			return nil, err
		}
		schema := core.GetDbms().Schema(table)
		return schemaOutput{Schema: schema}, nil
	},
})

type schemaOutput struct {
	Schema string `json:"schema" jsonschema:"description=Schema definition for the requested table"`
}

// query
var _ = addTool(toolSpec{
	name:        "suneido_query",
	description: "Execute a Suneido database query and return the results as Suneido-format text (Value.String) in a simple row/column array format (limit 100)",
	params: []stringParam{
		{name: "query", description: "Suneido query (e.g. 'tables sort table')", required: true},
	},
	outputSchema: mcp.WithOutputSchema[queryOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		qs, err := requireString(args, "query")
		if err != nil {
			return nil, err
		}
		return queryTool(qs)
	},
})

type queryOutput struct {
	Query   string `json:"query" jsonschema:"description=Query string that was executed"`
	Results string `json:"results" jsonschema:"description=Formatted row/column output"`
	HasMore bool   `json:"has_more,omitempty" jsonschema:"description=True when additional rows were truncated"`
}

// read_code
var _ = addTool(toolSpec{
	name:        "suneido_read_code",
	description: "Get the source code from a library for a specific name",
	params: []stringParam{
		{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
		{name: "name", description: "Name of the definition (e.g. 'Alert')", required: true, kind: paramString},
		{name: "start_line", description: "1-based line number to start from (default 1)", required: false, kind: paramNumber},
		{name: "plain", description: "If true, don't add line numbers (default false)", required: false, kind: paramBool},
	},
	outputSchema: mcp.WithOutputSchema[readCodeOutput](),
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
})

type readCodeOutput struct {
	Plain        bool   `json:"plain" jsonschema:"description=Whether line numbers were omitted"`
	Library      string `json:"library" jsonschema:"description=Library name the definition was loaded from"`
	Name         string `json:"name" jsonschema:"description=Definition name"`
	Text         string `json:"text" jsonschema:"description=The source code content"`
	StartLine    int    `json:"start_line" jsonschema:"description=1-based starting line number for the returned text"`
	TotalLines   int    `json:"total_lines" jsonschema:"description=Total number of lines in the definition"`
	HasMore      bool   `json:"has_more,omitempty" jsonschema:"description=True when additional lines remain past the returned text"`
	Modified  string `json:"modified,omitempty" jsonschema:"description=Date/time when the record was last modified"`
	Committed string `json:"committed,omitempty" jsonschema:"description=Date/time when the record was last committed to version control"`
}

// search_code
var _ = addTool(toolSpec{
	name:        "suneido_search_code",
	description: "Search library code by regex on library, name, and text",
	params: []stringParam{
		{name: "library", description: "Regular expression applied to library names", required: true, kind: paramString},
		{name: "name", description: "Regular expression applied to definition names (optional if code provided)", required: false, kind: paramString},
		{name: "code", description: "Regular expression applied to definition text (optional if name provided)", required: false, kind: paramString},
		{name: "case_sensitive", description: "If true, regex matching is case sensitive (default false)", required: false, kind: paramBool},
	},
	outputSchema: mcp.WithOutputSchema[searchCodeOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		libraryRx, err := requireString(args, "library")
		if err != nil {
			return nil, err
		}
		nameRx := optionalString(args, "name")
		codeRx := optionalString(args, "code")
		if strings.TrimSpace(nameRx) == "" && strings.TrimSpace(codeRx) == "" {
			return nil, fmt.Errorf("name or code is required")
		}
		caseSensitive, err := optionalBool(args, "case_sensitive", false)
		if err != nil {
			return nil, err
		}
		return searchTool(libraryRx, nameRx, codeRx, caseSensitive)
	},
})

type searchCodeOutput struct {
	Matches []codeMatch `json:"matches" jsonschema:"description=List of matching library/name pairs"`
	HasMore bool        `json:"has_more,omitempty" jsonschema:"description=True when additional matches were truncated"`
}

type codeMatch struct {
	Library string `json:"library" jsonschema:"description=Library name"`
	Name    string `json:"name" jsonschema:"description=Definition name"`
	Path    string `json:"path" jsonschema:"description=Folder path within the library"`
	Line    string `json:"line" jsonschema:"description=Matching line of source code with line number prefix"`
}

// read_book
var _ = addTool(toolSpec{
	name:        "suneido_read_book",
	description: "Read from a Suneido book (documentation) table. Returns a JSON object containing the page 'text' and a 'children' array of sub-topic names.",
	params: []stringParam{
		{name: "book", description: "Name of the book table (e.g. 'suneidoc')", required: true},
		{name: "path", description: "The path to the book page. If sub-topics are returned in 'children', append them to this path to dive deeper. (e.g. 'Database/Reference/Query'). Empty or omitted for root.", required: false},
	},
	outputSchema: mcp.WithOutputSchema[readBookOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		book, err := requireString(args, "book")
		if err != nil {
			return nil, err
		}
		path := optionalString(args, "path")
		return bookTool(book, path)
	},
})

type readBookOutput struct {
	Book     string   `json:"book" jsonschema:"description=Book table name"`
	Path     string   `json:"path" jsonschema:"description=Normalized page path"`
	Text     string   `json:"text" jsonschema:"description=Book page text"`
	Children []string `json:"children" jsonschema:"description=Child topic names at this path"`
}

// code_folders
var _ = addTool(toolSpec{
	name:        "suneido_code_folders",
	description: "List folders and code items under a library path",
	params: []stringParam{
		{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
		{name: "path", description: "Folder path within the library (e.g. 'Debugging/Tests', empty string for root)", required: true, kind: paramString},
	},
	outputSchema: mcp.WithOutputSchema[codeFoldersOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		library, err := requireString(args, "library")
		if err != nil {
			return nil, err
		}
		path, err := requireString(args, "path")
		if err != nil {
			return nil, err
		}
		return codeFoldersTool(library, path)
	},
})

type codeFoldersOutput struct {
	Library  string   `json:"library" jsonschema:"description=Library name the folders were loaded from"`
	Path     string   `json:"path" jsonschema:"description=Normalized folder path"`
	Children []string `json:"children" jsonschema:"description=Child items at this path (folders end with '/')"`
}

// execute
var _ = addTool(toolSpec{
	name: "suneido_execute",
	description: "Executes Suneido code for its result or side effects.\n" +
		"Use this for calculations, data manipulation, or system commands.\n" +
		"Note: A single returned object will appear as the first result (e.g., [[1,2]]), while multiple return values appear as separate elements (e.g., [1,2]).",
	params:       []stringParam{{name: "code", description: "Suneido code to execute (as the body of a function)", required: true}},
	outputSchema: mcp.WithOutputSchema[execOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		code, err := requireString(args, "code")
		if err != nil {
			return nil, err
		}
		return execTool(code)
	},
})

type execOutput struct {
	Code     string   `json:"code" jsonschema:"description=The code that was executed"`
	Warnings []string `json:"warnings" jsonschema:"description=Compiler warnings"`
	Results  string   `json:"results" jsonschema:"description=0, 1, or multiple return values as Suneido-format strings"`
}

// check_code
var _ = addTool(toolSpec{
	name:         "suneido_check_code",
	description:  "Checks Suneido code for syntax and compilation errors without executing it. Returns compiler warnings only.",
	params:       []stringParam{{name: "code", description: "Suneido code to check (as the body of a function)", required: true}},
	outputSchema: mcp.WithOutputSchema[checkCodeOutput](),
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		code, err := requireString(args, "code")
		if err != nil {
			return nil, err
		}
		return checkTool(code)
	},
})

type checkCodeOutput struct {
	Code     string   `json:"code" jsonschema:"description=The code that was checked"`
	Warnings []string `json:"warnings" jsonschema:"description=Compiler warnings"`
}
