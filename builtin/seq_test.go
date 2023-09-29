// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/race"
)

// NOTE: these tests depend on the race detector to find problems

func TestConcurrentSeq(t *testing.T) {
	if !race.Enabled {
		t.Skip("RACE NOT ENABLED")
	}
	to := 1_000_000
	if testing.Short() {
		to = 50_000
	}
	it := seqIter{by: 1, to: to}
	it.SetConcurrent()
	fn := func() {
		for it.Next() != nil {
			it.Infinite()
			it.Dup().Next()
		}
	}
	const nthreads = 6
	for i := 0; i < nthreads; i++ {
		go fn()
	}
	fn()
}

func TestConcurrentSequence(t *testing.T) {
	if !race.Enabled {
		t.Skip("RACE NOT ENABLED")
	}
	size := 100_000
	if testing.Short() {
		size = 10_000
	}
	sq := core.NewSuSequence(&seqIter{by: 1, to: 1_000})
	sq.SetConcurrent()
	var wg sync.WaitGroup
	const nthreads = 6
	for i := 0; i < nthreads; i++ {
		wg.Add(1)
		go func() {
			n := rand.Intn(size)
			for i := 0; i < n; i++ {
				sq.Infinite()
				sq.Instantiated()
				sq.Iter().Next()
				if i == n/2 {
					sq.ToContainer() // instantiates
				}
			}
			assert.That(sq.Instantiated())
			wg.Done()
		}()
	}
	wg.Wait()
}
