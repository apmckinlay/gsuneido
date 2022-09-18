// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows

package builtin

import (
	"encoding/binary"
	"syscall"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(SystemMemory, "()")

func SystemMemory() Value {
	s, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		panic(err)
	}
	var buf [8]byte
	copy(buf[:], s)
	m := binary.LittleEndian.Uint64(buf[:])
	return Int64Val(int64(m))
}
