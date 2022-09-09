// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable

package builtin

import (
	"io"
	"os"
	"runtime"
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	Interrupt = func() bool {
		select {
		case <-c:
			return true
		default:
			return false
		}
	}
}

func Run() {
}

var _ = builtin0("OperatingSystem()", func() Value { // deprecated
	return SuStr(runtime.GOOS)
})
var _ = builtin0("OSName()", func() Value {
	return SuStr(runtime.GOOS)
})

var _ = builtin0("GetComputerName()", func() Value {
	name, err := os.Hostname()
	if err != nil {
		panic("GetComputerName " + err.Error())
	}
	return SuStr(name)
})

var _ = builtin0("GetTempPath()",
	func() Value {
		return SuStr(os.TempDir())
	})

var _ = builtin3("CopyFile(from, to, failIfExists)",
	func(a, b, c Value) Value {
		from := ToStr(a)
		to := ToStr(b)
		failIfExists := ToBool(c)

		flags := os.O_WRONLY | os.O_CREATE
		if failIfExists {
			flags |= os.O_EXCL
		} else {
			flags |= os.O_TRUNC
		}

		srcFile, err := os.Open(from)
		if err != nil {
			return False
		}
		defer srcFile.Close()

		fi, err := srcFile.Stat()
		if err != nil {
			return False
		}

		destFile, err := os.OpenFile(to, flags, fi.Mode())
		if err != nil {
			return False
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return False
		}

		destFile.Close()
		os.Chtimes(to, time.Now(), fi.ModTime())

		return True
	})

func CallbacksCount() int {
	return 0
}

func WndProcCount() int {
	return 0
}

func GetGuiResources() (int, int) {
	return 0, 0
}

func ErrlogDir() string {
	return os.TempDir()
}

func OnUIThread() bool {
	return false
}
