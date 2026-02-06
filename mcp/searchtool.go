// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/regex"
)

const searchLimit = 100

func searchTool(libraryRx, nameRx, codeRx string, caseSensitive bool) (result searchCodeOutput, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("search code failed: %v", r)
		}
	}()
	if strings.TrimSpace(nameRx) == "" && strings.TrimSpace(codeRx) == "" {
		return searchCodeOutput{}, fmt.Errorf("name or code is required")
	}

	th := core.NewThread(core.MainThread)
	defer th.Close()
	libs, err := filterLibraries(th.Dbms().Libraries(), libraryRx, caseSensitive)
	if err != nil {
		return searchCodeOutput{}, err
	}
	if len(libs) == 0 {
		return searchCodeOutput{Matches: []codeMatch{}}, nil
	}

	nameQuery, err := searchQuery(nameRx, codeRx, caseSensitive)
	if err != nil {
		return searchCodeOutput{}, err
	}

	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	st := core.NewSuTran(tran, false)

	var matches []codeMatch
	for _, lib := range libs {
		q := tran.Query(lib+" "+nameQuery, nil)
		hdr := q.Header()
		for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
			if len(matches) >= searchLimit {
				return searchCodeOutput{Matches: matches, HasMore: true}, nil
			}
			nameVal := row.GetVal(hdr, "name", th, st)
			name, ok := nameVal.ToStr()
			if !ok {
				continue
			}
			path, err := codeItemPath(th, tran, st, lib, row, hdr)
			if err != nil {
				return searchCodeOutput{}, err
			}
			line, err := matchLine(row, hdr, th, st, codeRx, caseSensitive)
			if err != nil {
				return searchCodeOutput{}, err
			}
			matches = append(matches, codeMatch{Library: lib, Name: name, Path: path, Line: line})
		}
	}

	result = searchCodeOutput{Matches: matches}
	return result, nil
}

func searchQuery(nameRx, codeRx string, caseSensitive bool) (string, error) {
	nameRx = strings.TrimSpace(nameRx)
	codeRx = strings.TrimSpace(codeRx)
	nameRx = applyCaseSensitivity(nameRx, caseSensitive)
	codeRx = applyCaseSensitivity(codeRx, caseSensitive)
	if _, err := compileRegex(nameRx); err != nil {
		return "", fmt.Errorf("invalid name regex: %w", err)
	}
	if _, err := compileRegex(codeRx); err != nil {
		return "", fmt.Errorf("invalid code regex: %w", err)
	}
	return fmt.Sprintf("where group = -1 and name =~ %s and text =~ %s sort name",
		strconv.Quote(nameRx), strconv.Quote(codeRx)), nil
}

func filterLibraries(libs []string, libraryRx string, caseSensitive bool) ([]string, error) {
	libraryRx = strings.TrimSpace(libraryRx)
	if libraryRx == "" {
		return slices.Clone(libs), nil
	}
	libraryRx = applyCaseSensitivity(libraryRx, caseSensitive)
	pat, err := compileRegex(libraryRx)
	if err != nil {
		return nil, fmt.Errorf("invalid library regex: %w", err)
	}
	filtered := make([]string, 0, len(libs))
	for _, lib := range libs {
		if pat.Match(lib, nil) {
			filtered = append(filtered, lib)
		}
	}
	return filtered, nil
}

func applyCaseSensitivity(rx string, caseSensitive bool) string {
	if caseSensitive {
		return rx
	}
	if strings.HasPrefix(rx, "(?i)") {
		return rx
	}
	return "(?i)" + rx
}

func matchLine(row core.Row, hdr *core.Header, th *core.Thread, st *core.SuTran, codeRx string, caseSensitive bool) (string, error) {
	codeRx = strings.TrimSpace(codeRx)
	if codeRx == "" {
		return "", nil
	}
	codeRx = applyCaseSensitivity(codeRx, caseSensitive)
	pat, err := compileRegex(codeRx)
	if err != nil {
		return "", fmt.Errorf("invalid code regex: %w", err)
	}
	textVal := row.GetVal(hdr, "text", th, st)
	if textVal == nil {
		return "", fmt.Errorf("text column not found or null")
	}
	text, ok := textVal.ToStr()
	if !ok {
		return "", fmt.Errorf("text column is not a string")
	}
	var cap regex.Captures
	if !pat.Match(text, &cap) {
		return "", nil
	}
	start := int(cap[0])
	if start < 0 {
		return "", nil
	}
	if start > len(text) {
		start = len(text)
	}
	lineNum := 1 + strings.Count(text[:start], "\n")
	lineText := lineAt(text, start)
	if lineText == "" {
		return "", nil
	}
	return addLineNumbers(lineText, lineNum), nil
}

func lineAt(text string, pos int) string {
	if text == "" {
		return ""
	}
	if pos < 0 {
		pos = 0
	}
	if pos > len(text) {
		pos = len(text)
	}
	start := strings.LastIndex(text[:pos], "\n")
	if start == -1 {
		start = 0
	} else {
		start++
	}
	end := strings.Index(text[pos:], "\n")
	if end == -1 {
		end = len(text)
	} else {
		end = pos + end
	}
	return text[start:end]
}

func codeItemPath(th *core.Thread, tran core.ITran, st *core.SuTran, library string, row core.Row, hdr *core.Header) (string, error) {
	group, err := intValue(row.GetVal(hdr, "group", th, st), "group")
	if err != nil {
		return "", err
	}
	if group != -1 {
		return "", fmt.Errorf("expected leaf group -1")
	}
	parent, err := intValue(row.GetVal(hdr, "parent", th, st), "parent")
	if err != nil {
		return "", err
	}
	if parent == 0 {
		return "", nil
	}
	segments := []string{}
	for parent != 0 {
		folderArgs := core.SuObjectOf(core.SuStr(library))
		folderArgs.Set(core.SuStr("num"), core.IntVal(parent))
		folder, fhdr, _ := tran.Get(th, folderArgs, core.Only)
		if folder == nil {
			return "", fmt.Errorf("folder not found: %d", parent)
		}
		nameVal := folder.GetVal(fhdr, "name", th, st)
		name, ok := nameVal.ToStr()
		if !ok {
			return "", fmt.Errorf("folder name is not a string")
		}
		segments = append(segments, name)
		parent, err = intValue(folder.GetVal(fhdr, "parent", th, st), "parent")
		if err != nil {
			return "", err
		}
	}
	slices.Reverse(segments)
	return strings.Join(segments, "/"), nil
}

func compileRegex(rx string) (pat regex.Pattern, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	pat = regex.Compile(rx)
	return pat, nil
}
