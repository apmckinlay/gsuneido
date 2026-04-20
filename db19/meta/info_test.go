// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hamt"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInfo(t *testing.T) {
	assert := assert.T(t)
	one := &Info{
		Table:      "one",
		Nrows:      100,
		BtreeNrows: 100,
		Size:       1000,
		BtreeSize:  1000,
		Deltas:     []Delta{{}},
		lastMod:    -1,
		Indexes:    []*index.Overlay{index.OverlayStub()},
	}
	assert.That(!one.IsTomb())
	two := &Info{
		Table:      "two",
		Nrows:      200,
		BtreeNrows: 200,
		Size:       2000,
		BtreeSize:  2000,
		Deltas:     []Delta{{}},
		lastMod:    -1,
		Indexes:    []*index.Overlay{index.OverlayStub()},
	}
	tbl := InfoHamt{}.Mutable()
	tbl.Put(one)
	tbl.Put(two)
	tbl = tbl.Freeze()

	st := stor.HeapStor(8192)
	btree.CreateBtree(st) // OverlayStub has a btree root offset of 0
	off := tbl.Write(st, 0, hamt.All)

	ic := hamt.ReadChain(st, off, ReadInfo)
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
	btree.CreateBtree(st) // OverlayStub has a btree root offset of 0
	off := tbl.Freeze().Write(st, 0, hamt.All)

	tbl = hamt.ReadChain(st, off, ReadInfo).Hamt
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
		tbl.Put(&Info{Table: data[i], Nrows: i, BtreeNrows: i,
			Indexes: []*index.Overlay{index.OverlayStub()}})
	}
	return data
}

func TestSmallestKeyIndex(t *testing.T) {
	assert := assert.T(t)
	ov := func(treeLevels int) *index.Overlay {
		return index.OverlayFor(btree.BtreeWithLevels(treeLevels))
	}
	key := func(c ...string) schema.Index {
		return schema.Index{Mode: 'k', Columns: c}
	}
	idx := func(c ...string) schema.Index {
		return schema.Index{Mode: 'i', Columns: c}
	}

	// single key index
	info := &Info{Indexes: []*index.Overlay{ov(0)}}
	assert.This(info.SmallestKeyIndex([]schema.Index{key("a")})).Is(0)

	// same treeLevels, different column counts => fewer columns wins
	info = &Info{Indexes: []*index.Overlay{ov(2), ov(2), ov(2)}}
	assert.This(info.SmallestKeyIndex(
		[]schema.Index{key("a", "b", "c"), key("x"), key("y", "z")})).Is(1)

	// different treeLevels => fewer levels wins
	info = &Info{Indexes: []*index.Overlay{ov(3), ov(1), ov(2)}}
	assert.This(info.SmallestKeyIndex(
		[]schema.Index{key("a"), key("x", "y", "z"), key("y")})).Is(1)

	// treeLevels takes priority over columns
	info = &Info{Indexes: []*index.Overlay{ov(1), ov(2)}}
	assert.This(info.SmallestKeyIndex(
		[]schema.Index{key("a", "b", "c"), key("x")})).Is(0)

	// tie on both => first wins
	info = &Info{Indexes: []*index.Overlay{ov(2), ov(2)}}
	assert.This(info.SmallestKeyIndex(
		[]schema.Index{key("a"), key("b")})).Is(0)

	// non-key indexes are skipped
	info = &Info{Indexes: []*index.Overlay{ov(3), ov(1), ov(2)}}
	assert.This(info.SmallestKeyIndex(
		[]schema.Index{idx("a"), key("x"), idx("y")})).Is(1)

	// first key is not first index
	info = &Info{Indexes: []*index.Overlay{ov(3), ov(1), ov(2)}}
	assert.This(info.SmallestKeyIndex(
		[]schema.Index{idx("a"), key("x"), key("y")})).Is(1)
}
