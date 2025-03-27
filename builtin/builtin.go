// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"maps"
	"reflect"
	"runtime"
	"slices"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// builtin function names can end in Q or X for ? or !
func builtin(f any, p string) any {
	name := funcName(f)
	name = str.Capitalize(name)
	Global.Builtin(name, builtinVal(name, f, p))
	return nil
}

func builtinVal(name string, f any, p string) Value {
	ps := params(p)
	ps.Name = name
	switch f := f.(type) {
	case func() Value:
		assert.That(ps.Nparams == 0)
		return &SuBuiltin0{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(Value) Value:
		assert.That(ps.Nparams == 1)
		return &SuBuiltin1{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(Value, Value) Value:
		assert.That(ps.Nparams == 2)
		return &SuBuiltin2{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(Value, Value, Value) Value:
		assert.That(ps.Nparams == 3)
		return &SuBuiltin3{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(Value, Value, Value, Value) Value:
		assert.That(ps.Nparams == 4)
		return &SuBuiltin4{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(Value, Value, Value, Value, Value) Value:
		assert.That(ps.Nparams == 5)
		return &SuBuiltin5{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(Value, Value, Value, Value, Value, Value) Value:
		assert.That(ps.Nparams == 6)
		return &SuBuiltin6{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(Value, Value, Value, Value, Value, Value, Value) Value:
		assert.That(ps.Nparams == 7)
		return &SuBuiltin7{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(th *Thread, args []Value) Value:
		return &SuBuiltin{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(th *Thread, as *ArgSpec, args []Value) Value:
		// params are just for documentation
		return &SuBuiltinRaw{Fn: f, BuiltinParams: BuiltinParams{ParamSpec: ps}}
	default:
		panic("invalid builtin function: " + name)
	}
}

// curMethods and curPrefix are only used during initialization
// which is single threaded, so no locking is required
var curMethods = make(map[string]Methods)

func methods(prefix string) Methods {
	if _, ok := curMethods[prefix]; ok {
		Fatal("duplicate methods prefix:", prefix)
	}
	curMethods[prefix] = Methods{}
	return curMethods[prefix]
}

func exportMethods(m *Methods, prefix string) any {
	*m = methods(prefix)
	return nil
}

// staticMethod adds to curMethods, like method,
// but creates a standalone function like builtin
func staticMethod(f any, p string) any {
	name, meths := methodName(f)
	fn := builtinVal(name, f, p)
	meths[name] = fn
	return nil
}

// method adds to curMethods, which is set by methods.
// method function names must start with a prefix e.g. xyz_
func method(f any, p string) any {
	name, meths := methodName(f)
	ps := params(p)
	ps.Name = name
	switch f := f.(type) {
	case func(this Value) Value:
		assert.That(ps.Nparams == 0)
		meths[name] = &SuBuiltinMethod0{SuBuiltin1: SuBuiltin1{Fn: f,
			BuiltinParams: BuiltinParams{ParamSpec: ps}}}
	case func(this, a Value) Value:
		assert.That(ps.Nparams == 1)
		meths[name] = &SuBuiltinMethod1{SuBuiltin2: SuBuiltin2{Fn: f,
			BuiltinParams: BuiltinParams{ParamSpec: ps}}}
	case func(this, a, b Value) Value:
		assert.That(ps.Nparams == 2)
		meths[name] = &SuBuiltinMethod2{SuBuiltin3: SuBuiltin3{Fn: f,
			BuiltinParams: BuiltinParams{ParamSpec: ps}}}
	case func(this, a, b, c Value) Value:
		assert.That(ps.Nparams == 3)
		meths[name] = &SuBuiltinMethod3{SuBuiltin4: SuBuiltin4{Fn: f,
			BuiltinParams: BuiltinParams{ParamSpec: ps}}}
	case func(th *Thread, this Value, args []Value) Value:
		meths[name] = &SuBuiltinMethod{Fn: f,
			BuiltinParams: BuiltinParams{ParamSpec: ps}}
	case func(th *Thread, as *ArgSpec, this Value, args []Value) Value:
		// params are just for documentation
		meths[name] = &SuBuiltinMethodRaw{Fn: f, 
			BuiltinParams: BuiltinParams{ParamSpec: params(p)}}
	default:
		Fatal("invalid builtin function", name)
	}
	return nil
}

func methodName(f any) (string, Methods) {
	fname := funcName(f)
	prefix, name, _ := strings.Cut(fname, "_")
	if name == "" {
		Fatal("method missing prefix:", fname)
	}
	meths, ok := curMethods[prefix]
	if !ok {
		Fatal("unknown method prefix:", prefix)
	}
	if _, ok := meths[name]; ok {
		Fatal("duplicate method name:", fname)
	}
	return name, meths
}

func funcName(f any) string {
	s := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	s = str.AfterFirst(s, "builtin.")
	switch s[len(s)-1] {
	case 'Q':
		s = s[:len(s)-1] + "?"
	case 'X':
		s = s[:len(s)-1] + "!"
	}
	return s
}

// params builds a ParamSpec from a string like (a, b) or (@args)
func params(s string) ParamSpec {
	s = strings.ReplaceAll(s, "nil", "'nil'")
	fn := compile.Constant("function " + s + " {}").(*SuFunc)
	for i := range int(fn.ParamSpec.Ndefaults) {
		if fn.Values[i].Equal(SuStr("nil")) {
			fn.Values[i] = nil
		}
	}
	return fn.ParamSpec
}

type staticClass[E any] struct {
	ValueBase[E]
}

func (*staticClass[E]) SetConcurrent() {
	// read-only so ok
}

func methodList(m map[string]Value) Value {
	return SuObjectOfStrs(slices.AppendSeq(make([]string, 0, len(m)), maps.Keys(m)))
}
