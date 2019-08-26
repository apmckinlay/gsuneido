package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var comdlg32 = windows.NewLazyDLL("comdlg32.dll")

// dll long ComDlg32:CommDlgExtendedError()
var commDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")
var _ = builtin0("CommDlgExtendedError()",
	func() Value {
		rtn, _, _ := commDlgExtendedError.Call()
		return intRet(rtn)
	})
