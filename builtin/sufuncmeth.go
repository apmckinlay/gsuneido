package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	SuFuncMethods = Methods{
		"Disasm": method0(func(self Value) Value {
			fn := self.(*SuFunc)
			buf := &strings.Builder{}
			Disasm(buf, fn)
			return SuStr(buf.String())
		}),
		"Params": method0(func(self Value) Value {
			fn := self.(*SuFunc)
			return SuStr(fn.Params())
		}),
	}
	Params = method0(func (self Value) Value {
		ps := self.(Paramsable)
		return SuStr(ps.Params())
	})
}

type Paramsable interface {
	Params() string
}
