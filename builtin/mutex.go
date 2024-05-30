// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

type suMutex struct {
	ValueBase[suMutex]
	ch chan struct{}
}

var _ = builtin(Mutex, "()")

func Mutex() Value {
	return &suMutex{ch: make(chan struct{}, 1)}
}

var suMutexMethods = methods()

var _ = method(mu_Do, "(block)")

func mu_Do(th *Thread, this Value, args []Value) Value {
	sm := this.(*suMutex)
	sm.lock()
	defer sm.unlock()
	return th.Call(args[0])
}

// Value implementation

var _ Value = (*suMutex)(nil)

func (sm *suMutex) Equal(other any) bool {
	return sm == other
}

func (*suMutex) Lookup(_ *Thread, method string) Callable {
	return suMutexMethods[method]
}

func (*suMutex) SetConcurrent() {
	// ok for concurrent use
}

//-------------------------------------------------------------------

func (sm *suMutex) lock() {
	select {
	case sm.ch <- struct{}{}:
		// lock acquired
	case <-time.After(10 * time.Second):
		panic("Mutex: lock timeout")
	}
}

func (sm *suMutex) unlock() {
	select {
	case <-sm.ch:
		// lock released
	default:
		panic("Mutex: unlock failed")
	}
}
