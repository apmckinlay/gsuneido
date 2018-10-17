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

type Method struct {
	Name string
	ParamSpec
	Fn func(t *Thread, self Value, args ...Value) Value
}

type Methods = []*Method

var _ Callable = (*Method)(nil)

func (m *Method) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&m.ParamSpec, as)
	return m.Fn(t, t.stack[t.sp - as.Nargs() - 1], args...)
}

func lookupMethod(methods Methods, method string) Callable {
	for _,m := range methods {
		if m.Name == method {
			return m
		}
	}
	return nil
}
