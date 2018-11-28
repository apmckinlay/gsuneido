package runtime

// Methods is a map of method name strings to Values
type Methods = map[string]Value

// Builtin is a Value for a builtin function
type Builtin struct {
	Fn func(t *Thread, args ...Value) Value
	ParamSpec
}

var _ Value = (*Builtin)(nil)

func (*Builtin) TypeName() string {
	return "BuiltinFunction"
}

func (b *Builtin) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t, args...)
}

// Builtin0 is a Value for a builtin function with no arguments
type Builtin0 struct {
	Fn func() Value
	ParamSpec
}

var _ Value = (*Builtin0)(nil)

func (*Builtin0) TypeName() string {
	return "BuiltinFunction"
}

func (b *Builtin0) Call(t *Thread, as *ArgSpec) Value {
	t.Args(&b.ParamSpec, as)
	return b.Fn()
}

// Builtin1 is a Value for a builtin function with one argument
type Builtin1 struct {
	Fn func(a1 Value) Value
	ParamSpec
}

var _ Value = (*Builtin1)(nil)

func (*Builtin1) TypeName() string {
	return "BuiltinFunction"
}

func (b *Builtin1) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(args[0])
}

// Builtin2 is a Value for a builtin function with two arguments
type Builtin2 struct {
	Fn func(a1, a2 Value) Value
	ParamSpec
}

var _ Value = (*Builtin2)(nil)

func (*Builtin2) TypeName() string {
	return "BuiltinFunction"
}

func (b *Builtin2) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1])
}

// Builtin3 is a Value for a builtin function with three arguments
type Builtin3 struct {
	Fn func(a1, a2, a3 Value) Value
	ParamSpec
}

var _ Value = (*Builtin3)(nil)

func (*Builtin3) TypeName() string {
	return "BuiltinFunction"
}

func (b *Builtin3) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2])
}

// ------------------------------------------------------------------

// Method is a Value for a builtin method
type Method struct {
	Fn func(t *Thread, this Value, args ...Value) Value
	ParamSpec
}

var _ Value = (*Method)(nil)

func (*Method) TypeName() string {
	return "BuiltinFunction"
}

func (b *Method) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t, t.this, args...)
}

// Method0 is a Value for a builtin method with no arguments
type Method0 struct {
	Builtin1
}

func (b *Method0) Call(t *Thread, as *ArgSpec) Value {
	t.Args(&b.ParamSpec, as)
	return b.Fn(t.this)
}

// Method1 is a Value for a builtin method with one argument
type Method1 struct {
	Builtin2
}

func (b *Method1) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t.this, args[0])
}

// Method2 is a Value for a builtin method with two arguments
type Method2 struct {
	Builtin3
}

func (b *Method2) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t.this, args[0], args[1])
}

// RawMethod is a Value for a builtin function with no massage
type RawMethod struct {
	Fn func(t *Thread, as *ArgSpec, this Value, args ...Value) Value
	ParamSpec
}

var _ Value = (*RawMethod)(nil)

func (*RawMethod) TypeName() string {
	return "BuiltinFunction"
}

func (b *RawMethod) Call(t *Thread, as *ArgSpec) Value {
	base := t.sp - as.Nargs()
	args := t.stack[base:]
	return b.Fn(t, as, t.this, args...)
}
