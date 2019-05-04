package runtime

import "github.com/apmckinlay/gsuneido/runtime/types"

// SuInstance is an instance of an SuInstance
type SuInstance struct {
	MemBase
	class *SuClass
}

func NewInstance(class *SuClass) *SuInstance {
	return &SuInstance{NewMemBase(), class}
}

func (ob *SuInstance) Base() *SuClass {
	return ob.class
}

// ToString is used by Cat, Display, and Print
// to handle user defined ToString
func (ob *SuInstance) ToString(t *Thread) string {
	if f := ob.class.get2(t, "ToString"); f != nil {
		t.this = ob
		x := f.Call(t, ArgSpec0)
		if x != nil {
			if s, ok := x.IfStr(); ok {
				return s
			}
		}
		panic("ToString should return a string")
	}
	return ob.String()
}

func (ob *SuInstance) Copy() *SuInstance {
	return &SuInstance{ob.MemBase.Copy(), ob.class}
}

// Value interface --------------------------------------------------

var _ Value = (*SuInstance)(nil)

func (ob *SuInstance) String() string {
	if ob.class.Name != "" {
		return ob.class.Name + "()"
	}
	return "/* instance */"
}

func (*SuInstance) Type() types.Type {
	return types.Instance
}

func (ob *SuInstance) Get(t *Thread, m Value) Value {
	if m.Type() != types.String {
		return nil
	}
	ms := ToStr(m)
	if x, ok := ob.Data[ms]; ok {
		return x
	}
	x := ob.class.get1(t, ob, ms)
	if m, ok := x.(*SuMethod); ok {
		m.this = ob // fix 'this' to be instance, not method class
	}
	return x
}

func (ob *SuInstance) Put(_ *Thread, m Value, v Value) {
	ob.Data[ToStr(m)] = v
}

func (*SuInstance) RangeTo(int, int) Value {
	panic("instance does not support range")
}

func (*SuInstance) RangeLen(int, int) Value {
	panic("instance does not support range")
}

func (*SuInstance) Hash() uint32 {
	panic("instance hash not implemented") //TODO
}

func (*SuInstance) Hash2() uint32 {
	panic("instance hash not implemented")
}

// Equal returns true if two instances have the same class and data
func (ob *SuInstance) Equal(other interface{}) bool {
	o2, ok := other.(*SuInstance)
	if !ok {
		return false
	}
	if ob == o2 {
		return true
	}
	var stack [maxpairs]pair
	return siEqual(ob, o2, stack[:0])
}

func siEqual(ob, o2 *SuInstance, inProgress pairs) bool {
	if ob.class != o2.class || len(ob.Data) != len(o2.Data) {
		return false
	}
	if inProgress.contains(ob, o2) {
		return true
	}
	inProgress.push(ob, o2)
	for k, x := range ob.Data {
		if y, ok := o2.Data[k]; !ok || !equals3(x, y, inProgress) {
			return false
		}
	}
	return true
}

func (*SuInstance) Compare(Value) int {
	panic("instance compare not implemented")
}

// InstanceMethods is initialized by the builtin package
var InstanceMethods Methods

func (ob *SuInstance) Lookup(t *Thread, method string) Callable {
	if method == "*new*" {
		panic("can't create instance of instance")
	}
	if f, ok := InstanceMethods[method]; ok {
		return f
	}
	if f, ok := BaseMethods[method]; ok {
		return f
	}
	if x := ob.class.get2(t, "Default"); x != nil {
		return &defaultAdapter{x, method}
	}
	return ob.class.get2(t, method)
}

func (ob *SuInstance) Call(t *Thread, as *ArgSpec) Value {
	if f := ob.class.get2(t, "Call"); f != nil {
		t.this = ob
		return f.Call(t, as)
	}
	panic("method not found: Call")
}

// Finder implements Findable
func (ob *SuInstance) Finder(t *Thread, fn func(Value, *MemBase) Value) Value {
	if x := fn(ob, &ob.MemBase); x != nil {
		return x
	}
	return ob.class.Finder(t, fn)
}

var _ Findable = (*SuInstance)(nil)

func (ob *SuInstance) Delete(key Value) {
	m := IfStr(key)
	delete(ob.Data, m)
}

func (ob *SuInstance) Clear() {
	ob.Data = map[string]Value{}
}
