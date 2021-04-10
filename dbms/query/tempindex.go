// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sortlist"
	"github.com/apmckinlay/gsuneido/util/str"
)

type TempIndex struct {
	Query1
	order   []string
	tran    QueryTran
	iter    *sortlist.Iter
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

func (ti *TempIndex) Get(dir runtime.Dir) runtime.Row {
	if ti.iter == nil {
		ti.create()
	}
	var off uint64
	if dir == runtime.Next {
		off = ti.iter.Next()
	} else {
		off = ti.iter.Prev()
	}
	if off == 0 {
		return nil
	}
	dbrec := runtime.DbRec{Record: ti.tran.GetRecord(off), Off: off}
	return runtime.Row{dbrec}
}

func (ti *TempIndex) create() {
	b := sortlist.NewSorting(ti.tran.MakeCompare(ti.ixspec()))
	for {
		row := ti.source.Get(runtime.Next)
		if row == nil {
			break
		}
		b.Add(row[0].Off) //TODO handle multiple & derived
	}
	ti.iter = b.Finish().Iter()
}

func (ti *TempIndex) ixspec() *ixkey.Spec {
	fields := ti.source.Header().Fields[0] //TODO handle multiple
	flds := make([]int, len(fields))
	for i, f := range ti.order {
		fi := str.List(fields).Index(f)
		assert.That(fi >= 0)
		flds[i] = fi
	}
	return &ixkey.Spec{Fields: flds}
}
