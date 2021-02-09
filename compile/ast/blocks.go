// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
)

type set map[string]struct{}

var yes = struct{}{}

type blok struct {
	block  *Block
	parent *blok
	params set
	vars   set
	hasRet bool
}

type bloks struct {
	bloks []*blok
	cur   *blok
}

// Blocks sets compileAsFunction
// if a block can be compiled as a function
// i.e. isn't a closure, doesn't share any variables.
// NOTE: This is trickier than it seems.
// e.g. Have to handle nested blocks and sharing between peer blocks.
// Does not process nested functions (they're already codegen and not Ast);
// they are checked as constructed bottom up.
func Blocks(f *Function) {
	// first traverse the ast and collect outer variables
	// and a list of blocks, their params & variables, and their parent if nested.
	var b bloks
	vars := make(set)
	b.params(f.Params, vars)
	for _, stmt := range f.Body {
		b.statement(stmt, vars)
	}
	// then check for variable sharing
	for _, x := range b.bloks {
		x.block.CompileAsFunction = true
	}
	for i, x := range b.bloks {
		_, this := x.vars["this"]
		_, super := x.vars["super"]
		if this || super || x.hasRet || shares(x.vars, vars) ||
			(x.parent != nil && shares(x.vars, x.parent.params)) {
			closure(x)
			continue
		}
		for j := i + 1; j < len(b.bloks); j++ {
			y := b.bloks[j]
			if shares(x.vars, y.vars) {
				closure(x)
				closure(y)
			}
		}
	}
}

func shares(v1, v2 set) bool {
	for v := range v1 {
		if _, ok := v2[v]; ok {
			return true
		}
	}
	return false
}

func closure(x *blok) {
	for x != nil {
		x.block.CompileAsFunction = false
		x = x.parent
	}
}

func (b *bloks) params(params []Param, vars set) {
	for _, p := range params {
		name := p.Name.Name
		if name[0] == '.' {
			name = str.UnCapitalize(name[1:])
		} else if name[0] == '@' || name[0] == '_' {
			name = name[1:]
		}
		vars[name] = yes
	}
}

// statement processes one statement (and its children)
func (b *bloks) statement(stmt Statement, vars set) {
	if stmt == nil {
		return
	}
	switch stmt := stmt.(type) {
	case *Compound:
		for _, stmt := range stmt.Body {
			b.statement(stmt, vars)
		}
	case *ExprStmt:
		b.expr(stmt.E, vars)
	case *Return:
		if b.cur != nil {
			b.cur.hasRet = true
		}
		b.expr(stmt.E, vars)
	case *Throw:
		b.expr(stmt.E, vars)
	case *TryCatch:
		b.statement(stmt.Try, vars)
		if stmt.CatchVar.Name != "" {
			vars[stmt.CatchVar.Name] = yes
		}
		b.statement(stmt.Catch, vars)
	case *While:
		b.expr(stmt.Cond, vars)
		b.statement(stmt.Body, vars)
	case *Forever:
		b.statement(stmt.Body, vars)
	case *DoWhile:
		b.statement(stmt.Body, vars)
		b.expr(stmt.Cond, vars)
	case *If:
		b.expr(stmt.Cond, vars)
		b.statement(stmt.Then, vars)
		b.statement(stmt.Else, vars)
	case *Switch:
		b.expr(stmt.E, vars)
		for _, c := range stmt.Cases {
			for _, e := range c.Exprs {
				b.expr(e, vars)
			}
			for _, stmt := range c.Body {
				b.statement(stmt, vars)
			}
		}
		for _, d := range stmt.Default {
			b.statement(d, vars)
		}
	case *ForIn:
		vars[stmt.Var.Name] = yes
		b.expr(stmt.E, vars)
		b.statement(stmt.Body, vars)
	case *For:
		for _, expr := range stmt.Init {
			b.expr(expr, vars)
		}
		b.expr(stmt.Cond, vars)
		b.statement(stmt.Body, vars)
		for _, expr := range stmt.Inc {
			b.expr(expr, vars)
		}
	case *Break, *Continue:
		// nothing to do
	default:
		panic("unexpected statement type " + fmt.Sprintf("%T", stmt))
	}
}

func (b *bloks) expr(expr Expr, vars set) {
	if expr == nil {
		return
	}
	switch expr := expr.(type) {
	case *Binary:
		if expr.Tok == tok.Eq {
			if id, ok := expr.Lhs.(*Ident); ok {
				// assignment
				vars[id.Name] = yes
				b.expr(expr.Rhs, vars)
				break
			}
		}
		b.expr(expr.Lhs, vars)
		b.expr(expr.Rhs, vars)
	case *Ident:
		if !ascii.IsUpper(expr.Name[0]) {
			vars[expr.Name] = yes
		}
	case *Trinary:
		b.expr(expr.Cond, vars)
		b.expr(expr.T, vars)
		b.expr(expr.F, vars)
	case *Nary:
		if expr.Tok == tok.And || expr.Tok == tok.Or {
			b.expr(expr.Exprs[0], vars) // first is always done
			for _, e := range expr.Exprs[1:] {
				b.expr(e, vars) // rest are conditional
			}
		} else {
			expr.Children(func(e Node) Node {
				b.expr(e.(Expr), vars)
				return e
			})
		}
	case *Block:
		b.block(expr)
	default:
		expr.Children(func(e Node) Node {
			b.expr(e.(Expr), vars)
			return e
		})
	}
}

func (b *bloks) block(block *Block) {
	parent := b.cur
	params := make(set)
	b.params(block.Params, params)
	b.cur = &blok{block: block, parent: parent, params: params}
	blockVars := make(set)
	for _, stmt := range block.Body {
		b.statement(stmt, blockVars)
	}
	for v := range params {
		delete(blockVars, v)
	}
	b.cur.vars = blockVars
	b.bloks = append(b.bloks, b.cur)
	b.cur = parent
}
