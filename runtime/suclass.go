package runtime

import (
	"sort"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

// SuClass is a user defined (Suneido language) class
type SuClass struct {
	MemBase
	Name string
	Base Global
}

var _ Value = (*SuClass)(nil)

func (c *SuClass) String() string {
	s := ""
	if c.Name != "" {
		s = c.Name + " "
	}
	s += "/* class"
	if c.Base != 0 {
		s += " : " + GlobalName(c.Base)
	}
	s += " */"
	return s
}

func (c *SuClass) Show() string {
	s := ""
	if c.Base == 0 {
		s = "class"
	} else {
		s += GlobalName(c.Base)
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

func (*SuClass) Compare(Value) int {
	panic("class compare not implemented")
}

func (c *SuClass) parent() *SuClass {
	if c.Base == 0 {
		return nil
	}
	base := GetGlobal(c.Base)
	if baseClass, ok := base.(*SuClass); ok {
		return baseClass
	}
	panic("base must be class")
}

// ClassMethods is initialized by the builtin package
var ClassMethods Methods

func (c *SuClass) Lookup(method string) Value {
	if f, ok := ClassMethods[method]; ok {
		return f
	}
	if f := c.lookup(method); f != nil {
		return f
	}
	if method == "New" {
		return &Builtin0{func() Value { return nil },
			BuiltinParams{ParamSpec: ParamSpec0}}
	}
	return nil
}

func (c *SuClass) Call(t *Thread, as *ArgSpec) Value {
	if f := c.Data["CallClass"]; f != nil {
		t.this = c
		return f.Call(t, as)
	}
	// default for calling a class is to create an instance
	return c.New(t, as)
}

func (c *SuClass) New(t *Thread, as *ArgSpec) Value {
	ob := NewInstance(c)
	nu := c.Lookup("New")
	t.this = ob
	nu.Call(t, as)
	return ob
}

var _ Named = &SuClass{}

func (c *SuClass) SetName(name string) {
	c.Name = name
}

func (c *SuClass) GetName() string {
	return c.Name
}

// finder applies fn to ob and all its parents
// stopping if fn returns something other than nil, and returning that value
func (c *SuClass) finder(fn func(*MemBase) Value) Value {
	for i := 0; i < 100; i++ {
		if x := fn(&c.MemBase); x != nil {
			return x
		}
		c = c.parent()
	}
	panic("too many levels of derivation (>100)")
}
