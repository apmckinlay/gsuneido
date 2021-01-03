// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/race"
)

// NOTE: these tests depend on the race detector to find problems

func TestConcurrentAtomic(t *testing.T) {
	if !race.Enabled {
		t.Skip("RACE NOT ENABLED")
	}
	// This is equivalent to what the code does
	var concurrent bool
	var mu sync.Mutex
	go func() {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		mu.Unlock()
		assert.That(concurrent == true)
	}()
	time.Sleep(5 * time.Millisecond)
	mu.Lock()
	concurrent = true
	mu.Unlock()
	time.Sleep(10 * time.Millisecond)
}

func TestConcurrentMutex(t *testing.T) {
	if !race.Enabled {
		t.Skip("RACE NOT ENABLED")
	}
	// This demonstrates that atomic is also sufficient
	var x int
	var m int64
	go func() {
		time.Sleep(10 * time.Millisecond)
		assert.That(atomic.LoadInt64(&m) == 1)
		assert.That(x == 123)
	}()
	time.Sleep(5 * time.Millisecond)
	x = 123
	atomic.StoreInt64(&m, 1)
	time.Sleep(10 * time.Millisecond)
}

func TestConcurrentSuObjectIter(t *testing.T) {
	if !race.Enabled {
		t.Skip("RACE NOT ENABLED")
	}
	ob := SuObjectOf(One, True)
	ob.SetConcurrent()
	for i := 0; i < 4; i++ {
		go func() {
			time.Sleep(5 * time.Millisecond)
			iter := ob.Iter2(true, true)
			iter()
		}()
	}
	ob.Add(One)
}
