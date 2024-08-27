// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package queue

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestQueue(t *testing.T) {
	const nthreads = 10
	const nitemsperthread = 1000
	q := New[int](10, 20)
	var np atomic.Int32
	var rw sync.RWMutex
	for range nthreads {
		rw.RLock()
		go func() {
			defer rw.RUnlock()
			for i := range nitemsperthread {
				q.Put(i)
				if !q.Remove(func(x int) bool { return x == i-9 }) {
					np.Add(1)
				}
			}
		}()
	}
	n := 0
	for !rw.TryLock() {
		if _, ok := q.TryGet(); ok {
			n++
		}
	}
	for {
		if _, ok := q.TryGet(); !ok {
			break
		}
		n++
	}
	// fmt.Println(n)
	assert.T(t).This(n).Is(int(np.Load()))
}
