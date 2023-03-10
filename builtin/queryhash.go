// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"hash/adler32"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/hmap"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"golang.org/x/exp/slices"
)

var _ = builtin(QueryHash, "(query, details=false)")

type rowHash struct {
	row  Row
	hash uint32
}

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

	hfn := func(row rowHash) uint32 { return row.hash }
	eqfn := func(x, y rowHash) bool {
		return x.hash == y.hash && equalRow(x.row, y.row, hdr, fields)
	}
	rows := hmap.NewHmapFuncs[rowHash, struct{}](hfn, eqfn)

	n := 0
	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		rh := rowHash{row, hashRow(hdr, fields, row)}
		_, _, ok := rows.GetPut(rh, struct{}{})
		if ok {
			panic("QueryHash: duplicate row")
		}
		hash += rh.hash
		// fmt.Println("row", row)
		n++
		// fmt.Println(n, hash)
		// if n >= 10 {
		// 	break
		// }
	}
	if details {
		return SuStr(fmt.Sprintln("nrows", n, "hash", hash,
			"ncols", len(hdr.Columns), "hash", colhash))
	}
	return IntVal(int(hash))
}

func hashCols(hdr *Header) uint32 {
	cols := slices.Clone(hdr.Columns)
	slices.Sort(cols)
	hash := uint32(31)
	for _, col := range cols {
		// fmt.Println(col)
		hash = hash*31 + adler32.Checksum(hacks.Stobs(col))
	}
	return hash
}

func hashRow(hdr *Header, fields []string, row Row) uint32 {
	hash := uint32(0)
	// fmt.Print(">>> ")
	for _, fld := range fields {
		hash = hash*31 + hashPacked(row.GetRaw(hdr, fld))
		// fmt.Print(fld, ": ", Unpack(row.GetRaw(hdr, fld)), " ")
	}
	// fmt.Println()
	return hash
}

func hashPacked(p string) uint32 {
	if len(p) > 0 && p[0] >= PackObject {
		return hashObject(p)
	}
	return adler32.Checksum(hacks.Stobs(p))
}

func hashObject(p string) uint32 {
	hash := uint32(0)
	for i := 0; i < len(p); i++ {
		// use simple addition to be insensitive to member order
		hash += uint32(p[i])
	}
	return hash
}

func equalRow(x, y Row, hdr *Header, cols []string) bool {
	for _, col := range cols {
		if x.GetRaw(hdr, col) != y.GetRaw(hdr, col) {
			return false
		}
	}
	return true
}
