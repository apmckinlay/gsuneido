// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBitSet(t *testing.T) {
	bs := BitSet{}
	nums := []int16{0, 7, 13, 99, 111}
	for _, n := range nums {
		assert.False(bs.Has(n))
	}
	for _, n := range nums {
		bs.Add(n)
	}
	for _, n := range nums {
		assert.True(bs.Has(n))
	}
	bs.Clear()
	for _, n := range nums {
		assert.False(bs.Has(n))
	}
}

var bs BitSet

func BenchmarkBitSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs = BitSet{}
		for j := 0; j < 10; j++ {
			bs.Add(int16(rand.Intn(100)))
		}
		bs.Clear()
	}
}
