// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
)

func GoGen(src string) string {
	p := GogenParser(src)
	f := p.constant().(*SuFunc)
	if p.Token != tok.Eof {
		p.Error("did not consume all input")
	}
	return f.Code
}

// gogen compiles an ast.Function to Go source code placed in SuFunc.Code.
// Using SuFunc for output is for compatibility with byte code codegen.
func (*gogenAspects) codegen(_, _ string, f *ast.Function, _ Value) Value {
	if len(f.Final) > 0 {
		ast.PropFold(f)
	}
	var g ggen
	g.locals = make(map[string]struct{})
	g.function(f)
	init := ""
	if g.init.Len() > 0 {
		init = g.init.String()
	}
	return &SuFunc{Code: init + g.String()}
}

type ggen struct {
	str.Builder
	next   int
	init   strings.Builder
	locals map[string]struct{}
}

func (g *ggen) function(fn *ast.Function) {
	g.params(fn.Params)
	g.Add("{\n")
	stmts := fn.Body
	for si, stmt := range stmts {
		g.statement(stmt, si == len(stmts)-1)
	}
	if len(stmts) == 0 {
		g.Add("return nil\n")
	}
	g.Add("}")
}

func (g *ggen) params(ps []ast.Param) {
	g.Add("func(")
	sep := ""
	for _, p := range ps {
		g.Adds(sep, p.Name.Name)
		g.locals[p.Name.Name] = struct{}{}
		sep = ", "
	}
	if len(ps) > 0 {
		g.Add(" Value")
	}
	g.Add(") Value ")
}

func (g *ggen) statement(node ast.Statement, lastStmt bool) {
	if node == nil {
		return
	}
	switch node := node.(type) {
	case *ast.Compound:
		g.statements(node.Body)
		if lastStmt {
			g.Add("return nil\n")
		}
		return
	case *ast.Return:
		g.returnStmt(node.E)
		g.Add("\n")
		return
	case *ast.If:
		g.ifStmt(node)
	case *ast.Forever:
		g.foreverStmt(node)
	case *ast.While:
		g.whileStmt(node)
	case *ast.DoWhile:
		g.dowhileStmt(node)
	case *ast.For:
		g.forStmt(node)
	case *ast.ForIn:
		g.forinStmt(node)
	case *ast.Break:
		g.Add("break")
	case *ast.Continue:
		g.Add("continue")
	case *ast.Throw:
		g.throwStmt(node)
		g.Add("\n")
		return
	case *ast.TryCatch:
		g.trycatchStmt(node)
	case *ast.ExprStmt:
		g.exprStmt(node.E, lastStmt)
		g.Add("\n")
		return
	default:
		panic("unhandled statement " + fmt.Sprintf("%T", node))
	}
	g.Add("\n")
	if lastStmt {
		g.Add("return nil\n")
	}
}

func (g *ggen) statements(stmts []ast.Statement) {
	for _, stmt := range stmts {
		g.statement(stmt, false)
	}
}

func (g *ggen) returnStmt(expr ast.Expr) {
	g.Add("return ")
	if expr == nil {
		g.Add("nil")
	} else {
		g.expr(expr, suvalue)
	}
}

func (g *ggen) ifStmt(node *ast.If) {
	g.Add("if ")
	g.expr(node.Cond, gobool)
	g.Add(" {\n")
	g.statement(node.Then, false)
	if node.Else == nil {
		g.Add("}")
	} else {
		g.Add("} else {\n")
		g.statement(node.Else, false)
		g.Add("}")
	}
}

func (g *ggen) foreverStmt(node *ast.Forever) {
	g.Add("for {\n")
	g.statement(node.Body, false)
	g.Add("}")
}

func (g *ggen) whileStmt(node *ast.While) {
	g.Add("for ")
	g.expr(node.Cond, gobool)
	g.Add(" {\n")
	g.statement(node.Body, false)
	g.Add("}")
}

func (g *ggen) dowhileStmt(node *ast.DoWhile) {
	g.Add("for {\n")
	g.statement(node.Body, false)
	g.Add("if !")
	g.expr(node.Cond, gobool)
	g.Add(" { break }\n}")
}

func (g *ggen) forStmt(node *ast.For) {
	// init can be outside the for statement
	for _, e := range node.Init {
		g.expr(e, void)
		g.Add("\n")
	}
	g.Add("for ; ")
	if node.Cond != nil {
		g.expr(node.Cond, gobool)
	}
	g.Add("; ")
	// put increment in for statement so continue works
	if len(node.Inc) == 1 {
		g.expr(node.Inc[0], void)
	} else if len(node.Inc) > 1 {
		g.Add("func(){ ")
		sep := ""
		for _, e := range node.Inc {
			g.Add(sep)
			sep = "; "
			g.expr(e, void)
		}
		g.Add(" }()")
	}
	g.Add(" {\n")
	g.statement(node.Body, false)
	g.Add("}")
}

