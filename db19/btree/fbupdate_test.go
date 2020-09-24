// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/ixspec"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestUpdate(t *testing.T) {
	var nTimes = 10
	if testing.Short() {
		nTimes = 1
	}
	for j := 0; j < nTimes; j++ {
		const n = 1000
		var data [n]string
		GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string { return data[i] }
		defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
		MaxNodeSize = 44
		fb := CreateFbtree(nil, nil)
		mfb := fb.makeMutable()
		randKey := str.UniqueRandomOf(3, 6, "abcde")
		for i := 0; i < n; i++ {
			key := randKey()
			data[i] = key
			mfb.Insert(key, uint64(i))
		}
		mfb.checkData(t, data[:])
	}
}

func TestUnevenSplit(t *testing.T) {
	const n = 1000
	var data [n]string
	test := func() {
		GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string { return data[i] }
		defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
		MaxNodeSize = 128
		fb := CreateFbtree(nil, nil)
		mfb := fb.makeMutable()
		for i := 0; i < n; i++ {
			mfb.Insert(data[i], uint64(i))
		}
		count, size, nnodes := mfb.check(nil)
		assert.T(t).This(count).Is(n)
		full := float32(size) / float32(nnodes) / float32(MaxNodeSize)
		// print("count", count, "nnodes", nnodes, "size", size, "full", full)
		if full < .65 {
			t.Error("expected > .65 got", full)
		}
		mfb.checkData(t, data[:])
	}
	randKey := str.UniqueRandomOf(3, 6, "abcde")
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	test()
	sort.Strings(data[:])
	test()
	str.List(data[:]).Reverse()
	test()
}

func (fb *fbtree) checkData(t *testing.T, data []string) {
	t.Helper()
	count, _, _ := fb.check(nil)
	n := 0
	for i, k := range data {
		if data[i] == "" {
			continue
		}
		o := fb.Search(k)
		if o != uint64(i) {
			t.Log("checkData", k, "expect", i, "actual", o)
			t.FailNow()
		}
		n++
	}
	assert.T(t).This(count).Is(n)
}

func TestSampleData(t *testing.T) {
	var nShuffle = 12
	if testing.Short() {
		nShuffle = 4
	}
	test := func(file string) {
		data := fileData(file)
		// fmt.Println(len(data))
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string {
				return data[i]
			}
			defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
			MaxNodeSize = 256
			fb := CreateFbtree(nil, nil)
			mfb := fb.makeMutable()
			for i, d := range data {
				mfb.Insert(d, uint64(i))
			}
			mfb.checkData(t, data)
		}
	}
	test("../../../bizpartnername.txt")
	test("../../../bizpartnerabbrev.txt")
}

func fileData(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("can't open", filename)
	}
	defer file.Close()
	data := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return data
}

func TestFbdelete(t *testing.T) {
	var n = 1000
	if testing.Short() {
		n = 100
	}
	data := make([]string, n)
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 44
	fb := CreateFbtree(nil, nil)
	mfb := fb.makeMutable()
	randKey := str.UniqueRandomOf(3, 6, "abcdef")
	for i := 0; i < n; i++ {
		key := randKey()
		data[i] = key
		mfb.Insert(key, uint64(i))
	}
	mfb.checkData(t, data)
	// mfb.print()

	for i := 0; i < len(data); i++ {
		off := rand.Intn(len(data))
		for data[off] == "" {
			off = (off + 1) % len(data)
		}
		// print("================================= delete", data[off])
		mfb.Delete(data[off], uint64(off))
		// mfb.print()
		data[off] = ""
		if i%11 == 0 {
			mfb.checkData(t, data)
		}
	}
}

func TestFreeze(t *testing.T) {
	assert := assert.T(t).This
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string {
		return strconv.Itoa(int(i))
	}
	store := stor.HeapStor(8192)
	store.Alloc(1) // avoid offset 0
	fb := CreateFbtree(store, nil)
	assert(fb.redirs.Len()).Is(1)
	fb = fb.Update(func(mfb *fbtree) {
		mfb.Insert("1", 1)
	})
	assert(fb.redirs.Len()).Is(1)
	assert(fb.list()).Is("1")
	fb = fb.Update(func(mfb *fbtree) {
		mfb.Insert("2", 2)
	})
	assert(fb.redirs.Len()).Is(1)
	assert(fb.list()).Is("1 2")

	fb = fb.Save()
	fb = OpenFbtree(store, fb.root, fb.treeLevels, fb.redirsOff)
	assert(fb.redirs.Len()).Is(0)
	assert(fb.list()).Is("1 2")
}

func (fb *fbtree) list() string {
	s := ""
	iter := fb.Iter()
	for _, o, ok := iter(); ok; _, o, ok = iter() {
		s += strconv.Itoa(int(o)) + " "
	}
	return strings.TrimSpace(s)
}

