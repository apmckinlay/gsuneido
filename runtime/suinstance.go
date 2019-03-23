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

// Value interface --------------------------------------------------

var _ Value = (*SuInstance)(nil)

func (ob *SuInstance) String() string {
	if ob.class.Name != "" {
		return ob.class.Name + "()"
	}
	return "/* instance */"
}

// Display is used by Display and Print
// to handle user defined ToString
func (ob *SuInstance) Display(t *Thread) string {
	if f := ob.class.get2("ToString"); f != nil {
		return ToStr(f.Call(t, ArgSpec0))
	}
	return ob.String()
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
	x := ob.class.get1(t, ms)
	if m, ok := x.(*SuMethod); ok {
		m.this = ob // fix 'this' to be instance, not method class
	}
	return x
}

func (ob *SuInstance) Put(m Value, v Value) {
	ob.Data[ToStr(m)] = v
}

func (SuInstance) RangeTo(int, int) Value {
	panic("instance does not support range")
}

func (SuInstance) RangeLen(int, int) Value {
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

func (ob *SuInstance) Lookup(method string) Value {
	if method == "*new*" {
		panic("can't create instance of instance")
	}
	if f, ok := InstanceMethods[method]; ok {
		return f
	}
	return ob.class.get2(method)
}

func (ob *SuInstance) Call(t *Thread, as *ArgSpec) Value {
	if f := ob.class.get2("Call"); f != nil {
		t.this = ob
		return f.Call(t, as)
	}
	panic("method not found: Call")
}

// finder applies fn to ob and all its parents
// stopping if fn returns something other than nil, and returning that value
func (ob *SuInstance) finder(fn func(*MemBase) Value) Value {
	if x := fn(&ob.MemBase); x != nil {
		return x
	}
	return ob.class.finder(fn)
}
