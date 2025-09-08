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
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Check struct {
	th *Thread
	// AllInit is the set of variables assigned to, including conditionally
	AllInit map[string]int
	// AllUsed is the set of variables read from, including conditionally
	AllUsed   map[string]struct{}
	results   []string
	resultPos []int
	// pos is used to store the position of the current statement
	pos int
}

// New returns a Check instance
func New(th *Thread) *Check {
	return &Check{th: th}
}

// CheckFunc is the main entry point.
// It can be called more than once (for nested functions).
func (ck *Check) CheckFunc(f *ast.Function) {
	ck.CheckFunc2(f)
}

func (ck *Check) CheckFunc2(f *ast.Function) set {
	ck.AllInit = make(map[string]int)
	ck.AllUsed = make(map[string]struct{})
	var init set = make([]string, 0, 8)
	init = ck.check(f, init)
	ck.process(f.Params, init)
	return init
}

// CheckGlobal checks if a global name is defined.
// It is also called by compile constant to check class base.
func (ck *Check) CheckGlobal(name string, pos int) {
	if nil == Global.Find(ck.th, Global.Num(name)) {
		ck.CheckResult(pos, "ERROR: can't find: "+name)
	}
}

// CheckResults returns the results sorted by code position
func (ck *Check) CheckResults() []string {
	sort.Sort(resultsByPos{ck})
	return ck.results
}

//-------------------------------------------------------------------

func (ck *Check) check(f *ast.Function, init set) set {
	init = ck.params(f.Params, init)
	init, _ = ck.statements(f.Body, init, true)
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

func (ck *Check) statements(
	stmts []ast.Statement, init set, fnBody bool) (initOut set, exit bool) {
	for si, stmt := range stmts {
		if exit {
			ck.CheckResult(stmt.Position(), "ERROR: unreachable code")
		}
		init, exit = ck.statement(stmt, init, fnBody && si == len(stmts)-1)
	}
	return init, exit
}

// statement processes one statement (and its children)
// Conditional statements are assumed to run for used, and not to run for init.
// So we accumulate used, but not init.
func (ck *Check) statement(
	stmt ast.Statement, init set, last bool) (initOut set, exit bool) {
	if stmt == nil {
		return init, exit
	}
	var effects bool
	ck.pos = stmt.Position()
	switch stmt := stmt.(type) {
	case *ast.Compound:
		init, exit = ck.statements(stmt.Body, init, false)
	case *ast.ExprStmt:
		init, effects = ck.expr(stmt.E, init)
		if !last && !effects {
			ck.CheckResult(stmt.Position(), "ERROR: useless expression")
		}
	case *ast.Return:
		for _, e := range stmt.Exprs {
			init, _ = ck.expr(e, init)
		}
		exit = true
	case *ast.MultiAssign:
		for _, e := range stmt.Lhs {
			id := e.(*ast.Ident)
			if id.Name != "unused" {
				init = ck.initVar(init, id.Name, int(id.Pos))
			}
		}
		init, _ = ck.expr(stmt.Rhs, init)
	case *ast.Throw:
		init, _ = ck.expr(stmt.E, init)
		exit = true
	case *ast.TryCatch:
		var exit1 bool
		if expr, ok := ck.exprStmt(stmt.Try); ok {
			// allow useless expression as try statement
			init, _ = ck.expr(expr, init)
		} else {
			init, exit1 = ck.statement(stmt.Try, init, false)
		}
		if stmt.CatchVar.Name != "" && stmt.CatchVar.Name != "unused" &&
			!stmt.CatchVarUnused {
			init = ck.initVar(init, stmt.CatchVar.Name, int(stmt.CatchVar.Pos))
		}
		_, exit2 := ck.statement(stmt.Catch, init, false)
		if exit1 && exit2 {
			exit = true
		}
	case *ast.While:
		initTrue, initFalse := ck.cond(stmt.Cond, init)
		ck.statement(stmt.Body, initTrue, false)
		init = initFalse
	case *ast.Forever:
		init, _ = ck.statement(stmt.Body, init, false)
		// Forever loop exits only if there are no break statements
		if !ck.hasBreak(stmt.Body) {
			exit = true
		}
	case *ast.DoWhile:
		init, _ = ck.statement(stmt.Body, init, false)
		init, _ = ck.expr(stmt.Cond, init)
	case *ast.If:
		initTrue, initFalse := ck.cond(stmt.Cond, init)
		thenInit, ex1 := ck.statement(stmt.Then, initTrue, false)
		elseInit, ex2 := ck.statement(stmt.Else, initFalse.cow(), false)
		if ck.isReturn(stmt.Then) {
			init = elseInit
		} else if ck.isReturn(stmt.Else) {
			init = thenInit
		} else {
			init = init.unionIntersect(thenInit, elseInit)
		}
		if ex1 && ex2 {
			exit = true
		}
	case *ast.Switch:
		// there will always be at least a default default that throws
		exAll := true
		init, _ = ck.expr(stmt.E, init)
		var initInAll set
		for _, c := range stmt.Cases {
			in := init
			for _, e := range c.Exprs {
				in, _ = ck.expr(e, in)
			}
			in, ex := ck.statements(c.Body, in, false)
			exAll = exAll && ex
			if initInAll == nil {
				initInAll = in.copy()
			} else {
				initInAll = initInAll.intersect(in)
			}
		}
		if stmt.Default != nil { // specifically nil and not len 0
			in, ex := ck.statements(stmt.Default, init, false)
			exAll = exAll && ex
			if initInAll == nil {
				initInAll = in.copy()
			} else {
				initInAll = initInAll.intersect(in)
			}
		}
		if exAll {
			exit = true
		}
		init = init.union(initInAll)
	case *ast.ForIn:
		if stmt.Var.Name != "" {
			init = ck.initVar(init, stmt.Var.Name, int(stmt.Var.Pos))
		}
		if stmt.Var2.Name != "" {
			init = ck.initVar(init, stmt.Var2.Name, int(stmt.Var2.Pos))
		}
		init, _ = ck.expr(stmt.E, init)
		init, _ = ck.expr(stmt.E2, init)
		ck.statement(stmt.Body, init, false)
	case *ast.For:
		for _, expr := range stmt.Init {
			init, _ = ck.expr(expr, init)
		}
		initTrue, initFalse := ck.cond(stmt.Cond, init)
		afterBody, _ := ck.statement(stmt.Body, initTrue, false)
		ck.pos = stmt.Position() // restore after statement has modified
		for _, expr := range stmt.Inc {
			ck.expr(expr, afterBody)
		}
		init = initFalse
		if stmt.Cond == nil && !ck.hasBreak(stmt.Body) {
			exit = true
		}
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
			first, _ := ck.expr(expr.Exprs[0], init) // first is always done
			rest := first
			for _, e := range expr.Exprs[1:] {
				rest, _ = ck.expr(e, rest) // rest are conditional
			}
			if expr.Tok == tok.And {
				return rest, first
			}
			return first, rest
		}
	}
	init, _ = ck.expr(expr, init)
	return init, init
}

