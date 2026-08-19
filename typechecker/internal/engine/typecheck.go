// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

func TypeCheckPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for name, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				checkOperatorsIn(stmt, name, env)
			}
		}
	}
	return false
}

func checkOperatorsIn(n ast.Node, method string, env TypeEnv) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		if ep.Expr != nil {
			checkOperatorsIn(ep.Expr, method, env)
		}
		return
	}
	switch e := n.(type) {
	case *ast.Binary:
		checkBinary(e, method, env)
	case *ast.Nary:
		checkNary(e, method, env)
	case *ast.Unary:
		checkUnary(e, method, env)
	}

	n.Children(func(c ast.Node) ast.Node {
		checkOperatorsIn(c, method, env)
		return c
	})
}

var operatorSpecs = map[tok.Token]struct {
	required DynType
	label    string
}{
	tok.Add:      {TNumber, "+"},
	tok.Sub:      {TNumber, "-"},
	tok.Mul:      {TNumber, "*"},
	tok.Div:      {TNumber, "/"},
	tok.Mod:      {TNumber, "%"},
	tok.BitAnd:   {TNumber, "&"},
	tok.BitOr:    {TNumber, "|"},
	tok.BitXor:   {TNumber, "^"},
	tok.LShift:   {TNumber, "<<"},
	tok.RShift:   {TNumber, ">>"},
	tok.AddEq:    {TNumber, "+="},
	tok.SubEq:    {TNumber, "-="},
	tok.MulEq:    {TNumber, "*="},
	tok.DivEq:    {TNumber, "/="},
	tok.ModEq:    {TNumber, "%="},
	tok.BitAndEq: {TNumber, "&="},
	tok.BitOrEq:  {TNumber, "|="},
	tok.BitXorEq: {TNumber, "^="},
	tok.LShiftEq: {TNumber, "<<="},
	tok.RShiftEq: {TNumber, ">>="},
	tok.Inc:      {TNumber, "++"},
	tok.PostInc:  {TNumber, "++"},
	tok.Dec:      {TNumber, "--"},
	tok.PostDec:  {TNumber, "--"},
	tok.BitNot:   {TNumber, "~"},
	tok.Cat:      {TString, "$"},
	tok.CatEq:    {TString, "$="},
}

func operatorContract(op tok.Token) (required DynType, label string, ok bool) {
	spec, found := operatorSpecs[op]
	if !found {
		return nil, "", false
	}
	return spec.required, spec.label, true
}

func checkBinary(b *ast.Binary, method string, env TypeEnv) {
	required, opLabel, ok := operatorContract(b.Tok)
	if ok {
		if d, ok := classifyOperand(method, opLabel, b.Lhs, env, required); ok {
			env.Report(&d)
		}
		if d, ok := classifyOperand(method, opLabel, b.Rhs, env, required); ok {
			env.Report(&d)
		}
		return
	}
	checkCompareCrossType(b, method, env)
}

// warns when `< <= > >=` compares statically-known singletons of
// different Suneido ranks. Suneido defines a total order
// (bool < number < string < date < object) which is a real feature
// (db indexes, generic sort) - hence warning not error.
//
// ```suneido
// if 5 < "x"   { ... }   // Number vs String - warning (cross-rank)
//
//	^   ^^^
//
// if a < b     { ... }   // optional (a: False | Number) - silent
// if 5 < 10    { ... }   // both Number - silent
// ```
// only fires when both operands are primitive singletons of different
// rank. the `False | T` "optional" case is where the tiebreaker earns
// its keep, so a warning there would be noise.
func checkCompareCrossType(b *ast.Binary, method string, env TypeEnv) {
	var opLabel string
	switch b.Tok {
	case tok.Lt:
		opLabel = "<"
	case tok.Lte:
		opLabel = "<="
	case tok.Gt:
		opLabel = ">"
	case tok.Gte:
		opLabel = ">="
	default:
		return
	}
	lt := env.GetType(b.Lhs)
	rt := env.GetType(b.Rhs)
	lo, lok := typeRankOrd(lt)
	ro, rok := typeRankOrd(rt)
	if !lok || !rok || lo == ro {
		return
	}
	pos := nodePos(b.Lhs)
	if pos == 0 {
		pos = nodePos(b.Rhs)
	}
	env.Report(&Diagnostic{
		Severity: SeverityWarning,
		Method:   method,
		Pos:      pos,
		Flag:     FlagStrictCrossTypeCompares,
		Msg: fmt.Sprintf("comparison %q between %v and %v falls back to Suneido's "+
			"cross-type ordering (bool < number < string < date < object); "+
			"relying on this tiebreaker is error-prone and not recommended when "+
			"operand types are statically known - add runtime type checks on the "+
			"operands or use Suneido.StrictCompare",
			opLabel, lt, rt),
	})
}

