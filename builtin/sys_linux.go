// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"syscall"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("SystemMemory()", func() Value {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		panic(err)
	}
	return Int64Val(int64(info.Totalram))
})
