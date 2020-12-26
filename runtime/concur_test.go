// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestConcurrentAtomic(*testing.T) {
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

func TestConcurrentMutex(*testing.T) {
	// This demonstrates that atomic is also sufficient
	// although the code doesn't currently work like this
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
