// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
)

var _ = addTool(toolSpec{
	name: "suneido_upsert_code",
	description: "Output or update a library definition. " +
		"The definition must be valid Suneido code.",
	params: []stringParam{
		{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
		{name: "path", description: "Folder path within the library (e.g. 'Debugging/Tests', empty string for root)", required: true, kind: paramString},
		{name: "name", description: "Name of the definition (e.g. 'Alert')", required: true, kind: paramString},
		{name: "text", description: "The source code for the definition", required: true, kind: paramString},
	},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		library, err := requireString(args, "library")
		if err != nil {
			return nil, err
		}
		path, err := requireString(args, "path")
		if err != nil {
			return nil, err
		}
		name, err := requireString(args, "name")
		if err != nil {
			return nil, err
		}
		text, err := requireString(args, "text")
		if err != nil {
			return nil, err
		}
		return upsertCodeTool(library, path, name, text)
	},
})

type upsertCodeOutput struct {
	Library  string   `json:"library" jsonschema:"Library name"`
	Name     string   `json:"name" jsonschema:"Definition name"`
	Action   string   `json:"action" jsonschema:"Either 'inserted' or 'updated'"`
	Warnings []string `json:"warnings" jsonschema:"Compiler warnings"`
}

func upsertCodeTool(library, path, name, text string) (result upsertCodeOutput, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("upsert error: %v", r)
		}
	}()

	if !isValidName(name) {
		return upsertCodeOutput{}, fmt.Errorf("invalid name: %s", name)
	}
	path = normalizeFolderPath(path)

	th := core.NewThread(core.MainThread)
	defer th.Close()

	libs := th.Dbms().Libraries()
	if !slices.Contains(libs, library) {
		return upsertCodeOutput{}, fmt.Errorf("library not found: %s", library)
	}

	parent, err := ensurePathParent(th, library, path)
	if err != nil {
		return upsertCodeOutput{}, err
	}

	// Validate the code by compiling it
	warnings, err := validateLibCode(th, text)
	if err != nil {
		return upsertCodeOutput{}, err
	}
	if warnings == nil {
		warnings = []string{}
	}

	now := core.Now()

	// Look up existing record in a read transaction
	query := fmt.Sprintf("%s where group = -1 and name = %q", library, name)
	rtran := th.Dbms().Transaction(false)
	rq := rtran.Query(query, nil)
	hdr := rq.Header()
	row, _ := rq.Get(th, core.Next)

	var action string
	if row == nil {
		// Get max num to assign a unique num to the new record
		maxNum := maxLibNum(th, rtran, library)
		rtran.Complete()
		// Insert: use Output on the table
		action = "inserted"
		utran := th.Dbms().Transaction(true)
		iq := utran.Query(library, nil)
		ihdr := iq.Header()
		rec := buildRecord(ihdr, map[string]core.Value{
			"name":            core.SuStr(name),
			"text":            core.SuStr(text),
			"lib_before_text": core.SuStr(""),
			"lib_modified":    now,
			"group":           core.SuInt(-1),
			"num":             core.SuInt(maxNum + 1),
			"parent":          core.SuInt(parent),
		})
		iq.Output(th, rec)
		if conflict := utran.Complete(); conflict != "" {
			return upsertCodeOutput{}, fmt.Errorf("transaction conflict: %s", conflict)
		}
	} else {
		// Update existing record
		action = "updated"
		off := row[0].Off
		st := core.NewSuTran(rtran, false)
		// Collect all existing field values
		fields := hdr.Fields[0]
		vals := make(map[string]core.Value, len(fields))
		for _, f := range fields {
			vals[f] = row.GetVal(hdr, f, th, st)
		}
		// Check lib_before_text: if empty, set to original text
		existingBefore := core.ToStr(vals["lib_before_text"])
		if existingBefore == "" {
			vals["lib_before_text"] = vals["text"]
		}
		// Apply changes
		vals["text"] = core.SuStr(text)
		vals["lib_modified"] = now
		vals["parent"] = core.SuInt(parent)
		rtran.Complete()

		// Build new record
		utran := th.Dbms().Transaction(true)
		uq := utran.Query(query, nil)
		uhdr := uq.Header()
		newRec := buildRecord(uhdr, vals)
		utran.Update(th, library, off, newRec)
		if conflict := utran.Complete(); conflict != "" {
			return upsertCodeOutput{}, fmt.Errorf("transaction conflict: %s", conflict)
		}
		core.Global.Unload(name)
	}

	result = upsertCodeOutput{
		Library:  library,
		Name:     name,
		Action:   action,
		Warnings: warnings,
	}
	return result, nil
}

