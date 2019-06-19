package runtime

import (
	"sort"
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

// SuClass is a user defined (Suneido language) class
type SuClass struct {
	MemBase
	Name     string
	Base     Gnum
	noGetter bool
}

var _ Value = (*SuClass)(nil)

func (c *SuClass) String() string {
	s := ""
	if !anonymous(c.Name) {
		s = c.Name + " "
	}
	s += "/* class"
	if c.Base != 0 {
		s += " : " + Global.Name(c.Base)
	}
	s += " */"
	return s
}

func anonymous(s string) bool {
	return s == "" || s == "?" ||
		(strings.HasPrefix(s, "Class") && ascii.IsDigit(s[len(s)-1]))
}

func (c *SuClass) Show() string {
	s := ""
	if c.Base == 0 {
		s = "class"
	} else {
		s += Global.Name(c.Base)
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

func (*SuClass) Type() types.Type {
	return types.Class
}

func (c *SuClass) Get(t *Thread, m Value) Value {
	if m.Type() != types.String {
		return nil
	}
	return c.get1(t, c, AsStr(m))
}

func (c *SuClass) get1(t *Thread, this Value, mem string) Value {
	val := c.get2(t, mem)
	if val != nil {
		if _, ok := val.(*SuFunc); ok {
			return &SuMethod{fn: val, this: c}
		}
		return val
	}
	if !c.noGetter {
		if getter := c.get2(t, "Getter_"); getter != nil {
			return t.CallThis(getter, this, SuStr(mem))
		}
		c.noGetter = true
	}
	getterName := "Getter_" + mem
	if getter := c.get2(t, getterName); getter != nil {
		return getter.Call(t, this, ArgSpec0)
	}
	return nil
}

func (c *SuClass) get2(t *Thread, m string) Value {
	for {
		if x, ok := c.Data[m]; ok {
			return x
		}
		if c = c.Parent(t); c == nil {
			break
		}
	}
	return nil
}

func (*SuClass) Put(*Thread, Value, Value) {
	panic("class does not support put")
}

func (*SuClass) RangeTo(int, int) Value {
	panic("class does not support range")
}

func (*SuClass) RangeLen(int, int) Value {
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

func (c *SuClass) Parent(t *Thread) *SuClass {
	if c.Base <= 0 {
		return nil
	}
	base := Global.Get(t, c.Base)
	if baseClass, ok := base.(*SuClass); ok {
		return baseClass
	}
	panic("base must be class")
}

// BaseMethods is initialized by the builtin package
var BaseMethods Methods

// ClassMethods is initialized by the builtin package
var ClassMethods Methods

var DefaultNewMethod = &SuBuiltin0{func() Value { return nil },
	BuiltinParams{ParamSpec: ParamSpec0}}

func (c *SuClass) Lookup(t *Thread, method string) Callable {
	if f, ok := ClassMethods[method]; ok {
		return f
	}
	if f, ok := BaseMethods[method]; ok {
		return f
	}
	if x := c.get2(t, method); x != nil {
		return x
	}
	if method == "New" {
		return DefaultNewMethod
	}
	if x := UserDef(t, gnObjects, method); x != nil {
		return x
	}
	if x := c.get2(t, "Default"); x != nil {
		return &defaultAdapter{x, method}
	}
	return nil
}

// defaultAdapter wraps a Default method to insert the method argument
type defaultAdapter struct {
	fn     Callable
	method string
}

func (d *defaultAdapter) Call(t *Thread, this Value, as *ArgSpec) Value {
	method := SuStr(d.method)
	if as.Each >= EACH {
		args := ToContainer(t.Pop()).Slice(int(as.Each) - 1)
		args.Insert(0, method)
		t.Push(args)
		as = ArgSpecEach
	} else if as.Nargs == 0 {
		t.Push(method)
		as = ArgSpec1
	} else {
		t.Push(nil) // allow for another value
		base := t.sp - 1 - int(as.Nargs)
		copy(t.stack[base+1:], t.stack[base:]) // shift over
		t.stack[base] = method
		as2 := *as
		as2.Nargs++
		as = &as2
	}
	return d.fn.Call(t, this, as)
}

func (c *SuClass) Call(t *Thread, _ Value, as *ArgSpec) Value {
	if f := c.get2(t, "CallClass"); f != nil {
		return f.Call(t, c, as)
	}
	// default for calling a class is to create an instance
	return c.New(t, as)
}

func (c *SuClass) New(t *Thread, as *ArgSpec) Value {
	ob := NewInstance(c)
	nu := c.Lookup(t, "New")
	nu.Call(t, ob, as)
	return ob
}

var _ Named = &SuClass{}

func (c *SuClass) GetName() string {
	return c.Name
}

// Finder implements Findable
func (c *SuClass) Finder(t *Thread, fn func(v Value, mb *MemBase) Value) Value {
	for i := 0; i < 100; i++ {
		if x := fn(c, &c.MemBase); x != nil {
			return x
		}
		c = c.Parent(t)
		if c == nil {
			return nil
		}
	}
	panic("too many levels of derivation (>100)")
}

var _ Findable = (*SuClass)(nil)
