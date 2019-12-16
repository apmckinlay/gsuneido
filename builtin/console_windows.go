// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func connectToConsole(repl bool) {
	attachStdIn := repl && !redirected(windows.STD_INPUT_HANDLE)
	attachStdOut := !redirected(windows.STD_OUTPUT_HANDLE)
	attachStdErr := !redirected(windows.STD_ERROR_HANDLE)

	if attachStdIn || attachStdOut || attachStdErr {
		// If not directed to a file, try attaching to console,
		// or if repl, create a new one
		attachedToConsole := AttachParentConsole()
		if !attachedToConsole && repl {
			attachedToConsole = AllocConsole()
		}

		if attachedToConsole {
			if attachStdIn {
				if f, e := os.Open("CONIN$"); e == nil {
					os.Stdin = f
				}
			}
			if attachStdOut || attachStdErr {
				if f, _ := os.OpenFile("CONOUT$", os.O_WRONLY, 0644); f != nil {
					if attachStdOut {
						os.Stdout = f
					}
					if attachStdErr {
						os.Stderr = f
					}
				}
			}
		}
	}
}

func redirected(which uint32) bool {
	handle, _ := windows.GetStdHandle(which)
	if handle != 0 {
	dwFileType, _ := windows.GetFileType(handle)
	return dwFileType == windows.FILE_TYPE_DISK ||
		dwFileType == windows.FILE_TYPE_PIPE
	}
	return false
}

var allocConsole = kernel32.MustFindProc("AllocConsole").Addr()

func AllocConsole() bool {
	rtn, _, _ := syscall.Syscall(allocConsole, 0, 0, 0, 0)
	return rtn != 0
}

var attachConsole = kernel32.MustFindProc("AttachConsole").Addr()

func AttachParentConsole() bool {
	rtn, _, _ := syscall.Syscall(attachConsole, 1, ATTACH_PARENT_PROCESS, 0, 0)
	return rtn != 0
}

const ATTACH_PARENT_PROCESS = uintptrMinusOne