func typeRankOrd(t DynType) (int, bool) {
	p, ok := t.(Primitive)
	if !ok {
		return 0, false
	}
	switch p {
	case TBoolean, TFalse, TTrue:
		return 0, true
	case TNumber:
		return 1, true
	case TString:
		return 2, true
	case TDate:
		return 3, true
	case TObject, TSequence:
		return 4, true
	}
	return 0, false
}

func checkNary(n *ast.Nary, method string, env TypeEnv) {
	required, opLabel, ok := operatorContract(n.Tok)
	if !ok {
		return
	}
	for _, operand := range n.Exprs {
		if d, ok := classifyOperand(method, opLabel, operand, env, required); ok {
			env.Report(&d)
		}
	}
}

func checkUnary(u *ast.Unary, method string, env TypeEnv) {
	required, opLabel, ok := operatorContract(u.Tok)
	if !ok {
		return
	}
	if d, ok := classifyOperand(method, opLabel, u.E, env, required); ok {
		env.Report(&d)
	}
}

func flagForOpLabel(opLabel string) Flag {
	switch opLabel {
	case "$", "$=":
		return FlagStrictStringConcat
	}
	return FlagNone
}

func catCrashMember(m DynType) bool {
	if _, ok := m.(Instance); ok {
		return true
	}
	switch m {
	case TObject, TSequence, TDate, TFunction, TBlock, TClass:
		return true
	}
	return false
}

func catCrashDiag(method, opLabel string, pos int, ty DynType, badMembers []DynType) (Diagnostic, bool) {
	for _, m := range badMembers {
		if !catCrashMember(m) {
			continue
		}
		msg := fmt.Sprintf("operator %q may fail at runtime: can't convert %v to String",
			opLabel, m)
		if !dynEqual(ty, m) {
			msg += fmt.Sprintf(" (operand type %v)", ty)
		}
		return Diagnostic{
			Severity: SeverityError,
			Method:   method,
			Pos:      pos,
			Flag:     FlagNone,
			Got:      []DynType{ty},
			Msg:      msg,
		}, true
	}
	return Diagnostic{}, false
}

func isCatOp(opLabel string) bool {
	return opLabel == "$" || opLabel == "$="
}

func classifyOperand(method, opLabel string, operand ast.Expr, env TypeEnv, required DynType) (Diagnostic, bool) {
	if operand == nil {
		return Diagnostic{}, false
	}
	ty := env.GetType(operand)
	if ty == nil || ty == TVoid {
		return Diagnostic{}, false
	}
	members, dirty := decomposeForCheck(ty)
	var badMembers []DynType
	for _, m := range members {
		if !memberFits(m, required) {
			badMembers = append(badMembers, m)
		}
	}
	pos := nodePos(operand)
	switch {
	case len(badMembers) > 0 && exprGuessed(operand, method, env):
		return Diagnostic{
			Severity: SeverityWarning,
			Method:   method,
			Pos:      pos,
			Flag:     flagForOpLabel(opLabel),
			Got:      []DynType{ty},
			Msg: fmt.Sprintf("operator %q expects %v, got %v "+
				"(type guessed from builtin overloads; receiver unknown)",
				opLabel, required, ty),
		}, true
	case len(badMembers) > 0:
		if isCatOp(opLabel) {
			if d, ok := catCrashDiag(method, opLabel, pos, ty, badMembers); ok {
				return d, true
			}
		}
		return Diagnostic{
			Severity: SeverityError,
			Method:   method,
			Pos:      pos,
			Flag:     flagForOpLabel(opLabel),
			Got:      []DynType{ty},
			Msg: fmt.Sprintf("operator %q expects %v, got %v",
				opLabel, required, ty),
		}, true
	case dirty:
		return Diagnostic{
			Severity: SeverityWarning,
			Method:   method,
			Pos:      pos,
			Got:      []DynType{ty},
			Msg: fmt.Sprintf("operator %q operand has type %v cannot prove it is %v",
				opLabel, ty, required),
		}, true
	}
	return Diagnostic{}, false
}

func nodePos(n ast.Node) int {
	if n == nil {
		return 0
	}
	type positioner interface{ GetPos() int }
	if p, ok := n.(positioner); ok {
		return p.GetPos()
	}
	if m, ok := n.(*ast.Mem); ok {
		return int(m.DotPos)
	}
	type stmtPositioner interface{ Position() int }
	if p, ok := n.(stmtPositioner); ok {
		return p.Position()
	}
	var found int
	n.Children(func(c ast.Node) ast.Node {
		if found == 0 {
			if p := nodePos(c); p > 0 {
				found = p
			}
		}
		return c
	})
	return found
}
