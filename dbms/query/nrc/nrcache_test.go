// Governed by the MIT license found in the LICENSE file.

package nrc

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"golang.org/x/exp/rand"
)

func TestAdd(t *testing.T) {
	for b := 0; b < 1000; b++ {
		var batch Batch
		for range 100 {
			batch.Add(randKey(), 0)
		}
		batch.Save()
	}
}

func BenchmarkHamt_Get(b *testing.B) {
	var ht hamt.Hamt[Hash, *item]
	m := ht.Mutable()
	for i := uint32(0); i < 8000; i++ {
		m.Put(&item{key: randKey(), n: i})
	}
	ht = m.Freeze()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ht.Get(randKey())
	}
}

func BenchmarkMap_Get(b *testing.B) {
	ht := map[Hash]uint32{}
	for i := uint32(0); i < 8000; i++ {
		ht[randKey()] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, _ = ht[randKey()]
	}
}

var n uint32

func randKey() Hash {
	var k Hash
	for i := range k {
		k[i] = byte(rand.Intn(256))
	}
	return k
}
