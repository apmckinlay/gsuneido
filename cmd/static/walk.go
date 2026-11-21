// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

// walkRefs traverses the AST returned by parse and calls fn with the type
// ("global" or "string") and value of each global identifier and string
// constant. A function result is traversed directly; a class result has each of
// its member functions traversed; other results are traversed as values so
// strings in containers are included.
func walkRefs(val core.Value, fn func(typ, value string)) {
	if fnAST, ok := val.(*ast.Function); ok {
		walkNode(fnAST, fn)
		return
	}
	if typ := getVal(val, core.StrVal("type")); typ != nil &&
		core.StrOrEmpty(typ) == "Class" {
		walkMethods(val, fn)
		return
	}
	walkValue(val, fn)
}

// walkMethods traverses each member value of a class.
func walkMethods(cls core.Value, fn func(typ, value string)) {
	sizeVal := getVal(cls, core.StrVal("size"))
	if sizeVal == nil {
		return
	}
	size, ok := sizeVal.ToInt()
	if !ok {
		return
	}
	for i := range size {
		member := getVal(cls, core.IntVal(i))
		walkValue(getVal(member, core.StrVal("value")), fn)
	}
}

func walkNode(node ast.Node, fn func(typ, value string)) {
	switch node := node.(type) {
	case *ast.ExprPos:
		walkNode(node.Expr, fn)
	case *ast.Ident:
		if node.Name != "" && ascii.IsUpper(node.Name[0]) {
			fn("global", node.Name)
		}
	case *ast.Mem:
		walkNode(node.E, fn)
		if member, ok := memberName(node.M); ok {
			// a.b.c may get a rule, but .x, this.x, and super.x are not record fields
			if name, ok := node.E.(*ast.Ident); !ok ||
				(name.Name != "this" && name.Name != "super") {
				if member != "" && ascii.IsLower(member[0]) {
					fn("global", "Rule_"+member)
				}
			}
		} else {
			walkNode(node.M, fn)
		}
	case *ast.Constant:
		walkValue(node.Val, fn)
	default:
		node.Children(func(child ast.Node) ast.Node {
			walkNode(child, fn)
			return child
		})
	}
}

// memberName returns the name of a constant member expression.
func memberName(expr ast.Expr) (string, bool) {
	c, ok := expr.(*ast.Constant)
	if !ok {
		return "", false
	}
	return c.Val.ToStr()
}

// walkValue traverses a constant value. Containers hold their members as
// keyVal values, which expose their key and value.
func walkValue(val core.Value, fn func(typ, value string)) {
	if val == nil {
		return
	}
	if nested, ok := val.(*ast.Function); ok {
		walkNode(nested, fn)
		return
	}
	if s, ok := val.ToStr(); ok {
		fn("string", s)
		return
	}
	typ := getVal(val, core.StrVal("type"))
	if typ == nil {
		return
	}
	switch core.StrOrEmpty(typ) {
	case "Member":
		walkValue(getVal(val, core.StrVal("key")), fn)
		walkValue(getVal(val, core.StrVal("value")), fn)
	case "Object", "Record", "Class":
		sizeVal := getVal(val, core.StrVal("size"))
		if sizeVal == nil {
			return
		}
		size, ok := sizeVal.ToInt()
		if !ok {
			return
		}
		for i := range size {
			walkValue(getVal(val, core.IntVal(i)), fn)
		}
	}
}

// getVal is a panic-safe Value.Get.
func getVal(val, key core.Value) (result core.Value) {
	defer func() { _ = recover() }()
	if val == nil {
		return nil
	}
	return val.Get(nil, key)
}
