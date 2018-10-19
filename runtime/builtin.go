package runtime

// Builtin is a Callable Value for a builtin function with massaged arguments
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
func (b *Builtin) Call0(t *Thread) Value {
	return b.Call(t, ArgSpec0)
}
func (b *Builtin) Call1(t *Thread, _ Value) Value {
	return b.Call(t, ArgSpec1)
}
func (b *Builtin) Call2(t *Thread, _, _ Value) Value {
	return b.Call(t, ArgSpec2)
}
func (b *Builtin) Call3(t *Thread, _, _, _ Value) Value {
	return b.Call(t, ArgSpec3)
}
func (b *Builtin) Call4(t *Thread, _, _, _, _ Value) Value {
	return b.Call(t, ArgSpec4)
}

// Builtin0 is a Callable Value for a builtin function with no arguments
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
func (b *Builtin0) Call0(*Thread) Value {
	return b.Fn() // fast path
}

// Builtin1 is a Callable Value for a builtin function with one argument
type Builtin1 struct {
	Fn func(a Value) Value
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
func (b *Builtin1) Call0(*Thread) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin1) Call1(_ *Thread, a Value) Value {
	return b.Fn(a) // fast path
}

// Builtin2 is a Callable Value for a builtin function with two arguments
type Builtin2 struct {
	Fn func(a, b Value) Value
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
func (b *Builtin2) Call0(*Thread) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin2) Call1(*Thread, Value) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin2) Call2(_ *Thread, a1, a2 Value) Value {
	return b.Fn(a1, a2) // fast path
}

// Builtin3 is a Callable Value for a builtin function with three arguments
type Builtin3 struct {
	Fn func(a, b, c Value) Value
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
func (b *Builtin3) Call0(*Thread) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin3) Call1(*Thread, Value) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin3) Call2(_ *Thread, _, _ Value) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin3) Call3(_ *Thread, a1, a2, a3 Value) Value {
	return b.Fn(a1, a2, a3) // fast path
}

// Builtin4 is a Callable Value for a builtin function with three arguments
type Builtin4 struct {
	Fn func(a, b, c, d Value) Value
	ParamSpec
}

var _ Value = (*Builtin4)(nil)

func (*Builtin4) TypeName() string {
	return "BuiltinFunction"
}

func (b *Builtin4) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&b.ParamSpec, as)
	return b.Fn(args[0], args[1], args[2], args[3])
}
func (b *Builtin4) Call0(*Thread) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin4) Call1(*Thread, Value) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin4) Call2(_ *Thread, _, _ Value) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin4) Call3(_ *Thread, _, _, _ Value) Value {
	// TODO use default if available
	panic("not enough arguments")
}
func (b *Builtin4) Call4(_ *Thread, a1, a2, a3, a4 Value) Value {
	return b.Fn(a1, a2, a3, a4) // fast path
}

// Methods is a map of method name strings to Callables
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
func (m *Method) Call0(*Thread) Value {
	panic("shouldn't get here")
}
func (m *Method) Call1(t *Thread, _ Value) Value {
	return m.Call(t, ArgSpec0)
}
func (m *Method) Call2(t *Thread, _, _ Value) Value {
	return m.Call(t, ArgSpec1)
}
func (m *Method) Call3(t *Thread, _, _, _ Value) Value {
	return m.Call(t, ArgSpec2)
}
func (m *Method) Call4(t *Thread, _, _, _, _ Value) Value {
	return m.Call(t, ArgSpec3)
}

// Method0 is a Callable for a builtin method with no arguments
type Method0 struct {
	ParamSpec
	Fn func(self Value) Value
}

var _ Callable = (*Method0)(nil)

func (m *Method0) Call(t *Thread, as *ArgSpec) Value {
	self := t.stack[t.sp-as.Nargs()-1]
	t.Args(&m.ParamSpec, as)
	return m.Fn(self)
}
func (m *Method0) Call0(*Thread) Value {
	panic("shouldn't get here")
}
func (m *Method0) Call1(_ *Thread, self Value) Value {
	return m.Fn(self) // fast path
}

