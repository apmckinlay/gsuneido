// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"sort"
	"strings"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// SuClass is a user defined (Suneido language) class.
// Classes are read-only so there is no locking.
type SuClass struct {
	ValueBase[*SuClass]
	parentsCache atomic.Pointer[SuClassChain] // used by SuInstance getParents
	Lib          string
	Name         string
	MemBase
	Base Gnum
}

// SetParents caches the parents chain on the class.
func (c *SuClass) SetParents(parents *SuClassChain) {
	c.parentsCache.Store(parents)
}

// GetParents returns the cached parents chain, or nil if not cached.
func (c *SuClass) GetParents() *SuClassChain {
	return c.parentsCache.Load()
}

func (c *SuClass) Class() *SuClass {
	return c
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
	sort.Strings(mems)
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

func (c *SuClass) Get(th *Thread, m Value) Value {
	return c.get1(th, c, m, nil)
}

func (c *SuClass) get1(th *Thread, this Value, m Value, parents []*SuClass) Value {
	ms := AsStr(m) //TODO not quite right - allows class { "4": }[4]
	val := c.get2(th, ms, parents)
	if val != nil {
		if _, ok := val.(*SuFunc); ok {
			return &SuMethod{fn: val, this: this}
		}
		return val
	}
	if getter := c.get2(th, "Getter_", parents); getter != nil {
		return th.CallThis(getter, this, m)
	}
	if getter := c.get2(th, "Getter_"+ms, parents); getter != nil {
		return th.CallThis(getter, this)
	}
	return nil
}

// get2 looks for m in the current Data and then the parent's
func (c *SuClass) get2(th *Thread, m string, parents []*SuClass) Value {
	if parents == nil {
		for n := 0; ; n++ {
			if x, ok := c.Data[m]; ok {
				assert.That(x != nil)
				return x
			}
			if c = c.Parent(th); c == nil {
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

func (c *SuClass) Equal(other any) bool {
	if cc, ok := other.(*SuClassChain); ok {
		return c == cc.Class()
	}
	return c == other
}

func (*SuClass) SetConcurrent() {
	// classes are immutable so no locking is required
}

func (*SuClass) IsConcurrent() Value {
	return EmptyStr
}

func (c *SuClass) Parent(th *Thread) *SuClass {
	if c.Base <= 0 {
		return nil
	}
	base := Global.Get(th, c.Base)
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

func (c *SuClass) Lookup(th *Thread, method string) Value {
	return c.lookup(th, method, nil)
}

func (c *SuClass) lookup(th *Thread, method string, parents []*SuClass) Value {
	if f, ok := ClassMethods[method]; ok {
		return f
	}
	if f, ok := BaseMethods[method]; ok {
		return f
	}
	if x := c.get2(th, method, parents); x != nil {
		return x
	}
	if method == "New" {
		return DefaultNewMethod
	}
	if x := UserDef(th, gnObjects, method); x != nil {
		return x
	}
	//TODO explicit CallClass doesn't go to Default in cSuneido or jSuneido
	if x := c.get2(th, "Default", parents); x != nil {
		return &defaultAdapter{fn: x, method: method}
	}
	return nil
}

// defaultAdapter wraps a Default method to insert the method argument
type defaultAdapter struct {
	ValueBase[*defaultAdapter]
	fn     Value
	method string
}

var _ Value = (*defaultAdapter)(nil)

func (d *defaultAdapter) Equal(other any) bool {
	return d == other
}

func (*defaultAdapter) SetConcurrent() {
	// immutable so ok
}

func (d *defaultAdapter) String() string {
	return "Default(" + d.method + " /* method */"
}

func (d *defaultAdapter) Call(th *Thread, this Value, as *ArgSpec) Value {
	method := SuStr(d.method)
	if as.Each >= EACH0 {
		args := ToContainer(th.Pop()).Slice(int(as.Each) - 1)
		args.Insert(0, method)
		th.Push(args)
		as = &ArgSpecEach0
	} else if as.Nargs == 0 {
		th.Push(method)
		as = &ArgSpec1
	} else {
		th.Push(nil) // allow for another value
		base := th.sp - 1 - int(as.Nargs)
		copy(th.stack[base+1:], th.stack[base:]) // shift over
		th.stack[base] = method
		as2 := *as
		as2.Nargs++
		as = &as2
	}
	return d.fn.Call(th, this, as)
}

func (c *SuClass) Call(th *Thread, this Value, as *ArgSpec) Value {
	if this == nil {
		this = c
	}
	if f := c.get2(th, "CallClass", nil); f != nil {
		return f.Call(th, this, as)
	}
	// default for calling a class is to create an instance
	return c.New(th, as)
}

func (c *SuClass) New(th *Thread, as *ArgSpec) *SuInstance {
	ob := NewInstance(th, c)
	nu := c.Lookup(th, "New")
	nu.Call(th, ob, as)
	return ob
}

var _ Named = (*SuClass)(nil)

func (c *SuClass) GetName() string {
	return c.Name
}

// Finder implements Findable
func (c *SuClass) Finder(th *Thread, fn func(v Value, mb *MemBase) Value) Value {
	for range inheritanceLimit {
		if x := fn(c, &c.MemBase); x != nil {
			return x
		}
		c = c.Parent(th)
		if c == nil {
			return nil
		}
	}
	panic("too many levels of inheritance")
}

var inheritanceLimit = 100

var _ Findable = (*SuClass)(nil)

// coverage ---------------------------------------------------------

func (c *SuClass) StartCoverage(count bool) {
	for _, v := range c.Data {
		if c2, ok := v.(*SuClass); ok {
			c2.StartCoverage(count) // RECURSE
		}
		if f, ok := v.(*SuFunc); ok {
			f.StartCoverage(count)
		}
	}
}

func (c *SuClass) StopCoverage() *SuObject {
	ob := &SuObject{}
	first := true
	count := false
	c.stopCoverage(ob, &first, &count)
	return ob
}

func (c *SuClass) stopCoverage(ob *SuObject, first, count *bool) {
	for _, v := range c.Data {
		if c2, ok := v.(*SuClass); ok {
			c2.stopCoverage(ob, first, count) // RECURSE
		}
		if f, ok := v.(*SuFunc); ok {
			if *first {
				*count = len(f.cover) >= len(f.Code)
				*first = false
			}
			f.getCoverage(ob, *count)
		}
	}
}
