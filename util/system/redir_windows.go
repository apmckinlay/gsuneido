// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package system

import (
	"log"
	"os"

	"golang.org/x/sys/windows"
)

// Redirect redirects stderr and stdout to a file
// (unless already redirected).
// This is to capture Go errors when running as a GUI program or service.
func Redirect(filename string) error {
	if redirected() {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	log.SetOutput(f) // redundant, but just to make sure
	err = windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(f.Fd()))
	if err != nil {
		return err
	}
	err = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(f.Fd()))
	if err != nil {
		return err
	}
	// need these because SetStdHandle does not affect prior references
	os.Stderr = f
	os.Stdout = f
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
