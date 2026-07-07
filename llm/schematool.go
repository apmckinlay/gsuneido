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
	summarize: func(args map[string]any) string {
		return mdSummary("Schema", argReqStr(args, "table"))
	},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		table, err := requireString(args, "table")
		if err != nil {
			return nil, err
		}
		return getSchema(table)
	},
})

func getSchema(table string) (any, error) {
	th := core.NewThread(core.MainThread)
	defer th.Close()
	dbms := th.Dbms()
	schema := dbms.Schema(table)
	if schema != "" {
		return schemaOutput{Schema: schema}, nil
	}

	tran := dbms.Transaction(false)
	defer tran.Complete()
	q := tran.Query(fmt.Sprintf("views where view_name = %q", table), nil)
	defer q.Close()
	row, _ := q.Get(th, core.Next)
	if row == nil {
		return schemaOutput{Schema: ""}, fmt.Errorf("table or view not found: %s", table)
	}
	hdr := q.Header()
	st := core.NewSuTran(tran, false)
	val := row.GetVal(hdr, "view_definition", th, st)
	return schemaOutput{Schema: "view " + table + " = " + core.ToStr(val)}, nil
}

type schemaOutput struct {
	Schema string `json:"schema" jsonschema:"Schema definition for the requested table"`
}
