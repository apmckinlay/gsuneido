package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var atl = windows.NewLazyDLL("atl.dll")

// dll bool atl:AtlAxWinInit()
var atlAxWinInit = atl.NewProc("AtlAxWinInit")
var _ = builtin0("AtlAxWinInit()",
	func() Value {
		rtn, _, _ := atlAxWinInit.Call()
		return boolRet(rtn)
	})
