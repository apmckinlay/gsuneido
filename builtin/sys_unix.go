// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build unix

package builtin

import (
	"io"
	"log"
	"os"
	"syscall"
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(GetDiskFreeSpace, "(dir = '.')")

func GetDiskFreeSpace(arg Value) Value {
	var stat syscall.Statfs_t
	syscall.Statfs(ToStr(arg), &stat)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	return Int64Val(int64(freeBytes))
}

var _ = builtin(CopyFile, "(from, to, failIfExists)")

func CopyFile(a, b, c Value) Value {
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
	// needed when the destination is on a Samba network drive
	if err := destFile.Chmod(fi.Mode()); err != nil {
		log.Println("WARN CopyFile Chmod", err)
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		log.Println("WARN CopyFile Copy", err)
		return False
	}

	destFile.Close()
	if err := os.Chtimes(to, time.Now(), fi.ModTime()); err != nil {
		log.Println("WARN CopyFile Chtimes", err)
	}

	return True
}
