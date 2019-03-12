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
	name, ps := paramSplit(s)
	Global.Add(name, &SuBuiltin{f, BuiltinParams{ParamSpec: ps}})
	return true
}

func builtin0(s string, f func() Value) bool {
	name, ps := paramSplit(s)
	Global.Add(name, &SuBuiltin0{f, BuiltinParams{ParamSpec: ps}})
	return true
}

func builtin1(s string, f func(a Value) Value) bool {
	name, ps := paramSplit(s)
	Global.Add(name, &SuBuiltin1{f, BuiltinParams{ParamSpec: ps}})
	return true
}

func builtin2(s string, f func(a, b Value) Value) bool {
	name, ps := paramSplit(s)
	Global.Add(name, &SuBuiltin2{f, BuiltinParams{ParamSpec: ps}})
	return true
}

func builtin3(s string, f func(a, b, c Value) Value) bool {
	name, ps := paramSplit(s)
	Global.Add(name, &SuBuiltin3{f, BuiltinParams{ParamSpec: ps}})
	return true
}

func builtinRaw(s string, f func(t *Thread, as *ArgSpec, args ...Value) Value) bool {
	name, ps := paramSplit(s)
	Global.Add(name, &SuBuiltinRaw{f, BuiltinParams{ParamSpec: ps}})
	return true
}

// paramSplit takes Foo(x, y) and returns name and ParamSpec
func paramSplit(s string) (string, ParamSpec) {
	i := strings.IndexByte(s, byte('('))
	name := s[:i]
	ps := params(s[i:])
	ps.Name = name
	return name, ps
}

func method(p string, f func(t *Thread, this Value, args ...Value) Value) Value {
	return &SuBuiltinMethod{Fn: f, ParamSpec: params(p)}
}

func method0(f func(this Value) Value) Value {
	return &SuBuiltinMethod0{SuBuiltin1{f, BuiltinParams{}}}
}

func method1(p string, f func(this, a1 Value) Value) Value {
	return &SuBuiltinMethod1{SuBuiltin2{f, BuiltinParams{ParamSpec: params(p)}}}
}

func methodRaw(p string,
	f func(t *Thread, as *ArgSpec, this Value, args ...Value) Value) Value {
	// params are just for documentation, SuBuiltinMethodRaw doesn't use them
	return &SuBuiltinMethodRaw{Fn: f, ParamSpec: params(p)}
}

// params builds a ParamSpec from a string like (a, b) or (@args)
func params(s string) ParamSpec {
	fn := compile.Constant("function " + s + " {}").(*SuFunc)
	for i := 0; i < int(fn.ParamSpec.Ndefaults); i++ {
		if fn.Values[i].Equal(SuStr("nil")) {
			fn.Values[i] = nil
		}
	}
	return fn.ParamSpec
}
