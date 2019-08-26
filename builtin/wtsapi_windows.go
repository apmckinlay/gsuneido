package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var wtsapi32 = windows.NewLazyDLL("wtsapi32.dll")

// dll void WTSAPI32:WTSFreeMemory(pointer adr)
var wtsFreeMemory = wtsapi32.NewProc("WTSFreeMemory")
var _ = builtin1("WTSFreeMemory(adr)",
	func(a Value) Value {
		wtsFreeMemory.Call(
			intArg(a))
			return nil
		})
