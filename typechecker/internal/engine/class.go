// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// astContainer in gsuneido is private so we redefine a bare version here
package engine

import (
	"sort"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
)

type ClassObject struct {
	Name       string
	Base       string
	Members    map[string]core.Value
	Methods    map[string]*ast.Function
	IsFunction bool
}

func NewClassObject(name string, val core.Value) *ClassObject {
	cls := &ClassObject{
		Name:    name,
		Members: make(map[string]core.Value, 10),
		Methods: make(map[string]*ast.Function, 10),
	}

	if fn, ok := val.(*ast.Function); ok {
		cls.Name = "thisIsNotAnActualSuneidoClassJustAShim"
		cls.IsFunction = true
		cls.Methods[name] = fn
	}

	if base := getValue(val, "base"); base != nil {
		b := core.ToStr(base)
		if b != "" && b != "class" {
			cls.Base = b
		}
	}

	var size int
	if v := getValue(val, "size"); v != nil {
		size, _ = v.ToInt()
	}
	for i := range size {
		member := val.Get(nil, core.IntVal(i))
		if member == nil {
			continue
		}
		key := getValue(member, "key")
		val := getValue(member, "value")
		if key == nil || val == nil {
			continue
		}

		name, isStr := key.ToStr()
		if !isStr {
			continue // a numeric key names no method or member
		}
		if fn, ok := val.(*ast.Function); ok {
			cls.Methods[name] = fn
		} else {
			cls.Members[name] = val
		}
	}

	return cls

}

func (c *ClassObject) SortedMethods(yield func(string, *ast.Function) bool) {
	names := make([]string, 0, len(c.Methods))
	for n := range c.Methods {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		if !yield(n, c.Methods[n]) {
			return
		}
	}
}

func getValue(val core.Value, key string) core.Value {
	defer func() { _ = recover() }()
	return val.Get(nil, core.SuStr(key))
}

func (c *ClassObject) Lineage(r ClassResolver) []*ClassObject {
	chain := []*ClassObject{c}
	if r == nil {
		return chain
	}
	for current := c.Base; current != ""; {
		parent, err := r.Resolve(current)
		if err != nil {
			break // unresolvable base
		}
		chain = append([]*ClassObject{parent}, chain...)
		current = parent.Base
	}
	return chain
}
