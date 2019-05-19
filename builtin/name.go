package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = builtin1("Name(value)",
	func(arg Value) Value {
		if named, ok := arg.(Named); ok {
			return SuStr(str.AfterFirst(named.GetName(), ":"))
		}
		return EmptyStr
	})
