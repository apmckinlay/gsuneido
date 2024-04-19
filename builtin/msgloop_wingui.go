// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"runtime"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/core"
	"golang.org/x/sys/windows"
)

var _ = builtin(MessageLoop, "(hwnd = 0)")

func MessageLoop(_ *Thread, args []Value) Value {
	goc.MessageLoop(uintptr(ToInt(args[0])))
	return nil
}

// uiThreadId is used to detect calls from other threads (not allowed)
var uiThreadId uint32

func init() {
	// inject dependencies
	goc.Callback2 = callback2
	goc.Callback3 = callback3
	goc.Callback4 = callback4
	goc.SunAPP = sunAPP
	goc.RunOnGoSide = runOnGoSide
	RunOnGoSide = runOnGoSide // runtime
	Interrupt = goc.Interrupt
	goc.Shutdown = shutdown

	runtime.LockOSThread()
	uiThreadId = windows.GetCurrentThreadId()

	goc.Init()
}

func Run() {
	goc.Run()
}

func shutdown(exitcode int) {
	Exit(exitcode)
}

func OnUIThread() bool {
	return windows.GetCurrentThreadId() == uiThreadId
}