func TestSave(t *testing.T) {
	var nSaves = 40
	if testing.Short() {
		nSaves = 10
	}
	const updatesPerSave = 3
	const insertsPerUpdate = 17
	data := make([]string, 0, nSaves*updatesPerSave*insertsPerUpdate)
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	st := stor.HeapStor(8192)
	st.Alloc(1) // avoid offset 0
	fb := CreateFbtree(st, nil)
	randKey := str.UniqueRandomOf(3, 7, "abcdef")
	for i := 0; i < nSaves; i++ {
		for j := 0; j < updatesPerSave; j++ {
			fb = fb.Update(func(mfb *fbtree) {
				for k := 0; k < insertsPerUpdate; k++ {
					key := randKey()
					mfb.Insert(key, uint64(len(data)))
					data = append(data, key)
				}
			})
			fb.checkData(t, data)
			fb.ckpaths()
		}
		fb = fb.Save()
		fb.checkData(t, data)
	}
}

func TestSplitDup(*testing.T) {
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string {
		return strconv.Itoa(int(i))
	}
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	data := []int{}
	for i := 3; i < 8; i++ {
		data = append(data, i)
	}
	for i := 53; i < 58; i++ {
		data = append(data, i)
	}
	for i := 553; i < 558; i++ {
		data = append(data, i)
	}
	for i := 5553; i < 5558; i++ {
		data = append(data, i)
	}
	for i := 55553; i < 55558; i++ {
		data = append(data, i)
	}
	n := 10000
	if testing.Short() {
		n = 1000
	}
	for i := 0; i < n; i++ {
		rand.Shuffle(len(data),
			func(i, j int) { data[i], data[j] = data[j], data[i] })
		fb := CreateFbtree(nil, nil)
		fb = fb.Update(func(mfb *fbtree) {
			for _, n := range data {
				key := strconv.Itoa(n)
				mfb.Insert(key, uint64(n))
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	assert := assert.T(t)
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string {
		return strconv.Itoa(int(i))
	}
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	const from, to = 10000, 10800
	inserted := map[int]bool{}
	var fb *fbtree

	build := func() {
		trace("==============================")
		MaxNodeSize = 96
		inserted = map[int]bool{}
		store := stor.HeapStor(8192)
		bldr := NewFbtreeBuilder(store)
		for i := from; i < to; i += 2 {
			key := strconv.Itoa(i)
			bldr.Add(key, uint64(i))
		}
		fb = bldr.Finish().base()
		assert.That(fb.treeLevels == 2)
		fb.redirs.tbl.ForEach(func(*redir) { panic("redir!") })
		fb.redirs.paths.ForEach(func(uint64) { panic("path!") })
	}
	check := func() {
		t.Helper()
		fb.check(nil)
		iter := fb.Iter()
		for i := from; i < to; i++ {
			if i%2 == 1 && !inserted[i] {
				continue
			}
			key := strconv.Itoa(i)
			k, o, ok := iter()
			assert.True(ok)
			assert.True(strings.HasPrefix(key, k))
			assert.This(o).Is(i)
			if o != uint64(i) {
				t.FailNow()
			}
		}
		_, _, ok := iter()
		assert.False(ok)
	}
	insert := func(i int) {
		fb = fb.Update(func(mfb *fbtree) {
			mfb.Insert(strconv.Itoa(i), uint64(i))
			inserted[i] = true
		})
		check()
	}
	maybeSave := func(save bool) {
		check()
		if save {
			fb = fb.Save()
			check()
			trace("---------------------------")
		}
	}
	flatten := func() {
		fb = fb.Update(func(mfb *fbtree) { mfb.flatten() })
		check()
	}

	for _, save := range []bool{false, true} {
		for _, mns := range []int{999, 90} {
			build()
			MaxNodeSize = mns // prevent or force splitting
			insert(10051)
			maybeSave(save)
			flatten()
		}
	}
	for _, save := range []bool{false, true} {
		build()
		MaxNodeSize = 999 // no split
		insert(10051)
		MaxNodeSize = 90 // split all the way
		insert(10551)
		maybeSave(save)
		flatten()
	}
}

//-------------------------------------------------------------------

// ckpaths checks that all the redirects can be reached by following the paths
// and that all paths are in the current tree
func (fb *fbtree) ckpaths() {
	var rset = make(map[uint64]bool)
	fb.redirs.tbl.ForEach(func(r *redir) {
		rset[r.offset] = true
	})
	var pset = make(map[uint64]bool)
	fb.redirs.paths.ForEach(func(o uint64) {
		pset[o] = true
	})

	delete(rset, fb.root)
	fb.ckpaths1(0, fb.root, rset, pset)
	if len(rset) != 0 {
		fb.print()
		fmt.Println("root", fb.root, "rset", rset)
		fb.printPaths("paths")
	}
	assert.This(len(rset)).Is(0)
	if len(pset) != 0 {
		fb.print()
		fmt.Println("pset", pset)
		fb.printPaths("paths")
	}
	assert.This(len(pset)).Is(0)
}

func (fb *fbtree) ckpaths1(depth int, nodeOff uint64,
	rset map[uint64]bool, pset map[uint64]bool) {
	delete(pset, nodeOff)
	if depth >= fb.treeLevels {
		if _, pathNode := fb.redirs.paths.Get(nodeOff); pathNode {
			fb.print()
			panic("leaf found in paths")
		}
		return // skip leaf nodes
	}
	// tree node
	if depth > 0 && !fb.pathNode(nodeOff) {
		return // skip if not root and not on path
	}
	node := fb.getNode(nodeOff)
	for it := node.iter(); it.next(); {
		off := it.offset
		delete(rset, off)
		fb.ckpaths1(depth+1, off, rset, pset) // RECURSE
	}
}
