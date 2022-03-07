// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// container allows using memberList etc. for both objects and classes
type container interface {
	Add(Value)
	HasKey(Value) bool
	Set(Value, Value)
}

//-------------------------------------------------------------------

type cgMaker struct{}

var _ maker = (*cgMaker)(nil)

func (cgMaker) mkObject() container {
	return &SuObject{}
}
func (cgMaker) mkRecord() container {
	return NewSuRecord()
}
func (cgMaker) mkRecOrOb(rec container) container {
	if rec.(*SuRecord).ListSize() > 0 {
		return rec.(*SuRecord).ToObject()
	}
	return rec
}
func (cgMaker) mkClass(_ string) container {
	return classBuilder{}
}
func (cgMaker) mkConcat(strs []string) Value {
	return SuStr(strings.Join(strs, ""))
}

type classBuilder map[string]Value

func (c classBuilder) Add(Value) {
	panic("class members must be named")
}

func (c classBuilder) HasKey(m Value) bool {
	if s, ok := m.(SuStr); ok {
		_, ok = c[string(s)]
		return ok
	}
	panic("class member names must be strings")
}

func (c classBuilder) Set(m, v Value) {
	c[string(m.(SuStr))] = v
}

//-------------------------------------------------------------------

type astMaker struct{}

var _ maker = (*astMaker)(nil)

func (astMaker) mkObject() container {
	return &astContainer{which: "Object"}
}
func (astMaker) mkRecord() container {
	return &astContainer{which: "Record"}
}
func (astMaker) mkRecOrOb(rec container) container {
	ac := rec.(*astContainer)
	for _, kv := range ac.kv {
		if kv.key == nil {
			ac.which = "Object"
			return rec
		}
	}
	return rec
}
func (astMaker) mkClass(base string) container {
	return &astContainer{which: "Class", base: base}
}
func (astMaker) mkConcat(strs []string) Value {
	exprs := make([]ast.Expr, len(strs))
	for i, s := range strs {
		exprs[i] = &ast.Constant{Val: SuStr(s)}
	}
	return &ast.Nary{Tok: tokens.Cat, Exprs: exprs}
}

type astContainer struct {
	ast.SuAstNode
	which string
	base  string
	kv    []keyVal
}

type keyVal struct {
	key Value
	val Value
}

func (c *astContainer) Add(val Value) {
	c.kv = append(c.kv, keyVal{val: val})
}

func (c *astContainer) Set(key Value, val Value) {
	c.kv = append(c.kv, keyVal{key: key, val: val})
}

func (c *astContainer) HasKey(key Value) bool {
	for _, kv := range c.kv {
		if key.Equal(kv.key) {
			return true
		}
	}
	return false
}

// Value interface

var _ Value = (*astContainer)(nil)

func (c *astContainer) String() string {
	style := []string{"#(", ", ", ")"}
	if c.which == "Record" {
		style = []string{"#{", ", ", "}"}
	} else if c.which == "Class" {
		style = []string{c.base + " {\n", "\n", "\n}"}
	}
	var sb strings.Builder
	sep := style[0]
	for _, kv := range c.kv {
		sb.WriteString(sep)
		sep = style[1]
		if kv.key != nil {
			if ks := Unquoted(kv.key); ks != "" {
				sb.WriteString(ks)
			} else {
				sb.WriteString(Display(nil, kv.key))
			}
			sb.WriteString(": ")
		}
		sb.WriteString(Display(nil, kv.val))
	}
	sb.WriteString(style[2])
	return sb.String()
}

func (c *astContainer) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr(c.which)
	case SuStr("base"):
		return SuStr(c.base)
	case SuStr("size"):
		return IntVal(len(c.kv))
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(c.kv) <= i {
			return False
		}
		kv := c.kv[i]
		ob := &SuObject{}
		ob.Add(kv.val)
		if kv.key != nil {
			ob.Add(kv.key)
		}
		return ob
	}
	return nil
}
