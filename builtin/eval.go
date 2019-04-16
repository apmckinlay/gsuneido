package builtin

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/regex"
)

// Eval executes string containing Suneido code
func EvalString(t *Thread, s string) Value {
	s = strings.Trim(s, " \t\r\n")
	if isGlobal(s) {
		// optimize if just a global name
		return Global.GetName(s)
	}
	s = "function () {\n" + s + "\n}"
	fn := compile.NamedConstant("eval", s).(*SuFunc)
	return t.Call(fn)
}

var rxGlobal = regex.Compile("^[A-Z][_a-zA-Z0-9]*?[!?]?$")

func isGlobal(s string) bool {
	return rxGlobal.Matches(s)
}

// EvalAsMethod runs a function as if it were a method of an object
// implements object.Eval
func EvalAsMethod(t *Thread, as *ArgSpec, ob *SuObject, args []Value) Value {
	// first argument is function
	k, f := NewArgsIter(as, args)()
	if k != nil || f == nil {
		panic("usage: object.Eval(function, ...)")
	}

	as2 := *as
	if as.Each == EACH {
		as2.Each = EACH1
	} else if as.Each == EACH1 {
		panic("object.Eval does not support @+1")
	} else {
		as2.Nargs--
	}

	if m, ok := f.(*SuMethod); ok {
		f = m.GetFn()
	}

	return CallMethod(t, ob, f, &as2)
}
