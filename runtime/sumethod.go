package runtime

// SuMethod is a bound method originating from an SuClass or SuInstance
// when called, it sets 'this' to the origin
type SuMethod struct {
	SuFunc
	this Value
}
var _ Value = (*SuMethod)(nil) // verify *SuFunc satisfies Value

func (m *SuMethod) Call(t *Thread, as *ArgSpec) Value {
	t.this = m.this
	return m.SuFunc.Call(t, as)
}
