// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/shmap"
)

var _ = builtin(QueryHash, "(query, details=false)")

type rowHash struct {
	row  Row
	hash uint64
}

func QueryHash(th *Thread, args []Value) Value {
	query := ToStr(args[0]) + `
		/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */
		/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */
		/* CHECKQUERY SUPPRESS: JOIN MANY TO MANY */`
	details := ToBool(args[1])
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	q := tran.Query(query, nil)
	qh := qry.NewQueryHasher(q.Header())

	hfn := func(row rowHash) uint64 { return row.hash }
	eqfn := func(x, y rowHash) bool {
		return x.hash == y.hash && equalRow(x.row, y.row, qh.Hdr, qh.Fields)
	}
	rows := shmap.NewMapFuncs[rowHash, struct{}](hfn, eqfn)

	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		rh := rowHash{row: row, hash: qh.Row(row)}
		if _, exists := rows.GetInit(rh); exists {
			panic("QueryHash: duplicate row")
		}
	}
	return qh.Result(details)
}

func equalRow(x, y Row, hdr *Header, cols []string) bool {
	for _, col := range cols {
		if x.GetRaw(hdr, col) != y.GetRaw(hdr, col) {
			return false
		}
	}
	return true
}
