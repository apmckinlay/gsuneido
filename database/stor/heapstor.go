package stor

type heapStor struct {
	chunksize int
}

// HeapStor returns an in-memory stor for testing.
func HeapStor(chunksize int) *stor {
	return &stor{chunksize: int64(chunksize), impl: &heapStor{chunksize}}
}

func (hs heapStor) Get(chunk int) []byte {
	return make([]byte, hs.chunksize, hs.chunksize)
}
func (hs heapStor) Close() {
}
