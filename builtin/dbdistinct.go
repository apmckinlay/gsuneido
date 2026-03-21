// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/hll"
)

var _ = builtin(DbDistinct, "(table)")

func DbDistinct(th *Thread, args []Value) Value {
	table := ToStr(args[0])
	t := th.Dbms().Transaction(false)
	defer t.Complete()
	rt := t.(*dbms.ReadTranLocal).ReadTran
	cols := indexedColumns(rt.GetSchema(table).Indexes)
	hdr := SimpleHeader(cols)
	sketches := make([]*hll.HLL, len(cols))
	for i := range cols {
		sketches[i] = hll.New()
	}
	q := t.Query(table, nil)
	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		for i, col := range cols {
			sketches[i].Add(row.GetRawVal(hdr, col, nil, nil))
		}
	}
	ob := &SuObject{}
	for i, col := range cols {
		ob.Set(SuStr(col), Int64Val(int64(sketches[i].Count())))
	}
	return ob
}

func indexedColumns(indexes []schema.Index) []string {
	cols := make([]string, 0, len(indexes))
	seen := make(map[string]struct{}, len(indexes))
	for _, ix := range indexes {
		for _, col := range ix.Columns {
			if strings.HasSuffix(col, "_lower!") {
				continue
			}
			if _, ok := seen[col]; ok {
				continue
			}
			seen[col] = struct{}{}
			cols = append(cols, col)
		}
	}
	return cols
}
