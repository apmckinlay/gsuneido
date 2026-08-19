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
	description: "Execute a Suneido database query and return the results as Suneido-format text (Value.String) in a simple row/column array format",
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
const querySizeLimit = 10000

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
	truncated := false
	st := core.NewSuTran(tran, false)
	nrows := 0
	qr := querier{hdr: hdr, th: th, st: st}
	qr.sb.WriteString("[\n")
	qr.formatQueryHeader()
	for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
		nrows++
		if nrows > queryLimit || !qr.formatQueryRow(row) {
			qr.sb.WriteString("// too large, truncated\n")
			truncated = true
			break
		}
	}
	qr.sb.WriteString("]\n")
	result = queryOutput{
		Query:   query,
		Results: qr.sb.String(),
		HasMore: truncated,
	}
	return result, nil
}

type querier struct {
	hdr *core.Header
	th  *core.Thread
	st  *core.SuTran
	sb  strings.Builder
}

func (qr *querier) formatQueryHeader() {
	sep := "["
	for _, col := range qr.hdr.Columns {
		qr.sb.WriteString(sep)
		sep = ", "
		qr.sb.WriteString(strconv.Quote(col))
	}
	qr.sb.WriteString("]\n")
}

func (qr *querier) formatQueryRow(row core.Row) bool {
	sep := "["
	for _, col := range qr.hdr.Columns {
		v := row.GetVal(qr.hdr, col, qr.th, qr.st)
		s := ""
		if v == nil {
			s = "null"
		} else {
			s = v.String()
		}
		if qr.sb.Len()+len(s) > querySizeLimit {
			qr.sb.WriteString("]\n")
			return false
		}
		qr.sb.WriteString(sep)
		sep = ", "
		qr.sb.WriteString(s)
	}
	qr.sb.WriteString("]\n")
	return true
}
