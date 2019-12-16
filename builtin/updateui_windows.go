// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"log"
	"syscall"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var _ = builtin("UpdateUI(block)",
	func(t *Thread, args []Value) Value {
		if windows.GetCurrentThreadId() == uiThreadId {
			synchronized(t, args)
		} else {
			block := args[0]
			block.SetConcurrent()
			ChanUI <- block
			// NOTE: this has to be the Go Syscall, not goc.Syscall
			r, _, _ := syscall.Syscall6(postThreadMessage, 4,
				goc.CThreadId(), WM_USER, 0xffffffff, 0, 0, 0)
			if r == 0 {
				log.Panicln("UpdateUI failed")
			}
		}
		return nil
	})

var updateThread *Thread

func updateUI() {
	select {
	case block := <-ChanUI:
		defer func() {
			if e := recover(); e != nil {
				log.Println("error in UpdateUI:", e)
			}
		}()
		if updateThread == nil {
			updateThread = UIThread.SubThread()
		}
		updateThread.Call(block)
	default: // non-blocking
	}
}

// ChanUI is used for cross thread UpdateUI (e.g. Print)
// Need buffer of 1 so UpdateUI can send to channel and then SendThreadMessage
var ChanUI = make(chan Value, 1)
