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

	"github.com/apmckinlay/gsuneido/database/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMemOffsets(t *testing.T) {
	mo := newMemOffsets()
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
		mo := newMemOffsets()
		up := fbupdate{
			fb:          fbtree{root: mo.add(fNode{})},
			moffs:       mo,
			getLeafKey:  func(i uint64) string { return data[i] },
			maxNodeSize: 44,
		}
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
		mo := newMemOffsets()
		up := fbupdate{
			fb:          fbtree{root: mo.add(fNode{})},
			moffs:       mo,
			getLeafKey:  func(i uint64) string { return data[i] },
			maxNodeSize: 128,
		}
		for i := 0; i < n; i++ {
			up.Insert(data[i], uint64(i))
		}
		count, size, nnodes := up.check()
		Assert(t).That(count, Equals(n))
		full := float32(size) / float32(nnodes) / float32(up.maxNodeSize)
		// print("count", count, "nnodes", nnodes, "size", size, "full", full)
		if full < .65 {
			t.Error("expected > .65 got", full)
		}
		up.checkData(t, data[:])
	}
	randKey := str.UniqueRandomOf(3, 6, "abcd")
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	test()
	sort.Strings(data[:])
	test()
	str.ListReverse(data[:])
	test()
}

func (up *fbupdate) checkData(t *testing.T, data []string) {
	count, _, _ := up.check()
	Assert(t).That(count, Equals(len(data)))
	for i, k := range data {
		o := up.Search(k)
		if o != uint64(i) {
			t.Error("checkData", k, "expected", i, "got", o)
		}
	}
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
			mo := newMemOffsets()
			up := fbupdate{
				fb:          fbtree{root: mo.add(fNode{})},
				moffs:       mo,
				getLeafKey:  func(i uint64) string { return data[i] },
				maxNodeSize: 256,
			}
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

func TestSave(t *testing.T) {
	const n = 1000
	data := make([]string, 0, n)
	mo := newMemOffsets()
	st, err := stor.MmapStor("tmp.db", stor.CREATE)
	if err != nil {
		log.Fatalln(err)
	}
	up := fbupdate{
		fb:          fbtree{root: mo.add(fNode{}), store: st},
		moffs:       mo,
		getLeafKey:  func(i uint64) string { return data[i] },
		maxNodeSize: 44,
	}
	up.save()
	randKey := str.UniqueRandomOf(3, 6, "abcd")
	for i := 0; i < n; i++ {
		key := randKey()
		data = append(data, key)
		up.Insert(key, uint64(i))
		if (i+1)%17 == 0 {
			up.checkData(t, data)
			up.save()
			up.checkData(t, data)
		}
	}
	st.Close()
	os.Remove("tmp.db")
}