func (*Check) isReturn(stmt ast.Statement) bool {
	if cmpd, ok := stmt.(*ast.Compound); ok && len(cmpd.Body) == 1 {
		stmt = cmpd.Body[0]
	}
	_, ok := stmt.(*ast.Return)
	return ok
}

func (*Check) exprStmt(stmt ast.Statement) (expr ast.Expr, ok bool) {
	if cmpd, ok := stmt.(*ast.Compound); ok && len(cmpd.Body) == 1 {
		stmt = cmpd.Body[0]
	}
	// TODO also allow <exprstmt>, <return>
	if es, ok := stmt.(*ast.ExprStmt); ok {
		return es.E, true
	}
	return nil, false
}

func (ck *Check) expr(expr ast.Expr, init set) (initOut set, effects bool) {
	if expr == nil {
		return init, false
	}
	var ef1, ef2 bool
	switch expr := expr.(type) {
	case *ast.Unary:
		init, effects = ck.expr(expr.E, init)
		switch expr.Tok {
		case tok.Inc, tok.Dec, tok.PostInc, tok.PostDec:
			effects = true
		}
	case *ast.Binary:
		if expr.Tok == tok.Eq {
			effects = true
			if id, ok := expr.Lhs.(*ast.Ident); ok {
				init, _ = ck.expr(expr.Rhs, init)
				init = ck.initVar(init, id.Name, int(id.Pos))
				break
			}
		}
		if tok.AssignStart < expr.Tok && expr.Tok < tok.AssignEnd {
			init, _ = ck.expr(expr.Rhs, init)
			init, _ = ck.expr(expr.Lhs, init)
			effects = true
		} else {
			init, ef1 = ck.expr(expr.Lhs, init)
			init, ef2 = ck.expr(expr.Rhs, init)
			effects = ef1 || ef2
		}
	case *ast.Mem:
		effects = true // accessing member of record can have rule effects
		expr.Children(func(e ast.Node) ast.Node {
			init, _ = ck.expr(e.(ast.Expr), init)
			return e
		})
	case *ast.Ident:
		if ascii.IsLower(expr.Name[0]) {
			init = ck.usedVar(init, expr.Name, int(expr.Pos))
		}
		if ascii.IsUpper(expr.Name[0]) {
			ck.CheckGlobal(expr.Name, int(expr.Pos))
		}
	case *ast.Trinary:
		initTrue, initFalse := ck.cond(expr.Cond, init)
		tInit, ef1 := ck.expr(expr.T, initTrue)
		fInit, ef2 := ck.expr(expr.F, initFalse.cow())
		init = init.unionIntersect(tInit, fInit)
		effects = ef1 || ef2
	case *ast.Nary:
		if expr.Tok == tok.And || expr.Tok == tok.Or {
			init, effects = ck.expr(expr.Exprs[0], init) // first is always done
			in := init
			for _, e := range expr.Exprs[1:] {
				in, ef1 = ck.expr(e, in) // rest are conditional
				effects = effects || ef1
			}
		} else {
			expr.Children(func(e ast.Node) ast.Node {
				init, _ = ck.expr(e.(ast.Expr), init)
				return e
			})
		}
	case *ast.Call:
		effects = true
		expr.Children(func(e ast.Node) ast.Node {
			init, _ = ck.expr(e.(ast.Expr), init)
			return e
		})
	case *ast.Block:
		init = ck.block(expr, init)

	default:
		expr.Children(func(e ast.Node) ast.Node {
			init, _ = ck.expr(e.(ast.Expr), init)
			return e
		})
	}
	return init, effects
}

