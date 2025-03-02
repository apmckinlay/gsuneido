// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hash"
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
	qh := NewQueryHasher(q.Header())

	hfn := func(row rowHash) uint64 { return row.hash }
	eqfn := func(x, y rowHash) bool {
		return x.hash == y.hash && equalRow(x.row, y.row, qh.hdr, qh.fields)
	}
	rows := shmap.NewMapFuncs[rowHash, struct{}](hfn, eqfn)

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
	hdr      *Header
	fields   []string
	ncols    int
	colsHash uint64
	nrows    int
	hash     uint64
}

func NewQueryHasher(hdr *Header) *queryHasher {
	qh := queryHasher{}
	qh.hdr = hdr
	qh.fields = slc.Without(hdr.GetFields(), "-")
	slices.Sort(qh.fields)
	cols := slc.Clone(hdr.Columns)
	slices.Sort(cols)
	h := uint64(17)
	for _, col := range cols {
		h = h*31 + hash.String(col)
	}
	qh.ncols = len(cols)
	qh.colsHash = h
	qh.hash = h
	return &qh
}

func (qh *queryHasher) Row(row Row) uint64 {
	hash := uint64(17)
	for _, fld := range qh.fields {
		s := row.GetRaw(qh.hdr, fld)
		if len(s) > 0 && s[0] == PackForward {
			s = ""
		}
		hash = hash*31 + hashPacked(s)
	}
	//TODO order sensitive if sorted
	qh.hash += hash // '+' to ignore order
	qh.nrows++
	return hash
}

func hashPacked(p string) uint64 {
	if len(p) > 0 && (p[0] == PackObject || p[0] == PackRecord) {
		return hashObject(p)
	}
	return hash.FullString(p)
}

func hashObject(p string) uint64 {
	hash := uint64(17)
	for i := range len(p) {
		// use simple addition to be insensitive to member order
		hash += uint64(p[i])
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
