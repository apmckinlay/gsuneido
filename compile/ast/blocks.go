// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

// scope tracks block scope information during Blocks processing
type scope struct {
	parent   *scope    // nil for outer function
	block    *Block    // nil for outer function (scopes[0])
	fn       *Function // the function or block's function
	hasThis  bool
	hasSuper bool
	hasRet   bool
}

// scopes manages the collection of scopes during Blocks processing
type scopes struct {
	cur    *scope
	scopes []*scope // scopes[0] is the outer function
}

// Blocks sets CompileAsFunction if a block can be compiled as a function
// i.e. isn't a closure, doesn't share any variables.
// It also assigns slot indexes to variables:
//   - Local stack frame slots:  < SharedSlotStart
//   - Shared (heap) slots: >= SharedSlotStart
//
// NOTE: This is trickier than it seems.
// due to e.g. nested blocks and parameter shadowing.
// Blocks does not process nested *functions* (they're already codegen and not Ast);
// they are handled as constructed bottom up.
func Blocks(f *Function) {
	var ss scopes
	// Phase 1: Collect scopes and variables
	ss.collect(f)
	// Phase 2: Assign shared slots >= SharedSlotStart and set CompileAsFunction
	ss.assignShared()
	// Phase 3: Assign local slots < SharedSlotStart
	ss.assignLocal()
}

// collect scopes and variables.
func (ss *scopes) collect(f *Function) {
	// Create scope for outer function (scopes[0])
	ss.cur = &scope{block: nil, parent: nil, fn: f}
	ss.scopes = append(ss.scopes, ss.cur)
	// Initialize Vars on the outer function
	f.Vars = make(map[string]int16)
	// Add parameters with slot indexes starting at 0
	AddParams(f.Params, f.Vars)
	// Collect variables from the function body
	vars := f.Vars
	for _, stmt := range f.Body {
		ss.statement(stmt, vars, f)
	}
}

// AddParams adds parameters to vars
// assigning slot indexes starting at 0
func AddParams(params []Param, vars map[string]int16) {
	for i, p := range params {
		if i >= SharedSlotStart {
			panic("too many local variables")
		}
		vars[p.Name.ParamName()] = int16(i)
	}
}

// assignShared slots (>= SharedSlotStart)
// Variables shared between a block and its parent chain get shared slots
// Also sets CompileAsFunction
// NOTE: assignShared requires scopes to be top down (outer first)
func (ss *scopes) assignShared() {
	sharedIdx := int16(SharedSlotStart) // shared slots start at SharedSlotStart
	// For each block scope (skip scopes[0] which is the outer function)
	for _, s := range ss.scopes[1:] {
		if s.block == nil {
			continue
		}
		blockVars := s.block.Function.Vars
		hasShared := false
		// Check each variable in the block
		for vname, vidx := range blockVars {
			// Skip parameters (already assigned starting at 0)
			if vidx >= 0 && vidx < SharedSlotStart {
				continue
			}
			// Dynamic vars (_name) are never shared across scopes.
			if vname != "" && vname[0] == '_' {
				continue
			}
			// Walk up the parent chain to find shared variables
			// but do NOT go past a parameter in the CURRENT scope
			// that shadows an outer variable with the same name
			// Note: we check s (current scope), not parent, because we want to
			// stop if the variable we're searching for is a parameter in the
			// block that captured it (creating a new binding)
			for parent := s.parent; parent != nil; parent = parent.parent {
				parentVars := parent.fn.Vars
				if parentVidx, ok := parentVars[vname]; ok {
					// Variable exists in parent
					if parentVidx >= SharedSlotStart {
						// Already assigned as shared, use same slot
						blockVars[vname] = parentVidx
						hasShared = true
					} else if parentVidx >= 0 {
						// Parameter in parent - make it shared
						if sharedIdx > 255 {
							panic("too many shared variables")
						}
						blockVars[vname] = sharedIdx
						parentVars[vname] = sharedIdx
						sharedIdx++
						hasShared = true
					} else {
						// Variable in parent is unassigned - make it shared
						if sharedIdx > 255 {
							panic("too many shared variables")
						}
						blockVars[vname] = sharedIdx
						parentVars[vname] = sharedIdx
						sharedIdx++
						hasShared = true
					}
					break // found in parent, no need to continue up chain
				}
			}
		}
		// Block can be compiled as function if:
		// - no shared variables
		// - no this/super references
		// - no return statement
		// Only set to false (never back to true, since inner blocks may have set it)
		if hasShared || s.hasThis || s.hasSuper || s.hasRet {
			s.block.CompileAsFunction = false
			// Parents must be closures too
			for parent := s.parent; parent != nil; parent = parent.parent {
				if parent.block != nil {
					if !parent.block.CompileAsFunction {
						break // already false, no need to continue
					}
					parent.block.CompileAsFunction = false
				}
			}
		}
	}
	root := ss.scopes[0].fn
	root.SharedCount = sharedIdx - SharedSlotStart
}

// assignLocal slots (< SharedSlotStart)
// Unshared variables get local slots starting at len(params)
func (ss *scopes) assignLocal() {
	for _, s := range ss.scopes {
		vars := s.fn.Vars
		nextSlot := int16(len(s.fn.Params))
		for vname, vidx := range vars {
			if vidx == -1 { // unassigned
				if nextSlot >= SharedSlotStart {
					panic("too many local variables")
				}
				vars[vname] = nextSlot
				nextSlot++
			}
		}
	}
}

// ast traversal (used by scopes.collect) ---------------------------

