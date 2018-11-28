package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	SuFuncMethods = Methods{
		"Disasm": method0(func(this Value) Value {
			fn := this.(*SuFunc)
			buf := &strings.Builder{}
			Disasm(buf, fn)
			return SuStr(buf.String())
		}),
		"Params": method0(func(this Value) Value {
			fn := this.(*SuFunc)
			return SuStr(fn.Params())
		}),
	}
	Params = method0(func(this Value) Value {
		ps := this.(Paramsable)
		return SuStr(ps.Params())
	})
}

type Paramsable interface {
	Params() string
}
