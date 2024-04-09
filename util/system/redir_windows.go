// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package system

import (
	"log"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

// Redirect redirects stderr and stdout to a file
// (unless already redirected).
// This is to capture Go errors when running as a service or GUI program.
func Redirect(filename string) error {
	if redirected() {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	wh := windows.Handle(f.Fd())
	err = windows.SetStdHandle(windows.STD_ERROR_HANDLE, wh)
	if err != nil {
		return err
	}
	err = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, wh)
	if err != nil {
		return err
	}
	// redo initialization
	os.Stdout = f
	os.Stderr = f
	log.SetOutput(f)
	syscall.Stdout = syscall.Handle(wh)
	syscall.Stderr = syscall.Handle(wh)
	return nil
}

func redirected() bool {
	handle, _ := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if handle != 0 {
		dwFileType, _ := windows.GetFileType(handle)
		return dwFileType == windows.FILE_TYPE_DISK ||
			dwFileType == windows.FILE_TYPE_PIPE
	}
	return false
}
