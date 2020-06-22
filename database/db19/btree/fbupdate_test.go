// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMemOffsets(t *testing.T) {
	mo := memOffsets{nextOff: stor.MaxSmallOffset, nodes: make(map[uint64]fNode)}
	Assert(t).True(mo.get(0) == nil)
	Assert(t).True(mo.get(123) == nil)
	n := fNode{123}
	o1 := mo.add(n)
	n1 := mo.get(o1)
	Assert(t).That(n1[0], Equals(123))
	n = fNode{99}
	o2 := mo.add(n)
	n2 := mo.get(o2)
	Assert(t).That(n2[0], Equals(99))
	n1 = mo.get(o1)
	Assert(t).That(n1[0], Equals(123))
}

func TestFbupdate(t *testing.T) {
	var nTimes = 10
	if testing.Short() {
		nTimes = 1
	}
	for j := 0; j < nTimes; j++ {
		const n = 1000
		var data [n]string
		getLeafKey := func(i uint64) string { return data[i] }
		fb := CreateFbtree(nil, getLeafKey, 44)
		up := newFbupdate(fb)
		randKey := str.UniqueRandomOf(3, 6, "abcde")
		for i := 0; i < n; i++ {
			key := randKey()
			data[i] = key
			up.Insert(key, uint64(i))
			// count, _, _ := up.check()
			// Assert(t).That(count, Equals(i+1))
		}
		up.checkData(t, data[:])
	}
}

func TestUnevenSplit(t *testing.T) {
	const n = 1000
	var data [n]string
	test := func() {
		getLeafKey := func(i uint64) string { return data[i] }
		fb := CreateFbtree(nil, getLeafKey, 128)
		up := newFbupdate(fb)
		for i := 0; i < n; i++ {
			up.Insert(data[i], uint64(i))
		}
		count, size, nnodes := up.check()
		Assert(t).That(count, Equals(n))
		full := float32(size) / float32(nnodes) / float32(up.fb.maxNodeSize)
		// print("count", count, "nnodes", nnodes, "size", size, "full", full)
		if full < .65 {
			t.Error("expected > .65 got", full)
		}
		up.checkData(t, data[:])
	}
	randKey := str.UniqueRandomOf(3, 6, "abcde")
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	test()
	sort.Strings(data[:])
	test()
	str.ListReverse(data[:])
	test()
}

func (fb *fbtree) checkData(t *testing.T, data []string) {
	count, _, _ := fb.check()
	Assert(t).That(count, Equals(len(data)))
	for i, k := range data {
		o := fb.Search(k)
		if o != uint64(i) {
			t.Error("checkData", k, "expect", i, "actual", o)
		}
	}
}

func (up *fbupdate) checkData(t *testing.T, data []string) {
	count, _, _ := up.check()
	n := 0
	for i, k := range data {
		if data[i] == "" {
			continue
		}
		o := up.Search(k)
		if o != uint64(i) {
			t.Error("checkData", k, "expect", i, "actual", o)
		}
		n++
	}
	Assert(t).That(count, Equals(n))
}

func TestSampleData(t *testing.T) {
	var nShuffle = 16
	if testing.Short() {
		nShuffle = 4
	}
	test := func(file string) {
		data := fileData(file)
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			getLeafKey := func(i uint64) string { return data[i] }
			fb := CreateFbtree(nil, getLeafKey, 256)
			up := newFbupdate(fb)
			for i, d := range data {
				up.Insert(d, uint64(i))
			}
			up.checkData(t, data)
		}
	}
	test("../../../bizpartnername.txt")
	test("../../../bizpartnerabbrev.txt")
}

func fileData(filename string) []string {
	file, _ := os.Open(filename)
	defer file.Close()
	data := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return data
}

func TestFbdelete(t *testing.T) {
	var n = 2000
	if testing.Short() {
		n = 100
	}
	data := make([]string, n)
	getLeafKey := func(i uint64) string { return data[i] }
	fb := CreateFbtree(nil, getLeafKey, 44)
	up := newFbupdate(fb)
	randKey := str.UniqueRandomOf(3, 6, "abcde")
	for i := 0; i < n; i++ {
		key := randKey()
		data[i] = key
		up.Insert(key, uint64(i))
	}
	up.checkData(t, data)
	// up.print()

	for i := 0; i < len(data); i++ {
		off := rand.Intn(len(data))
		for data[off] == "" {
			off = (off + 1) % len(data)
		}
		// print("================================= delete", data[off])
		up.Delete(data[off], uint64(off))
		// up.print()
		data[off] = ""
		if i%11 == 0 {
			up.checkData(t, data)
		}
	}
}

func TestSave(t *testing.T) {
	var nSaves = 50
	if testing.Short() {
		nSaves = 10
	}
	const freezesPerSave = 3
	const insertsPerFreeze = 17
	data := make([]string, 0, nSaves*freezesPerSave*insertsPerFreeze)
	getLeafKey := func(i uint64) string { return data[i] }
	st, err := stor.MmapStor("tmp.db", stor.CREATE)
	if err != nil {
		log.Fatalln(err)
	}
	fb := CreateFbtree(st, getLeafKey, 64)
	randKey := str.UniqueRandomOf(3, 7, "abcdef")
	for i := 0; i < nSaves; i++ {
		for j := 0; j < freezesPerSave; j++ {
			fb = fb.Update(func(up *fbupdate) {
				for k := 0; k < insertsPerFreeze; k++ {
					key := randKey()
					up.Insert(key, uint64(len(data)))
					data = append(data, key)
				}
				up.checkData(t, data)
			})
			fb.checkData(t, data)
		}
		fb = fb.save()
		fb.checkData(t, data)
	}
	st.Close()
	os.Remove("tmp.db")
}
