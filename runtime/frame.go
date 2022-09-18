// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

// Frame is the context for a function/method/block invocation.
type Frame struct {
	// fn is the Function being executed
	fn *SuFunc

	// ip is the current index into the Function's code
	ip int

	// locals are the local variables (including arguments)
	// Normally they are on the thread stack
	// but for closure blocks they are moved to the heap.
	locals locals

	// this is the instance if we're running a method
	this Value

	// blockParent is used for block returns
	blockParent *Frame
}

type locals struct {
	v []Value
	// onHeap is true when locals have been moved from the stack to the heap
	onHeap bool
}

func (ls *locals) moveToHeap() {
	if ls.onHeap {
		return
	}
	// not concurrent at this point
	oldlocals := ls.v
	ls.v = make([]Value, len(oldlocals))
	copy(ls.v, oldlocals)
	ls.onHeap = true
}
