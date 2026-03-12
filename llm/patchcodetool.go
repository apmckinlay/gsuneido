// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
)

var _ = addTool(toolSpec{
	name: "suneido_patch_code",
	description: `Modify an existing library definition by replacing a line range with new text.
The updated definition must be valid Suneido code.

Line range rules:
- from and to are 1-based line numbers.
- to is exclusive.
- from == to inserts without deleting.
- text may be empty to delete lines.
`,
	params: []stringParam{
		{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
		{name: "name", description: "Name of the definition (e.g. 'Alert')", required: true, kind: paramString},
		{name: "from", description: "First line number to replace (inclusive, 1-based)", required: true, kind: paramNumber},
		{name: "to", description: "Line number after the replaced range (exclusive, 1-based)", required: true, kind: paramNumber},
		{name: "text", description: "Replacement text", required: true, kind: paramString},
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
		from, err := requireInt(args, "from")
		if err != nil {
			return nil, err
		}
		to, err := requireInt(args, "to")
		if err != nil {
			return nil, err
		}
		text, err := requireString(args, "text")
		if err != nil {
			return nil, err
		}
		return patchCodeTool(ctx, library, name, from, to, text)
	},
})

type patchCodeOutput struct {
	Library  string   `json:"library" jsonschema:"Library name"`
	Name     string   `json:"name" jsonschema:"Definition name"`
	Warnings []string `json:"warnings" jsonschema:"Compiler warnings"`
}

func patchCodeTool(ctx context.Context, library, name string, from, to int, text string) (patchCodeOutput, error) {
	if !isValidName(name) {
		return patchCodeOutput{}, fmt.Errorf("invalid name: %s", name)
	}

	th := core.NewThread(core.MainThread)
	defer th.Close()

	if err := validateLibrary(th, library); err != nil {
		return patchCodeOutput{}, err
	}

	query := fmt.Sprintf("%s where group = -1 and name = %q", library, name)
	rtran := th.Dbms().Transaction(false)
	rq := rtran.Query(query, nil)
	hdr := rq.Header()
	row, _ := rq.Get(th, core.Next)
	if row == nil {
		rtran.Complete()
		return patchCodeOutput{}, fmt.Errorf("code not found for: %s in %s", name, library)
	}

	off := row[0].Off
	st := core.NewSuTran(rtran, false)
	vals := make(map[string]string, len(hdr.Fields[0]))
	for _, f := range hdr.Fields[0] {
		vals[f] = row.GetRawVal(hdr, f, th, st)
	}
	oldText := core.ToStr(core.Unpack(vals["text"]))
	rtran.Complete()

	newText, err := applyLineEdit(oldText, from, to, text)
	if err != nil {
		return patchCodeOutput{}, err
	}

	warnings, err := validateLibCode(th, newText)
	if err != nil {
		return patchCodeOutput{}, err
	}
	if warnings == nil {
		warnings = []string{}
	}

	if err := requireApproval(ctx, "patchCodeTool"); err != nil {
		return patchCodeOutput{}, err
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
		return patchCodeOutput{}, fmt.Errorf("transaction conflict: %s", conflict)
	}

	core.Global.Unload(name)
	return patchCodeOutput{
		Library:  library,
		Name:     name,
		Warnings: warnings,
	}, nil
}

func requireInt(args map[string]any, name string) (int, error) {
	v, ok := args[name]
	if !ok || v == nil {
		return 0, errors.New(name + " must be an integer")
	}
	return optionalInt(args, name, 0)
}

func applyLineEdit(oldText string, from, to int, insert string) (string, error) {
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
		return 0, 0, fmt.Errorf("from and to must be >= 1")
	}
	if to < from {
		return 0, 0, fmt.Errorf("to must be >= from")
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
		return 0, 0, fmt.Errorf("line range [%d,%d) out of bounds for %d lines", from, to, line)
	}
	return startOff, endOff, nil
}

func normalizeCRLF(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.ReplaceAll(s, "\n", "\r\n")
}
