// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/fbtree"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/index/testdata"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestEmptyOverlay(*testing.T) {
	var data []string
	defer func(mns int) { fbtree.MaxNodeSize = mns }(fbtree.MaxNodeSize)
	fbtree.MaxNodeSize = 64
	fb := fbtree.CreateFbtree(stor.HeapStor(8192), nil)
	mut := &ixbuf.T{}
	u := &ixbuf.T{}
	ov := &Overlay{fb: fb, layers: []*ixbuf.T{u}, mut: mut}
	checkIter(data, ov)

	const n = 100
	randKey := str.UniqueRandom(3, 7)

	data = insert(data, n, randKey, mut)
	checkIter(data, ov)

	data = insert(data, n, randKey, u)
	checkIter(data, ov)

	for i := 0; i < n/2; i++ {
		j := rand.Intn(len(data))
		if data[j] != "" {
			ov.Delete(data[j], key2off(data[j]))
			data[j] = ""
		}
	}
	checkIter(data, ov)
}

func insert(data []string, n int, randKey func() string, dest *ixbuf.T) []string {
	for i := 0; i < n; i++ {
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
	it := ov.Iter(true)
	for _, k := range data {
		if k == "" {
			continue
		}
		k2, o2, ok := it()
		assert.True(ok)
		assert.This(k2).Is(k)
		assert.This(o2).Is(key2off(k))
	}
	_, _, ok := it()
	assert.False(ok)
}

func TestOverlayBug(*testing.T) {
	d := testdata.New()
	fbtree.GetLeafKey = d.GetLeafKey
	defer func(mns int) { fbtree.MaxNodeSize = mns }(fbtree.MaxNodeSize)
	fbtree.MaxNodeSize = 64
	const n = 100

	fb := fbtree.CreateFbtree(stor.HeapStor(8192), nil)
	ov := &Overlay{fb: fb}
	d.CheckIter(ov.Iter(false))

	u := &ixbuf.T{}
	insertTestData(d, n, u)
	ov.layers = []*ixbuf.T{u}
	d.CheckIter(ov.Iter(false))

	ov.fb = ov.fb.MergeAndSave(u.Iter(false))
	ov.layers[0] = &ixbuf.T{}
	d.CheckIter(ov.Iter(false))
}

func insertTestData(dat *testdata.T, n int, dest *ixbuf.T) {
	for i := 0; i < n; i++ {
		dest.Insert(dat.Gen())
	}
}

func TestOverlayMerge(t *testing.T) {
	randKey := str.UniqueRandomOf(3, 10, "abcdef")
	var data []string
	randIxbuf := func() *ixbuf.T {
		const n = 300
		mut := &ixbuf.T{}
		for i := 0; i < n; i++ {
			key := randKey()
			off := uint64(len(data))
			data = append(data, key)
			mut.Insert(key, off)
		}
		return mut
	}
	mut := randIxbuf()
	fbtree.GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return data[i]
	}
	defer func(mns int) { fbtree.MaxNodeSize = mns }(fbtree.MaxNodeSize)
	fbtree.MaxNodeSize = 64
	fb := fbtree.CreateFbtree(stor.HeapStor(8192), nil)
	bi := &ixbuf.T{}
	ov := Overlay{fb: fb, layers: []*ixbuf.T{bi, mut}}
	bi = ov.Merge(1)
	checkData(t, bi, data)

	mut = randIxbuf()
	ov = Overlay{fb: fb, layers: []*ixbuf.T{bi, mut}}
	bi = ov.Merge(1)
	checkData(t, bi, data)
}

func checkData(t *testing.T, bi *ixbuf.T, data []string) {
	t.Helper()
	assert.T(t).This(bi.Len()).Is(len(data))
	sort.Strings(data)
	i := 0
	n := 0
	it := bi.Iter(false)
	for key, _, ok := it(); ok; key, _, ok = it() {
		assert.T(t).This(key).Is(data[i])
		i++
		n++
	}
	assert.T(t).This(n).Is(len(data))
}

func TestOverlayLookup(*testing.T) {
	dat := testdata.New()
	fbtree.GetLeafKey = dat.GetLeafKey
	defer func(mns int) { fbtree.MaxNodeSize = mns }(fbtree.MaxNodeSize)
	fbtree.MaxNodeSize = 128
	store := stor.HeapStor(8192)
	randFbtree := func(nkeys int) *fbtree.T {
		for i := 0; i < nkeys; i++ {
			dat.Gen()
		}
		sort.Strings(dat.Keys)
		b := fbtree.Builder(store)
		for _,k := range dat.Keys {
			b.Add(k, dat.K2o[k])
		}
		return b.Finish()
	}
	randIxbuf := func(nkeys int) *ixbuf.T {
		ib := &ixbuf.T{}
		for i := 0; i < nkeys; i++ {
			ib.Insert(dat.Gen())
		}
		return ib
	}
	fb := randFbtree(10000)
	layers := []*ixbuf.T{randIxbuf(1000), randIxbuf(100)}
	mut := randIxbuf(100)
	ov := &Overlay{fb: fb, layers: layers, mut: mut}
	for _,k := range dat.Keys {
		assert.This(ov.Lookup(k)).Is(dat.K2o[k])
		assert.This(ov.Lookup(k + "0")).Is(0) // nonexistent
	}
}
