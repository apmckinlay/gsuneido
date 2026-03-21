// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"cmp"
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/ss"
)

var _ = builtin(DbTop10, "(table, column)")

func DbTop10(th *Thread, args []Value) Value {
	table := ToStr(args[0])
	column := ToStr(args[1])

	tran := th.Dbms().Transaction(false)
	defer tran.Complete()

	sk := ss.New(128)
	q := tran.Query(table, nil)
	hdr := q.Header()
	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		sk.Add(row.GetRawVal(hdr, column, nil, nil))
	}

	top := sk.Top()
	slices.SortFunc(top, func(a, b ss.Entry) int {
		return -cmp.Compare(a.Count-a.Error, b.Count-b.Error)
	})
	if len(top) > 10 {
		top = top[:10]
	}

	result := &SuObject{}
	for _, e := range top {
		result.Set(Unpack(e.Value), IntVal(e.Count-e.Error))
	}
	return result
}
