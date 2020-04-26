// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package check processes an ast.Function
// and finds used but (possibly) not initialized,
// and initialized but (possibly) not used.
// "possibly" meaning not on all code paths.
// Does not check nested functions (they're already codegen and not Ast)
// they are checked as constructed bottom up.
package check

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Check struct {
	t *Thread
	// pos is used to store the position of the current statement
	pos int
	// AllInit is the set of variables assigned to, including conditionally
	AllInit map[string]int
	// AllUsed is the set of variables read from, including conditionally
	AllUsed   map[string]struct{}
	results   []string
	resultPos []int
}

// New returns a Check instance
func New(t *Thread) *Check {
	return &Check{t: t}
}

// Check is the main entry point.
// It can be called more than once (for nested functions).
func (ck *Check) Check(f *ast.Function) set {
	ck.AllInit = make(map[string]int)
	ck.AllUsed = make(map[string]struct{})
	var init set
	init = ck.check(f, init)
	ck.process(f.Params, init)
	return init
}

// CheckGlobal checks if a global name is defined.
// It is also called by compile constant to check class base.
func (ck *Check) CheckGlobal(name string, pos int) {
	if nil == Global.FindName(ck.t, name) {
		ck.addResult(pos, "ERROR: can't find: "+name)
	}
}

// Results returns the results sorted by code position
func (ck *Check) Results() []string {
	sort.Sort(resultsByPos{ck})
	return ck.results
}

//-------------------------------------------------------------------

func (ck *Check) check(f *ast.Function, init set) set {
	init = ck.params(f.Params, init)
	init, _ = ck.statements(f.Body, init)
	return init
}

func (ck *Check) params(params []ast.Param, init set) set {
	for _, p := range params {
		name := p.Name.Name
		if name[0] == '.' {
			name = str.UnCapitalize(name[1:])
			ck.AllUsed[name] = struct{}{}
		} else if name[0] == '@' || name[0] == '_' {
			name = name[1:]
		}
		if !p.Unused {
			init = append(init, name)
		}
	}
	return init
}

func (ck *Check) used(id string) bool {
	_, ok := ck.AllUsed[id]
	return ok
}

func (ck *Check) statements(stmts []ast.Statement, init set) (initOut set, exit bool) {
	for _, stmt := range stmts {
		if exit {
			ck.addResult(stmt.Position(), "ERROR: unreachable code")
		}
		init, exit = ck.statement(stmt, init)
	}
	return init, exit
}

