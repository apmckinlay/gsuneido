// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core_test

import (
	"testing"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSuClassChain_instance(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	def := func(name, src string) *SuClass {
		c := compile.Constant(src)
		if Global.Exists(name) {
			Global.TestDef(name, c)
		} else {
			Global.Add(name, c)
		}
		return c.(*SuClass)
	}
	call := func(x Value, src string) Value {
		fn := compile.Constant("function(x) { " + src + " }").(*SuFunc)
		return th.Call(fn, x)
	}
	c := def("C", "class { F() { 123 } }")
	b := def("B", "class : C { }")
	a := def("A", "class : B { }")
	i := NewInstance(th, a)
	assert.This(i.Parents()).Is([]*SuClass{a, b, c})
	i = NewInstance(th, a)
	assert.This(i.Parents()).Is([]*SuClass{a, b, c})
	b = def("B", "class { F() { false }}")
	x := call(i, "x.F()")
	assert.This(x).Is(SuInt(123))

	m := call(i, "x.F") // bound method
	x = call(m, "x()")
	assert.This(x).Is(SuInt(123))

	i = NewInstance(th, a) // new capture with new B
	assert.This(i.Parents()).Is([]*SuClass{a, b})
	assert.This(call(i, "x.F()")).Is(False)
}

func TestSuClassChain_overload_instance(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	def := func(lib, name, src string, prev Value) *SuClass {
		c := compile.NamedConstant(lib, name, src, prev)
		return c.(*SuClass)
	}
	call := func(x Value, src string) Value {
		fn := compile.Constant("function(x) { " + src + " }").(*SuFunc)
		return th.Call(fn, x)
	}
	c := def("stdlib", "A", "class { F() { 123 } }", nil)
	b := def("axonlib", "A", "class : _A { }", c)
	a := def("etalib", "A", "class : _A { }", b)
	i := NewInstance(th, a)
	assert.This(i.Parents()).Is([]*SuClass{a, b, c})
	assert.This(call(i, "x.F()")).Is(SuInt(123))
	Global.UnloadAll()
	assert.This(call(i, "x.F()")).Is(SuInt(123))

	m := call(i, "x.F") // bound method
	x := call(m, "x()")
	assert.This(x).Is(SuInt(123))
}

func TestSuClassChain_overload_method(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	def := func(lib, name, src string, prev Value) *SuClass {
		c := compile.NamedConstant(lib, name, src, prev)
		return c.(*SuClass)
	}
	Libload = func(_ *Thread, name string) (Value, any) {
		if name != "A" {
			return nil, nil
		}
		c := def("stdlib", "A", "class { F() { 123 } }", nil)
		b := def("axonlib", "A", "class : _A { }", c)
		a := def("etalib", "A",
			`class : _A { 
				M() {
					Unload()
					x = .F()
					Unload()
					y = .F()
					return x + y
				}
			}`, b)
		return a, nil
	}
	x := compile.EvalString(th, "A.M()")
	assert.This(x).Is(SuInt(246))
}

func TestSuClassChain_overload_bound(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	def := func(lib, name, src string, prev Value) *SuClass {
		c := compile.NamedConstant(lib, name, src, prev)
		return c.(*SuClass)
	}
	Libload = func(_ *Thread, name string) (Value, any) {
		if name != "A" {
			return nil, nil
		}
		c := def("stdlib", "A", "class { F() { 123 } }", nil)
		b := def("axonlib", "A", "class : _A { }", c)
		a := def("etalib", "A",
			`class : _A { 
				M() {
					f = .F
					Unload()
					x = f()
					Unload()
					y = f()
					return x + y
				}
			}`, b)
		return a, nil
	}
	x := compile.EvalString(th, "A.M()")
	assert.This(x).Is(SuInt(246))
}

func TestSuClassChain_equal(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	def := func(name, src string) *SuClass {
		c := compile.Constant(src)
		if Global.Exists(name) {
			Global.TestDef(name, c)
		} else {
			Global.Add(name, c)
		}
		return c.(*SuClass)
	}
	a := def("ChainEqA", "class { }")
	b := def("ChainEqB", "class { }")
	_ = NewInstance(th, a)
	_ = NewInstance(th, b)
	ca := a.GetParents()
	cb := b.GetParents()
	assert.That(ca != nil)
	assert.That(cb != nil)
	assert.This(a).Is(ca)
	assert.This(ca).Is(a)
	assert.This(ca).Isnt(cb)
}

func TestSuClassChain_getter_default(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	def := func(name, src string) *SuClass {
		c := compile.Constant(src)
		if Global.Exists(name) {
			Global.TestDef(name, c)
		} else {
			Global.Add(name, c)
		}
		return c.(*SuClass)
	}
	call := func(x Value, src string) Value {
		fn := compile.Constant("function(x) { " + src + " }").(*SuFunc)
		return th.Call(fn, x)
	}
	_ = def("ChainGetB",
		`class {
			Getter_Foo() { return 111 }
			Default(m) { return "d-" $ m }
		}`)
	a := def("ChainGetA", "class : ChainGetB { }")
	_ = NewInstance(th, a)
	ccFoo := a.GetParents()
	assert.That(ccFoo != nil)
	assert.This(call(ccFoo, "x.Foo")).Is(SuInt(111))
	assert.This(call(ccFoo, "x.Baz()")).Is(SuStr("d-Baz"))

	_ = def("ChainGetB",
		`class {
			Getter_(m) { return "g-" $ m }
			Default(m) { return "d2-" $ m }
		}`)
	_ = NewInstance(th, a)
	ccGet := a.GetParents()
	assert.That(ccGet != nil)
	assert.This(call(ccGet, "x.Bar")).Is(SuStr("g-Bar"))
	assert.This(call(ccGet, "x.Baz()")).Is(SuStr("d2-Baz"))
	assert.This(call(ccFoo, "x.Foo")).Is(SuInt(111))
	assert.This(call(ccFoo, "x.Baz()")).Is(SuStr("d-Baz"))

	_ = def("ChainGetB",
		`class {
			Getter_(m) { return "new-" $ m }
			Default(m) { return "new-" $ m }
		}`)
	assert.This(call(ccGet, "x.Bar")).Is(SuStr("g-Bar"))
	assert.This(call(ccGet, "x.Baz()")).Is(SuStr("d2-Baz"))
}

func TestSuClassChain_super(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	prevLibload := Libload
	t.Cleanup(func() { Libload = prevLibload })
	def := func(lib, name, src string, prev Value) *SuClass {
		c := compile.NamedConstant(lib, name, src, prev)
		return c.(*SuClass)
	}
	Libload = func(_ *Thread, name string) (Value, any) {
		if name != "A" {
			return nil, nil
		}
		c := def("stdlib", "A", "class { F() { 123 } }", nil)
		b := def("axonlib", "A", "class : _A { }", c)
		a := def("etalib", "A",
			`class : _A { 
				M() {
					x = super.F()
					Unload()
					y = super.F()
					return x + y
				}
			}`, b)
		return a, nil
	}
	x := compile.EvalString(th, "A.M()")
	assert.This(x).Is(SuInt(246))
}

func TestSuClassChain_callclass(t *testing.T) {
	assert := assert.T(t)
	th := &Thread{}
	prevLibload := Libload
	t.Cleanup(func() { Libload = prevLibload })
	def := func(lib, name, src string, prev Value) *SuClass {
		c := compile.NamedConstant(lib, name, src, prev)
		return c.(*SuClass)
	}
	Libload = func(_ *Thread, name string) (Value, any) {
		if name != "A" {
			return nil, nil
		}
		c := def("stdlib", "A", "class { CallClass() { return 7 } }", nil)
		b := def("axonlib", "A", "class : _A { }", c)
		a := def("etalib", "A",
			`class : _A { 
				M() {
					x = this()
					Unload()
					y = this()
					return x + y
				}
			}`, b)
		return a, nil
	}
	x := compile.EvalString(th, "A.M()")
	assert.This(x).Is(SuInt(14))
}
