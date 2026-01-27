// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"fmt"
	"sync/atomic"
	"syscall"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/util/queue"
)

const timerInterval = 10 // milliseconds //???

func init() {
	// a Windows timer so it will get called on the main ui thread
	// and even if a Windows message loop is running e.g. in MessageBox
	// although timers are the lowest priority,
	// so it will only be called when there are no other messages to process.
	timerFn := syscall.NewCallback(func(a, b, c, d uintptr) uintptr {
		runDefer()
		return 0
	})
	syscall.SyscallN(setTimer, 0, 0, timerInterval, timerFn)
}

var deferQueue = queue.New[dqitem](32, 64)

type dqitem struct {
	id int
	fn Value
}

var dqid atomic.Int32

// dqMustPut panics if the queue is full
func dqMustPut(fn Value) int {
	id := int(dqid.Add(1))
	deferQueue.MustPut(dqitem{id: id, fn: fn})
	return id
}

func dqGet() Value {
	if it, ok := deferQueue.TryGet(); ok {
		return it.fn
	}
	return nil
}

func dqRemove(id int) bool {
	return deferQueue.Remove(func(it dqitem) bool { return it.id == id })
}

// runDefer runs all the pending deferred functions.
// It is called by a timer on the main UI thread.
func runDefer() {
	state := MainThread.GetState()
	defer MainThread.RestoreState(state)
	defer func() {
		if e := recover(); e != nil {
			handler(MainThread, e)
		}
	}()
	// Only run the ones currently in the queue, not the ones added by these.
	// Otherwise chaining runs continuously, blocking the GUI
	for range deferQueue.Size() {
		fn := dqGet()
		if fn == nil {
			break
		}
		MainThread.Call(fn)
		MainThread.RestoreState(state)
	}
}

var _ = AddInfo("windows.nDefer", deferQueue.Size)

//-------------------------------------------------------------------

var _ = builtin(Defer, "(block)")

func Defer(th *Thread, args []Value) Value {
	trace.Defer.Println("Defer", args[0])
	if th != MainThread {
		// because it will be executed on the main thread
		args[0].SetConcurrent() 
	}
	id := dqMustPut(args[0]) // can't block because MainThread is the consumer
	return &killer{kill: func() { dqRemove(id) }}
}

var _ = builtin(Delay, "(delayMs, block)")

func Delay(th *Thread, args []Value) Value {
	if th != MainThread {
		panic("Delay can only be called from the main GUI thread")
	}
	delay := ToInt(args[0])
	const minDelay = 100 // ms
	if delay < minDelay {
		panic(fmt.Sprint("Delay minimum is ", minDelay, " (ms)"))
	}
	fn := args[1]
	id := -1
	timer := time.AfterFunc(time.Duration(delay)*time.Millisecond,
		func() {
			trace.Defer.Println("Delay", fn)
			id = dqMustPut(fn)
		})
	return &killer{kill: func() {
		timer.Stop()
		if id >= 0 {
			dqRemove(id)
		}
	}}
}
