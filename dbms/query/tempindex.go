// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sortlist"
	"github.com/apmckinlay/gsuneido/util/strs"
)

type TempIndex struct {
	Query1
	order   []string
	tran    QueryTran
	hdr     *Header
	iter    rowIter
	rewound bool
	selOrg  string
	selEnd  string
}

func (ti *TempIndex) String() string {
	return parenQ2(ti.source) + " TEMPINDEX" + strs.Join("(,)", ti.order)
}

func (ti *TempIndex) Transform() Query {
	return ti
}

// execution --------------------------------------------------------

func (ti *TempIndex) Rewind() {
	if ti.iter != nil {
		ti.iter.Rewind()
	}
	ti.rewound = true
}

func (ti *TempIndex) Select(cols, vals []string) {
	encode := len(ti.order) > 1
	ti.selOrg, ti.selEnd = selKeys(encode, ti.order, cols, vals)
	ti.Rewind()
}

func (ti *TempIndex) Lookup(key string) Row {
	if ti.iter == nil {
		ti.iter = ti.makeIter()
	}
	row := ti.iter.Seek(key)
	if row == nil || ti.rowKey(row) != key {
		return nil
	}
	return row
}

func (ti *TempIndex) Get(dir Dir) Row {
	if ti.iter == nil {
		ti.iter = ti.makeIter()
		if ti.selEnd == "" {
			ti.selEnd = ixkey.Max
		}
		ti.rewound = true
	}
	var row Row
	if ti.rewound {
		if dir == Next {
			row = ti.iter.Seek(ti.selOrg)
		} else { // Prev
			row = ti.iter.Seek(ti.selEnd)
			if row == nil {
				ti.iter.Rewind()
			}
			row = ti.iter.Get(dir)
		}
		ti.rewound = false
	} else {
		row = ti.iter.Get(dir)
	}
	if row == nil || !ti.selected(row) {
		return nil
	}
	return row
}

type rowIter interface {
	Get(Dir) Row
	Rewind()
	Seek(key string) Row
}

func (ti *TempIndex) makeIter() rowIter {
	ti.hdr = ti.source.Header()
	if ti.source.SingleTable() {
		return ti.single()
	}
	return ti.multi()
}

func (ti *TempIndex) selected(row Row) bool {
	key := ti.rowKey(row)
	return ti.selOrg <= key && key < ti.selEnd
}

func (ti *TempIndex) rowKey(row Row) string {
	if len(ti.order) == 1 {
		return row.GetRaw(ti.hdr, ti.order[0])
	}
	var enc ixkey.Encoder
	for _, col := range ti.order {
		enc.Add(row.GetRaw(ti.hdr, col))
	}
	return enc.String()
}

//-------------------------------------------------------------------

type singleIter struct {
	tran QueryTran
	iter *sortlist.Iter
}

func (ti *TempIndex) single() rowIter {
	spec := ti.ixspec()
	b := sortlist.NewSorting(ti.tran.MakeLess(spec))
	for {
		row := ti.source.Get(Next)
		if row == nil {
			break
		}
		b.Add(row[0].Off)
	}
	less := func(off uint64, key string) bool {
		rec := ti.tran.GetRecord(off)
		return spec.Key(rec) < key
	}
	return &singleIter{tran: ti.tran, iter: b.Finish().Iter(less)}
}

func (ti *TempIndex) ixspec() *ixkey.Spec {
	fields := ti.hdr.Fields[0]
	flds := make([]int, len(ti.order))
	for i, f := range ti.order {
		fi := strs.Index(fields, f)
		assert.That(fi >= 0)
		flds[i] = fi
	}
	return &ixkey.Spec{Fields: flds}
}

func (it *singleIter) Get(dir Dir) Row {
	if dir == Next {
		it.iter.Next()
	} else {
		it.iter.Prev()
	}
	if it.iter.Eof() {
		return nil
	}
	return it.get()
}

