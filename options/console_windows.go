// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package options

import (
	"os"
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
)

var traceConOnce sync.Once
var traceCon *os.File

func Console(s string) {
	traceConOnce.Do(func() {
		traceCon = os.Stdout
		if Mode == "gui" {
			AllocConsole()
			f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0644)
			if err == nil && f != nil {
				traceCon = f
			}
		}
	})
	traceCon.WriteString(s)
}

var kernel32 = windows.MustLoadDLL("kernel32.dll")

var allocConsole = kernel32.MustFindProc("AllocConsole").Addr()

func AllocConsole() bool {
	rtn, _, _ := syscall.Syscall(allocConsole, 0, 0, 0, 0)
	return rtn != 0
}
