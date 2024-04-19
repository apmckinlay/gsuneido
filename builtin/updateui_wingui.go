// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"log"
	"syscall"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"golang.org/x/sys/windows"
)

// rogsChan is used by other threads to Run code On the Go Side UI thread
// Need buffer so we can send to channel and then notifyCside
var rogsChan = make(chan func(), 1)

// UpdateUI runs the block on the main UI thread
var _ = builtin(UpdateUI, "(block)")

func UpdateUI(th *Thread, args []Value) Value {
	block := args[0]
	if windows.GetCurrentThreadId() == uiThreadId {
		th.Call(block)
	} else {
		block.SetConcurrent()
		rogsChan <- func() { runUI(block) }
		notifyCside()
	}
	return nil
}

const notifyWparam = 0xffffffff

// notifyCside is used by UpdateUI, SetTimer, and KillTimer
func notifyCside() {
	// NOTE: this has to be the Go Syscall, not goc.Syscall
	r, _, _ := syscall.SyscallN(postMessage,
		goc.CHelperHwnd(), WM_USER, notifyWparam, 0)
	if r == 0 {
		log.Panicln("notifyCside PostMessage failed")
	}
}

// runOnGoSide is used by runtime.RunOnGoSide (called by interp)
// and goc.RunOnGoSide (called by cside)
func runOnGoSide() {
	for {
		select {
		case fn := <-rogsChan:
			fn()
		default: // non-blocking
			return
		}
	}
}

func runUI(block Value) {
	state := MainThread.GetState()
	defer func() {
		if e := recover(); e != nil {
			log.Println("ERROR in UpdateUI:", e)
			MainThread.PrintStack()
			dbg.PrintStack()
		}
		MainThread.RestoreState(state)
	}()
	MainThread.Call(block)
}