func (g *ggen) forinStmt(node *ast.ForIn) {
	v := node.Var.Name
	if _, ok := g.locals[v]; !ok {
		g.Adds("var ", v, " Value\n")
	}
	g.Add("for _it_ := OpIter(")
	g.expr(node.E, suvalue)
	g.Adds("); ; {\n",
		v, " = _it_.Next()\n"+
			"if ", v, " == nil { break }\n")
	g.statement(node.Body, false)
	g.Add("}")
}

func (g *ggen) throwStmt(node *ast.Throw) {
	g.Add("panic(")
	g.expr(node.E, suvalue)
	g.Add(")")
}

func (g *ggen) trycatchStmt(node *ast.TryCatch) {
	g.Add("func() {\n" +
		"defer func() {\n" +
		"if _e_ := recover(); _e_ != nil {\n")
	if node.CatchVar.Name != "" {
		g.Adds(node.CatchVar.Name, " = ") //TODO := declare var ?
	}
	g.Add(fmt.Sprintf("OpCatch(t, _e_, %q)\n", node.CatchFilter))
	g.statement(node.Catch, false)
	g.Add("}\n}()\n")
	g.statement(node.Try, false)
	g.Add("}()") // end outer func
}

func (g *ggen) exprStmt(expr ast.Expr, lastStmt bool) {
	if lastStmt {
		g.Add("return ")
		g.expr(expr, suvalue)
	} else {
		if _, ok := expr.(*ast.Constant); ok {
			return
		}
		g.expr(expr, void)
	}
}

// expressions -----------------------------------------------------------------

type resultType int

const (
	void resultType = iota
	suvalue
	gobool
)

func (g *ggen) expr(node ast.Expr, want resultType) {
	result := suvalue
	pos := g.Len()
	switch node := node.(type) {
	case *ast.Constant:
		g.value(node.Val)
	case *ast.Ident:
		g.ident(node)
	case *ast.Unary:
		result = g.unary(node, want)
	case *ast.Binary:
		result = g.binary(node, want)
	case *ast.Nary:
		result = g.nary(node, want)
	case *ast.Trinary:
		g.trinary(node, want)
	default:
		panic("unhandled expression: " + node.String())
	}
	if result == suvalue && want == gobool {
		g.Insert(pos, "OpBool(")
		g.Add(")")
	} else if result == gobool && want == suvalue {
		g.Insert(pos, "SuBool(")
		g.Add(")")
	} else if want != void && result != want {
		panic("didn't get desired expression type")
	}
}

func (g *ggen) ident(node *ast.Ident) {
	c := node.Name[0]
	if ascii.IsLower(c) {
		g.Add(node.Name) // local
	} else {
		panic("unhandled identifier " + node.Name)
	}
}

func (g *ggen) value(val Value) {
	if val == True {
		g.Add("True")
	} else if val == False {
		g.Add("False")
	} else if val == Zero {
		g.Add("Zero")
	} else if val == One {
		g.Add("One")
	} else if val == MinusOne {
		g.Add("MinusOne")
	} else if val == EmptyStr {
		g.Add("EmptyStr")
	} else if f, ok := val.(*SuFunc); ok {
		g.Add(f.Code) //TODO needs to be SuBuiltin with ParamSpec
	} else {
		g.Add(g.pack64(val))
	}
}

func (g *ggen) pack64(v Value) string {
	name := fmt.Sprintf("_c%d_", g.next)
	g.next++
	data := Pack2(v.(Packable)).Buffer()
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(buf, data)
	g.init.WriteString("var " + name + " = Unpack64(`")
	for ; len(buf) >= 64; buf = buf[64:] {
		g.init.Write(buf[:64])
		g.init.WriteByte('\n')
	}
	if len(buf) > 0 {
		g.init.Write(buf)
		g.init.WriteString("`)\n")
	}
	return name
}

var cOne = &ast.Constant{Val: One}

