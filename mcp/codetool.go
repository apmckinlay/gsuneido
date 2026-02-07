// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"fmt"
	"slices"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/aymanbagabas/go-udiff"
)

const codeLineLimit = 400

func codeTool(library, name string, startLine int, plain bool) (readCodeOutput, error) {
	if !isValidName(name) {
		return readCodeOutput{}, fmt.Errorf("invalid name: %s", name)
	}
	if startLine < 1 {
		return readCodeOutput{}, fmt.Errorf("start_line must be >= 1")
	}

	th := core.NewThread(core.MainThread)
	defer th.Close()
	libs := th.Dbms().Libraries()
	if !slices.Contains(libs, library) {
		return readCodeOutput{}, fmt.Errorf("library not found: %s", library)
	}

	query := fmt.Sprintf("%s where group = -1 and name = '%s'", library, name)
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	q := tran.Query(query, nil)
	hdr := q.Header()
	row, _ := q.Get(th, core.Next)
	if row == nil {
		return readCodeOutput{}, fmt.Errorf("code not found for: %s in %s", name, library)
	}

	st := core.NewSuTran(tran, false)
	text := core.ToStr(row.GetVal(hdr, "text", th, st))

	beforeText := core.ToStr(row.GetVal(hdr, "lib_before_text", th, st))
	var diff *string
	if beforeText != "" && text != "" {
		ud := udiff.Unified("old", "new", beforeText, text)
		diff = &ud
	}

	snippet, totalLines, hasMore := sliceCode(text, startLine, codeLineLimit)
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
