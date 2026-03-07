// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"

	"github.com/apmckinlay/gsuneido/core"
)

// schema
var _ = addTool(toolSpec{
	name:        "suneido_schema",
	description: "Get the schema for a Suneido database table, or the definition for a view",
	params:      []stringParam{{name: "table", description: "Name of the table or view to get schema for", required: true}},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		table, err := requireString(args, "table")
		if err != nil {
			return nil, err
		}
		schema := core.GetDbms().Schema(table)
		if schema == "" {
			if viewDef := getViewDefinition(table); viewDef != "" {
				return schemaOutput{Schema: "view " + table + " = " + viewDef}, nil
			}
		}
		return schemaOutput{Schema: schema}, nil
	},
})

func getViewDefinition(name string) string {
	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	q := tran.Query(fmt.Sprintf("views where view_name = %q", name), nil)
	defer q.Close()
	row, _ := q.Get(th, core.Next)
	if row == nil {
		return ""
	}
	hdr := q.Header()
	st := core.NewSuTran(tran, false)
	val := row.GetVal(hdr, "view_definition", th, st)
	return core.ToStr(val)
}

type schemaOutput struct {
	Schema string `json:"schema" jsonschema:"Schema definition for the requested table"`
}
