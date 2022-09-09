// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build unix

package builtin

import (
	"io"
	"os"
	"syscall"
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("GetDiskFreeSpace(dir = '.')", func(arg Value) Value {
	var stat syscall.Statfs_t
	syscall.Statfs(ToStr(arg), &stat)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	return Int64Val(int64(freeBytes))
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