// statement processes one statement (and its children)
// Conditional statements are assumed to run for used, and not to run for init.
// So we accumulate used, but not init.
func (ck *Check) statement(stmt ast.Statement, init set) (initOut set, exit bool) {
	if stmt == nil {
		return init, exit
	}
	ck.pos = stmt.Position()
	switch stmt := stmt.(type) {
	case *ast.Compound:
		init, exit = ck.statements(stmt.Body, init)
	case *ast.ExprStmt:
		init = ck.expr(stmt.E, init)
	case *ast.Return:
		init = ck.expr(stmt.E, init)
		exit = true
	case *ast.Throw:
		init = ck.expr(stmt.E, init)
		exit = true
	case *ast.TryCatch:
		init, _ = ck.statement(stmt.Try, init)
		if stmt.CatchVar.Name != "" && stmt.CatchVar.Name != "unused" &&
			!stmt.CatchVarUnused {
			init = ck.initVar(init, stmt.CatchVar.Name, int(stmt.CatchVar.Pos))
		}
		ck.statement(stmt.Catch, init)
	case *ast.While:
		initTrue, initFalse := ck.cond(stmt.Cond, init)
		ck.statement(stmt.Body, initTrue)
		init = initFalse
	case *ast.Forever:
		init, _ = ck.statement(stmt.Body, init)
	case *ast.DoWhile:
		init, _ = ck.statement(stmt.Body, init)
		init = ck.expr(stmt.Cond, init)
	case *ast.If:
		initTrue, initFalse := ck.cond(stmt.Cond, init)
		thenInit, ex1 := ck.statement(stmt.Then, initTrue)
		elseInit, ex2 := ck.statement(stmt.Else, initFalse)
		if _, ok := stmt.Then.(*ast.Return); ok {
			init = initFalse
		} else {
			init = init.unionIntersect(thenInit, elseInit)
		}
		if ex1 && ex2 {
			exit = true
		}
	case *ast.Switch:
		// there will always be at least a default default that throws
		exAll := true
		init = ck.expr(stmt.E, init)
		for _, c := range stmt.Cases {
			in := init
			for _, e := range c.Exprs {
				in = ck.expr(e, in)
			}
			in, ex := ck.statements(c.Body, in)
			exAll = exAll && ex
		}
		if stmt.Default != nil { // specifically nil and not len 0
			_, ex := ck.statements(stmt.Default, init)
			exAll = exAll && ex
		}
		if exAll {
			exit = true
		}
	case *ast.ForIn:
		init = ck.initVar(init, stmt.Var.Name, int(stmt.Var.Pos))
		init = ck.expr(stmt.E, init)
		ck.statement(stmt.Body, init)
	case *ast.For:
		for _, expr := range stmt.Init {
			init = ck.expr(expr, init)
		}
		initTrue, initFalse := ck.cond(stmt.Cond, init)
		afterBody, _ := ck.statement(stmt.Body, initTrue)
		ck.pos = stmt.Pos // restore after statement has modified
		for _, expr := range stmt.Inc {
			ck.expr(expr, afterBody)
		}
		init = initFalse
	case *ast.Break, *ast.Continue:
		exit = true
	default:
		panic("unexpected statement type " + fmt.Sprintf("%T", stmt))
	}
	return init, exit
}

func (ck *Check) cond(expr ast.Expr, init set) (initTrue set, initFalse set) {
	if u, ok := expr.(*ast.Unary); ok && u.Tok == tok.LParen {
		expr = u.E
	}
	if expr, ok := expr.(*ast.Nary); ok {
		if expr.Tok == tok.And || expr.Tok == tok.Or {
			first := ck.expr(expr.Exprs[0], init) // first is always done
			rest := first
			for _, e := range expr.Exprs[1:] {
				rest = ck.expr(e, rest) // rest are conditional
			}
			if expr.Tok == tok.And {
				return rest, first
			}
			return first, rest
		}
	}
	init = ck.expr(expr, init)
	return init, init
}

func (ck *Check) expr(expr ast.Expr, init set) set {
	if expr == nil {
		return init
	}
	switch expr := expr.(type) {
	case *ast.Binary:
		if expr.Tok == tok.Eq {
			if id, ok := expr.Lhs.(*ast.Ident); ok {
				init = ck.expr(expr.Rhs, init)
				init = ck.initVar(init, id.Name, int(id.Pos))
				break
			}
		}
		if tok.AssignStart < expr.Tok && expr.Tok < tok.AssignEnd {
			init = ck.expr(expr.Rhs, init)
			init = ck.expr(expr.Lhs, init)
		} else {
			init = ck.expr(expr.Lhs, init)
			init = ck.expr(expr.Rhs, init)
		}
	case *ast.Ident:
		if ascii.IsLower(expr.Name[0]) {
			init = ck.usedVar(init, expr.Name, int(expr.Pos))
		}
		if ascii.IsUpper(expr.Name[0]) {
			ck.CheckGlobal(expr.Name, int(expr.Pos))
		}
	case *ast.Trinary:
		initTrue, initFalse := ck.cond(expr.Cond, init)
		tInit := ck.expr(expr.T, initTrue)
		fInit := ck.expr(expr.F, initFalse)
		init = init.unionIntersect(tInit, fInit)
	case *ast.Nary:
		if expr.Tok == tok.And || expr.Tok == tok.Or {
			init = ck.expr(expr.Exprs[0], init) // first is always done
			in := init
			for _, e := range expr.Exprs[1:] {
				in = ck.expr(e, in) // rest are conditional
			}
		} else {
			expr.Children(func(e ast.Node) ast.Node {
				init = ck.expr(e.(ast.Expr), init)
				return e
			})
		}
	case *ast.Block:
		init = ck.block(expr, init)

	default:
		expr.Children(func(e ast.Node) ast.Node {
			init = ck.expr(e.(ast.Expr), init)
			return e
		})
	}
	return init
}

