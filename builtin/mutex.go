// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

type suMutex struct {
	ValueBase[suMutex]
	mut MutexT
}

var _ = builtin(Mutex, "()")

func Mutex() Value {
	return &suMutex{mut: MakeMutexT()}
}

var suMutexMethods = methods()

var _ = method(mu_Do, "(block)")

func mu_Do(th *Thread, this Value, args []Value) Value {
	sm := this.(*suMutex)
	sm.mut.Lock()
	defer sm.mut.Unlock()
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
