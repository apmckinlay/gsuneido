// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !portable
// +build !portable

package trace

import (
	"os"
	"sync"
	"syscall"

	"github.com/apmckinlay/gsuneido/options"
	"golang.org/x/sys/windows"
)

var traceConOnce sync.Once
var traceCon *os.File

func consolePrintln(s string) {
	traceConOnce.Do(func() {
		traceCon = os.Stdout
		if options.Mode == "gui" {
			allocConsole()
			f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0644)
			if err == nil && f != nil {
				traceCon = f
			}
		}
	})
	traceCon.WriteString(s)
}

var kernel32 = windows.MustLoadDLL("kernel32.dll")

var allocConsoleAddr = kernel32.MustFindProc("AllocConsole").Addr()

func allocConsole() bool {
	rtn, _, _ := syscall.Syscall(allocConsoleAddr, 0, 0, 0, 0)
	return rtn != 0
}