func (ck *Check) block(b *ast.Block, init set) set {
	// save & remove variables shadowed by params
	allInit := map[string]int{}
	allUsed := map[string]struct{}{}
	nUsedParams := 0
	for _, p := range b.Function.Params {
		id := p.Name.Name
		if !p.Unused {
			nUsedParams++
		}
		if n, ok := ck.AllInit[id]; ok {
			allInit[id] = n
			delete(ck.AllInit, id)
		}
		if _, ok := ck.AllUsed[id]; ok {
			allUsed[id] = struct{}{}
			delete(ck.AllUsed, id)
		}
	}

	// assume that blocks are executed at point of definition
	// this is not necessarily true
	// they may be called elsewhere or not at all
	// but too many spurious warnings otherwise
	before := init
	after := ck.check(&b.Function, init)
	// remove params from init
	init = append(before, after[len(before)+nUsedParams:]...)

	// detect unused params
	for _, p := range b.Function.Params {
		if !p.Unused {
			id := p.Name.Name
			if _, ok := ck.AllUsed[id]; !ok {
				ck.addResult(int(p.Name.Pos),
					"WARNING: initialized but not used: "+id)
			}
		}
	}

	// remove params
	for _, p := range b.Function.Params {
		id := p.Name.Name
		delete(ck.AllInit, id)
		delete(ck.AllUsed, id)
	}
	// restore shadowed variables
	for id, n := range allInit {
		ck.AllInit[id] = n
	}
	for id := range allUsed {
		ck.AllUsed[id] = struct{}{}
	}

	return init
}

func (ck *Check) initVar(init set, id string, pos int) set {
	if strings.HasPrefix(id, "_") || id == "unused" {
		return init
	}
	ck.AllInit[id] = pos
	return init.with(id)
}

func (ck *Check) usedVar(init set, id string, pos int) set {
	if strings.HasPrefix(id, "_") {
		return init
	}
	if id != "this" && id != "super" && !init.has(id) {
		p := "ERROR: used but"
		if _, ok := ck.AllInit[id]; ok {
			p = "WARNING: used but possibly"
		}
		ck.addResult(pos, p+" not initialized: "+id)
	}
	ck.AllUsed[id] = struct{}{}
	return init
}

//-------------------------------------------------------------------

func (ck *Check) process(params []ast.Param, init set) {
	for _, id := range init {
		if !ck.used(id) {
			var at int
			if pos := paramPos(params, id); pos >= 0 {
				at = pos
			} else if pos, ok := ck.AllInit[id]; ok {
				at = int(pos)
			}
			ck.addResult(at, "WARNING: initialized but not used: "+id)
		}
	}
	for id, pos := range ck.AllInit {
		if _, ok := ck.AllUsed[id]; !ok && !init.has(id) {
			ck.addResult(pos, "WARNING: initialized but not used: "+id)
		}
	}
}
func paramPos(params []ast.Param, id string) int {
	for _, p := range params {
		if p.Name.Name == id {
			return int(p.Name.Pos)
		}
	}
	return -1
}

//-------------------------------------------------------------------

func (ck *Check) addResult(pos int, str string) {
	ck.resultPos = append(ck.resultPos, pos)
	ck.results = append(ck.results, str+" @"+strconv.Itoa(pos))
}

// resultByPos is used to sort the results
type resultsByPos struct {
	*Check
}

func (r resultsByPos) Len() int {
	return len(r.results)
}
func (r resultsByPos) Swap(i, j int) {
	r.results[i], r.results[j] = r.results[j], r.results[i]
	r.resultPos[i], r.resultPos[j] = r.resultPos[j], r.resultPos[i]
}
func (r resultsByPos) Less(i, j int) bool {
	return r.resultPos[i] < r.resultPos[j]
}
