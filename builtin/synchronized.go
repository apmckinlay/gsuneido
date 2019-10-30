package builtin

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/runtime"
)

// Synchronized is used recursively
// and Go doesn't have a reentrant mutex
// so we have to implement our own

type synchInfo struct {
	// lockThread is the current owner of the lock
	lockThread *Thread // guarded by threadLock
	threadLock sync.Mutex
	lock       sync.Mutex
}

var sy synchInfo

var _ = builtin("Synchronized(block)",
	synchronized)

func synchronized(t *Thread, args []Value) Value {
	sy.threadLock.Lock()
	reentry := sy.lockThread == t
	sy.threadLock.Unlock()
	if reentry {
		return t.Call(args[0])
	}
	sy.lock.Lock()
	sy.threadLock.Lock()
	sy.lockThread = t
	sy.threadLock.Unlock()
	defer func() {
		sy.threadLock.Lock()
		sy.lockThread = nil
		sy.threadLock.Unlock()
		sy.lock.Unlock()
	}()
	return t.Call(args[0])
}
