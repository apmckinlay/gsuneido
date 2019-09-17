// Package check processes an ast.Function
// and finds used but (possibly) not initialized,
// and initialized but (possibly) not used.
// "possibly" meaning not on all code paths.
// Does not check nested functions (they're already codegen and not Ast)
// they are checked as constructed bottom up.
package check

// TODO if/for/while where condition is: x and y
// should know y will always run for the body
// e.g. if i > 0 and false isnt x = Next() { ... x ...}

// TODO add a way to indicate that a block is always called e.g. Transaction
// to avoid spurious "possibly not initialized"

import (
	"fmt"
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
	AllInit map[string]int32
	// AllUsed is the set of variables read from, including conditionally
	AllUsed map[string]struct{}
	Results []string
}

func New(t *Thread) *Check {
	return &Check{t: t}
}

func (ck *Check) Check(f *ast.Function) set {
	ck.AllInit = make(map[string]int32)
	ck.AllUsed = make(map[string]struct{})
	var init set
	init = ck.check(f, init)
	for i, id := range init {
		if !ck.used(id) {
			at := ""
			if i < len(f.Params) {
				at = " @" + strconv.Itoa(int(f.Pos))
			} else if pos, ok := ck.AllInit[id]; ok {
				at = " @" + strconv.Itoa(int(pos))
			}
			ck.Results = append(ck.Results,
				"WARNING: initialized but not used: "+id+at)
		}
	}

	// for _, s := range ck.Results {
	// 	fmt.Println(s)
	// }
	return init
}

func (ck *Check) check(f *ast.Function, init set) set {
	init = ck.params(f.Params, init)
	for _, stmt := range f.Body {
		init = ck.statement(stmt, init)
	}
	return init
}

func (ck *Check) params(params []ast.Param, init set) set {
	for _, p := range params {
		name := p.Name
		if name[0] == '.' {
			name = str.UnCapitalize(name[1:])
			ck.AllUsed[name] = struct{}{}
		} else if name[0] == '@' || name[0] == '_' {
			name = name[1:]
		}
		if !p.Unused {
			init = init.with(name)
		}
	}
	return init
}

func (ck *Check) used(id string) bool {
	if id == "unused" {
		return true
	}
	_, ok := ck.AllUsed[id]
	return ok
}

// statement processes one statement (and its children)
// Conditional statements are assumed to run for used, and not to run for init.
// So we accumulate used, but not init.
func (ck *Check) statement(stmt ast.Statement, init set) set {
	if stmt == nil {
		return init
	}
	ck.pos = stmt.Position()
	switch stmt := stmt.(type) {
	case *ast.Compound:
		for _, stmt := range stmt.Body {
			init = ck.statement(stmt, init)
		}
	case *ast.ExprStmt:
		init = ck.expr(stmt.E, init)
	case *ast.Return:
		init = ck.expr(stmt.E, init)
	case *ast.Throw:
		init = ck.expr(stmt.E, init)
	case *ast.TryCatch:
		init = ck.statement(stmt.Try, init)
		if stmt.CatchVar != "" {
			init = ck.initVar(init, stmt.CatchVar)
		}
		ck.statement(stmt.Catch, init)
	case *ast.While:
		init = ck.expr(stmt.Cond, init)
		ck.statement(stmt.Body, init)
	case *ast.Forever:
		init = ck.statement(stmt.Body, init)
	case *ast.DoWhile:
		init = ck.statement(stmt.Body, init)
		init = ck.expr(stmt.Cond, init)
	case *ast.If:
		init = ck.expr(stmt.Cond, init)
		thenInit := ck.statement(stmt.Then, init)
		elseInit := ck.statement(stmt.Else, init)
		init = init.unionIntersect(thenInit, elseInit)
	case *ast.Switch:
		init = ck.expr(stmt.E, init)
		for _, c := range stmt.Cases {
			in := init
			for _, e := range c.Exprs {
				in = ck.expr(e, in)
			}
			for _, b := range c.Body {
				in = ck.statement(b, in)
			}
		}
		in := init
		for _, d := range stmt.Default {
			in = ck.statement(d, in)
		}
	case *ast.ForIn:
		init = ck.initVar(init, stmt.Var)
		init = ck.expr(stmt.E, init)
		ck.statement(stmt.Body, init)
	case *ast.For:
		for _, expr := range stmt.Init {
			init = ck.expr(expr, init)
		}
		init = ck.expr(stmt.Cond, init)
		ck.statement(stmt.Body, init)
		for _, expr := range stmt.Inc {
			ck.expr(expr, init)
		}
	case *ast.Break, *ast.Continue:
		// nothing to do
	default:
		panic("unexpected statement type " + fmt.Sprintf("%T", stmt))
	}
	return init
}

func (ck *Check) expr(expr ast.Expr, init set) set {
	if expr == nil {
		return init
	}
	switch expr := expr.(type) {
	case *ast.Binary:
		if expr.Tok == tok.Eq {
			if id, ok := expr.Lhs.(*ast.Ident); ok {
				init = ck.initVar(init, id.Name)
				init = ck.expr(expr.Rhs, init)
				break
			}
		}
		init = ck.expr(expr.Lhs, init)
		init = ck.expr(expr.Rhs, init)
	case *ast.Ident:
		if ascii.IsLower(expr.Name[0]) {
			init = ck.usedVar(init, expr.Name)
		}
		if ascii.IsUpper(expr.Name[0]) {
			if nil == Global.FindName(ck.t, expr.Name) {
				ck.Results = append(ck.Results,
					"WARNING: can't find: "+expr.Name+" @"+strconv.Itoa(ck.pos))
			}
		}
	case *ast.Trinary:
		init = ck.expr(expr.Cond, init)
		tInit := ck.expr(expr.T, init)
		fInit := ck.expr(expr.F, init)
		init = init.unionIntersect(tInit, fInit)
	case *ast.Nary:
		if expr.Tok == tok.And || expr.Tok == tok.Or {
			init = ck.expr(expr.Exprs[0], init) // first is always done
			in := init
			for _, e := range expr.Exprs[1:] {
				in = ck.expr(e, in) // rest are conditional
			}
		} else {
			expr.Children(func(e ast.Node) {
				init = ck.expr(e.(ast.Expr), init)
			})
		}
	case *ast.Block:
		ck.check(&expr.Function, init.with("it"))
	default:
		expr.Children(func(e ast.Node) {
			if e != nil {
				init = ck.expr(e.(ast.Expr), init)
			}
		})
	}
	return init
}

func (ck *Check) initVar(init set, id string) set {
	if strings.HasPrefix(id, "_") {
		return init
	}
	ck.AllInit[id] = int32(ck.pos)
	return init.with(id)
}

func (ck *Check) usedVar(init set, id string) set {
	if strings.HasPrefix(id, "_") {
		return init
	}
	if id != "this" && id != "super" && !init.has(id) {
		p := "ERROR: used but"
		if _, ok := ck.AllInit[id]; ok {
			p = "WARNING: used but possibly"
		}
		ck.Results = append(ck.Results,
			p+" not initialized: "+id+" @"+strconv.Itoa(ck.pos))
	}
	ck.AllUsed[id] = struct{}{}
	return init
}