func (it *singleIter) get() Row {
	if it.iter.Eof() {
		return nil
	}
	off := it.iter.Cur()
	dbrec := DbRec{Record: it.tran.GetRecord(off), Off: off}
	return Row{dbrec}
}

func (it *singleIter) Seek(key string) Row {
	it.iter.Seek(key)
	if it.iter.Eof() {
		return nil
	}
	return it.get()
}

func (it *singleIter) Rewind() {
	it.iter.Rewind()
}

//-------------------------------------------------------------------

type multiIter struct {
	tran  QueryTran
	nrecs int
	heap  *stor.Stor
	hdr   *Header
	order []string
	iter  *sortlist.Iter
}

func (ti *TempIndex) multi() rowIter {
	it := multiIter{tran: ti.tran, hdr: ti.hdr, order: ti.order,
		nrecs: len(ti.hdr.Fields), heap: stor.HeapStor(8192)}
	it.heap.Alloc(1) // avoid offset 0
	b := sortlist.NewSorting(it.multiLess)
	for {
		row := ti.source.Get(Next)
		if row == nil {
			break
		}
		assert.That(len(row) == it.nrecs)
		n := it.nrecs * stor.SmallOffsetLen
		for _, dbrec := range row {
			if dbrec.Off == 0 { // derived record e.g. from extend
				n += len(dbrec.Record)
			}
		}
		off, buf := it.heap.Alloc(n)
		for _, dbrec := range row {
			if dbrec.Off > 0 {
				stor.WriteSmallOffset(buf, dbrec.Off)
				buf = buf[stor.SmallOffsetLen:]
			} else { // derived
				size := len(dbrec.Record)
				stor.WriteSmallOffset(buf, multiMask|uint64(size))
				buf = buf[stor.SmallOffsetLen:]
				copy(buf, dbrec.Record)
				buf = buf[size:]
			}
		}
		b.Add(off)
	}
	less := func(off uint64, key string) bool {
		row := it.get(off)
		return ti.rowKey(row) < key
	}
	it.iter = b.Finish().Iter(less)
	return &it
}

const multiMask = 0xffff000000

func (it *multiIter) Seek(key string) Row {
	it.iter.Seek(key)
	if it.iter.Eof() {
		return nil
	}
	return it.get(it.iter.Cur())
}

func (it *multiIter) Rewind() {
	it.iter.Rewind()
}

func (it *multiIter) Get(dir Dir) Row {
	if dir == Next {
		it.iter.Next()
	} else {
		it.iter.Prev()
	}
	if it.iter.Eof() {
		return nil
	}
	return it.get(it.iter.Cur())
}

func (it *multiIter) get(off uint64) Row {
	row := make([]DbRec, it.nrecs)
	buf := it.heap.Data(off)
	var rec Record
	for i := 0; i < it.nrecs; i++ {
		buf, rec, off = it.getRec(buf)
		row[i] = DbRec{Record: rec, Off: off}
	}
	return row
}

func (it *multiIter) getRec(buf []byte) ([]byte, Record, uint64) {
	off := stor.ReadSmallOffset(buf)
	buf = buf[stor.SmallOffsetLen:]
	var rec Record
	if off < multiMask {
		rec = it.tran.GetRecord(off)
	} else { // derived
		size := off &^ multiMask
		rec = Record(buf[:size])
		buf = buf[size:]
		off = 0
	}
	return buf, rec, off
}

func (it *multiIter) multiLess(x, y uint64) bool {
	xrow := make(Row, it.nrecs)
	yrow := make(Row, it.nrecs)
	xbuf := it.heap.Data(x)
	ybuf := it.heap.Data(y)
	for i := 0; i < it.nrecs; i++ {
		xbuf, xrow[i].Record, _ = it.getRec(xbuf)
		ybuf, yrow[i].Record, _ = it.getRec(ybuf)
	}
	for _, col := range it.order {
		x := xrow.GetRaw(it.hdr, col)
		y := yrow.GetRaw(it.hdr, col)
		if x != y {
			if x < y {
				return true
			}
			return false // >
		}
	}
	return false
}
