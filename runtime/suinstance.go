package runtime

import "github.com/apmckinlay/gsuneido/util/dnum"

// SuInstance is an instance of an SuInstance
type SuInstance struct {
	MemBase
	class *SuClass
}

func NewInstance(class *SuClass) *SuInstance {
	return &SuInstance{NewMemBase(), class}
}

var _ Value = (*SuInstance)(nil)

func (ob *SuInstance) String() string {
	if ob.class.Name != "" {
		return ob.class.Name + "()"
	}
	return "/* instance */"
}

func (*SuInstance) TypeName() string {
	return "Instance"
}

func (*SuInstance) ToInt() int {
	panic("cannot convert instance to integer")
}

func (*SuInstance) ToDnum() dnum.Dnum {
	panic("cannot convert instance to number")
}

func (*SuInstance) ToStr() string {
	panic("cannot convert instance to string")
}

func (ob *SuInstance) Get(m Value) Value {
	if m.TypeName() != "String" {
		return nil
	}
	if x, ok := ob.Data[m.ToStr()]; ok {
		return x
	}
	return ob.class.Get(m)
}

func (ob *SuInstance) Put(m Value, v Value) {
	ob.Data[m.ToStr()] = v
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

// func (ob *SuInstance) parent() *SuClass {
// 	return ob.class
// }

// InstanceMethods is initialized by the builtin package
var InstanceMethods Methods

func (ob *SuInstance) Lookup(method string) Value {
	if f, ok := InstanceMethods[method]; ok {
		return f
	}
	return ob.lookup(method)
}

func (ob *SuInstance) Call(t *Thread, as *ArgSpec) Value {
	if f := ob.lookup("Call"); f != nil {
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