// statement processes one statement (and its children)
func (ss *scopes) statement(stmt Statement, vars map[string]int16, fn *Function) {
	if stmt == nil {
		return
	}
	switch stmt := stmt.(type) {
	case *Compound:
		for _, stmt := range stmt.Body {
			ss.statement(stmt, vars, fn)
		}
	case *ExprStmt:
		ss.expr(stmt.E, vars, fn)
	case *Return:
		if ss.cur != nil && ss.cur.block != nil {
			ss.cur.hasRet = true
		}
		for _, expr := range stmt.Exprs {
			ss.expr(expr, vars, fn)
		}
	case *MultiAssign:
		// see expr Binary Eq
		for _, expr := range stmt.Lhs {
			id := expr.(*Ident)
			ss.addVar(id.Name, vars)
		}
		ss.expr(stmt.Rhs, vars, fn)
	case *Throw:
		ss.expr(stmt.E, vars, fn)
	case *TryCatch:
		ss.statement(stmt.Try, vars, fn)
		if stmt.CatchVar.Name != "" {
			ss.addVar(stmt.CatchVar.Name, vars)
		}
		ss.statement(stmt.Catch, vars, fn)
	case *While:
		ss.expr(stmt.Cond, vars, fn)
		ss.statement(stmt.Body, vars, fn)
	case *Forever:
		ss.statement(stmt.Body, vars, fn)
	case *DoWhile:
		ss.statement(stmt.Body, vars, fn)
		ss.expr(stmt.Cond, vars, fn)
	case *If:
		ss.expr(stmt.Cond, vars, fn)
		ss.statement(stmt.Then, vars, fn)
		ss.statement(stmt.Else, vars, fn)
	case *Switch:
		ss.expr(stmt.E, vars, fn)
		for _, c := range stmt.Cases {
			for _, e := range c.Exprs {
				ss.expr(e, vars, fn)
			}
			for _, stmt := range c.Body {
				ss.statement(stmt, vars, fn)
			}
		}
		for _, d := range stmt.Default {
			ss.statement(d, vars, fn)
		}
	case *ForIn:
		if stmt.Var.Name != "" {
			ss.addVar(stmt.Var.Name, vars)
		}
		if stmt.Var2.Name != "" {
			ss.addVar(stmt.Var2.Name, vars)
		}
		ss.expr(stmt.E, vars, fn)
		ss.expr(stmt.E2, vars, fn)
		ss.statement(stmt.Body, vars, fn)
	case *For:
		for _, expr := range stmt.Init {
			ss.expr(expr, vars, fn)
		}
		ss.expr(stmt.Cond, vars, fn)
		ss.statement(stmt.Body, vars, fn)
		for _, expr := range stmt.Inc {
			ss.expr(expr, vars, fn)
		}
	case *Break, *Continue:
		// nothing to do
	default:
		panic("unexpected statement type " + fmt.Sprintf("%T", stmt))
	}
}

func (ss *scopes) expr(expr Expr, vars map[string]int16, fn *Function) {
	if expr == nil {
		return
	}
	switch expr := expr.(type) {
	case *Binary:
		if expr.Tok == tok.Eq {
			if id, ok := expr.Lhs.(*Ident); ok {
				// assignment
				ss.addVar(id.Name, vars)
				ss.expr(expr.Rhs, vars, fn)
				break
			}
		}
		ss.expr(expr.Lhs, vars, fn)
		ss.expr(expr.Rhs, vars, fn)
	case *Ident:
		if expr.Name == "this" {
			if ss.cur != nil && ss.cur.block != nil {
				ss.cur.hasThis = true
			}
			break
		}
		if expr.Name == "super" {
			if ss.cur != nil && ss.cur.block != nil {
				ss.cur.hasSuper = true
			}
			break
		}
		if isLocalVarName(expr.Name) {
			ss.addVar(expr.Name, vars)
		}
	case *Trinary:
		ss.expr(expr.Cond, vars, fn)
		ss.expr(expr.T, vars, fn)
		ss.expr(expr.F, vars, fn)
	case *Nary:
		if expr.Tok == tok.And || expr.Tok == tok.Or {
			ss.expr(expr.Exprs[0], vars, fn) // first is always done
			for _, e := range expr.Exprs[1:] {
				ss.expr(e, vars, fn) // rest are conditional
			}
		} else {
			expr.Children(func(e Node) Node {
				ss.expr(e.(Expr), vars, fn)
				return e
			})
		}
	case *Block:
		ss.block(expr)
	default:
		expr.Children(func(e Node) Node {
			ss.expr(e.(Expr), vars, fn)
			return e
		})
	}
}

// addVar adds a variable to the vars map if not already present
// Variables are initialized with slot index -1 (unassigned)
func (ss *scopes) addVar(name string, vars map[string]int16) {
	if _, ok := vars[name]; !ok {
		vars[name] = -1 // unassigned
	}
}

func isLocalVarName(name string) bool {
	return ascii.IsLower(name[0]) ||
		(name[0] == '_' && len(name) > 1 && ascii.IsLower(name[1]))
}

func (ss *scopes) block(block *Block) {
	parent := ss.cur
	ss.cur = &scope{block: block, parent: parent, fn: &block.Function}
	// append before recursing so scopes are top down
	ss.scopes = append(ss.scopes, ss.cur)
	block.Function.Vars = make(map[string]int16)
	block.CompileAsFunction = true
	AddParams(block.Params, block.Function.Vars)
	blockVars := block.Function.Vars
	for _, stmt := range block.Body {
		ss.statement(stmt, blockVars, &block.Function)
	}
	ss.cur = parent
}
