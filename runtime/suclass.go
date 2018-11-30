package runtime

import (
	"sort"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

// SuClass is a user defined (Suneido language) class
type SuClass struct {
	Name string
	Base string
	Data map[string]Value // or SuStr instead of string ???
}

var _ Value = (*SuClass)(nil)

func (c *SuClass) String() string {
	s := ""
	if c.Name != "" {
		s = c.Name + " "
	}
	s += "/* class"
	if c.Base != "" {
		s += " : " + c.Base
	}
	s += " */"
	return s
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

func (c *SuClass) Get(m Value) Value {
	if s, ok := m.(SuStr); ok {
		return c.Data[string(s)]
	}
	return nil
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

func (c *SuClass) Lookup(method string) Value {
	if f := c.lookup(method); f != nil {
		return f
	}
	return ClassMethods[method]
}

func (c *SuClass) Call(t *Thread, as *ArgSpec) Value {
	if f := c.lookup("CallClass"); f != nil {
		t.this = c
		return f.Call(t, as)
	}
	panic("CallClass not found")
}
func (c *SuClass) lookup(method string) Value {
	if x, ok := c.Data[method]; ok {
		return x
	}
	return nil // could make dummy Value with Call's doing panic
}

var _ Named = &SuClass{}

func (c *SuClass) SetName(name string) {
	c.Name = name
}

func (c *SuClass) GetName() string {
	return c.Name
}
