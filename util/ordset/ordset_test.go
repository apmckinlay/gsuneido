// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ordset

import (
	"math/rand"
	"sort"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestRandom(t *testing.T) {
	var nGenerate = 8
	var nShuffle = 8
	if testing.Short() {
		nGenerate = 2
		nShuffle = 2
	}
	const n = nodeSize * 80
	data := make([]string, n)
	for g := 0; g < nGenerate; g++ {
		randKey := str.UniqueRandom(3, 10)
		for i := 0; i < n; i++ {
			data[i] = randKey()
		}
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			var x Set
			for _, k := range data {
				x.Insert(k)
			}
			x.checkData(t, data)
		}
	}
}

func TestUnevenSplit(t *testing.T) {
	const n = nodeSize * 87 // won't fit without uneven splits
	data := make([]string, n)
	randKey := str.UniqueRandom(3, 10)
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	sort.Strings(data)
	m := Set{}
	for _, k := range data {
		m.Insert(k)
	}
	m.checkData(t, data)
	m = Set{}
	for i := len(data) - 1; i >= 0; i-- {
		m.Insert(data[i])
	}
	m.checkData(t, data)
}

func TestAnyInRange(t *testing.T) {
	const n = nodeSize * 80
	data := make([]string, n)
	randKey := str.UniqueRandom(3, 10)
	var x Set
	for i := 0; i < n; i++ {
		data[i] = randKey()
		x.Insert(data[i])
	}
	yes := func(d int) {
		Assert(t).True(x.AnyInRange(data[d], data[d]))
		Assert(t).True(x.AnyInRange(smaller(data[d]), data[d]))
		Assert(t).True(x.AnyInRange(data[d], bigger(data[d])))
		Assert(t).True(x.AnyInRange(smaller(data[d]), bigger(data[d])))
	}
	sort.Strings(data)
	yes(0)
	yes(n - 1)
	const ntimes = 1000
	for i := 0; i < ntimes; i++ {
		d := rand.Intn(n - 1)
		Assert(t).False(x.AnyInRange(bigger(data[d]), smaller(data[d+1])))
		yes(rand.Intn(n))
		d = rand.Intn(n - 10)
		e := d + rand.Intn(10)
		Assert(t).True(x.AnyInRange(smaller(data[d]), bigger(data[e])))
	}
}

//-------------------------------------------------------------------

func (set *Set) checkData(t *testing.T, data []string) {
	for _, key := range data {
		Assert(t).True(set.Contains(key))
		Assert(t).False(set.Contains(bigger(key)))
		Assert(t).False(set.Contains(smaller(key)))
	}
}

func bigger(s string) string {
	return s + " "
}

func smaller(s string) string {
	last := len(s) - 1
	return s[:last] + string(s[last]-1) + " "
}
