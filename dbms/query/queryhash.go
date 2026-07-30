// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/shmap"
	"github.com/apmckinlay/gsuneido/util/slc"
)

type QueryHash struct {
	Hdr      *Header
	Fields   []string
	rows     *shmap.Map[rowHash, struct{}, shmap.Funcs[rowHash]]
	ncols    int
	colsHash uint64
	nrows    int
	hash     uint64
}

func NewQueryHasher(hdr *Header) *QueryHash {
	qh := QueryHash{}
	qh.Hdr = hdr
	qh.Fields = slc.Clone(hdr.Physical())
	slices.Sort(qh.Fields)
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

func (qh *QueryHash) CheckDups() *QueryHash {
	// Note: hash and equal do NOT evaluate rules
	hfn := func(row rowHash) uint64 { return row.hash }
	eqfn := func(x, y rowHash) bool {
		return x.hash == y.hash && equalRow(x.row, y.row, qh.Hdr, qh.Fields)
	}
	qh.rows = shmap.NewMapFuncs[rowHash, struct{}](hfn, eqfn)
	return qh
}

func (qh *QueryHash) Row(row Row) uint64 {
	hash := uint64(17)
	for _, fld := range qh.Fields {
		s := row.GetRaw(qh.Hdr, fld)
		// ignore PackForward
		if len(s) > 0 && s[0] != PackForward {
			hash = hash*31 + hashPacked(s)
		}
	}
	if qh.rows != nil {
		rh := rowHash{row: row, hash: hash}
		if _, exists := qh.rows.GetInit(rh); exists {
			panic("QueryHash: duplicate row")
		}
	}
	//TODO order sensitive if sorted
	qh.hash += hash // '+' to ignore order
	qh.nrows++
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

func (qh *QueryHash) Result(details bool) Value {
	if details {
		return SuStr(fmt.Sprintln("nrows", qh.nrows, "hash", qh.hash,
			"ncols", qh.ncols, "hash", qh.colsHash))
	}
	return IntVal(int(qh.hash))
}

func queryHashAll(db *db19.Database, query string) {
	tran := db.NewReadTran()
	q := ParseQuery(query, tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	th := &Thread{}

	// fmt.Println("=== Simple ===")
	h1 := NewQueryHasher(q.Header()).CheckDups()
	for _, row := range q.Simple(th) {
		// fmt.Println(row)
		h1.Row(row)
	}

	// fmt.Println("== Get ===")
	h2 := NewQueryHasher(q.Header()).CheckDups()
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		// fmt.Println(row)
		h2.Row(row)
	}

	// fmt.Println("optimized:", String(q))
	// fmt.Println("Fields:", q.Header().Fields)
	// fmt.Println("Columns:", q.Header().Columns)
	assert.This(h2.Result(true)).Is(h1.Result(true))
}
