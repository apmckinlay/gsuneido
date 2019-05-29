package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Construct(what, suffix='')",
	func(t *Thread, args []Value) Value {
		what := args[0]
		suffix := ToStr(args[1])
		c, ok := what.ToContainer()
		if ok {
			what = c.ListGet(0)
			if what == nil {
				panic("Construct: object requires member 0")
			}
		}
		if s, ok := what.ToStr(); ok {
			if !strings.HasSuffix(s, suffix) {
				s += suffix
			}
			what = Global.GetName(t, s)
		}
		if c == nil {
			return t.CallLookup(what, "*new*")
		}
		return t.CallLookupEach1(what, "*new*", c)
	})
