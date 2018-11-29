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

func builtin2(s string, f func(a, b Value) Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &Builtin2{Fn: f, ParamSpec: params(p)})
	return true
}

func rawbuiltin(s string, f func(t *Thread, as *ArgSpec, args ...Value) Value) bool {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	p := s[i:]
	AddGlobal(name, &RawBuiltin{Fn: f, ParamSpec: params(p)})
	return true
}

// params builds a ParamSpec from a string like (a, b) or (@args)
func params(s string) ParamSpec {
	fn := compile.Constant("function " + s + " {}").(*SuFunc)
	return fn.ParamSpec
}

func method(p string, f func(t *Thread, this Value, args ...Value) Value) Value {
	return &Method{Fn: f, ParamSpec: params(p)}
}

func method0(f func(this Value) Value) Value {
	return &Method0{Builtin1{Fn: f, ParamSpec: ParamSpec{}}}
}

func method1(p string, f func(this, a1 Value) Value) Value {
	return &Method1{Builtin2{Fn: f, ParamSpec: params(p)}}
}

func rawmethod(p string,
	f func(t *Thread, as *ArgSpec, this Value, args ...Value) Value) Value {
	// params are just for documentation, RawMethod doesn't use them
	return &RawMethod{Fn: f, ParamSpec: params(p)}
}
