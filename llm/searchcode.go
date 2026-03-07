// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var _ = addTool(toolSpec{
	name:        "suneido_search_code",
	description: "Search library code by regex on library, name, and text",
	params: []stringParam{
		{name: "library", description: "Regular expression applied to library names", required: true, kind: paramString},
		{name: "name", description: "Regular expression applied to definition names (optional if code provided)", required: false, kind: paramString},
		{name: "code", description: "Regular expression applied to definition text (optional if name provided)", required: false, kind: paramString},
		{name: "case_sensitive", description: "If true, regex matching is case sensitive (default false)", required: false, kind: paramBool},
		{name: "modified", description: "If true, only return results where the code has been modified", required: false, kind: paramBool},
	},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		libraryRx, err := requireString(args, "library")
		if err != nil {
			return nil, err
		}
		nameRx := optionalString(args, "name")
		codeRx := optionalString(args, "code")
		caseSensitive, err := optionalBool(args, "case_sensitive", false)
		if err != nil {
			return nil, err
		}
		modified, err := optionalBool(args, "modified", false)
		if err != nil {
			return nil, err
		}
		return searchCode(libraryRx, nameRx, codeRx, caseSensitive, modified)
	},
})

type searchCodeOutput struct {
	Matches []codeMatch `json:"matches" jsonschema:"List of matching library/name pairs"`
	HasMore bool        `json:"has_more,omitempty" jsonschema:"True when additional matches were truncated"`
}

type codeMatch struct {
	Library string   `json:"library" jsonschema:"Library name"`
	Name    string   `json:"name" jsonschema:"Definition name"`
	Path    string   `json:"path" jsonschema:"Folder path within the library"`
	Lines   []string `json:"lines" jsonschema:"Matching lines of source code with line number prefixes"`
	HasMore bool     `json:"has_more,omitempty" jsonschema:"True when additional matching lines were truncated"`
}

const searchLimit = 100

func searchCode(libraryRx, nameRx, codeRx string, caseSensitive, modified bool) (result searchCodeOutput, err error) {
	nameRx = strings.TrimSpace(nameRx)
	if nameRx == "" && codeRx == "" && !modified {
		return searchCodeOutput{}, fmt.Errorf("name or code is required (unless modified is true)")
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

	nameQuery, err := searchQuery(nameRx, codeRx, caseSensitive, modified)
	if err != nil {
		return searchCodeOutput{}, err
	}

	hasCodeRx := false
	var codePat regex.Pattern
	if codeRx != "" {
		hasCodeRx = true
		codeRx = applyCaseSensitivity(codeRx, caseSensitive)
		if codePat, err = compileRegex(codeRx); err != nil {
			return searchCodeOutput{}, fmt.Errorf("invalid code regex: %w", err)
		}
	}

	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	st := core.NewSuTran(tran, false)
	folderCache := make(map[int]folderInfo)

	var matches []codeMatch
	for _, lib := range libs {
		clear(folderCache)
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
			path, err := codeItemPath(th, tran, st, lib, row, hdr, folderCache)
			if err != nil {
				return searchCodeOutput{}, err
			}
			var lines []string
			var hasMore bool
			if hasCodeRx {
				lines, hasMore = matchLines(row, hdr, th, st, codePat)
			}
			matches = append(matches, codeMatch{Library: lib, Name: name,
				Path: path, Lines: lines, HasMore: hasMore})
		}
	}

	if len(matches) == 0 && isValidName(nameRx) && core.Global.IsBuiltin(nameRx) {
		return searchCodeOutput{}, fmt.Errorf("%s is built-in", nameRx)
	}

	result = searchCodeOutput{Matches: matches}
	return result, nil
}

func searchQuery(nameRx, codeRx string, caseSensitive, modified bool) (string, error) {
	var query strings.Builder
	query.WriteString("where group = -1")
	if nameRx != "" {
		nameRx = applyCaseSensitivity(nameRx, caseSensitive)
		if _, err := compileRegex(nameRx); err != nil {
			return "", fmt.Errorf("invalid name regex: %w", err)
		}
		fmt.Fprintf(&query, " where name =~ %q", nameRx)
	}
	if codeRx != "" {
		codeRx = applyCaseSensitivity(codeRx, caseSensitive)
		if _, err := compileRegex(codeRx); err != nil {
			return "", fmt.Errorf("invalid code regex: %w", err)
		}
		fmt.Fprintf(&query, " where text =~ %q", codeRx)
	}
	if modified {
		query.WriteString(" where lib_modified isnt \"\"")
	}
	query.WriteString(" sort name")
	return query.String(), nil
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

func matchLines(row core.Row, hdr *core.Header, th *core.Thread, st *core.SuTran,
	pat regex.Pattern) (lines []string, hasMore bool) {
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

type folderInfo struct {
	name   string
	parent int
}

func codeItemPath(th *core.Thread, tran core.ITran, st *core.SuTran, library string, row core.Row, hdr *core.Header, folderCache map[int]folderInfo) (string, error) {
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
		info, err := cachedFolderInfo(th, tran, st, library, parent, folderCache)
		if err != nil {
			return "", err
		}
		segments = append(segments, info.name)
		parent = info.parent
	}
	slices.Reverse(segments)
	return strings.Join(segments, "/"), nil
}

func cachedFolderInfo(th *core.Thread, tran core.ITran, st *core.SuTran, library string, num int, folderCache map[int]folderInfo) (folderInfo, error) {
	if folderCache != nil {
		if info, ok := folderCache[num]; ok {
			return info, nil
		}
		info, err := fetchFolderInfo(th, tran, st, library, num)
		if err != nil {
			return folderInfo{}, err
		}
		folderCache[num] = info
		return info, nil
	}
	return fetchFolderInfo(th, tran, st, library, num)
}

func fetchFolderInfo(th *core.Thread, tran core.ITran, st *core.SuTran, library string, num int) (folderInfo, error) {
	folderArgs := core.SuObjectOf(core.SuStr(library))
	folderArgs.Set(core.SuStr("num"), core.IntVal(num))
	folder, fhdr, _ := tran.Get(th, folderArgs, core.Only)
	if folder == nil {
		return folderInfo{}, fmt.Errorf("folder not found: %d", num)
	}
	nameVal := folder.GetVal(fhdr, "name", th, st)
	name, ok := nameVal.ToStr()
	if !ok {
		return folderInfo{}, fmt.Errorf("folder name is not a string")
	}
	parent, err := intValue(folder.GetVal(fhdr, "parent", th, st), "parent")
	if err != nil {
		return folderInfo{}, err
	}
	return folderInfo{name: name, parent: parent}, nil
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
