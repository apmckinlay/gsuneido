// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/core"
	"golang.org/x/sys/windows"
)

func init() {
	uiThreadId = windows.GetCurrentThreadId()
	goc.SunAPP = sunAPP
	core.Interrupt = goc.Interrupt
}

var uiThreadId uint32

func OnUIThread() bool {
	return windows.GetCurrentThreadId() == uiThreadId
}

func Run() int {
	return goc.Run()
}

var _ = builtin(MessageLoop, "(hwnd = 0)")

func MessageLoop(a Value) Value {
	goc.MessageLoop(uintptr(ToInt(a)))
	return nil
}
