// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"sync"
	"time"
)

// total and self come from examining the frame stack
// once per millisecond (1000 times/second).
// They are not exact and will only be meaningful in aggregate
// for a large enough sample period.

type profile struct {
	enabled bool
	stop    chan struct{}
	// profileLock is used to guard access to the frame stack
	lock  sync.Mutex
	// total is the samples in functions and all the functions they call
	total map[string]int32
	// self is the samples in functions themselves
	self  map[string]int32
	// ops is the number of operations/instructions executed in the interpreter
	ops   map[string]int32
	// calls is the number of times the function is called
	calls map[string]int32
}

func (t *Thread) StartProfile() {
	t.profile.stop = make(chan struct{})
	t.profile.total = make(map[string]int32)
	t.profile.self = make(map[string]int32)
	t.profile.ops = make(map[string]int32)
	t.profile.calls = make(map[string]int32)
	t.profile.enabled = true
	go t.profiler()
}

func (t *Thread) profiler() {
	for {
		time.Sleep(1 * time.Millisecond)
		select {
		case <-t.profile.stop:
			return
		default:
		}
		t.sample()
	}
}

func (t *Thread) sample() {
	t.profile.lock.Lock()
	defer t.profile.lock.Unlock()
	if t.profile.enabled && t.fp > 0 {
		t.profile.self[t.frames[t.fp-1].fn.Name]++
		for i := t.fp - 1; i >= 0; i-- {
			fn := t.frames[i].fn
			t.profile.total[fn.Name]++
		}
	}
}

func (t *Thread) StopProfile() (total, self, ops, calls map[string]int32) {
	t.profile.lock.Lock()
	defer t.profile.lock.Unlock()
	if !t.profile.enabled {
		return
	}
	t.profile.enabled = false
	close(t.profile.stop)
	t.profile.stop = nil
	total = t.profile.total
	self = t.profile.self
	ops = t.profile.ops
	calls = t.profile.calls
	t.profile.total, t.profile.self, t.profile.ops, t.profile.calls =
		nil, nil, nil, nil
	return total, self, ops, calls
}
