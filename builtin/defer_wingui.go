// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"fmt"
	"sync/atomic"
	"syscall"

	. "github.com/apmckinlay/gsuneido/core"
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

func dqPut(fn Value) int {
	id := int(dqid.Add(1))
	deferQueue.Put(dqitem{id: id, fn: fn})
	return id
}

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

// runDefer runs all the pending deferred functions
func runDefer() {
	state := MainThread.GetState()
	defer func() {
		if e := recover(); e != nil {
			handler(e, state)
		}
		MainThread.RestoreState(state)
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
	if th != MainThread {
		panic("Defer can only be used from the main GUI thread")
	}
	id := dqMustPut(args[0]) // can't block because MainThread is the consumer
	return &killer{kill: func() { dqRemove(id) }}
}

var _ = builtin(RunOnGui, "(block)")

func RunOnGui(th *Thread, args []Value) Value {
	if th == MainThread {
		panic("RunOnGui can only be used from other threads")
	}
	id := dqPut(args[0]) // blocks if queue is full
	return &killer{kill: func() { dqRemove(id) }}
}

var _ = builtin(Delay, "(delayMs, block)")

func Delay(th *Thread, args []Value) Value {
	if th != MainThread {
		panic("Delay can only be called from the main GUI thread")
	}
	const minDelay = 100 // ms
	if ToInt(args[0]) < minDelay {
		panic(fmt.Sprint("Delay minimum is ", minDelay, " (ms)"))
	}
	tf := &timerFn{callback: args[1]}
	tf.timerid = SetTimer(Zero, Zero, args[0], tf)
	if tf.timerid == Zero {
		panic("Delay SetTimer failed")
	}
	return &killer{kill: func() { tf.kill() }}
}

func (tf *timerFn) kill() bool {
	if tf.timerid == Zero {
		return false
	}
	tid := tf.timerid
	tf.timerid = Zero
	KillTimer(Zero, tid)
	clearCallback(tf)
	return true
}

type timerFn struct {
	ValueBase[*timerFn]
	timerid  Value
	callback Value
}

var _ Value = (*timerFn)(nil)

func (tf *timerFn) Equal(other any) bool {
	return tf == other
}

func (tf *timerFn) Call(th *Thread, this Value, _ *ArgSpec) Value {
	if tf.kill() {
		tf.callback.Call(th, this, &ArgSpec0)
	}
	return nil
}
