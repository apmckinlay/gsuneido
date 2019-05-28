package runtime

import "github.com/apmckinlay/gsuneido/runtime/types"

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

// methodForce ignores the "this" passed to Call.
// It is used by Lookup to ensure Params and Disasm get the right "this".
type methodForce struct {
	SuMethod
}

func (m *methodForce) Call(t *Thread, _ Value, as *ArgSpec) Value {
	return m.fn.Call(t, m.this, as)
}

// Value interface --------------------------------------------------

var _ Value = (*SuMethod)(nil)

func (m *SuMethod) Call(t *Thread, this Value, as *ArgSpec) Value {
	if this == nil {
		this = m.this
	}
	return m.fn.Call(t, this, as)
}

// Lookup is used for .Params or .Disasm
func (m *SuMethod) Lookup(t *Thread, method string) Callable {
	if f := m.fn.Lookup(t, method); f != nil {
		return &methodForce{SuMethod{fn: f.(Value), this: m.fn}}
	}
	return nil
}

func (*SuMethod) Type() types.Type {
	return types.Method
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

func (*SuMethod) Put(*Thread, Value, Value) {
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

// Named interface --------------------------------------------------

var _ Named = (*SuMethod)(nil)

func (m *SuMethod) GetName() string {
	if n, ok := m.fn.(Named); ok {
		return n.GetName()
	}
	return ""
}
