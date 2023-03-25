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

func (cgMaker) set(c container, key, val Value, pos, end int32) {
	if key == nil {
		c.Add(val)
	} else {
		c.Set(key, val)
	}
}

func (cgMaker) setPos(container, int32, int32) {
	// ignore
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

func (astMaker) set(c container, key, val Value, pos, end int32) {
	ac := c.(*astContainer)
	ac.kv = append(ac.kv, keyVal{key: key, val: val, pos: pos, end: end})
}

func (astMaker) setPos(c container, pos1, pos2 int32) {
	c.(*astContainer).pos1 = pos1
	c.(*astContainer).pos2 = pos2
}

type astContainer struct {
	ast.SuAstNode
	which string
	base  string
	kv    []keyVal
	ast.TwoPos
	pos1 int32
	pos2 int32
}

type keyVal struct {
	ast.SuAstNode
	key Value
	val Value
	pos int32
	end int32
}

func (c *astContainer) Add(val Value) {
	panic("should not reach here")
}

func (c *astContainer) Set(key Value, val Value) {
	panic("should not reach here")
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
	sb.WriteString(style[0])
	sep := ""
	for _, kv := range c.kv {
		sb.WriteString(sep)
		sep = style[1]
		sb.WriteString(kv.String())
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
	case SuStr("pos"):
		return IntVal(c.GetPos())
	case SuStr("pos1"):
		return IntVal(int(c.pos1))
	case SuStr("pos2"):
		return IntVal(int(c.pos2))
	case SuStr("end"):
		return IntVal(c.GetEnd())
	case SuStr("children"):
		list := make([]Value, len(c.kv))
		for i := range c.kv {
			list[i] = &c.kv[i]
		}
		return ast.NewChildren(list)
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(c.kv) <= i {
			return False
		}
		return &c.kv[i]
	}
	return nil
}

func (kv *keyVal) String() string {
	var sb strings.Builder
	if kv.key != nil {
		if ks := Unquoted(kv.key); ks != "" {
			sb.WriteString(ks)
		} else {
			sb.WriteString(Display(nil, kv.key))
		}
		sb.WriteString(": ")
	}
	sb.WriteString(Display(nil, kv.val))
	return sb.String()
}

func (kv *keyVal) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Member")
	case SuStr("named"):
		return SuBool(kv.key != nil)
	case SuStr("key"):
		return kv.key
	case SuStr("value"):
		return kv.val
	case SuStr("pos"):
		return IntVal(int(kv.pos))
	case SuStr("end"):
		return IntVal(int(kv.end))
	case SuStr("children"):
		if kv.key == nil {
			return ast.NewChildren([]Value{kv.val})
		} else {
			return ast.NewChildren([]Value{kv.key, kv.val})
		}
	}
	return nil
}
