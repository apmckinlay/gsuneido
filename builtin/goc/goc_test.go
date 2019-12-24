// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build interactive

package goc

import (
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

var user32 = MustLoadDLL("user32.dll")
var messageBox = user32.MustFindProc("MessageBoxA").Addr()

func TestGoc(*testing.T) {
	Init()
	text := []byte("hello\x00")
	caption := []byte("world\x00")
	Syscall4(messageBox,
		0,
		uintptr(unsafe.Pointer(&text[0])),
		uintptr(unsafe.Pointer(&caption[0])),
		0)
}
