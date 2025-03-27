// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import "github.com/apmckinlay/gsuneido/core/types"

// Methods is a map of method name strings to Values
type Methods = map[string]Value

type BuiltinParams struct {
	ParamSpec
}

func (ps *BuiltinParams) String() string {
	s := "/* builtin function */"
	if ps.Name == "" {
		return s
	}
	return ps.Name + " " + s
}

func (*BuiltinParams) Type() types.Type {
	return types.BuiltinFunction
}

// SuBuiltin is a built in function
type SuBuiltin struct {
	Fn func(th *Thread, args []Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin)(nil)

func (b *SuBuiltin) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(th, args)
}

// SuBuiltin0 is a builtin function with no arguments
type SuBuiltin0 struct {
	Fn func() Value
	BuiltinParams
}

var _ Value = (*SuBuiltin0)(nil)

func (b *SuBuiltin0) Call(th *Thread, _ Value, as *ArgSpec) Value {
	th.Args(&b.ParamSpec, as)
	return b.Fn()
}

// SuBuiltin1 is a builtin function with one argument
type SuBuiltin1 struct {
	Fn func(a1 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin1)(nil)

func (b *SuBuiltin1) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(args[0])
}

// SuBuiltin2 is a builtin function with two arguments
type SuBuiltin2 struct {
	Fn func(a1, a2 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin2)(nil)

func (b *SuBuiltin2) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1])
}

// SuBuiltin3 is a builtin function with three arguments
type SuBuiltin3 struct {
	Fn func(a1, a2, a3 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin3)(nil)

func (b *SuBuiltin3) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2])
}

// SuBuiltin4 is a builtin function with four arguments
type SuBuiltin4 struct {
	Fn func(a1, a2, a3, a4 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin4)(nil)

func (b *SuBuiltin4) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2], args[3])
}

// SuBuiltin5 is a builtin function with five arguments
type SuBuiltin5 struct {
	Fn func(a1, a2, a3, a4, a5 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin5)(nil)

func (b *SuBuiltin5) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2], args[3], args[4])
}

// SuBuiltin6 is a builtin function with six arguments
type SuBuiltin6 struct {
	Fn func(a1, a2, a3, a4, a5, a6 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin6)(nil)

func (b *SuBuiltin6) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2], args[3], args[4], args[5])
}

// SuBuiltin7 is a builtin function with seven arguments
var _ Value = (*SuBuiltin7)(nil)

type SuBuiltin7 struct {
	Fn func(a1, a2, a3, a4, a5, a6, a7 Value) Value
	BuiltinParams
}

func (b *SuBuiltin7) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2], args[3], args[4], args[5], args[6])
}

// SuBuiltinRaw is a builtin function with no massage
type SuBuiltinRaw struct {
	Fn func(th *Thread, as *ArgSpec, args []Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltinRaw)(nil)

func (b *SuBuiltinRaw) Call(th *Thread, _ Value, as *ArgSpec) Value {
	base := th.sp - int(as.Nargs)
	args := th.stack[base : base+int(as.Nargs)]
	return b.Fn(th, as, args)
}

// ------------------------------------------------------------------

// SuBuiltinMethod is a builtin method
type SuBuiltinMethod struct {
	Fn func(th *Thread, this Value, args []Value) Value
	BuiltinParams
}

func (b *SuBuiltinMethod) Call(th *Thread, this Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(th, this, args)
}

// SuBuiltinMethod0 is a builtin method with no arguments
type SuBuiltinMethod0 struct {
	SuBuiltin1
}

func (b *SuBuiltinMethod0) Call(th *Thread, this Value, as *ArgSpec) Value {
	th.Args(&b.ParamSpec, as)
	return b.Fn(this)
}

// SuBuiltinMethod1 is a builtin method with one argument
type SuBuiltinMethod1 struct {
	SuBuiltin2
}

func (b *SuBuiltinMethod1) Call(th *Thread, this Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(this, args[0])
}

// SuBuiltinMethod2 is a builtin method with two arguments
type SuBuiltinMethod2 struct {
	SuBuiltin3
}

func (b *SuBuiltinMethod2) Call(th *Thread, this Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(this, args[0], args[1])
}

// SuBuiltinMethod3 is a builtin method with two arguments
type SuBuiltinMethod3 struct {
	SuBuiltin4
}

func (b *SuBuiltinMethod3) Call(th *Thread, this Value, as *ArgSpec) Value {
	args := th.Args(&b.ParamSpec, as)
	return b.Fn(this, args[0], args[1], args[2])
}

// SuBuiltinMethodRaw is a builtin function with no massage
type SuBuiltinMethodRaw struct {
	Fn func(th *Thread, as *ArgSpec, this Value, args []Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltinMethodRaw)(nil)

func (b *SuBuiltinMethodRaw) Call(th *Thread, this Value, as *ArgSpec) Value {
	base := th.sp - int(as.Nargs)
	args := th.stack[base : base+int(as.Nargs)]
	return b.Fn(th, as, this, args)
}

func (b *SuBuiltinMethodRaw) Type() types.Type {
	return types.Method
}
