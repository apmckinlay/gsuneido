// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/core"
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

var _ = builtin(Synchronized, "(block)")

func Synchronized(th *Thread, args []Value) Value {
	sy.threadLock.Lock()
	reentry := sy.lockThread == th
	sy.threadLock.Unlock()
	if reentry {
		return th.Call(args[0])
	}
	sy.lock.Lock()
	sy.threadLock.Lock()
	sy.lockThread = th
	sy.threadLock.Unlock()
	defer func() {
		sy.threadLock.Lock()
		sy.lockThread = nil
		sy.threadLock.Unlock()
		sy.lock.Unlock()
	}()
	return th.Call(args[0])
}
