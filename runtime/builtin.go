package runtime

type Builtin struct {
	Fn func(t *Thread, args ...Value) Value
	ParamSpec
}

var _ Value = (*Builtin)(nil)

func (b Builtin) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(t, args...)
}

func (*Builtin) TypeName() string {
	return "BuiltinFunction"
}