func (g *ggen) unary(node *ast.Unary, want resultType) resultType {
	switch node.Tok {
	case tok.LParen:
		g.expr(node.E, want)
		return want
	case tok.Not:
		if want == gobool {
			g.Add("!")
			g.expr(node.E, gobool)
			return gobool
		}
		fallthrough
	case tok.Add, tok.Sub:
		g.Add(uopfn[node.Tok])
		g.expr(node.E, want)
		g.Add(")")
		return suvalue
	case tok.Div:
		g.Add("OpDiv(One, ")
		g.expr(node.E, suvalue)
		g.Add(")")
		return suvalue
	case tok.PostInc:
		if want == suvalue {
			return g.postIncDec(node.E, "OpAdd")
		}
		fallthrough
	case tok.Inc:
		return g.opeq(node.E, cOne, tok.AddEq, want)
	case tok.PostDec:
		if want == suvalue {
			return g.postIncDec(node.E, "OpSub")
		}
		fallthrough
	case tok.Dec:
		return g.opeq(node.E, cOne, tok.SubEq, want)
	}
	panic("unhandled unary " + node.String())
}

func (g *ggen) postIncDec(x ast.Expr, op string) resultType {
	if id := g.localVar(x); id != "" {
		g.Adds("func(){ _r_ := ", id, "; ",
			id, " = ", op, "(_r_, One); "+
				"return _r_ }()")
		return suvalue
	}
	// else member
	lhs := x.(*ast.Mem)
	g.expr(lhs.E, suvalue)
	g.Add(".GetPut(t, ")
	g.expr(lhs.M, suvalue)
	g.Adds(", One, ", op, ", true)") // true means return original value
	return suvalue
}

var uopfn = map[tok.Token]string{
	tok.Add:    "OpUnaryPlus(",
	tok.Sub:    "OpUnaryMinus(",
	tok.Not:    "OpNot(",
	tok.BitNot: "OpBitNot(",
}

func (g *ggen) binary(node *ast.Binary, want resultType) resultType {
	switch node.Tok {
	case tok.Eq:
		if want == void {
			g.lvalue(node.Lhs)
			g.expr(node.Rhs, suvalue)
			g.Add(")")
			return void
		}
		g.Add("func(){ _r_ := ")
		g.expr(node.Rhs, suvalue)
		g.Add("; ")
		g.lvalue(node.Lhs)
		g.Add("_r_); return _r_ }()")
		return suvalue
	case tok.AddEq, tok.SubEq, tok.CatEq, tok.MulEq, tok.DivEq, tok.ModEq,
		tok.LShiftEq, tok.RShiftEq, tok.BitOrEq, tok.BitAndEq, tok.BitXorEq:
		return g.opeq(node.Lhs, node.Rhs, node.Tok, want)
	case tok.Is, tok.Isnt, tok.Lt, tok.Lte, tok.Gt, tok.Gte:
		s2 := cmpfn[node.Tok]
		g.Add("(")
		g.expr(node.Lhs, suvalue)
		g.Add(s2.s)
		g.expr(node.Rhs, suvalue)
		g.Adds(s2.t, ")")
		return gobool
	case tok.Match, tok.MatchNot, tok.Mod, tok.LShift, tok.RShift:
		g.Add(opfn[node.Tok])
		g.expr(node.Lhs, suvalue)
		g.Add(", ")
		g.expr(node.Rhs, suvalue)
		g.Add(")")
		return suvalue
	}
	panic("unhandled binary " + node.Tok.String())
}

func (g *ggen) opeq(x, y ast.Expr, op tok.Token, want resultType) resultType {
	if id := g.localVar(x); id != "" {
		if want == void {
			g.Adds(id, " = ", opfn[op], "(", id, ", ")
			g.expr(y, suvalue)
			g.Add(")")
			return void
		}
		g.Adds("func(){ _r_ := ", opfn[op], "(")
		g.expr(x, suvalue)
		g.Add(", ")
		g.expr(y, suvalue)
		g.Adds("); ", id, " = _r_; return _r_ }()")
		return suvalue
	}
	// else member
	lhs := x.(*ast.Mem)
	g.expr(lhs.E, suvalue)
	g.Add(".GetPut(t, ")
	g.expr(lhs.M, suvalue)
	g.Add(", ")
	g.expr(y, suvalue)
	g.Adds(", ", opfn[op], ", false)") // false means return new result
	return suvalue
}

func (g *ggen) localVar(node ast.Expr) string {
	if ident, ok := node.(*ast.Ident); ok && isLocal(ident.Name) {
		return ident.Name
	}
	return ""
}

