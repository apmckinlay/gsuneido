// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"math/rand"
	"sync"
	"testing"
)

var nThreads = 11
var nIterations = 1000000
const allocSize = 32
const chunkSize = 1024

func TestStress(*testing.T) {
	if testing.Short() {
		nThreads = 2
		nIterations = 10000
	}
	var wg sync.WaitGroup
	s := HeapStor(chunkSize)
	for i := 0; i < nThreads; i++ {
		go thread(&wg, s)
		wg.Add(1)
	}
	wg.Wait()
}

func thread(wg *sync.WaitGroup, s *Stor) {
	for i := 0; i < nIterations; i++ {
		n := rand.Intn(allocSize) + 1
		s.Alloc(n)
	}
	wg.Done()
}
