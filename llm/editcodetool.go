// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
)

var _ = addTool(toolSpec{
	name: "suneido_edit_code",
	description: `Modify a Suneido definition by inserting or replacing lines.
This tool is the preferred way to edit existing code.
- Lines are 1-based (matching the output of suneido_read_code)
- Modes:
  - "insert_before": Insert lines of code before the specified line
  - "insert_after": Insert lines of code after the specified line
  - "replace_lines": Replace 'count' lines of code starting at 'line'
- For deletions with replace_lines: Set 'code' to an empty string
- Always call suneido_read_code before this to ensure line numbers are current
- Do NOT include line numbers in the replacement code, just the code itself
`,
	params: []stringParam{
		{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
		{name: "name", description: "Name of the definition (e.g. 'Alert')", required: true, kind: paramString},
		{name: "mode", description: "Operation mode: 'insert_before', 'insert_after', or 'replace_lines'", required: true, kind: paramString},
		{name: "line", description: "Line number (1-based)", required: true, kind: paramNumber},
		{name: "count", description: "Number of lines to replace (only for replace_lines mode)", required: false, kind: paramNumber},
		{name: "code", description: "Replacement code", required: true, kind: paramString},
	},
	summarize: func(args map[string]any) string {
		line := argInt(args, "line", 0)
		count := argInt(args, "count", 0)
		mode := argString(args, "mode")
		code := argString(args, "code")
		var lineInfo string
		if mode == "replace_lines" && count > 0 {
			endLine := line + count - 1
			lineInfo = "line: `" + strconv.Itoa(line) + " to " + strconv.Itoa(endLine) + "`"
		} else {
			lineInfo = "line: `" + strconv.Itoa(line) + "`"
		}
		return mdSummary("Edit Code",
			argReqStr(args, "library"),
			argReqStr(args, "name"),
			mdInline(strings.ReplaceAll(mode, "_", "-")),
			lineInfo) + "\n" +
			summarizeCodeBlock(code)
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
		mode, err := requireString(args, "mode")
		if err != nil {
			return nil, err
		}
		line, err := requireInt(args, "line")
		if err != nil {
			return nil, err
		}
		count, err := optionalInt(args, "count", 0)
		if err != nil {
			return nil, err
		}
		if err := validateEditModeArgs(mode, count); err != nil {
			return nil, err
		}
		code, err := requireString(args, "code")
		if err != nil {
			return nil, err
		}
		return editCodeTool(ctx, library, name, mode, line, count, code)
	},
})

type editCodeOutput struct {
	Library       string   `json:"library" jsonschema:"Library name"`
	Name          string   `json:"name" jsonschema:"Definition name"`
	Warnings      []string `json:"warnings" jsonschema:"Compiler warnings"`
	ResultContext string   `json:"resultContext" jsonschema:"Text around the edit (4 lines before to 4 lines after)"`
	LineCount     int      `json:"lineCount" jsonschema:"Total number of lines in the resulting code"`
}

func editCodeTool(ctx context.Context, library, name, mode string, line, count int, text string) (editCodeOutput, error) {
	if !isValidName(name) {
		return editCodeOutput{}, fmt.Errorf("invalid name: %s", name)
	}

	th := core.NewThread(core.MainThread)
	defer th.Close()

	if err := validateLibrary(th, library); err != nil {
		return editCodeOutput{}, err
	}

	query := fmt.Sprintf("%s where group = -1 and name = %q", library, name)
	rtran := th.Dbms().Transaction(false)
	rq := rtran.Query(query, nil)
	hdr := rq.Header()
	row, _ := rq.Get(th, core.Next)
	if row == nil {
		rtran.Complete()
		return editCodeOutput{}, fmt.Errorf("code not found for: %s in %s", name, library)
	}

	off := row[0].Off
	st := core.NewSuTran(rtran, false)
	vals := make(map[string]string, len(hdr.Fields[0]))
	for _, f := range hdr.Fields[0] {
		vals[f] = row.GetRawVal(hdr, f, th, st)
	}
	oldText := core.ToStr(core.Unpack(vals["text"]))
	rtran.Complete()

	newText, err := applyLineEdit(oldText, mode, line, count, text)
	if err != nil {
		// Return context from oldText on error
		ctxLines := extractContext(oldText, line, line+count-1)
		lineCount := countLines(oldText)
		return editCodeOutput{
			Library:       library,
			Name:          name,
			ResultContext: ctxLines,
			LineCount:     lineCount,
		}, err
	}

	// Calculate edit location in new text for context
	normalizedText := normalizeCRLF(text)
	insertedLines := countLines(normalizedText)
	if insertedLines == 0 {
		insertedLines = 1
	}
	editStart, editEnd := editLineRange(mode, line, insertedLines)
	ctxLines := extractContext(newText, editStart, editEnd)
	lineCount := countLines(newText)

	warnings, err := validateLibCode(th, newText)
	if err != nil {
		return editCodeOutput{
			Library:       library,
			Name:          name,
			ResultContext: ctxLines,
			LineCount:     lineCount,
		}, err
	}
	if err := errorWarnings(warnings); err != nil {
		return editCodeOutput{
			Library:       library,
			Name:          name,
			ResultContext: ctxLines,
			LineCount:     lineCount,
			Warnings:      warnings,
		}, err
	}
	if warnings == nil {
		warnings = []string{}
	}

	if err := requireApproval(ctx, "editCodeTool", oldText, newText); err != nil {
		return editCodeOutput{
			Library:       library,
			Name:          name,
			ResultContext: ctxLines,
			LineCount:     lineCount,
		}, err
	}

	if core.ToStr(core.Unpack(vals["lib_before_text"])) == "" {
		vals["lib_before_text"] = core.PackValue(core.SuStr(oldText))
	}
	vals["text"] = core.PackValue(core.SuStr(newText))
	vals["lib_modified"] = core.PackValue(core.Now())

	utran := th.Dbms().Transaction(true)
	newRec := buildRecord(hdr, vals)
	utran.Update(th, library, off, newRec)
	if conflict := utran.Complete(); conflict != "" {
		return editCodeOutput{
			Library:       library,
			Name:          name,
			ResultContext: ctxLines,
			LineCount:     lineCount,
		}, fmt.Errorf("transaction conflict: %s", conflict)
	}

	core.Global.Unload(name)
	return editCodeOutput{
		Library:       library,
		Name:          name,
		Warnings:      warnings,
		ResultContext: ctxLines,
		LineCount:     lineCount,
	}, nil
}

func errorWarnings(warnings []string) error {
	errWarnings := make([]string, 0, len(warnings))
	for _, w := range warnings {
		if strings.HasPrefix(w, "ERROR:") {
			errWarnings = append(errWarnings, w)
		}
	}
	if len(errWarnings) == 0 {
		return nil
	}
	return fmt.Errorf("compile errors: %s", strings.Join(errWarnings, "; "))
}

func requireInt(args map[string]any, name string) (int, error) {
	v, ok := args[name]
	if !ok || v == nil {
		return 0, errors.New(name + " must be an integer")
	}
	return optionalInt(args, name, 0)
}

func validateEditModeArgs(mode string, count int) error {
	switch mode {
	case "replace_lines":
		if count < 1 {
			return errors.New("count must be >= 1 for replace_lines")
		}
	case "insert_before", "insert_after":
		if count != 0 {
			return errors.New("count is only valid for replace_lines")
		}
	default:
		return fmt.Errorf("invalid mode: %s (must be 'insert_before', 'insert_after', or 'replace_lines')", mode)
	}
	return nil
}

func applyLineEdit(oldText string, mode string, line, count int, insert string) (string, error) {
	var from, to int
	switch mode {
	case "insert_before":
		from = line
		to = line
	case "insert_after":
		from = line + 1
		to = line + 1
	case "replace_lines":
		from = line
		to = line + count
	}

	startOff, endOff, err := findFromTo(from, to, oldText)
	if err != nil {
		return "", err
	}

	insert = normalizeCRLF(insert)

	var b strings.Builder
	b.Grow(len(oldText) - (endOff - startOff) + len(insert))
	b.WriteString(oldText[:startOff])
	b.WriteString(insert)
	if endOff < len(oldText) {
		if insert != "" && !strings.HasSuffix(insert, "\n") && !strings.HasSuffix(insert, "\r\n") {
			b.WriteString("\r\n")
		}
		b.WriteString(oldText[endOff:])
	}
	return b.String(), nil
}

func findFromTo(from int, to int, oldText string) (int, int, error) {
	if from < 1 || to < 1 {
		return 0, 0, fmt.Errorf("line must be >= 1")
	}
	if to < from {
		return 0, 0, fmt.Errorf("line + count exceeds valid range")
	}

	// line 1 starts at offset 0; subsequent lines start after each '\n'
	startOff := -1
	endOff := -1
	line := 1
	if from == 1 {
		startOff = 0
	}
	if to == 1 {
		endOff = 0
	}
	pos := 0
	for startOff == -1 || endOff == -1 {
		i := strings.IndexByte(oldText[pos:], '\n')
		if i == -1 {
			break
		}
		line++
		pos += i + 1
		if line == from {
			startOff = pos
		}
		if line == to {
			endOff = pos
		}
	}
	// allow from/to to point just past the last line
	if startOff == -1 && from == line+1 {
		startOff = len(oldText)
	}
	if endOff == -1 && to == line+1 {
		endOff = len(oldText)
	}
	if startOff == -1 || endOff == -1 {
		return 0, 0, fmt.Errorf("line %d out of bounds for %d lines", from, line)
	}
	return startOff, endOff, nil
}

func normalizeCRLF(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.ReplaceAll(s, "\n", "\r\n")
	if s != "" && !strings.HasSuffix(s, "\r\n") {
		s += "\r\n"
	}
	return s
}

// countLines returns the number of lines in the text
func countLines(text string) int {
	return strings.Count(text, "\n")
}

// editLineRange returns the start and end line numbers (1-based, inclusive)
// of the edit context in the new text
func editLineRange(mode string, line, insertedLines int) (start, end int) {
	switch mode {
	case "insert_before":
		return line, line + insertedLines - 1
	case "insert_after":
		return line + 1, line + insertedLines
	case "replace_lines":
		return line, line + insertedLines - 1
	}
	return line, line
}

// extractContext returns lines from 4 lines before start to 4 lines after end
// with line numbers prefixed
func extractContext(text string, start, end int) string {
	lines := strings.Split(strings.TrimRight(text, "\r\n"), "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}

	// Calculate context range (4 lines before and after)
	contextStart := start - 4
	if contextStart < 1 {
		contextStart = 1
	}
	contextEnd := end + 4
	if contextEnd > len(lines) {
		contextEnd = len(lines)
	}

	// Build result with line numbers
	var b strings.Builder
	for i := contextStart; i <= contextEnd; i++ {
		fmt.Fprintf(&b, "%4d: %s\n", i, lines[i-1])
	}
	return b.String()
}