var opfn = map[tok.Token]string{
	tok.Add:      "OpAdd(",
	tok.Sub:      "OpSub(",
	tok.Cat:      "OpCat(",
	tok.BitAnd:   "OpBitAnd(",
	tok.BitOr:    "OpBitOr(",
	tok.BitXor:   "OpBitXor(",
	tok.And:      " && ",
	tok.Or:       " || ",
	tok.Match:    "OpMatch(t, ",
	tok.MatchNot: "!OpMatch(t, ",
	tok.Mod:      "OpMod(",
	tok.LShift:   "OpLShift(",
	tok.RShift:   "OpRShift(",
	tok.AddEq:    "OpAdd",
	tok.SubEq:    "OpSub",
	tok.CatEq:    "t.Cat",
	tok.MulEq:    "OpMul",
	tok.DivEq:    "OpDiv",
	tok.ModEq:    "OpMod",
	tok.LShiftEq: "OpLShift",
	tok.RShiftEq: "OpRShift",
	tok.BitAndEq: "OpBitAndEq",
	tok.BitOrEq:  "OpBitOrEq",
	tok.BitXorEq: "OpBitXorEq",
}

type string2 struct {
	s, t string
}

var cmpfn = map[tok.Token]string2{
	tok.Is:   {".Equal(", ")"},
	tok.Isnt: {".Equal(", ") != true"},
	tok.Lt:   {".Compare(", ") < 0"},
	tok.Lte:  {".Compare(", ") <= 0"},
	tok.Gt:   {".Compare(", ") > 0"},
	tok.Gte:  {".Compare(", ") >= 0"},
}

// lvalue caller must add closing parenthesis
func (g *ggen) lvalue(node ast.Expr) {
	switch node := node.(type) {
	case *ast.Ident:
		if isLocal(node.Name) {
			g.Adds("(", node.Name)
			if _, ok := g.locals[node.Name]; ok {
				g.Add(" = ")
			} else {
				g.Add(" := ")
				g.locals[node.Name] = struct{}{}
			}
			return
		}
	case *ast.Mem:
		g.expr(node.E, suvalue)
		g.Add(".Put(")
		g.expr(node.M, suvalue)
		g.Add(", ")
		return
	}
	panic("unhandled lvalue: " + node.String())
}

func (g *ggen) nary(node *ast.Nary, want resultType) resultType {
	if node.Tok == tok.And || node.Tok == tok.Or {
		g.andorExpr(node, want)
		return want
	}
	if node.Tok == tok.Mul {
		g.muldivExpr(node)
		return suvalue
	}
	// else Add, Sub, Cat, BitOr, BitAnd, BitXor
	// left associative, so work backwards to generate operations
	for i := len(node.Exprs) - 1; i > 0; i-- {
		if isUnary(node.Exprs[i], tok.Sub) {
			g.Add("OpSub(")
		} else {
			g.Add("OpAdd(")
		}
	}
	// now work forwards to generate operands
	g.expr(node.Exprs[0], suvalue)
	for _, e := range node.Exprs[1:] {
		g.Add(", ")
		if u, ok := e.(*ast.Unary); ok && u.Tok == tok.Sub {
			e = u.E
		}
		g.expr(e, suvalue)
		g.Add(")")
	}
	return suvalue
}

func (g *ggen) muldivExpr(node *ast.Nary) {
	var top []ast.Expr
	var bot []ast.Expr
	for _, e := range node.Exprs {
		if isUnary(e, tok.Div) {
			bot = append(bot, e.(*ast.Unary).E)
		} else {
			top = append(top, e)
		}
	}
	if len(bot) > 0 {
		g.Add("OpDiv(")
	}
	g.mul(top)
	if len(bot) > 0 {
		g.Add(", ")
		g.mul(bot)
		g.Add(")")
	}
}

func (g *ggen) mul(exprs []ast.Expr) {
	// left associative, so work backwards to generate operations
	for i := len(exprs) - 1; i > 0; i-- {
		g.Add("OpMul(")
	}
	// now work forwards to generate operands
	g.expr(exprs[0], suvalue)
	for _, e := range exprs[1:] {
		g.Add(", ")
		g.expr(e, suvalue)
		g.Add(")")
	}
}

func (g *ggen) andorExpr(node *ast.Nary, want resultType) {
	if want == suvalue {
		g.Add("SuBool")
	}
	g.Add("(")
	g.expr(node.Exprs[0], gobool)
	for _, e := range node.Exprs[1:] {
		g.Add(opfn[node.Tok])
		g.expr(e, gobool)
	}
	g.Add(")")
}

func (g *ggen) trinary(node *ast.Trinary, want resultType) resultType {
	if want == void {
		g.Add("if ")
		g.expr(node.Cond, gobool)
		g.Add(" { ")
		g.expr(node.T, void)
		g.Add(" } else { ")
		g.expr(node.F, void)
		g.Add(" }")
		return void
	}
	g.Add("func() { if ")
	g.expr(node.Cond, gobool)
	g.Add(" { return ")
	g.expr(node.T, void)
	g.Add(" } else { return ")
	g.expr(node.F, void)
	g.Add(" } }()")
	return suvalue
}
