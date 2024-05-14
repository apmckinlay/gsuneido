// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/check"
	"github.com/apmckinlay/gsuneido/core"
)

// Aspects allows tailoring the parser for different purposes.
// e.g. codegen, ast generation, queries, go gen
type Aspects interface {
	ast.Builder
	checker
	maker

	privatize(name, className string) string
	codegen(lib, name string, fn *ast.Function, prevDef core.Value) core.Value
	// exprPos is used to wrap expressions with their position for astvalue
	exprPos(expr ast.Expr, pos, end int32) ast.Expr
}

type checker interface {
	CheckFunc(*ast.Function)
	CheckGlobal(name string, pos int)
	CheckResult(pos int, str string)
	CheckResults() []string
}

var _ checker = (*check.Check)(nil)

type maker interface {
	mkObject() container
	mkRecord() container
	mkRecOrOb(container) container
	mkClass(base string) container
	mkConcat([]string) core.Value
	set(c container, key, val core.Value, pos, end int32)
	setPos(container, int32, int32)
}

// cgAspects is used to compile code --------------------------------
type cgAspects struct {
	cgAspectsBase
	nilChecker
}

type cgAspectsBase struct {
	ast.Folder
	cgMaker
}

func (*cgAspectsBase) privatize(name, className string) string {
	if strings.HasPrefix(name, "getter_") {
		return "Getter_" + className + "_" + name[7:]
	}
	return className + "_" + name
}

func (*cgAspectsBase) codegen(lib, name string, fn *ast.Function, prevDef core.Value) core.Value {
	return codegen(lib, name, fn, prevDef)
}

func (*cgAspectsBase) exprPos(expr ast.Expr, _, _ int32) ast.Expr {
	return expr
}

type cgckAspects struct {
	cgAspectsBase
	*check.Check
}

// astAspects is used to generate an AST ----------------------------
type astAspects struct {
	ast.Factory
	astMaker
	nilChecker
}

func (*astAspects) Symbol(s core.SuStr) ast.Expr {
	return &ast.Symbol{Constant: ast.Constant{Val: s}}
}

func (*astAspects) privatize(name, _ string) string {
	return name
}

func (*astAspects) codegen(_, _ string, fn *ast.Function, _ core.Value) core.Value {
	return fn
}

func (astMaker) exprPos(expr ast.Expr, pos, end int32) ast.Expr {
	return &ast.ExprPos{Expr: expr, Pos: pos, End: end}
}

type nilChecker struct{}

func (nilChecker) CheckFunc(*ast.Function) {
}
func (nilChecker) CheckGlobal(string, int) {
}
func (nilChecker) CheckResult(int, string) {
}
func (nilChecker) CheckResults() []string {
	return nil
}
