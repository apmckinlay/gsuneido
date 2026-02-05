// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
)

func tablesTool(prefix string) (tables []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tables failed: %v", r)
		}
	}()

	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := core.GetDbms().Transaction(false)
	defer tran.Complete()

	q := tran.Query(tablesQuery(prefix), nil)
	hdr := q.Header()
	for row, _ := q.Get(th, core.Next); row != nil; row, _ = q.Get(th, core.Next) {
		if len(tables) >= queryLimit {
			tables = append(tables, "...")
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
		tables = append(tables, s)
	}
	return tables, nil
}

func tablesQuery(prefix string) string {
	if prefix == "" {
		return "tables sort table"
	}
	lo := strconv.Quote(prefix)
	return "tables where table >= " + lo + " project table sort table"
}
