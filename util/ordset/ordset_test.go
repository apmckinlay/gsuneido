// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ordset

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestRandom(t *testing.T) {
	var nGenerate = 250
	var nShuffle = 4
	if testing.Short() {
		nGenerate = 2
		nShuffle = 2
	}
	const n = nodeSize * 75
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
				assert.That(x.Insert(k))
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
	assert := assert.T(t)
	const n = nodeSize * 80
	data := make([]string, n)
	randKey := str.UniqueRandom(3, 10)
	var x Set
	for i := 0; i < n; i++ {
		data[i] = randKey()
		x.Insert(data[i])
	}
	yes := func(d int) {
		t.Helper()
		assert.True(x.AnyInRange(data[d], data[d]))
		assert.True(x.AnyInRange(smaller(data[d]), data[d]))
		assert.True(x.AnyInRange(data[d], bigger(data[d])))
		assert.True(x.AnyInRange(smaller(data[d]), bigger(data[d])))
	}
	sort.Strings(data)
	yes(0)
	yes(n - 1)
	const ntimes = 1000
	for i := 0; i < ntimes; i++ {
		d := rand.Intn(n - 1)
		assert.False(x.AnyInRange(bigger(data[d]), smaller(data[d+1])))
		yes(rand.Intn(n))
		d = rand.Intn(n - 10)
		e := d + rand.Intn(10)
		assert.True(x.AnyInRange(smaller(data[d]), bigger(data[e])))
	}
}

func TestOverflow(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var x Set
	randKey := str.UniqueRandom(3, 10)
	for i := 0; i < 20000; i++ {
		if !x.Insert(randKey()) {
			fmt.Println(i)
			break
		}
	}
}

func TestDups(t *testing.T) {
	for _, n := range []int{50, 5000} {
		var x Set
		seed := time.Now().UnixNano()
		rand.Seed(seed)
		randKey := str.UniqueRandom(3, 10)
		for i := 0; i < n; i++ {
			assert.That(x.Insert(randKey()))
		}
		assert.This(x.count()).Is(n)
		rand.Seed(seed)
		randKey = str.UniqueRandom(3, 10)
		for i := 0; i < n; i++ {
			assert.That(x.Insert(randKey()))
		}
		assert.This(x.count()).Is(n)
	}
}

func (set *Set) count() int {
	if set.tree == nil {
		return set.leaf.size
	}
	n := 0
	for i := 0; i < set.tree.size; i++ {
		n += set.tree.slots[i].leaf.size
	}
	return n
}

//-------------------------------------------------------------------

func (set *Set) checkData(t *testing.T, data []string) {
	for _, key := range data {
		assert.T(t).True(set.Contains(key))
		assert.T(t).False(set.Contains(smaller(key)))
		bgr := bigger(key)
		assert.T(t).False(set.Contains(bgr))
		assert.T(t).False(set.AnyInRange(bgr, bigger(bgr)))
	}
}

func bigger(s string) string {
	return s + " "
}

func smaller(s string) string {
	last := len(s) - 1
	return s[:last] + string(s[last]-1) + "~"
}
