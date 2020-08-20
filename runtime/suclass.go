// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"sort"
	"strings"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// SuClass is a user defined (Suneido language) class
type SuClass struct {
	MemBase
	Lib          string
	Name         string
	Base         Gnum
	parentsCache atomic.Value // used by SuInstance getParents
	noGetter     bool
}

// NOTE: the parents argument on some SuClass methods is used by SuInstance

func (c *SuClass) SetParents(parents []*SuClass) {
	c.parentsCache.Store(parents)
}

func (c *SuClass) GetParents() []*SuClass {
	x, _ := c.parentsCache.Load().([]*SuClass) // allow nil
	return x
}

var _ Value = (*SuClass)(nil)

func (c *SuClass) String() string {
	s := ""
	if !anonymous(c.Name) {
		s = c.Name + " "
	}
	s += "/* " + str.Opt(c.Lib, " ") + "class"
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
	return c.get1(t, c, m, nil)
}

func (c *SuClass) get1(t *Thread, this Value, m Value, parents []*SuClass) Value {
	ms := AsStr(m) //TODO not quite right - allows class { "4": }[4]
	val := c.get2(t, ms, parents)
	if val != nil {
		if _, ok := val.(*SuFunc); ok {
			return &SuMethod{fn: val, this: c}
		}
		return val
	}
	if !c.noGetter {
		if getter := c.get2(t, "Getter_", parents); getter != nil {
			return t.CallThis(getter, this, m)
		}
		c.noGetter = true
	}
	getterName := "Getter_" + ms
	if getter := c.get2(t, getterName, parents); getter != nil {
		return t.CallThis(getter, this)
	}
	return nil
}

// get2 looks for m in the current Data and then the parent's
func (c *SuClass) get2(t *Thread, m string, parents []*SuClass) Value {
	if parents == nil {
		for n := 0; ; n++ {
			if x, ok := c.Data[m]; ok {
				assert.That(x != nil)
				return x
			}
			if c = c.Parent(t); c == nil {
				break
			}
			if n > inheritanceLimit {
				panic("too many levels of inheritance")
			}
		}
	} else {
		assert.That(parents[0] == c)
		for _, p := range parents {
			if x, ok := p.Data[m]; ok {
				assert.That(x != nil)
				return x
			}
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
	c2, ok := other.(*SuClass)
	return ok && c == c2
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
	return c.lookup(t, method, nil)
}

func (c *SuClass) lookup(t *Thread, method string, parents []*SuClass) Callable {
	if f, ok := ClassMethods[method]; ok {
		return f
	}
	if f, ok := BaseMethods[method]; ok {
		return f
	}
	if x := c.get2(t, method, parents); x != nil {
		return x
	}
	if method == "New" {
		return DefaultNewMethod
	}
	if x := UserDef(t, gnObjects, method); x != nil {
		return x
	}
	//TODO explicit CallClass doesn't go to Default in cSuneido or jSuneido
	if x := c.get2(t, "Default", parents); x != nil {
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
	if as.Each >= EACH0 {
		args := ToContainer(t.Pop()).Slice(int(as.Each) - 1)
		args.Insert(0, method)
		t.Push(args)
		as = &ArgSpecEach0
	} else if as.Nargs == 0 {
		t.Push(method)
		as = &ArgSpec1
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

func (c *SuClass) Call(t *Thread, this Value, as *ArgSpec) Value {
	if this == nil {
		this = c
	}
	if f := c.get2(t, "CallClass", nil); f != nil {
		return f.Call(t, this, as)
	}
	// default for calling a class is to create an instance
	return c.New(t, as)
}

func (c *SuClass) New(t *Thread, as *ArgSpec) Value {
	ob := NewInstance(t, c)
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
	for i := 0; i < inheritanceLimit; i++ {
		if x := fn(c, &c.MemBase); x != nil {
			return x
		}
		c = c.Parent(t)
		if c == nil {
			return nil
		}
	}
	panic("too many levels of inheritance")
}

var inheritanceLimit = 100

var _ Findable = (*SuClass)(nil)
