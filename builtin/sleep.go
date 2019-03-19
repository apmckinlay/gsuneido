package builtin

import (
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Sleep(ms)",
	func(arg Value) Value {
		ms := ToInt(arg)
		time.Sleep(time.Duration(int64(ms) * 1000000))
		return nil
	})
