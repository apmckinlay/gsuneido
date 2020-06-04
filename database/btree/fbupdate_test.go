// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"bufio"
	"math/rand"
	"os"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMemOffsets(t *testing.T) {
	var mo memOffsets
	Assert(t).That(mo.get(0), Equals(nil))
	Assert(t).That(mo.get(123), Equals(nil))
	n := fNode{123}
	o1 := mo.add(n)
	p := mo.get(o1)
	Assert(t).That([]byte(*p)[0], Equals(123))
	n = fNode{99}
	o2 := mo.add(n)
	p = mo.get(o2)
	Assert(t).That([]byte(*p)[0], Equals(99))
	p = mo.get(o1)
	Assert(t).That([]byte(*p)[0], Equals(123))
}

func TestFbupdate(t *testing.T) {
	var nTimes = 10
	if testing.Short() {
		nTimes = 1
	}
	for j := 0; j < nTimes; j++ {
		var data []string
		mo := memOffsets{}
		up := fbupdate{
			bt:          fbtree{root: mo.add(fNode{})},
			moffs:       mo,
			getLeafKey:  func(i uint64) string { return data[i] },
			maxNodeSize: 44,
		}
		const n = 1000
		randKey := str.UniqueRandomOf(4, 5, "abcd")
		for i := 0; i < n; i++ {
			key := randKey()
			off := len(data)
			data = append(data, key)
			up.Insert(key, uint64(off))
			count, _, _ := up.check()
			Assert(t).That(count, Equals(i+1))
		}
	}
}

func TestSampleData(t *testing.T) {
	const nShuffle = 8
	test := func(file string) {
		data := fileData(file)
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			mo := memOffsets{}
			up := fbupdate{
				bt:          fbtree{root: mo.add(fNode{})},
				moffs:       mo,
				getLeafKey:  func(i uint64) string { return data[i] },
				maxNodeSize: 128,
			}
			for i, d := range data {
				up.Insert(d, uint64(i))
			}
			count, _, _ := up.check()
			Assert(t).That(count, Equals(len(data)))
			// print("count", count, "nnodes", nnodes, "size", size,
			// 	"full", float32(size) / float32(nnodes) / float32(up.maxNodeSize))
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
