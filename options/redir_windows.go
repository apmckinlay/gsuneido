// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !portable
// +build !portable

package options

import "golang.org/x/sys/windows"

func Redirected() bool {
	handle, _ := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if handle != 0 {
		dwFileType, _ := windows.GetFileType(handle)
		return dwFileType == windows.FILE_TYPE_DISK ||
			dwFileType == windows.FILE_TYPE_PIPE
	}
	return false
}
