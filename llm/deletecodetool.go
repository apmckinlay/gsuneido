// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"slices"

	"github.com/apmckinlay/gsuneido/core"
)

var _ = addTool(toolSpec{
	name:        "suneido_delete_code",
	description: "Delete a library definition by name",
	params: []stringParam{
		{name: "library", description: "Name of the library (e.g. 'stdlib')", required: true, kind: paramString},
		{name: "name", description: "Name of the definition (e.g. 'Alert')", required: true, kind: paramString},
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
		return deleteCodeTool(ctx, library, name)
	},
})

type deleteCodeOutput struct {
	Library string `json:"library" jsonschema:"Library name"`
	Name    string `json:"name" jsonschema:"Definition name"`
	Action  string `json:"action" jsonschema:"Always 'deleted'"`
}

func deleteCodeTool(ctx context.Context, library, name string) (deleteCodeOutput, error) {
	if !isValidName(name) {
		return deleteCodeOutput{}, fmt.Errorf("invalid name: %s", name)
	}

	th := core.NewThread(core.MainThread)
	defer th.Close()

	if err := validateLibrary(th, library); err != nil {
		return deleteCodeOutput{}, err
	}

	query := fmt.Sprintf("%s where group = -1 and name = %q", library, name)
	rtran := th.Dbms().Transaction(false)
	rq := rtran.Query(query, nil)
	hdr := rq.Header()
	row, _ := rq.Get(th, core.Next)
	if row == nil {
		rtran.Complete()
		return deleteCodeOutput{}, fmt.Errorf("code not found for: %s in %s", name, library)
	}

	st := core.NewSuTran(rtran, false)
	vals := make(map[string]string, len(hdr.Fields[0]))
	for _, f := range hdr.Fields[0] {
		vals[f] = row.GetRawVal(hdr, f, th, st)
	}

	softDelete := false
	if slices.Contains(hdr.Fields[0], "lib_committed") {
		if v := vals["lib_committed"]; v != "" {
			softDelete = true
		}
	}

	off := row[0].Off
	rtran.Complete()

	if err := requireApproval(ctx, "deleteCodeTool"); err != nil {
		return deleteCodeOutput{}, err
	}

	utran := th.Dbms().Transaction(true)
	action := "deleted"
	if softDelete {
		vals["group"] = core.PackValue(core.SuInt(-2))
		vals["lib_modified"] = core.PackValue(core.Now())

		if slices.Contains(hdr.Fields[0], "lib_before_text") {
			if v := vals["lib_before_text"]; v == "" {
				vals["lib_before_text"] = vals["text"]
			}
		}
		if slices.Contains(hdr.Fields[0], "lib_before_path") {
			if v := vals["lib_before_path"]; v == "" {
				vals["lib_before_path"] = vals["path"]
			}
		}

		newRec := buildRecord(hdr, vals)
		utran.Update(th, library, off, newRec)
		action = "soft-deleted"
	} else {
		utran.Delete(th, library, off)
	}
	if conflict := utran.Complete(); conflict != "" {
		return deleteCodeOutput{}, fmt.Errorf("transaction conflict: %s", conflict)
	}

	core.Global.Unload(name)
	return deleteCodeOutput{Library: library, Name: name, Action: action}, nil
}
