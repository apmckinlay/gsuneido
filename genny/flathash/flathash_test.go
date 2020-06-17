// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package flathash

import (
	"math/rand"
	"testing"
)

func TestFlatHash(t *testing.T) {
	const n = 1200
	h := NewKVMap(0)
	data := make(map[int]int)
	check := func() {
		for k, v := range data {
			if h.Get(k) != v {
				t.Error("expected", v, "for", k)
			}
		}
	}
	addRand := func() {
		key := rand.Intn(10000)
		val := rand.Int()
		h.Set(key, val)
		data[key] = val
	}
	for i := 0; i < n; i++ {
		addRand()
		if i%11 == 0 {
			check()
		}
	}
	h = h.Dup()
	check()
}

func (h *KVMap) hash(k K) uint32 {
	return uint32(k.(int))
}
