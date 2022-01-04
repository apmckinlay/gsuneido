// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/check"
	"github.com/apmckinlay/gsuneido/runtime"
)

// Aspects allows tailoring the parser for different purposes.
// e.g. codegen, ast generation, queries, go gen
type Aspects interface {
	ast.Builder
	checker
	maker

	privatize(name, className string) string
	codegen(lib, name string, fn *ast.Function, prevDef runtime.Value) runtime.Value
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
	mkConcat([]string) runtime.Value
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
	return className + "_" + name
}

func (*cgAspectsBase) codegen(lib, name string, fn *ast.Function, prevDef runtime.Value) runtime.Value {
	return codegen(lib, name, fn, prevDef)
}

type cgckAspects struct {
	cgAspectsBase
	*check.Check
}

// gogenAspects is used when transpiling to Go ----------------------
type gogenAspects struct {
	cgAspectsBase
	nilChecker
}

// codegen defined in gogen.go

// astAspects is used to generate an AST ----------------------------
type astAspects struct {
	ast.Factory
	astMaker
	nilChecker
}

func (*astAspects) Symbol(s runtime.SuStr) ast.Expr {
	return &ast.Symbol{Constant: ast.Constant{Val: s}}
}

func (*astAspects) privatize(name, _ string) string {
	return name
}

func (*astAspects) codegen(_, _ string, fn *ast.Function, _ runtime.Value) runtime.Value {
	return fn
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