func (ck *Check) block(b *ast.Block, init set) set {
	// save & remove variables shadowed by params
	allInit := map[string]int{}
	allUsed := map[string]struct{}{}
	nUsedParams := 0
	for _, p := range b.Function.Params {
		id := p.Name.ParamName()
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
			id := p.Name.ParamName()
			if _, ok := ck.AllUsed[id]; !ok {
				ck.CheckResult(int(p.Name.Pos),
					"WARNING: initialized but not used: "+id)
			}
		}
	}

	// remove params
	for _, p := range b.Function.Params {
		id := p.Name.ParamName()
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
		ck.CheckResult(pos, p+" not initialized: "+id)
	}
	ck.AllUsed[id] = struct{}{}
	return init
}

// hasBreak checks if a statement contains any break statements
// that are not nested inside other loops.
func (ck *Check) hasBreak(stmt ast.Statement) bool {
	if stmt == nil {
		return false
	}
	switch stmt := stmt.(type) {
	case *ast.Break:
		return true
	case *ast.Compound:
		for _, s := range stmt.Body {
			if ck.hasBreak(s) {
				return true
			}
		}
	case *ast.If:
		return ck.hasBreak(stmt.Then) || ck.hasBreak(stmt.Else)
	case *ast.TryCatch:
		return ck.hasBreak(stmt.Try) || ck.hasBreak(stmt.Catch)
	case *ast.Switch:
		for _, c := range stmt.Cases {
			for _, s := range c.Body {
				if ck.hasBreak(s) {
					return true
				}
			}
		}
		if stmt.Default != nil {
			for _, s := range stmt.Default {
				if ck.hasBreak(s) {
					return true
				}
			}
		}
	// Don't recurse into nested loops - breaks in nested loops don't affect the outer loop
	case *ast.For, *ast.Forever, *ast.While, *ast.DoWhile, *ast.ForIn:
		return false
	}
	return false
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
			ck.CheckResult(at, "WARNING: initialized but not used: "+id)
		}
	}
	for id, pos := range ck.AllInit {
		if _, ok := ck.AllUsed[id]; !ok && !init.has(id) {
			ck.CheckResult(pos, "WARNING: initialized but not used: "+id)
		}
	}
}
func paramPos(params []ast.Param, id string) int {
	for _, p := range params {
		name := p.Name.Name
		if name == id ||
			((name[0] == '_' || name[0] == '@') && name[1:] == id) {
			return int(p.Name.Pos)
		}
	}
	return -1
}

//-------------------------------------------------------------------

func (ck *Check) CheckResult(pos int, str string) {
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
