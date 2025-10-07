// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"math/rand"
	"sort"
	"testing"

	btree "github.com/apmckinlay/gsuneido/db19/index/btree3"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/testdata"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestOverlay(*testing.T) {
	var data []string
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	bt.SetSplit(64)
	mut := &ixbuf.T{}
	u := &ixbuf.T{}
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{u}, mut: mut}
	checkIter(data, ov)

	const n = 100
	randKey := str.UniqueRandom(3, 7)

	data = insert(data, n, randKey, mut)
	checkIter(data, ov)

	data = insert(data, n, randKey, u)
	checkIter(data, ov)

	for range n / 2 {
		j := rand.Intn(len(data))
		if data[j] != "" {
			ov.Delete(data[j], key2off(data[j]))
			data[j] = ""
		}
	}
	checkIter(data, ov)
}

func insert(data []string, n int, randKey func() string, dest *ixbuf.T) []string {
	for range n {
		key := randKey()
		off := key2off(key)
		data = append(data, key)
		dest.Insert(key, off)
	}
	return data
}

func key2off(key string) uint64 {
	off := uint64(0)
	for _, c := range key {
		off = off<<8 + uint64(c)
	}
	assert.That(off>>62 == 0)
	return off
}

func checkIter(data []string, ov *Overlay) {
	sort.Strings(data)
	it := NewOverIter("", 0)
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	for _, k := range data {
		if k == "" {
			continue
		}
		it.Next(tran)
		assert.False(it.Eof())
		k2, o2 := it.Cur()
		assert.This(k2).Is(k)
		assert.This(o2).Is(key2off(k))
	}
	it.Next(tran)
	assert.True(it.Eof())
}

func TestOverlayBug(*testing.T) {
	d := testdata.New()
	const n = 100

	st := stor.HeapStor(8192)
	st.Alloc(1)
	bt := btree.CreateBtree(st, nil)
	bt.SetSplit(64)
	ov := &Overlay{bt: bt}
	checkOver(d, ov)

	u := &ixbuf.T{}
	insertTestData(d, n, u)
	ov.layers = []*ixbuf.T{u}
	checkOver(d, ov)
	
	ov.bt = ov.bt.MergeAndSave(u.Iter())
	ov.layers[0] = &ixbuf.T{}
	checkOver(d, ov)
}

func insertTestData(dat *testdata.T, n int, dest *ixbuf.T) {
	for range n {
		dest.Insert(dat.Gen())
	}
}

func checkOver(d *testdata.T, ov *Overlay) {
	t := &testTran{getIndex: func() *Overlay { return ov }}
	sort.Strings(d.Keys)
	i := 0
	it := NewOverIter("", 0)
	for it.Next(t); !it.Eof(); it.Next(t) {
		k, o := it.Cur()
		assert.This(k).Is(d.Keys[i])
		assert.This(d.O2k[o]).Is(d.Keys[i])
		i++
	}
	assert.This(i).Is(len(d.Keys))
}

func TestOverlayMerge(t *testing.T) {
	randKey := str.UniqueRandomOf(3, 10, "abcdef")
	var data []string
	randIxbuf := func() *ixbuf.T {
		const n = 300
		mut := &ixbuf.T{}
		for range n {
			key := randKey()
			off := uint64(len(data))
			data = append(data, key)
			mut.Insert(key, off+1)
		}
		return mut
	}
	mut := randIxbuf()
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	bt.SetSplit(64)
	bi := &ixbuf.T{}
	ov := Overlay{bt: bt, layers: []*ixbuf.T{bi, mut}}
	bi = ov.Merge(1)
	checkData(t, bi, data)

	mut = randIxbuf()
	ov = Overlay{bt: bt, layers: []*ixbuf.T{bi, mut}}
	bi = ov.Merge(1)
	checkData(t, bi, data)
}

func checkData(t *testing.T, bi *ixbuf.T, data []string) {
	t.Helper()
	assert.T(t).This(bi.Len()).Is(len(data))
	sort.Strings(data)
	i := 0
	n := 0
	it := bi.Iter()
	for key, _, ok := it(); ok; key, _, ok = it() {
		assert.T(t).This(key).Is(data[i])
		i++
		n++
	}
	assert.T(t).This(n).Is(len(data))
}

func TestOverlayLookup(*testing.T) {
	dat := testdata.New()
	store := stor.HeapStor(8192)
	randBtree := func(nkeys int) iface.Btree {
		for range nkeys {
			dat.Gen()
		}
		sort.Strings(dat.Keys)
		b := btree.Builder(store)
		b.SetSplit(64)
		for _, k := range dat.Keys {
			assert.That(b.Add(k, dat.K2o[k]))
		}
		return b.Finish()
	}
	randIxbuf := func(nkeys int) *ixbuf.T {
		ib := &ixbuf.T{}
		for range nkeys {
			ib.Insert(dat.Gen())
		}
		return ib
	}
	bt := randBtree(10000)
	layers := []*ixbuf.T{randIxbuf(1000), randIxbuf(100)}
	mut := randIxbuf(100)
	ov := &Overlay{bt: bt, layers: layers, mut: mut}
	for _, k := range dat.Keys {
		assert.This(ov.Lookup(k)).Is(dat.K2o[k])
		assert.This(ov.Lookup(k + "0")).Is(0) // nonexistent
	}
}
