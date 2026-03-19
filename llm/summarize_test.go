// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSummarizeOutput(t *testing.T) {
	assert.TestOnlyIndividually(t)
	var b bytes.Buffer
	b.WriteString("# Tool Summarize Output Examples\n\n")
	b.WriteString("This document shows the output of summarize functions for all tools.\n\n")

	for _, spec := range toolSpecs {
		b.WriteString("## " + spec.name + "\n\n")
		b.WriteString(spec.description + "\n\n")
		b.WriteString("### Parameters\n\n")
		for _, p := range spec.params {
			req := ""
			if p.required {
				req = " (required)"
			}
			b.WriteString("- `" + p.name + "`" + req + ": " + p.description + "\n")
		}
		b.WriteString("\n### Summarize Output Examples\n\n")

		// Generate example args based on tool type
		examples := generateExamples(spec.name, spec.params)
		for _, args := range examples {
			b.WriteString("**Args:** `" + formatArgs(args) + "`\n\n")
			b.WriteString("=> " + spec.summarize(args) + "\n\n")
		}
		b.WriteString("---\n\n")
	}

	filename := "summarize_output.md"
	err := os.WriteFile(filename, b.Bytes(), 0644)
	if err != nil {
		t.Fatal("failed to write markdown file:", err)
	}
	t.Logf("Wrote %s", filename)
}

func generateExamples(toolName string, params []stringParam) []map[string]any {
	var examples []map[string]any

	switch toolName {
	case "suneido_read_book":
		examples = []map[string]any{
			{"book": "suneidoc"},
			{"book": "suneidoc", "path": "/Database"},
			{"book": "suneidoc", "path": "/Database/Reference/Query"},
		}
	case "suneido_read_code":
		examples = []map[string]any{
			{"library": "stdlib", "name": "Alert"},
			{"library": "stdlib", "name": "Alert", "start_line": 1},
			{"library": "stdlib", "name": "Alert", "start_line": 10, "plain": true},
		}
	case "suneido_create_code":
		examples = []map[string]any{
			{"library": "stdlib", "path": "", "name": "TestFunc",
				"code": "TestFunc()\n\t{\n\t123\n\t}"},
			{"library": "stdlib", "path": "Debugging/Tests", "name": "TestFunc", "code": "TestFunc()\n\t{\n\t123\n\t}"},
		}
	case "suneido_delete_code":
		examples = []map[string]any{
			{"library": "stdlib", "name": "OldFunc"},
		}
	case "suneido_edit_code":
		examples = []map[string]any{
			{"library": "stdlib", "name": "Alert", "mode": "insert_before", "line": 5, "code": "new line"},
			{"library": "stdlib", "name": "Alert", "mode": "insert_after", "line": 10, "code": "new line"},
			{"library": "stdlib", "name": "Alert", "mode": "replace_lines", "line": 5, "count": 3, "code": "replacement"},
		}
	case "suneido_execute":
		examples = []map[string]any{
			{"code": "1 + 2"},
			{"code": "Date.Now()"},
			{"code": "result = []\nfor (i = 0; i < 10; i++)\n\tresult.Add(i)\nresult"},
		}
	case "suneido_check_code":
		examples = []map[string]any{
			{"code": "function foo() { return 123 }"},
			{"code": "function bar(x) { x + 1 }"},
		}
	case "suneido_code_folders":
		examples = []map[string]any{
			{"library": "stdlib", "path": ""},
			{"library": "stdlib", "path": "Debugging"},
			{"library": "stdlib", "path": "Debugging/Tests"},
		}
	case "suneido_libraries":
		examples = []map[string]any{
			{},
		}
	case "suneido_query":
		examples = []map[string]any{
			{"query": "tables sort table"},
			{"query": "columns\nwhere table = 'test'"},
		}
	case "suneido_schema":
		examples = []map[string]any{
			{"table": "tables"},
			{"table": "views"},
		}
	case "suneido_search_book":
		examples = []map[string]any{
			{"book": "suneidoc", "path": "Database"},
			{"book": "suneidoc", "text": "query"},
			{"book": "suneidoc", "path": "Database", "text": "query", "case_sensitive": true},
		}
	case "suneido_search_code":
		examples = []map[string]any{
			{"name": "Alert"},
			{"library": "stdlib", "name": "Alert"},
			{"library": "stdlib", "code": "function"},
			{"library": "", "name": "^A", "code": "", "case_sensitive": true},
			{"library": "stdlib", "modified": true},
		}
	case "suneido_tables":
		examples = []map[string]any{
			{"prefix": ""},
			{"prefix": "test"},
		}
	default:
		// Generic fallback
		args := map[string]any{}
		for _, p := range params {
			switch p.kind {
			case paramNumber:
				args[p.name] = 1
			case paramBool:
				args[p.name] = true
			default:
				args[p.name] = "example"
			}
		}
		examples = []map[string]any{args}
	}

	return examples
}

func formatArgs(args map[string]any) string {
	if len(args) == 0 {
		return "{}"
	}
	var parts []string
	for k, v := range args {
		parts = append(parts, k+": "+formatValue(v))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func formatValue(v any) string {
	switch x := v.(type) {
	case string:
		if strings.Contains(x, "\n") {
			return `"` + strings.ReplaceAll(x[:min(20, len(x))], "\n", "↩") + `…"`
		}
		return `"` + x + `"`
	case bool:
		if x {
			return "true"
		}
		return "false"
	case int:
		return fmt.Sprintf("%d", x)
	case float64:
		return fmt.Sprintf("%g", x)
	default:
		return `"unknown"`
	}
}
