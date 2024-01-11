// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/dbms"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
)

var _ = builtin(QueryAlt, "(@args)")

func QueryAlt(th *Thread, as *ArgSpec, args []Value) Value {
	query, _ := extractQuery(th, &queryParams, as, args)
	// this will only work if dbms is local
	t := th.Dbms().Transaction(false).(*dbms.ReadTranLocal).ReadTran
	q := qry.ParseQuery(query, t, nil)
	hdr := q.Header()
	rows := q.Simple(th)
	ob := &SuObject{}
	for _, row := range rows {
		ob.Add(SuRecordFromRow(row, hdr, "", nil))
	}
	return ob
}

var _ = builtin(QueryAltHash, "(query, details=false)")

func QueryAltHash(th *Thread, args []Value) Value {
	query := ToStr(args[0]) + `
		/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */
		/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */
		/* CHECKQUERY SUPPRESS: JOIN MANY TO MANY */`
	details := ToBool(args[1])
	t := th.Dbms().Transaction(false).(*dbms.ReadTranLocal).ReadTran
	q := qry.ParseQuery(query, t, nil)
	qh := NewQueryHasher(q.Header())
	rows := q.Simple(th)
	for _, row := range rows {
		qh.Row(row)
	}
	return qh.Result(details)
}
