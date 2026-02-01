// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	. "github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/slc"
)

// renameExpr is used by Where Transform on Rename.
// It renames identifiers in an expression.
// It does not modify the expression.
// If any renames are done, it returns a new expression.
func renameExpr(expr Expr, r *Rename) Expr {
	switch e := expr.(type) {
	case *Constant:
		return expr
	case *Ident:
		// this is the actual rename
		// the other cases are just traversal and path copying
		if to := r.renameRev([]string{e.Name})[0]; to != e.Name {
			return &Ident{Name: to}
		}
		return expr
	case *Unary:
		newExpr := renameExpr(e.E, r)
		if newExpr == expr {
			return expr
		}
		return &Unary{Tok: e.Tok, E: newExpr}
	case *Binary:
		lhs := renameExpr(e.Lhs, r)
		rhs := renameExpr(e.Rhs, r)
		if lhs == e.Lhs && rhs == e.Rhs {
			return expr
		}
		return &Binary{Tok: e.Tok, Lhs: lhs, Rhs: rhs}
	case *Mem:
		e2 := renameExpr(e.E, r)
		m := renameExpr(e.M, r)
		if e2 == e.E && m == e.M {
			return expr
		}
		return &Mem{E: e2, M: m}
	case *Trinary:
		cond := renameExpr(e.Cond, r)
		t := renameExpr(e.T, r)
		f := renameExpr(e.F, r)
		if cond == e.Cond && t == e.T && f == e.F {
			return expr
		}
		return &Trinary{Cond: cond, T: t, F: f}
	case *RangeTo:
		cond := renameExpr(e.E, r)
		f := renameExpr(e.From, r)
		t := renameExpr(e.To, r)
		if cond == e.E && f == e.From && t == e.To {
			return expr
		}
		return &RangeTo{E: cond, From: f, To: t}
	case *RangeLen:
		cond := renameExpr(e.E, r)
		f := renameExpr(e.From, r)
		n := renameExpr(e.Len, r)
		if cond == e.E && f == e.From && n == e.Len {
			return expr
		}
		return &RangeLen{E: cond, From: f, Len: n}
	case *Nary:
		exprs := renameExprs(e.Exprs, r)
		if exprs == nil {
			return expr
		}
		return &Nary{Tok: e.Tok, Exprs: exprs}
	case *Call:
		fn := renameExpr(e.Fn, r)
		args := renameArgs(e.Args, r)
		if fn == e.Fn && args == nil {
			return expr
		}
		if args == nil {
			args = e.Args
		}
		return &Call{Fn: fn, Args: args}
	case *In:
		e2 := renameExpr(e.E, r)
		exprs := renameExprs(e.Exprs, r)
		if e2 == e.E && exprs == nil {
			return expr
		}
		if exprs == nil {
			exprs = e.Exprs
		}
		return &In{E: e2, Exprs: exprs}
	case *InRange:
		e2 := renameExpr(e.E, r)
		if e2 == e.E {
			return expr
		}
		return &InRange{E: e2, OrgTok: e.OrgTok, Org: e.Org,
			EndTok: e.EndTok, End: e.End}
	default:
		panic(assert.ShouldNotReachHere())
	}
}

func renameExprs(exprs []Expr, r *Rename) []Expr {
	var newExprs []Expr
	for i, e := range exprs {
		e2 := renameExpr(e, r)
		if e2 != e {
			if newExprs == nil {
				newExprs = make([]Expr, len(exprs))
				copy(newExprs, exprs[:i])
			}
		}
		if newExprs != nil {
			newExprs[i] = e2
		}
	}
	return newExprs
}

func renameArgs(args []Arg, r *Rename) []Arg {
	var newArgs []Arg
	for i, a := range args {
		e2 := renameExpr(a.E, r)
		if e2 != a.E {
			if newArgs == nil {
				newArgs = make([]Arg, len(args))
				copy(newArgs, args[:i])
			}
		}
		if newArgs != nil {
			newArgs[i] = Arg{E: e2, Name: a.Name}
		}
	}
	return newArgs
}

var aFolder Folder

