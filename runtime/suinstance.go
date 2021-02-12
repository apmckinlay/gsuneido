// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// SuInstance is an instance of an SuClass
type SuInstance struct {
	MemBase
	class   *SuClass
	parents []*SuClass
}

func NewInstance(t *Thread, class *SuClass) *SuInstance {
	return &SuInstance{MemBase: NewMemBase(),
		class: class, parents: getParents(t, class)}
}

// getParents captures the inheritance chain (and caches it on the class).
// Otherwise, the chain via SuClass Base is indirect by global number,
// and if a parent is modified incompatibly or with an error
// then existing (running) instances can fail.
func getParents(t *Thread, class *SuClass) []*SuClass {
	if class == nil {
		return nil
	}
	// Use cached parents on class if valid.
	// Still have to follow inheritance chain to validate, but no allocation.
	parents := class.GetParents()
	c := class
	for _, p := range parents {
		if c != p {
			parents = nil // cached is invalid
			break
		}
		c = c.Parent(t)
	}
	if parents != nil {
		return parents // cached is valid
	}

	parents = make([]*SuClass, 0, 4)
	for c := class; c != nil; c = c.Parent(t) {
		parents = append(parents, c)
	}
	class.SetParents(parents) // cache on class
	return parents
}

func (ob *SuInstance) Base() *SuClass {
	return ob.class
}

// ToString is used by Cat, Display, and Print
// to handle user defined ToString
func (ob *SuInstance) ToString(t *Thread) string {
	if f := ob.class.get2(t, "ToString", ob.parents); f != nil && t != nil {
		x := t.CallThis(f, ob)
		if x != nil {
			if s, ok := x.ToStr(); ok {
				return s
			}
		}
		panic("ToString should return a string")
	}
	return ob.String()
}

func (ob *SuInstance) Copy() *SuInstance {
	return &SuInstance{MemBase: ob.MemBase.Copy(),
		class: ob.class, parents: ob.parents}
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
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.get(t, m)
}
func (ob *SuInstance) get(t *Thread, m Value) Value {
	if ms, ok := m.ToStr(); ok {
		if x, ok := ob.Data[ms]; ok {
			return x
		}
	}
	ob.Unlock() // can't hold lock because it may call getter
	defer ob.Lock()
	x := ob.class.get1(t, ob, m, ob.parents)
	if m, ok := x.(*SuMethod); ok {
		m.this = ob // fix 'this' to be instance, not method class
	}
	return x
}

func (ob *SuInstance) Put(_ *Thread, m Value, v Value) {
	if ob.Lock() {
		defer ob.Unlock()
		v.SetConcurrent()
	}
	ob.put(m, v)
}
func (ob *SuInstance) put(m Value, v Value) {
	ob.Data[AsStr(m)] = v
}

func (ob *SuInstance) GetPut(t *Thread, m Value, v Value,
	op func(x, y Value) Value, retOrig bool) Value {
	if ob.Lock() {
		defer ob.Unlock()
	}
	orig := ob.get(t, m)
	if orig == nil {
		panic("uninitialized member: " + m.String())
	}
	v = op(orig, v)
	ob.put(m, v)
	if retOrig {
		return orig
	}
	return v
}

func (*SuInstance) RangeTo(int, int) Value {
	panic("instance does not support range")
}

func (*SuInstance) RangeLen(int, int) Value {
	panic("instance does not support range")
}

func (*SuInstance) Hash() uint32 {
	panic("instance hash not implemented")
}

func (*SuInstance) Hash2() uint32 {
	panic("instance hash not implemented")
}

// Equal returns true if two instances have the same class and data
func (ob *SuInstance) Equal(other interface{}) bool {
	ob2, ok := other.(*SuInstance)
	return ok && deepEqual(ob, ob2)
}

func (*SuInstance) Compare(Value) int {
	panic("instance compare not implemented")
}

func (ob *SuInstance) SetConcurrent() {
	if ob.concurrent {
		return
	}
	ob.concurrent = true
	for _, v := range ob.Data {
		v.SetConcurrent() // recursive, deep
	}
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
	return ob.class.lookup(t, method, ob.parents)
}

func (ob *SuInstance) Call(t *Thread, _ Value, as *ArgSpec) Value {
	if f := ob.class.get2(t, "Call", ob.parents); f != nil {
		return f.Call(t, ob, as)
	}
	if f := ob.class.get2(t, "Default", ob.parents); f != nil {
		da := &defaultAdapter{f, "Call"}
		return da.Call(t, ob, as)
	}
	panic("method not found: Call")
}

// Finder implements Findable
func (ob *SuInstance) Finder(_ *Thread, fn func(Value, *MemBase) Value) Value {
	if x := fn(ob, &ob.MemBase); x != nil {
		return x
	}
	assert.That(ob.parents[0] == ob.class)
	for _, p := range ob.parents {
		if x := fn(p, &p.MemBase); x != nil {
			return x
		}
	}
	return nil
}

var _ Findable = (*SuInstance)(nil)

func (ob *SuInstance) Delete(key Value) {
	if ob.Lock() {
		defer ob.Unlock()
	}
	delete(ob.Data, ToStr(key))
}

func (ob *SuInstance) Clear() {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.Data = map[string]Value{}
}

func (ob *SuInstance) size() int {
	return len(ob.Data)
}
