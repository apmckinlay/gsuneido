// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"math/bits"

	"github.com/apmckinlay/gsuneido/util/assert"
)

type heapStor struct {
	chunksize int
}

// HeapStor returns an empty in-memory stor for testing.
func HeapStor(chunksize int) *Stor {
	assert.That(bits.OnesCount(uint(chunksize)) == 1)
	hs := NewStor(&heapStor{chunksize}, uint64(chunksize), 0)
	hs.chunks.Store([][]byte{make([]byte, chunksize)})
	return hs
}

func (hs heapStor) Get(int) []byte {
	return make([]byte, hs.chunksize)
}

func (hs heapStor) Close(int64) {
}
