// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows

package builtin

import (
	"encoding/binary"
	"syscall"
)

func systemMemory() uint64 {
	s, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		panic(err)
	}
	var buf [8]byte
	copy(buf[:], s)
	m := binary.LittleEndian.Uint64(buf[:])
	return m
}
