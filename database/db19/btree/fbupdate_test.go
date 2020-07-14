// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFbupdate(t *testing.T) {
	var nTimes = 10
	if testing.Short() {
		nTimes = 1
	}
	for j := 0; j < nTimes; j++ {
		const n = 1000
		var data [n]string
		GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string { return data[i] }
		defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
		MaxNodeSize = 44
		fb := CreateFbtree(nil)
		up := newFbupdate(fb)
		randKey := str.UniqueRandomOf(3, 6, "abcde")
		for i := 0; i < n; i++ {
			key := randKey()
			data[i] = key
			up.Insert(key, uint64(i))
		}
		up.checkData(t, data[:])
	}
}

func TestUnevenSplit(t *testing.T) {
	const n = 1000
	var data [n]string
	test := func() {
		GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string { return data[i] }
		defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
		MaxNodeSize = 128
		fb := CreateFbtree(nil)
		up := newFbupdate(fb)
		for i := 0; i < n; i++ {
			up.Insert(data[i], uint64(i))
		}
		count, size, nnodes := up.check()
		Assert(t).That(count, Equals(n))
		full := float32(size) / float32(nnodes) / float32(MaxNodeSize)
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
	t.Helper()
	count, _, _ := fb.check()
	Assert(t).That(count, Equals(len(data)))
	for i, k := range data {
		o := fb.Search(k)
		if o != uint64(i) {
			t.Log("checkData", k, "expect", i, "actual", o)
			t.FailNow()
		}
	}
}

func (up *fbupdate) checkData(t *testing.T, data []string) {
	t.Helper()
	count, _, _ := up.check()
	n := 0
	for i, k := range data {
		if data[i] == "" {
			continue
		}
		o := up.Search(k)
		if o != uint64(i) {
			t.Log("checkData", k, "expect", i, "actual", o)
			t.FailNow()
		}
		n++
	}
	Assert(t).That(count, Equals(n))
}

func TestSampleData(t *testing.T) {
	var nShuffle = 12
	if testing.Short() {
		nShuffle = 4
	}
	test := func(file string) {
		data := fileData(file)
		fmt.Println(len(data))
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string {
				return data[i]
			}
			defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
			MaxNodeSize = 256
			fb := CreateFbtree(nil)
			up := newFbupdate(fb)
			for i, d := range data {
				up.Insert(d, uint64(i))
			}
			up.checkData(t, data)
		}
	}
	test("../../../../bizpartnername.txt")
	test("../../../../bizpartnerabbrev.txt")
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
	GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 44
	fb := CreateFbtree(nil)
	up := newFbupdate(fb)
	randKey := str.UniqueRandomOf(3, 6, "abcdef")
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

func TestFreeze(t *testing.T) {
	GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string {
		return strconv.Itoa(int(i))
	}
	store := stor.HeapStor(8192)
	store.Alloc(1) // avoid offset 0
	fb := CreateFbtree(store)
	Assert(t).That(fb.moffs.Len(), Equals(1))
	fb = fb.Update(func(up *fbupdate) {
		up.Insert("1", 1)
	})
	Assert(t).That(fb.moffs.Len(), Equals(1))
	Assert(t).That(fb.list(), Equals("1"))
	fb = fb.Update(func(up *fbupdate) {
		up.Insert("2", 2)
	})
	Assert(t).That(fb.moffs.Len(), Equals(1))
	Assert(t).That(fb.list(), Equals("1 2"))

	fb = fb.Save()
	fb = OpenFbtree(store, fb.root, fb.treeLevels, fb.redirs)
	Assert(t).That(fb.moffs.Len(), Equals(1))
	Assert(t).That(fb.list(), Equals("1 2"))
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
	GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	st, err := stor.MmapStor("tmp.db", stor.CREATE)
	st.Alloc(1) // avoid offset 0
	if err != nil {
		log.Fatalln(err)
	}
	fb := CreateFbtree(st)
	randKey := str.UniqueRandomOf(3, 7, "abcdef")
	for i := 0; i < nSaves; i++ {
		for j := 0; j < updatesPerSave; j++ {
			fb = fb.Update(func(up *fbupdate) {
				for k := 0; k < insertsPerUpdate; k++ {
					key := randKey()
					up.Insert(key, uint64(len(data)))
					data = append(data, key)
				}
			})
			fb.checkData(t, data)
		}
		fb = fb.Save()
		fb.checkData(t, data)
	}
	st.Close()
	os.Remove("tmp.db")
}

func TestSplitDup(*testing.T) {
	GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string {
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
		fb := CreateFbtree(nil)
		fb = fb.Update(func(up *fbupdate) {
			for _, n := range data {
				key := strconv.Itoa(n)
				up.Insert(key, uint64(n))
			}
		})
	}
}
