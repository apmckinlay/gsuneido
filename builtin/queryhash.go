// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
)

var _ = builtin(QueryHash, "(query, details=false)")

func QueryHash(th *Thread, args []Value) Value {
	query := ToStr(args[0]) + `
		/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */
		/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */
		/* CHECKQUERY SUPPRESS: JOIN MANY TO MANY */`
	details := ToBool(args[1])
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	q := tran.Query(query, nil)
	// fmt.Println(q.Strategy(true))
	qh := qry.NewQueryHasher(q.Header()).CheckDups()

	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		qh.Row(row)
	}
	return qh.Result(details)
}
