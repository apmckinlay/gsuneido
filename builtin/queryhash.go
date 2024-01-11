// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"hash/adler32"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/generic/hmap"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

var _ = builtin(QueryHash, "(query, details=false)")

type rowHash struct {
	row  Row
	hash uint32
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
	qh := NewQueryHasher(q.Header())

	hfn := func(row rowHash) uint32 { return row.hash }
	eqfn := func(x, y rowHash) bool {
		return x.hash == y.hash && equalRow(x.row, y.row, qh.hdr, qh.fields)
	}
	rows := hmap.NewHmapFuncs[rowHash, struct{}](hfn, eqfn)

	for row, _ := q.Get(th, Next); row != nil; row, _ = q.Get(th, Next) {
		rh := rowHash{row: row, hash: qh.Row(row)}
		if _, _, exists := rows.GetPut(rh, struct{}{}); exists {
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

//-------------------------------------------------------------------

type queryHasher struct {
	hdr *Header
	fields []string
	ncols int
	colsHash uint32
	nrows int
	hash uint32
}

func NewQueryHasher(hdr *Header) *queryHasher {
	qh := queryHasher{}
	qh.hdr = hdr
	qh.fields = slc.Without(hdr.GetFields(), "-")
	slices.Sort(qh.fields)
	cols := slices.Clone(hdr.Columns)
	slices.Sort(cols)
	hash := uint32(31)
	for _, col := range cols {
		hash = hash*31 + adler32.Checksum(hacks.Stobs(col))
	}
	qh.ncols = len(cols)
	qh.colsHash = hash
	qh.hash = hash
	return &qh
}


func (qh *queryHasher) Row(row Row) uint32 {
	hash := uint32(0)
	for _, fld := range qh.fields {
		hash = hash*31 + hashPacked(row.GetRaw(qh.hdr, fld))
	}
	//TODO order sensitive if sorted
	qh.hash += hash // '+' to ignore order
	qh.nrows++
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

func (qh *queryHasher) Result(details bool) Value {
	if details {
		return SuStr(fmt.Sprintln("nrows", qh.nrows, "hash", qh.hash,
			"ncols", qh.ncols, "hash", qh.colsHash))
	}
	return IntVal(int(qh.hash))
}
