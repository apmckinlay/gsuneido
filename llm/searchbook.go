// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/lexer"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var _ = addTool(toolSpec{
	name:        "suneido_search_book",
	description: "Search book pages by regex on path and text",
	params: []stringParam{
		{name: "book", description: "Name of the book table (e.g. 'suneidoc')", required: true, kind: paramString},
		{name: "path", description: "Regular expression applied to the full page path (path + name)", required: false, kind: paramString},
		{name: "text", description: "Regular expression applied to page text (optional if path provided)", required: false, kind: paramString},
		{name: "case_sensitive", description: "If true, regex matching is case sensitive (default false)", required: false, kind: paramBool},
	},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		book, err := requireString(args, "book")
		if err != nil {
			return nil, err
		}
		pathRx := optionalString(args, "path")
		textRx := optionalString(args, "text")
		caseSensitive, err := optionalBool(args, "case_sensitive", false)
		if err != nil {
			return nil, err
		}
		return searchBook(book, pathRx, textRx, caseSensitive)
	},
})

type searchBookOutput struct {
	Matches []bookMatch `json:"matches" jsonschema:"List of matching book pages"`
	HasMore bool        `json:"has_more,omitempty" jsonschema:"True when additional matches were truncated"`
}

type bookMatch struct {
	Path    string   `json:"path" jsonschema:"Full path to the book page"`
	Lines   []string `json:"lines" jsonschema:"Matching lines of text with line number prefixes"`
	HasMore bool     `json:"has_more,omitempty" jsonschema:"True when additional matching lines were truncated"`
}

const linesLimit = 5

func searchBook(book, pathRx, textRx string, caseSensitive bool) (searchBookOutput, error) {
	if strings.TrimSpace(pathRx) == "" && strings.TrimSpace(textRx) == "" {
		return searchBookOutput{}, fmt.Errorf("path or text is required")
	}
	if !lexer.IsIdentifier(book) {
		return searchBookOutput{}, fmt.Errorf("invalid book name: %s", book)
	}

	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	st := core.NewSuTran(tran, false)

	pathRx = strings.TrimSpace(pathRx)
	textRx = strings.TrimSpace(textRx)
	hasTextRx := textRx != ""
	var textPat regex.Pattern
	if pathRx != "" {
		pathRx = applyCaseSensitivity(pathRx, caseSensitive)
		if _, err := compileRegex(pathRx); err != nil {
			return searchBookOutput{}, fmt.Errorf("invalid path regex: %w", err)
		}
	}
	if textRx != "" {
		textRx = applyCaseSensitivity(textRx, caseSensitive)
		var err error
		if textPat, err = compileRegex(textRx); err != nil {
			return searchBookOutput{}, fmt.Errorf("invalid text regex: %w", err)
		}
	}

	query := buildBookQuery(book, pathRx, textRx)
	q := tran.Query(query, nil)
	hdr := q.Header()

	var matches []bookMatch
	for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
		if len(matches) >= searchLimit {
			return searchBookOutput{Matches: matches, HasMore: true}, nil
		}
		pathVal := row.GetVal(hdr, "path", th, st)
		path := core.ToStr(pathVal)
		nameVal := row.GetVal(hdr, "name", th, st)
		name := core.ToStr(nameVal)
		fullPath := buildBookPath(path, name)
		var lines []string
		var hasMore bool
		if hasTextRx {
			lines, hasMore = bookMatchLines(row, hdr, th, st, textPat)
		}
		matches = append(matches, bookMatch{Path: fullPath, Lines: lines, HasMore: hasMore})
	}

	return searchBookOutput{Matches: matches}, nil
}

func bookMatchLines(row core.Row, hdr *core.Header, th *core.Thread, st *core.SuTran, pat regex.Pattern) (lines []string, hasMore bool) {
	textVal := row.GetVal(hdr, "text", th, st)
	if textVal == nil {
		return nil, false
	}
	text, ok := textVal.ToStr()
	if !ok {
		return nil, false
	}
	seenLines := make(map[int]bool)
	for cap := range pat.All(text) {
		start := int(cap[0])
		if start < 0 {
			continue
		}
		if start > len(text) {
			start = len(text)
		}
		lineNum := 1 + strings.Count(text[:start], "\n")
		if seenLines[lineNum] {
			continue
		}
		seenLines[lineNum] = true
		if len(lines) >= linesLimit {
			hasMore = true
			break
		}
		lineText := lineAt(text, start)
		if lineText == "" {
			continue
		}
		lines = append(lines, addLineNumbers(lineText, lineNum))
	}
	return lines, hasMore
}

func buildBookPath(path, name string) string {
	if path == "" {
		return "/" + name
	}
	return path + "/" + name
}

func buildBookQuery(book, pathRx, textRx string) string {
	query := book
	if pathRx != "" {
		query += fmt.Sprintf(" where (path $ '/' $ name) =~ %q", pathRx)
	}
	if textRx != "" {
		query += fmt.Sprintf(" where text =~ %q", textRx)
	}
	return query + " sort path, name"
}
