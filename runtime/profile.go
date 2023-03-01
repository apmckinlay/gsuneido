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
	// lock is used to guard access to the frame stack
	lock sync.Mutex
	// total is the samples in functions and all the functions they call
	total map[string]int32
	// self is the samples in functions themselves
	self map[string]int32
	// ops is the number of operations/instructions executed in the interpreter
	ops map[string]int32
	// calls is the number of times the function is called
	calls map[string]int32
}

func (th *Thread) StartProfile() {
	th.profile.stop = make(chan struct{})
	th.profile.total = make(map[string]int32)
	th.profile.self = make(map[string]int32)
	th.profile.ops = make(map[string]int32)
	th.profile.calls = make(map[string]int32)
	th.profile.enabled = true
	go th.profiler()
}

func (th *Thread) profiler() {
	for {
		time.Sleep(1 * time.Millisecond)
		select {
		case <-th.profile.stop:
			return
		default:
		}
		th.sample()
	}
}

func (th *Thread) sample() {
	th.profile.lock.Lock()
	defer th.profile.lock.Unlock()
	if th.profile.enabled && th.fp > 0 {
		th.profile.self[th.frames[th.fp-1].fn.Name]++
		for i := th.fp - 1; i >= 0; i-- {
			fn := th.frames[i].fn
			th.profile.total[fn.Name]++
		}
	}
}

func (th *Thread) StopProfile() (total, self, ops, calls map[string]int32) {
	th.profile.lock.Lock()
	defer th.profile.lock.Unlock()
	if !th.profile.enabled {
		return
	}
	th.profile.enabled = false
	close(th.profile.stop)
	th.profile.stop = nil
	total = th.profile.total
	self = th.profile.self
	ops = th.profile.ops
	calls = th.profile.calls
	th.profile.total, th.profile.self, th.profile.ops, th.profile.calls =
		nil, nil, nil, nil
	return total, self, ops, calls
}
