// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

type heapStor struct {
	chunksize int
}

// HeapStor returns an empty in-memory stor for testing.
func HeapStor(chunksize int) *stor {
	return &stor{chunksize: uint64(chunksize), impl: &heapStor{chunksize},
		chunks: [][]byte{make([]byte, chunksize, chunksize)}}
}

func (hs heapStor) Get(int) []byte {
	return make([]byte, hs.chunksize, hs.chunksize)
}

func (hs heapStor) Close(int64) {
}
