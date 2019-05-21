package builtin

import (
	"encoding/binary"
	"syscall"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

var _ = builtin0("SystemMemory()", func() Value {
	s, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		panic(err)
	}
	var buf [8]byte
	copy(buf[:], s)
	m := binary.LittleEndian.Uint64(buf[:])
	return SuDnum{Dnum: dnum.FromInt(int64(m))}
})

var _ = builtin0("OperatingSystem()", func() Value {
	return SuStr("darwin")
})
