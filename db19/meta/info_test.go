// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInfo(t *testing.T) {
	assert := assert.T(t)
	one := &Info{
		Table:   "one",
		Nrows:   100,
		Size:    1000,
		lastMod: -1,
		Indexes: []*index.Overlay{index.OverlayStub()},
	}
	assert.That(!one.IsTomb())
	two := &Info{
		Table:   "two",
		Nrows:   200,
		Size:    2000,
		lastMod: -1,
		Indexes: []*index.Overlay{index.OverlayStub()},
	}
	tbl := InfoHamt{}.Mutable()
	tbl.Put(one)
	tbl.Put(two)
	tbl = tbl.Freeze()

	st := stor.HeapStor(8192)
	st.Alloc(1) // avoid offset 0
	off := tbl.Write(st, 0, hamt.All)

	ic := hamt.ReadChain[string](st, off, ReadInfo)
	assert.This(ic.Ages[0]).Is(ic.MustGet("one").lastMod)
	tbl = ic.Hamt
	x, _ := tbl.Get("one")
	x.Indexes = nil
	one.Indexes = nil
	assert.This(*x).Is(*one)
	x, _ = tbl.Get("two")
	x.Indexes = nil
	two.Indexes = nil
	assert.This(*x).Is(*two)
}

func TestInfo2(t *testing.T) {
	tbl := InfoHamt{}.Mutable()
	const n = 1000
	data := mkdata(tbl, n)
	st := stor.HeapStor(32 * 1024)
	st.Alloc(1) // avoid offset 0
	off := tbl.Freeze().Write(st, 0, hamt.All)

	tbl = hamt.ReadChain[string](st, off, ReadInfo).Hamt
	for i, s := range data {
		ti, _ := tbl.Get(s)
		assert.T(t).Msg("table").This(ti.Table).Is(s)
		assert.T(t).Msg("nrows").This(ti.Nrows).Is(i)
		_, ok := tbl.Get(s + "Z")
		assert.T(t).That(!ok)
	}
}

func mkdata(tbl InfoHamt, n int) []string {
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := range n {
		data[i] = randStr()
		tbl.Put(&Info{Table: data[i], Nrows: i,
			Indexes: []*index.Overlay{index.OverlayStub()}})
	}
	return data
}
