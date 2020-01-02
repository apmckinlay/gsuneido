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

// uuiChan is used for cross thread UpdateUI (e.g. Print)
// Need buffer of 1 so UpdateUI can send to channel and then SendThreadMessage
var uuiChan = make(chan Value, 1)

// UpdateUI runs the block on the main UI thread
var _ = builtin("UpdateUI(block)",
	func(t *Thread, args []Value) Value {
		if windows.GetCurrentThreadId() == uiThreadId {
			synchronized(t, args)
		} else {
			block := args[0]
			block.SetConcurrent()
			uuiChan <- block
			notifyMessageLoop()
		}
		return nil
	})

func notifyMessageLoop() {
	// NOTE: this has to be the Go Syscall, not goc.Syscall
	r, _, _ := syscall.Syscall6(postThreadMessage, 4,
		goc.CThreadId(), WM_USER, 0xffffffff, 0, 0, 0)
	if r == 0 {
		log.Panicln("PostThreadMessage failed")
	}
}

func updateUI2() {
	for {
		select {
		case block := <-uuiChan:
			runUI(block)
		case t := <-setTimerChan:
			timerIdChan <- gocSetTimer(t.hwnd, t.id, t.ms, t.cb)
		default: // non-blocking
			return
		}
	}
}

func updateUI() {
	for {
		select {
		case block := <-uuiChan:
			runUI(block)
		default: // non-blocking
			return
		}
	}
}

var updateThread *Thread

func runUI(block Value) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("error in UpdateUI:", e)
		}
	}()
	if updateThread == nil {
		updateThread = UIThread.SubThread()
	}
	updateThread.Call(block)
}
