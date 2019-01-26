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
	locals []Value
	// this is the instance if we're running a method
	this Value
	// localsOnHeap is true when locals have been moved from the stack to the heap
	// for blocks
	localsOnHeap bool
}

func (fr *Frame) moveLocalsToHeap() {
	if fr.localsOnHeap {
		return
	}
	oldlocals := fr.locals
	fr.locals = make([]Value, len(oldlocals))
	copy(fr.locals, oldlocals)
	fr.localsOnHeap = true
}
