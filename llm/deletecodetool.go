// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"

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
	row, _ := rq.Get(th, core.Next)
	if row == nil {
		rtran.Complete()
		return deleteCodeOutput{}, fmt.Errorf("code not found for: %s in %s", name, library)
	}
	off := row[0].Off
	rtran.Complete()

	if approvalFn, ok := ctx.Value(approvalFnKey{}).(func() (bool, error)); ok {
		allowed, err := approvalFn()
		if err != nil {
			return deleteCodeOutput{}, err
		}
		if !allowed {
			return deleteCodeOutput{}, fmt.Errorf("DENIED")
		}
	}

	utran := th.Dbms().Transaction(true)
	utran.Delete(th, library, off)
	if conflict := utran.Complete(); conflict != "" {
		return deleteCodeOutput{}, fmt.Errorf("transaction conflict: %s", conflict)
	}

	core.Global.Unload(name)
	return deleteCodeOutput{Library: library, Name: name, Action: "deleted"}, nil
}
