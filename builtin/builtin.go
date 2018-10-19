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

func builtin0(s string, f func() Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &Builtin0{Fn: f, ParamSpec: params(p)})
	return true
}

func builtin1(s string, f func(a Value) Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &Builtin1{Fn: f, ParamSpec: params(p)})
	return true
}

func builtin2(s string, f func(a,b Value) Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &Builtin2{Fn: f, ParamSpec: params(p)})
	return true
}

func builtin3(s string, f func(a,b,c Value) Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &Builtin3{Fn: f, ParamSpec: params(p)})
	return true
}

func builtin4(s string, f func(a,b,c,d Value) Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &Builtin4{Fn: f, ParamSpec: params(p)})
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

func method0(f func(self Value) Value) Callable {
	return &Method0{Fn: f, ParamSpec: ParamSpec{}}
}

func method1(p string, f func(self, a1 Value) Value) Callable {
	return &Method1{Fn: f, ParamSpec: params(p)}
}

func method2(p string, f func(self, a1, a2 Value) Value) Callable {
	return &Method2{Fn: f, ParamSpec: params(p)}
}

func method3(p string, f func(self, a1, a2, a3 Value) Value) Callable {
	return &Method3{Fn: f, ParamSpec: params(p)}
}

func rawmethod(p string,
	f func(t *Thread, self Value, as *ArgSpec, args ...Value) Value) Callable {
	return &RawMethod{Fn: f, ParamSpec: params(p)}
}
