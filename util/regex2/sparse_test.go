// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex2

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSparse(t *testing.T) {
	ss := SparseSet{}
	nums := []int16{0, 5, 7, 13}
	for _, n := range nums {
		assert.False(ss.Has(n))
	}
	for _, n := range nums {
		ss.Add(n)
	}
	for _, n := range nums {
		assert.True(ss.Has(n))
	}
	ss.Clear()
	for _, n := range nums {
		assert.False(ss.Has(n))
	}
}

var ss = SparseSet{}

func BenchmarkSparseSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ss = SparseSet{}
		for j := 0; j < 10; j++ {
			ss.Add(int16(rand.Intn(100)))
		}
		ss.Clear()
	}
}