// replaceExpr is used by Where Transform on Extend.
// It replaces identifiers in an expression with expressions.
// It does not modify the original expression.
// If any replacements are done, it returns a new expression.
func replaceExpr(expr Expr, from []string, to []Expr, clone bool) Expr {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *Constant:
		// Constant is not dependent on context so no need to clone
		return expr
	case *Ident:
		// this is the actual replace
		// the other cases are just traversal and path copying
		if i := slc.LastIndex(from, e.Name); i != -1 {
			// need to clone regardless of changes
			return replaceExpr(to[i], from[:i], to[:i], true)
		}
		return expr
	case *Unary:
		newExpr := replaceExpr(e.E, from, to, clone)
		if newExpr == expr && !clone {
			return expr
		}
		return aFolder.Unary(e.Tok, newExpr)
	case *Binary:
		lhs := replaceExpr(e.Lhs, from, to, clone)
		rhs := replaceExpr(e.Rhs, from, to, clone)
		if lhs == e.Lhs && rhs == e.Rhs && !clone {
			return expr
		}
		// if it could be evaluated raw then we need to make a copy
		return aFolder.Binary(lhs, e.Tok, rhs)
	case *Mem:
		e2 := replaceExpr(e.E, from, to, clone)
		m := replaceExpr(e.M, from, to, clone)
		if e2 == e.E && m == e.M {
			return expr
		}
		return &Mem{E: e2, M: m}
	case *Trinary:
		cond := replaceExpr(e.Cond, from, to, clone)
		t := replaceExpr(e.T, from, to, clone)
		f := replaceExpr(e.F, from, to, clone)
		if cond == e.Cond && t == e.T && f == e.F && !clone {
			return expr
		}
		return aFolder.Trinary(cond, t, f)
	case *RangeTo:
		cond := replaceExpr(e.E, from, to, clone)
		f := replaceExpr(e.From, from, to, clone)
		t := replaceExpr(e.To, from, to, clone)
		if cond == e.E && f == e.From && t == e.To {
			return expr
		}
		return &RangeTo{E: cond, From: f, To: t}
	case *RangeLen:
		cond := replaceExpr(e.E, from, to, clone)
		f := replaceExpr(e.From, from, to, clone)
		n := replaceExpr(e.Len, from, to, clone)
		if cond == e.E && f == e.From && n == e.Len {
			return expr
		}
		return &RangeLen{E: cond, From: f, Len: n}
	case *Nary:
		exprs := replaceExprs(e.Exprs, from, to, clone)
		if exprs == nil && !clone {
			return expr
		}
		if exprs == nil {
			exprs = e.Exprs
		}
		return aFolder.Nary(e.Tok, exprs)
	case *Call:
		fn := replaceExpr(e.Fn, from, to, clone)
		args := replaceArgs(e.Args, from, to, clone)
		if fn == e.Fn && args == nil && !clone {
			return expr
		}
		if args == nil {
			args = e.Args
		}
		return aFolder.Call(fn, args, 0)
	case *In:
		e2 := replaceExpr(e.E, from, to, clone)
		exprs := replaceExprs(e.Exprs, from, to, clone)
		if e2 == e.E && exprs == nil && !clone {
			return expr
		}
		if exprs == nil {
			exprs = e.Exprs
		}
		// if it could be evaluated raw then we need to make a copy
		return aFolder.In(e2, exprs)
	case *InRange:
		e2 := replaceExpr(e.E, from, to, clone)
		if e2 == e.E && !clone {
			return expr
		}
		return aFolder.Nary(tok.And, []Expr{
			aFolder.Binary(e2, e.OrgTok, e.Org),
			aFolder.Binary(e2, e.EndTok, e.End)})
	default:
		panic(assert.ShouldNotReachHere())
	}
}

// replaceExprs returns nil if nothing was replaced,
// otherwise it returns a modified copy of the expression list
func replaceExprs(exprs []Expr, from []string, to []Expr, clone bool) []Expr {
	var newExprs []Expr
	for i, e := range exprs {
		e2 := replaceExpr(e, from, to, clone)
		if e2 != e {
			if newExprs == nil {
				newExprs = make([]Expr, len(exprs))
				copy(newExprs, exprs[:i])
			}
		}
		if newExprs != nil {
			newExprs[i] = e2
		}
	}
	return newExprs
}

func replaceArgs(args []Arg, from []string, to []Expr, clone bool) []Arg {
	var newArgs []Arg
	for i, a := range args {
		e2 := replaceExpr(a.E, from, to, clone)
		if e2 != a.E {
			if newArgs == nil {
				newArgs = make([]Arg, len(args))
				copy(newArgs, args[:i])
			}
		}
		if newArgs != nil {
			newArgs[i] = Arg{E: e2, Name: a.Name}
		}
	}
	return newArgs
}
