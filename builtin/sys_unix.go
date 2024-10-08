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

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(GetDiskFreeSpace, "(dir = '.')")

func GetDiskFreeSpace(arg Value) Value {
	var stat syscall.Statfs_t
	err := syscall.Statfs(ToStr(arg), &stat)
	if err != nil {
		panic("GetDiskFreeSpace: " + err.Error())
	}
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	return Int64Val(int64(freeBytes))
}

var _ = builtin(CopyFile, "(from, to, failIfExists)")

func CopyFile(th *Thread, args []Value) Value {
	from := ToStr(args[0])
	to := ToStr(args[1])
	failIfExists := ToBool(args[2])
	th.ReturnThrow = true

	srcFile, err := os.Open(from)
	if err != nil {
		return SuStr("CopyFile: " + err.Error())
	}
	defer srcFile.Close()

	fi, err := srcFile.Stat()
	if err != nil {
		return SuStr("CopyFile: " + err.Error())
	}

	flags := os.O_WRONLY | os.O_CREATE
	if failIfExists {
		flags |= os.O_EXCL
	} else {
		flags |= os.O_TRUNC
	}
	destFile, err := os.OpenFile(to, flags, fi.Mode())
	if err != nil {
		return SuStr("CopyFile: " + err.Error())
	}
	defer destFile.Close()
	// needed when the destination is on a Samba network drive
	if err := destFile.Chmod(fi.Mode()); err != nil {
		log.Println("WARNING: CopyFile Chmod", err)
	}

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return SuStr("CopyFile: " + err.Error())
	}

	destFile.Close()
	if err := os.Chtimes(to, time.Now(), fi.ModTime()); err != nil {
		log.Println("WARNING: CopyFile Chtimes", err)
	}

	return True
}

// Go os.Remove handles both files and directories.
// We only want to remove files.

func deleteFile(path string) error {
	// loop on EINTR the same as os.Remove
	for {
		err := syscall.Unlink(path)
		if err != syscall.EINTR {
			return err
		}
	}
}
