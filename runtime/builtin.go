package runtime

// Builtin is a Callable Value for a builtin function with massaged arguments
type Builtin struct {
	Fn func(t *Thread, args ...Value) Value
	ParamSpec
}

var _ Value = (*Builtin)(nil)

func (b *Builtin) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t, args...)
}

func (*Builtin) TypeName() string {
	return "BuiltinFunction"
}

type Methods = map[string]Callable

// Method is a Callable for a builtin method with massaged arguments
type Method struct {
	ParamSpec
	Fn func(t *Thread, self Value, args ...Value) Value
}

var _ Callable = (*Method)(nil)

func (m *Method) Call(t *Thread, as *ArgSpec) Value {
	self := t.stack[t.sp-as.Nargs()-1]
	args := t.Args(&m.ParamSpec, as)
	return m.Fn(t, self, args...)
}

// RawMethod is a Callable for a builtin method with raw arguments
type RawMethod struct {
	ParamSpec // ???
	Fn        func(t *Thread, self Value, as *ArgSpec, args ...Value) Value
}

var _ Callable = (*Method)(nil)

func (m *RawMethod) Call(t *Thread, as *ArgSpec) Value {
	self := t.stack[t.sp-as.Nargs()-1]
	args := t.stack[t.sp-as.Nargs():]
	return m.Fn(t, self, as, args...)
}