// Method1 is a Callable for a builtin method with one argument
type Method1 struct {
	ParamSpec
	Fn func(self, a1 Value) Value
}

var _ Callable = (*Method1)(nil)

func (m *Method1) Call(t *Thread, as *ArgSpec) Value {
	self := t.stack[t.sp-as.Nargs()-1]
	args := t.Args(&m.ParamSpec, as)
	return m.Fn(self, args[0])
}
func (m *Method1) Call0(*Thread) Value {
	panic("shouldn't get here")
}
func (m *Method1) Call1(_ *Thread, _ Value) Value {
	panic("not enough arguments")
}
func (m *Method1) Call2(_ *Thread, self, a1 Value) Value {
	return m.Fn(self, a1) // fast path
}

// Method2 is a Callable for a builtin method with two arguments
type Method2 struct {
	ParamSpec
	Fn func(self, a1, a2 Value) Value
}

var _ Callable = (*Method2)(nil)

func (m *Method2) Call(t *Thread, as *ArgSpec) Value {
	self := t.stack[t.sp-as.Nargs()-1]
	args := t.Args(&m.ParamSpec, as)
	return m.Fn(self, args[0], args[1])
}
func (m *Method2) Call0(*Thread) Value {
	panic("shouldn't get here")
}
func (m *Method2) Call1(_ *Thread, _ Value) Value {
	panic("not enough arguments")
}
func (m *Method2) Call2(_ *Thread, _, _ Value) Value {
	panic("not enough arguments")
}
func (m *Method2) Call3(_ *Thread, self, a1, a2 Value) Value {
	return m.Fn(self, a1, a2) // fast path
}

// Method3 is a Callable for a builtin method with two arguments
type Method3 struct {
	ParamSpec
	Fn func(self, a1, a2, a3 Value) Value
}

var _ Callable = (*Method3)(nil)

func (m *Method3) Call(t *Thread, as *ArgSpec) Value {
	self := t.stack[t.sp-as.Nargs()-1]
	args := t.Args(&m.ParamSpec, as)
	return m.Fn(self, args[0], args[1], args[2])
}
func (m *Method3) Call0(*Thread) Value {
	panic("shouldn't get here")
}
func (m *Method3) Call1(_ *Thread, _ Value) Value {
	panic("not enough arguments")
}
func (m *Method3) Call2(_ *Thread, _, _ Value) Value {
	panic("not enough arguments")
}
func (m *Method3) Call3(_ *Thread, _, _, _ Value) Value {
	panic("not enough arguments")
}
func (m *Method3) Call4(_ *Thread, self, a1, a2, a3 Value) Value {
	return m.Fn(self, a1, a2, a3) // fast path
}

// RawMethod is a Callable for a builtin method with raw arguments
type RawMethod struct {
	ParamSpec // ???
	Fn        func(t *Thread, self Value, as *ArgSpec, args ...Value) Value
}

var _ Callable = (*RawMethod)(nil)

func (m *RawMethod) Call(t *Thread, as *ArgSpec) Value {
	self := t.stack[t.sp-as.Nargs()-1]
	args := t.stack[t.sp-as.Nargs():]
	return m.Fn(t, self, as, args...)
}
func (m *RawMethod) Call0(*Thread) Value {
	panic("shouldn't get here")
}
func (m *RawMethod) Call1(t *Thread, _ Value) Value {
	return m.Call(t, ArgSpec0)
}
func (m *RawMethod) Call2(t *Thread, _, _ Value) Value {
	return m.Call(t, ArgSpec1)
}
func (m *RawMethod) Call3(t *Thread, _, _, _ Value) Value {
	return m.Call(t, ArgSpec2)
}
func (m *RawMethod) Call4(t *Thread, _, _, _, _ Value) Value {
	return m.Call(t, ArgSpec3)
}
