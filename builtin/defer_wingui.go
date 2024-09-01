// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"fmt"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/util/queue"
	"golang.org/x/sys/windows"
)

func init() {
	trace.SetupConsole = func() {
		if windows.GetCurrentThreadId() == uiThreadId {
			goc.SetupConsole()
		} else {
			dqMustPut(builtinVal("SetupConsole",
				func() Value { goc.SetupConsole(); return nil }, "()"))
		}
	}
}

var deferQueue = queue.New[dqitem](32, 64)

type dqitem struct {
	id int
	fn Callable
}

var dqid atomic.Int32

func dqPut(fn Callable) int {
	id := int(dqid.Add(1))
	deferQueue.Put(dqitem{id: id, fn: fn})
	return id
}

func dqMustPut(fn Callable) int {
	id := int(dqid.Add(1))
	deferQueue.MustPut(dqitem{id: id, fn: fn})
	return id
}

func dqGet() Callable {
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
	for fn := dqGet(); fn != nil; fn = dqGet() {
		MainThread.Call(fn)
		MainThread.RestoreState(state)
	}
}
var _ = AddInfo("windows.nDefer", deferQueue.Size)

//-------------------------------------------------------------------

var _ = builtin(Defer, "(block)")

func Defer(th *Thread, args []Value) Value {
	if th != MainThread {
		panic("ERROR: Defer can only be used from the main GUI thread")
	}
	id := dqMustPut(args[0]) // can't block because MainThread is the consumer
	return &killer{kill: func() { dqRemove(id) }}
}

var _ = builtin(RunOnGui, "(block)")

func RunOnGui(th *Thread, args []Value) Value {
	if th == MainThread {
		panic("ERROR: RunOnGui can only be used from other threads")
	}
	id := dqPut(args[0]) // blocks if queue is full
	return &killer{kill: func() { dqRemove(id) }}
}

var _ = builtin(Delay, "(delayMs, block)")

func Delay(th *Thread, args []Value) Value {
	if th != MainThread {
		panic("ERROR: Delay can only be called from the main GUI thread")
	}
	const minDelay = 100 // ms
	if ToInt(args[0]) < minDelay {
		panic(fmt.Sprint("ERROR: Delay minimum is ", minDelay, " (ms)"))
	}
	tf := &timerFn{callback: args[1]}
	tf.timerid = gocSetTimer(Zero, Zero, args[0], tf)
	if tf.timerid == Zero {
		panic("ERROR: Delay SetTimer failed")
	}
	return &killer{kill: func() { tf.kill() }}
}

func (tf *timerFn) kill() bool {
	if tf.timerid == Zero {
		return false
	}
	tid := tf.timerid
	tf.timerid = Zero
	gocKillTimer(Zero, tid)
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

//-------------------------------------------------------------------

type killer struct {
	ValueBase[*killer]
	kill func()
}

var _ Value = (*killer)(nil)

func (k *killer) Equal(other any) bool {
	return k == other
}

func (k *killer) SetConcurrent() {
	// need to allow this because of saving in concurrent places
	// still shouldn't be calling it from other threads
}

func (k *killer) Lookup(_ *Thread, method string) Callable {
	return killerMethods[method]
}

var killerMethods = methods()

var _ = method(killer_Kill, "()")

func killer_Kill(this Value) Value {
	this.(*killer).kill()
    return nil
}
