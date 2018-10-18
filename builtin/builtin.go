package builtin

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
)

/* builtin defines a built in function in globals
for example:
var _ = builtin("Foo(a,b)", func(t *Thread, args ...Value) Value {
		...
	}))
*/
func builtin(s string, f func(t *Thread, args ...Value) Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &Builtin{Fn: f, ParamSpec: params(p)})
	return true
}

// params builds a ParamSpec from a string like (a, b) or (@args)
func params(s string) ParamSpec {
	fn := compile.Constant("function " + s + " {}").(*SuFunc)
	return fn.ParamSpec
}

func method(p string, f func(t *Thread, self Value, args ...Value) Value) Callable {
	return &Method{Fn: f, ParamSpec: params(p)}
}

func rawmethod(p string,
	f func(t *Thread, self Value, as *ArgSpec, args ...Value) Value) Callable {
	return &RawMethod{Fn: f, ParamSpec: params(p)}
}
