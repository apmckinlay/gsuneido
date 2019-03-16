package runtime

// SuMethod is a bound method originating from an SuClass or SuInstance
// when called, it sets 'this' to the origin
type SuMethod struct {
	fn   Value
	this Value
	CantConvert
}

func (m *SuMethod) GetFn() Value {
	return m.fn
}

// Value interface --------------------------------------------------

var _ Value = (*SuMethod)(nil) // verify *SuMethod satisfies Value

func (m *SuMethod) Call(t *Thread, as *ArgSpec) Value {
	t.this = m.this
	return m.fn.Call(t, as)
}

// Lookup is used for .Params or .Disasm
func (m *SuMethod) Lookup(method string) Value {
	if f := m.fn.Lookup(method); f != nil {
		return &SuMethod{fn: f, this: m.fn}
	}
	return nil
}

func (*SuMethod) TypeName() string {
	return "Method"
}

func (m *SuMethod) String() string {
	return m.fn.String()
}

// Equal returns true if two methods have the same fn and this
func (m *SuMethod) Equal(other interface{}) bool {
	m2, ok := other.(*SuMethod)
	if !ok {
		return false
	}
	if m == m2 {
		return true
	}
	return m.fn == m2.fn && m.this == m2.this
}

func (*SuMethod) Get(*Thread, Value) Value {
	panic("method does not support get")
}

func (*SuMethod) Put(Value, Value) {
	panic("method does not support put")
}

func (*SuMethod) RangeTo(int, int) Value {
	panic("method does not support range")
}

func (*SuMethod) RangeLen(int, int) Value {
	panic("method does not support range")
}

func (*SuMethod) Hash() uint32 {
	panic("method hash not implemented")
}

func (*SuMethod) Hash2() uint32 {
	panic("method hash not implemented")
}

func (*SuMethod) Compare(Value) int {
	panic("method compare not implemented")
}

func (*SuMethod) Order() Ord {
	return OrdOther
}

// Named interface --------------------------------------------------

var _ Named = (*SuMethod)(nil)

func (*SuMethod) SetName(string) {
	panic("SuMethod SetName shouldn't be called")
}

func (m *SuMethod) GetName() string {
	if n, ok := m.fn.(Named); ok {
		return n.GetName()
	}
	return ""
}
