package runtime

import (
	"sort"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

// SuClass is a user defined (Suneido language) class
type SuClass struct {
	Base string
	Data map[string]Value // or SuStr instead of string ???
}

var _ Value = &SuClass{}

func (c *SuClass) String() string {
	s := "/* class"
	if c.Base != "" {
		s += " : " + c.Base
	}
	return s + " */"
}

func (c *SuClass) Show() string {
	s := c.Base
	if s == "" {
		s = "class"
	}
	s += "{"
	sep := ""
	mems := []string{}
	for k := range c.Data {
		mems = append(mems, k)
	}
	sort.Sort(sort.StringSlice(mems))
	for _, k := range mems {
		s += sep + k
		v := c.Data[k]
		if f, ok := v.(*SuFunc); ok {
			s += f.Params()
		} else {
			s += ": " + v.String()
		}
		sep = "; "
	}
	s += "}"
	return s
}

func (*SuClass) TypeName() string {
	return "Class"
}

func (*SuClass) ToInt() int {
	panic("cannot convert class to integer")
}

func (*SuClass) ToDnum() dnum.Dnum {
	panic("cannot convert class to number")
}

func (*SuClass) ToStr() string {
	panic("cannot convert class to string")
}

func (*SuClass) Get(Value) Value {
	panic("class get not implemented") //TODO
}

func (*SuClass) Put(Value, Value) {
	panic("class does not support put")
}

func (SuClass) RangeTo(int, int) Value {
	panic("class does not support range")
}

func (SuClass) RangeLen(int, int) Value {
	panic("class does not support range")
}

func (*SuClass) Hash() uint32 {
	panic("class hash not implemented") //TODO
}

func (*SuClass) Hash2() uint32 {
	panic("class hash not implemented")
}

func (c *SuClass) Equal(other interface{}) bool {
	if c2, ok := other.(*SuClass); ok {
		return c == c2
	}
	return false
}

func (*SuClass) Order() Ord {
	return OrdOther
}

func (*SuClass) Compare(Value) int {
	panic("class compare not implemented")
}

// ClassMethods is initialized by the builtin package
var ClassMethods Methods

func (*SuClass) Lookup(method string) Callable {
	return ClassMethods[method]
}

func (*SuClass) Call0(_ *Thread) Value {
	panic("call class not implemented") //TODO
}
func (*SuClass) Call1(_ *Thread, _ Value) Value {
	panic("call class not implemented") //TODO
}
func (*SuClass) Call2(_ *Thread, _, _ Value) Value {
	panic("call class not implemented") //TODO
}
func (*SuClass) Call3(_ *Thread, _, _, _ Value) Value {
	panic("call class not implemented") //TODO
}
func (*SuClass) Call4(_ *Thread, _, _, _, _ Value) Value {
	panic("call class not implemented") //TODO
}
func (*SuClass) Call(*Thread, *ArgSpec) Value {
	panic("call class not implemented") //TODO
}
