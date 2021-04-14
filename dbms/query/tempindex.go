// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sortlist"
	"github.com/apmckinlay/gsuneido/util/str"
)

type TempIndex struct {
	Query1
	order []string
	tran  QueryTran
	iter  rowIter
}

func (ti *TempIndex) String() string {
	return parenQ2(ti.source) + " TEMPINDEX" + str.Join("(,)", ti.order)
}

func (ti *TempIndex) Transform() Query {
	return ti
}

// execution --------------------------------------------------------

func (ti *TempIndex) Rewind() {
	if ti.iter != nil {
		ti.iter.Rewind()
	}
	ti.source.Rewind()
}

func (ti *TempIndex) Get(dir Dir) Row {
	if ti.iter == nil {
		if ti.source.SingleTable() {
			ti.iter = ti.single()
		} else {
			ti.iter = ti.multi()
		}
	}
	return ti.iter.Get(dir)
}

type rowIter interface {
	Get(Dir) Row
	Rewind()
}

type singleIter struct {
	tran QueryTran
	iter *sortlist.Iter
}

func (ti *TempIndex) single() rowIter {
	b := sortlist.NewSorting(ti.tran.MakeLess(ti.ixspec()))
	for {
		row := ti.source.Get(Next)
		if row == nil {
			break
		}
		b.Add(row[0].Off)
	}
	return &singleIter{tran: ti.tran, iter: b.Finish().Iter()}
}

func (ti *TempIndex) ixspec() *ixkey.Spec {
	fields := ti.source.Header().Fields[0]
	flds := make([]int, len(fields))
	for i, f := range ti.order {
		fi := str.List(fields).Index(f)
		assert.That(fi >= 0)
		flds[i] = fi
	}
	return &ixkey.Spec{Fields: flds}
}

func (it singleIter) Get(dir Dir) Row {
	var off uint64
	if dir == Next {
		off = it.iter.Next()
	} else {
		off = it.iter.Prev()
	}
	if off == 0 {
		return nil
	}
	dbrec := DbRec{Record: it.tran.GetRecord(off), Off: off}
	return Row{dbrec}
}

func (it singleIter) Rewind() {
	it.iter.Rewind()
}

type multiIter struct {
	tran   QueryTran
	nrecs  int
	heap   *stor.Stor
	fields []RowAt
	iter   *sortlist.Iter
}

func (ti *TempIndex) multi() rowIter {
	it := multiIter{tran: ti.tran}
	hdr := ti.source.Header()
	it.fields = make([]RowAt, len(ti.order))
	for i, f := range ti.order {
		it.fields[i] = hdr.Map[f]
	}
	it.nrecs = len(hdr.Fields)
	it.heap = stor.HeapStor(8192)
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
	it.iter = b.Finish().Iter()
	return &it
}

const multiMask = 0xffff000000

func (it *multiIter) Rewind() {
	it.iter.Rewind()
}

func (it *multiIter) Get(dir Dir) Row {
	var off uint64
	if dir == Next {
		off = it.iter.Next()
	} else {
		off = it.iter.Prev()
	}
	if off == 0 {
		return nil
	}
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
	xrow := make([]Record, it.nrecs)
	yrow := make([]Record, it.nrecs)
	xbuf := it.heap.Data(x)
	ybuf := it.heap.Data(y)
	for i := 0; i < it.nrecs; i++ {
		xbuf, xrow[i], _ = it.getRec(xbuf)
		ybuf, yrow[i], _ = it.getRec(ybuf)
	}
	for _, at := range it.fields {
		x := xrow[at.Reci].GetRaw(int(at.Fldi))
		y := yrow[at.Reci].GetRaw(int(at.Fldi))
		if x != y {
			if x < y {
				return true
			}
			return false // >
		}
	}
	return false
}
