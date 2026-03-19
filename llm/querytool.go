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

var _ = addTool(toolSpec{
	name:        "suneido_query",
	description: "Execute a Suneido database query and return the results as Suneido-format text (Value.String) in a simple row/column array format (limit 100)",
	params: []stringParam{
		{name: "query", description: "Suneido query (e.g. 'tables sort table')", required: true},
	},
	summarize: func(args map[string]any) string {
		query := argString(args, "query")
		trimmed := strings.TrimSpace(query)
		if strings.Contains(trimmed, "\n") || strings.Contains(trimmed, "\r") {
			return mdSummary("Query") + "\n" + summarizeCodeBlock(query)
		}
		return mdSummary("Query", mdInline(trimmed))
	},
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		qs, err := requireString(args, "query")
		if err != nil {
			return nil, err
		}
		return queryTool(qs)
	},
})

type queryOutput struct {
	Query   string `json:"query" jsonschema:"Query string that was executed"`
	Results string `json:"results" jsonschema:"Formatted row/column output"`
	HasMore bool   `json:"has_more,omitempty" jsonschema:"True when additional rows were truncated"`
}

const queryLimit = 100

func queryTool(query string) (result queryOutput, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("query failed: %v", r)
		}
	}()
	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	q := tran.Query(query, nil)
	hdr := q.Header()
	cols := hdr.Columns
	var rows [][]core.Value
	truncated := false
	st := core.NewSuTran(tran, false)
	for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
		if len(rows) >= queryLimit {
			truncated = true
			break
		}
		vals := make([]core.Value, len(cols))
		for i, col := range cols {
			vals[i] = row.GetVal(hdr, col, th, st)
		}
		rows = append(rows, vals)
	}
	result = queryOutput{
		Query:   query,
		Results: formatQueryResult(cols, rows, truncated),
		HasMore: truncated,
	}
	return result, nil
}

func formatQueryResult(cols []string, rows [][]core.Value, truncated bool) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	sb.WriteString(formatQueryHeader(cols))
	sb.WriteString("\n")
	for _, row := range rows {
		sb.WriteString(formatQueryRow(row))
		sb.WriteString("\n")
	}
	if truncated {
		sb.WriteString("// truncated\n")
	}
	sb.WriteString("]\n")
	return sb.String()
}

func formatQueryHeader(cols []string) string {
	qs := make([]string, len(cols))
	for i, col := range cols {
		qs[i] = strconv.Quote(col)
	}
	return "[" + strings.Join(qs, ", ") + "]"
}

func formatQueryRow(row []core.Value) string {
	ss := make([]string, len(row))
	for i, v := range row {
		if v == nil {
			ss[i] = "null"
		} else {
			ss[i] = v.String()
		}
	}
	return "[" + strings.Join(ss, ", ") + "]"
}
