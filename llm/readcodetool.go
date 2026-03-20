// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/aymanbagabas/go-udiff"
)

var _ = addTool(toolSpec{
	name:        "suneido_read_code",
	description: "Get the source code from a library for a specific name",
	params: []stringParam{
		{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
		{name: "name", description: "Name of the definition (e.g. 'Alert')", required: true, kind: paramString},
		{name: "start_line", description: "1-based line number to start from (default 1)", required: false, kind: paramNumber},
		{name: "num_lines", description: "Maximum number of lines to return (default 400)", required: false, kind: paramNumber},
		{name: "plain", description: "If true, don't add line numbers (default false)", required: false, kind: paramBool},
	},
	summarize: func(args map[string]any) string {
		return mdSummary("Read Code",
			argReqStr(args, "library"),
			argReqStr(args, "name"),
			argOptInt(args, "start_line", 1),
			argOptInt(args, "num_lines", codeLineLimit),
			argOptBool(args, "plain"))
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
		numLines, err := optionalInt(args, "num_lines", codeLineLimit)
		if err != nil {
			return nil, err
		}
		plain, err := optionalBool(args, "plain", false)
		if err != nil {
			return nil, err
		}
		return readCodeTool(library, name, startLine, numLines, plain)
	},
})

type readCodeOutput struct {
	Plain      bool    `json:"plain" jsonschema:"Whether line numbers were omitted"`
	Library    string  `json:"library" jsonschema:"Library name the definition was loaded from"`
	Name       string  `json:"name" jsonschema:"Definition name"`
	Text       string  `json:"text" jsonschema:"The source code content"`
	Diff       *string `json:"diff,omitempty" jsonschema:"Unified diff when lib_before_text is available"`
	StartLine  int     `json:"start_line" jsonschema:"1-based starting line number for the returned text"`
	TotalLines int     `json:"total_lines" jsonschema:"Total number of lines in the definition"`
	HasMore    bool    `json:"has_more,omitempty" jsonschema:"True when additional lines remain past the returned text"`
	Modified   string  `json:"modified,omitempty" jsonschema:"Date/time when the record was last modified"`
	Committed  string  `json:"committed,omitempty" jsonschema:"Date/time when the record was last committed to version control"`
}

const codeLineLimit = 400

func readCodeTool(library, name string, startLine int, numLines int, plain bool) (readCodeOutput, error) {
	if !isValidName(name) {
		return readCodeOutput{}, fmt.Errorf("invalid name: %s", name)
	}
	if startLine < 1 {
		return readCodeOutput{}, fmt.Errorf("start_line must be >= 1")
	}
	if numLines < 1 {
		return readCodeOutput{}, fmt.Errorf("num_lines must be >= 1")
	}
	if numLines > codeLineLimit {
		return readCodeOutput{}, fmt.Errorf("num_lines must be <= %d", codeLineLimit)
	}

	th := core.NewThread(core.MainThread)
	defer th.Close()
	if err := validateLibrary(th, library); err != nil {
		return readCodeOutput{}, err
	}

	query := fmt.Sprintf("%s where group = -1 and name = %q", library, name)
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	q := tran.Query(query, nil)
	hdr := q.Header()
	row, _ := q.Get(th, core.Next)
	if row == nil {
		extra := ""
		if core.Global.IsBuiltin(name) {
			extra = " (it is built-in)"
		}
		return readCodeOutput{}, fmt.Errorf("code not found for: %s in %s%s",
			name, library, extra)
	}

	st := core.NewSuTran(tran, false)
	text := core.ToStr(row.GetVal(hdr, "text", th, st))

	beforeText := core.ToStr(row.GetVal(hdr, "lib_before_text", th, st))
	var diff *string
	if beforeText != "" && text != "" {
		ud := udiff.Unified("old", "new", beforeText, text)
		diff = &ud
	}

	snippet, totalLines, hasMore := sliceCode(text, startLine, numLines)
	if !plain {
		snippet = addLineNumbers(snippet, startLine)
	}
	result := readCodeOutput{
		Plain:      plain,
		Library:    library,
		Name:       name,
		Text:       snippet,
		Diff:       diff,
		StartLine:  startLine,
		TotalLines: totalLines,
		HasMore:    hasMore,
		Modified:   formatDateVal(row.GetVal(hdr, "lib_modified", th, st)),
		Committed:  formatDateVal(row.GetVal(hdr, "lib_committed", th, st)),
	}
	return result, nil
}

func formatDateVal(val core.Value) string {
	if val == nil {
		return ""
	}
	if d, ok := core.AsDate(val); ok {
		return d.Format("yyyy-MM-dd HH:mm:ss")
	}
	return ""
}

func sliceCode(text string, startLine int, limit int) (string, int, bool) {
	if text == "" {
		return "", 0, false
	}
	if limit < 1 {
		return "", 0, false
	}
	line := 1
	startIdx := 0
	startFound := startLine == 1
	endIdx := len(text)
	endLine := startLine + limit - 1
	for i := 0; ; {
		nl := strings.Index(text[i:], "\n")
		if nl == -1 {
			break
		}
		idx := i + nl
		if !startFound && line == startLine-1 {
			startIdx = idx + 1
			startFound = true
		}
		if line == endLine {
			endIdx = idx
		}
		line++
		i = idx + 1
	}
	total := line
	if !startFound {
		return "", total, true
	}
	snippet := text[startIdx:endIdx]
	full := startLine == 1 && endIdx == len(text)
	return snippet, total, !full
}

func addLineNumbers(text string, startLine int) string {
	if text == "" {
		return text
	}
	lines := strings.Split(text, "\n")
	var sb strings.Builder

	// Pre-allocate memory to avoid multiple re-allocations.
	// 4 digits + colon + space = 6 extra chars per line.
	sb.Grow(len(text) + (len(lines) * 6))

	for i, line := range lines {
		// Calculate the actual line number based on the offset
		actualLineNum := startLine + i

		// Use %04d for fixed 4-digit zero padding
		fmt.Fprintf(&sb, "%04d: %s", actualLineNum, line)

		// Add newline back except for the very last line
		if i < len(lines)-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func isValidName(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '?' || c == '!') {
			return false
		}
	}
	return true
}

// validateLibrary checks if a library exists in the database.
// It returns an error if the library is not found.
func validateLibrary(th *core.Thread, library string) error {
	libs := th.Dbms().Libraries()
	if !slices.Contains(libs, library) {
		return fmt.Errorf("library not found: %s", library)
	}
	return nil
}
