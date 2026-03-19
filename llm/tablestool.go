// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
)

// tables
var _ = addTool(toolSpec{
	name:        "suneido_tables",
	description: "Get a list of database table names that start with the given prefix (limit of 100)",
	params:      []stringParam{{name: "prefix", description: "Only return tables whose names start with this prefix (empty string for all)", required: true}},
	summarize: func(args map[string]any) string {
		return mdSummary("Tables", argReqStr(args, "prefix"))
	},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		prefix, err := requireString(args, "prefix")
		if err != nil {
			return nil, err
		}
		return tablesTool(prefix)
	},
})

type tablesOutput struct {
	Tables  []string `json:"tables" jsonschema:"Table names matching the requested prefix"`
	HasMore bool     `json:"has_more,omitempty" jsonschema:"True when additional tables were truncated"`
}

func tablesTool(prefix string) (output tablesOutput, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tables failed: %v", r)
		}
	}()

	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()

	q := tran.Query(tablesQuery(prefix), nil)
	hdr := q.Header()
	for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
		if len(output.Tables) >= queryLimit {
			output.HasMore = true
			break
		}
		v := row.GetVal(hdr, "table", nil, nil)
		s, ok := v.ToStr()
		if !ok {
			panic("tablesTool: expected table name string")
		}
		if prefix != "" && !strings.HasPrefix(s, prefix) {
			break
		}
		output.Tables = append(output.Tables, s)
	}
	return output, nil
}

func tablesQuery(prefix string) string {
	if prefix == "" {
		return "tables sort table"
	}
	lo := strconv.Quote(prefix)
	return "tables where table >= " + lo + " project table sort table"
}
