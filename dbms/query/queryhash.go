// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"

	"slices"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/shmap"
	"github.com/apmckinlay/gsuneido/util/slc"
)

type QueryHash struct {
	Hdr      *core.Header
	Fields   []string
	rows     *shmap.Map[rowHash, struct{}, shmap.Funcs[rowHash]]
	ncols    int
	colsHash uint64
	nrows    int
	hash     uint64
}

func NewQueryHasher(hdr *core.Header) *QueryHash {
	qh := QueryHash{}
	qh.Hdr = hdr
	qh.Fields = slc.Without(hdr.GetFields(), "-")
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
	hfn := func(row rowHash) uint64 { return row.hash }
	eqfn := func(x, y rowHash) bool {
		return x.hash == y.hash && equalRow(x.row, y.row, qh.Hdr, qh.Fields)
	}
	qh.rows = shmap.NewMapFuncs[rowHash, struct{}](hfn, eqfn)
	return qh
}

func (qh *QueryHash) Row(row core.Row) uint64 {
	hash := uint64(17)
	for _, fld := range qh.Fields {
		s := row.GetRawVal(qh.Hdr, fld, nil, nil)
		hash = hash*31 + hashPacked(s)
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

func equalRow(x, y core.Row, hdr *core.Header, cols []string) bool {
	for _, col := range cols {
		if x.GetRaw(hdr, col) != y.GetRaw(hdr, col) {
			return false
		}
	}
	return true
}

func hashPacked(p string) uint64 {
	if len(p) > 0 && (p[0] == core.PackObject || p[0] == core.PackRecord) {
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

func (qh *QueryHash) Result(details bool) core.Value {
	if details {
		return core.SuStr(fmt.Sprintln("nrows", qh.nrows, "hash", qh.hash,
			"ncols", qh.ncols, "hash", qh.colsHash))
	}
	return core.IntVal(int(qh.hash))
}
