package runtime

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

// SuBuiltin is a Value for a built in function
type SuBuiltin struct {
	Fn func(t *Thread, args ...Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin)(nil)

func (*SuBuiltin) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltin) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t, args...)
}

// SuBuiltin0 is a Value for a builtin function with no arguments
type SuBuiltin0 struct {
	Fn func() Value
	BuiltinParams
}

var _ Value = (*SuBuiltin0)(nil)

func (*SuBuiltin0) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltin0) Call(t *Thread, as *ArgSpec) Value {
	t.Args(&b.ParamSpec, as)
	return b.Fn()
}

// SuBuiltin1 is a Value for a builtin function with one argument
type SuBuiltin1 struct {
	Fn func(a1 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin1)(nil)

func (*SuBuiltin1) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltin1) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(args[0])
}

// SuBuiltin2 is a Value for a builtin function with two arguments
type SuBuiltin2 struct {
	Fn func(a1, a2 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin2)(nil)

func (*SuBuiltin2) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltin2) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1])
}

// SuBuiltin3 is a Value for a builtin function with three arguments
type SuBuiltin3 struct {
	Fn func(a1, a2, a3 Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltin3)(nil)

func (*SuBuiltin3) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltin3) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2])
}

// SuBuiltinRaw is a Value for a builtin function with no massage
type SuBuiltinRaw struct {
	Fn func(t *Thread, as *ArgSpec, args ...Value) Value
	BuiltinParams
}

var _ Value = (*SuBuiltinRaw)(nil)

func (*SuBuiltinRaw) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltinRaw) Call(t *Thread, as *ArgSpec) Value {
	base := t.sp - int(as.Nargs)
	args := t.stack[base:base + int(as.Nargs)]
	return b.Fn(t, as, args...)
}

// ------------------------------------------------------------------

// SuBuiltinMethod is a Value for a builtin method
type SuBuiltinMethod struct {
	Fn func(t *Thread, this Value, args ...Value) Value
	ParamSpec
}

var _ Value = (*SuBuiltinMethod)(nil)

func (*SuBuiltinMethod) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltinMethod) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t, t.this, args...)
}

// SuBuiltinMethod0 is a Value for a builtin method with no arguments
type SuBuiltinMethod0 struct {
	SuBuiltin1
}

func (b *SuBuiltinMethod0) Call(t *Thread, as *ArgSpec) Value {
	t.Args(&b.ParamSpec, as)
	return b.Fn(t.this)
}

// SuBuiltinMethod1 is a Value for a builtin method with one argument
type SuBuiltinMethod1 struct {
	SuBuiltin2
}

func (b *SuBuiltinMethod1) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t.this, args[0])
}

// SuBuiltinMethod2 is a Value for a builtin method with two arguments
type SuBuiltinMethod2 struct {
	SuBuiltin3
}

func (b *SuBuiltinMethod2) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t.this, args[0], args[1])
}

// SuBuiltinMethodRaw is a Value for a builtin function with no massage
type SuBuiltinMethodRaw struct {
	Fn func(t *Thread, as *ArgSpec, this Value, args ...Value) Value
	ParamSpec
}

var _ Value = (*SuBuiltinMethodRaw)(nil)

func (*SuBuiltinMethodRaw) TypeName() string {
	return "BuiltinFunction"
}

func (b *SuBuiltinMethodRaw) Call(t *Thread, as *ArgSpec) Value {
	base := t.sp - int(as.Nargs)
	args := t.stack[base:base + int(as.Nargs)]
	return b.Fn(t, as, t.this, args...)
}
