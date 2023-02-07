// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"hash/adler32"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

var _ = builtin(QueryHash, "(query, details=false)")

func QueryHash(th *Thread, args []Value) Value {
	query := ToStr(args[0])
	details := ToBool(args[1])
	tran := th.Dbms().Transaction(false)
	defer tran.Complete()
	q := tran.Query(query, nil)
	hdr := q.Header()
	// fmt.Println(hdr)
	fields := slc.Without(hdr.GetFields(), "-")
	slices.Sort(fields)
	// fmt.Println("Fields", fields)
	colhash := hashCols(hdr)
	hash := colhash
	// type fmtable interface{ Format() string }
	// fmt.Println(q.(fmtable).Format())
	n := 0
	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		hash += hashRow(hdr, fields, row)
		// fmt.Println("row", row)
		n++
		// fmt.Println(n, hash)
		// if n >= 10 {
		// 	break
		// }
	}
	if details {
		return SuStr(fmt.Sprintln("nrows", n, "hash", hash) +
			fmt.Sprintln(colhash, str.Join(", ", hdr.Columns)))
	}
	return IntVal(int(hash))
}

func hashCols(hdr *Header) uint32 {
	cols := slices.Clone(hdr.Columns)
	slices.Sort(cols)
	hash := uint32(31)
	for _, col := range cols {
		hash = hash*31 + adler32.Checksum(hacks.Stobs(col))
	}
	return hash
}

func hashRow(hdr *Header, fields []string, row Row) uint32 {
	hash := uint32(0)
	// fmt.Print(">>> ")
	for _, fld := range fields {
		hash = hash*31 + adler32.Checksum(hacks.Stobs(row.GetRaw(hdr, fld)))
		// fmt.Print(fld, ": ", Unpack(row.GetRaw(hdr, fld)), " ")
	}
	// fmt.Println()
	return hash
}