func ensurePathParent(th *core.Thread, library, path string) (int, error) {
	if path == "" {
		return 0, nil
	}

	utran := th.Dbms().Transaction(true)
	parent := 0
	nextNum := maxLibNum(th, utran, library) + 1

	for _, segment := range strings.Split(path, "/") {
		if segment == "" {
			continue
		}
		num, exists, leaf, err := lookupFolder(th, utran, library, parent, segment)
		if err != nil {
			utran.Complete()
			return 0, err
		}
		if leaf {
			utran.Complete()
			return 0, fmt.Errorf("path segment is not a folder: %s", segment)
		}
		if exists {
			parent = num
			continue
		}

		iq := utran.Query(library, nil)
		ihdr := iq.Header()
		rec := buildRecord(ihdr, map[string]core.Value{
			"name":            core.SuStr(segment),
			"text":            core.SuStr(""),
			"lib_before_text": core.SuStr(""),
			"lib_modified":    core.Now(),
			"group":           core.SuInt(parent),
			"num":             core.SuInt(nextNum),
			"parent":          core.SuInt(parent),
		})
		iq.Output(th, rec)
		parent = nextNum
		nextNum++
	}

	if conflict := utran.Complete(); conflict != "" {
		return 0, fmt.Errorf("transaction conflict: %s", conflict)
	}
	return parent, nil
}

func lookupFolder(th *core.Thread, tran core.ITran, library string, parent int, name string) (int, bool, bool, error) {
	folderArgs := core.SuObjectOf(core.SuStr(library))
	folderArgs.Set(core.SuStr("group"), core.IntVal(parent))
	folderArgs.Set(core.SuStr("name"), core.SuStr(name))
	row, hdr, _ := tran.Get(th, folderArgs, core.Only)
	if row != nil {
		num, err := intValue(row.GetVal(hdr, "num", th, nil), "num")
		if err != nil {
			return 0, false, false, err
		}
		return num, true, false, nil
	}

	leafArgs := core.SuObjectOf(core.SuStr(library))
	leafArgs.Set(core.SuStr("group"), core.IntVal(-1))
	leafArgs.Set(core.SuStr("parent"), core.IntVal(parent))
	leafArgs.Set(core.SuStr("name"), core.SuStr(name))
	leaf, _, _ := tran.Get(th, leafArgs, core.Only)
	return 0, false, leaf != nil, nil
}

// maxLibNum returns the maximum num value in the library, or 0 if none.
func maxLibNum(th *core.Thread, tran core.ITran, library string) int {
	q := tran.Query(library+" summarize max num", nil)
	hdr := q.Header()
	st := core.NewSuTran(tran, false)
	row, _ := q.Get(th, core.Next)
	if row == nil {
		return 0
	}
	v := row.GetVal(hdr, "max_num", th, st)
	if n, ok := v.IfInt(); ok {
		return n
	}
	return 0
}

// buildRecord builds a Record from a map of field values, using the header's
// field order.
func buildRecord(hdr *core.Header, vals map[string]core.Value) core.Record {
	fields := hdr.Fields[0]
	rb := core.RecordBuilder{}
	for _, f := range fields {
		if f == "-" {
			rb.AddRaw("") // deleted/dropped column placeholder
			continue
		}
		v, ok := vals[f]
		if !ok || v == nil {
			rb.AddRaw("")
		} else {
			rb.AddRaw(core.PackValue(v))
		}
	}
	return rb.Trim().Build()
}

// validateLibCode compiles the code as a library definition to check for errors.
func validateLibCode(th *core.Thread, text string) (warnings []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid code: %v", r)
		}
	}()
	_, warnings = compile.Checked(th, text)
	return warnings, nil
}
