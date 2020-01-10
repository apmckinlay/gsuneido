// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package console

import (
	"log"
	"os"
	"syscall"

	"github.com/apmckinlay/gsuneido/options"
	"golang.org/x/sys/windows"
)

func RedirOutput() {
	// try to ensure that output is captured, either to console or log
	attachedStdout, attachedStderr := outputToConsole()
	if !attachedStdout {
		stdoutToLog()
	}
	if !attachedStderr {
		stderrToLog()
		log.SetOutput(os.Stderr) // new Stderr
	} else {
		LogFileAlso()
	}
	if options.Repl {
		inputFromConsole()
	}
}

func stdoutToLog() {
	f, err := os.OpenFile(options.Outlog, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil && f != nil {
		os.Stdout = f
	}
}

func stderrToLog() {
	f, err := os.OpenFile(options.Errlog, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil && f != nil {
		os.Stderr = f
	}
}

var consoleAttached = GetConsoleWindow() != 0 || AttachParentConsole()
var stdoutRedirected = redirected(windows.STD_OUTPUT_HANDLE)
var stderrRedirected = redirected(windows.STD_ERROR_HANDLE)

func outputToConsole() (stdoutAttached bool, stderrAttached bool) {
	if !stdoutRedirected || !stderrRedirected {
		// If not directed to a file, try attaching to console,
		if consoleAttached {
			if f, _ := os.OpenFile("CONOUT$", os.O_WRONLY, 0644); f != nil {
				if !stdoutRedirected {
					os.Stdout = f
				}
				if !stderrRedirected {
					os.Stderr = f
				}
				return !stdoutRedirected, !stderrRedirected
			}
		}
	}
	return false, false
}

func inputFromConsole() {
	attachStdIn := !redirected(windows.STD_INPUT_HANDLE)
	if attachStdIn {
		OutputToConsole()
		if consoleAttached {
			if f, e := os.Open("CONIN$"); e == nil {
				os.Stdin = f
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

func OutputToConsole() {
	if !consoleAttached && (!stdoutRedirected || !stderrRedirected) {
		consoleAttached = AllocConsole()
		outputToConsole()
	}
}

func ConsoleAttached() bool {
	return consoleAttached
}

var kernel32 = windows.MustLoadDLL("kernel32.dll")

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

const ATTACH_PARENT_PROCESS = ^uintptr(0) // -1

var getConsoleWindow = kernel32.MustFindProc("GetConsoleWindow").Addr()

func GetConsoleWindow() uintptr {
	rtn, _, _ := syscall.Syscall(getConsoleWindow, 0, 0, 0, 0)
	return rtn
}
