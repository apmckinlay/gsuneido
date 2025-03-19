// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ast defines the node types
// used by the compiler to build syntax trees
//
//	Expr
//		Ident
//		Constant, Symbol, Function, Block
//		Unary, Binary, Trinary, Nary, In
//		Mem, RangeTo, RangeLen
//		Call
//	Statement
//		ExprStmt
//		Compound
//		Return, Throw, Break, Continue
//		TryCatch
//		If, Switch
//		For, ForIn, Forever, While, DoWhile
package ast

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Node is embedded by Expr and Statement
type Node interface {
	astNode()
	String() string
	// Children calls the given function for each child node
	Children(func(Node) Node)
	// SetPos is called by the parser
	SetPos(org, end int32)
	// Get is for the Value interface for Suneido.Parse
	Get(*Thread, Value) Value
}

type astNodeT struct {
	SuAstNode
}

func (*astNodeT) astNode() {}

func (*astNodeT) Children(func(Node) Node) {
}

type noPos struct {
}

func (noPos) SetPos(org, end int32) {
}

type TwoPos struct {
	org int32
	end int32
}

func (a *TwoPos) SetPos(org, end int32) {
	a.org = org
	a.end = end
}

func (a *TwoPos) GetPos() int {
	return int(a.org)
}

func (a *TwoPos) GetEnd() int {
	return int(a.end)
}

// Expr is implemented by expression nodes
type Expr interface {
	Node
	exprNode()
	Echo() string
	// Eval, CanEvalRaw, and Columns are used by queries
	Eval(*Context) Value
	EvalRaw(*Context) string
	// CanEvalRaw sets the Expr so future Eval's will (or will not) be raw.
	// It must be called before calling Eval.
	// It is primarily for Where expressions.
	// CanEvalRaw should call CanEvalRaw on all its children.
	CanEvalRaw(flds []string) bool
	Columns() []string
}
type exprNodeT struct {
	astNodeT
}

func (*exprNodeT) exprNode() {}

func (en *exprNodeT) CanEvalRaw(flds []string) bool {
	en.Children(func(node Node) Node {
		node.(Expr).CanEvalRaw(flds)
		return node
	})
	return false
}

func (*exprNodeT) EvalRaw(*Context) string {
	panic(assert.ShouldNotReachHere())
}

func (en *exprNodeT) Echo() string {
	panic("Echo not implemented")
}

type Ident struct {
	exprNodeT
	Name     string
	Pos      int32 // for check errors
	Implicit bool  // for implicit Record, Object, this
}

func (a *Ident) String() string {
	return a.Name
}

func (a *Ident) Echo() string {
	return a.Name
}

func (a *Ident) ParamName() string {
	name := a.Name
	if name[0] == '.' {
		name = str.UnCapitalize(name[1:])
	} else if name[0] == '@' || name[0] == '_' {
		name = name[1:]
	}
	return name
}

func (a *Ident) SetPos(org, end int32) {
	a.Pos = org
}

func (a *Ident) GetPos() int {
	return int(a.Pos)
}

func (a *Ident) GetEnd() int {
	if a.Implicit {
		return int(a.Pos)
	}
	return int(a.Pos) + len(a.Name)
}

type Constant struct {
	exprNodeT
	Val Value
	// Packed is used for queries. It is set by Binary.CanEvalRaw
	Packed string
	TwoPos
}

func (a *Constant) String() string {
	return a.Val.String()
}

func (a *Constant) Echo() string {
	return a.Val.String()
}

func (a *Constant) SetPos(org, end int32) {
	a.TwoPos.SetPos(org, end)
	if x, ok := a.Val.(SetPosAble); ok {
		x.SetPos(org, end)
	}
}

type SetPosAble interface {
	SetPos(org, end int32)
}

type Symbol struct {
	Constant
}

type Unary struct {
	exprNodeT
	E Expr
	TwoPos
	Tok     tok.Token
	evalRaw bool
}

func (a *Unary) String() string {
	return "Unary(" + a.Tok.String() + " " + a.E.String() + ")"
}

func (a *Unary) Echo() string {
	if a.Tok == tok.LParen {
		return "(" + a.E.Echo() + ")"
	}
	if a.Tok == tok.Not {
		if in, ok := a.E.(*In); ok {
			return in.E.Echo() + " not" + in.echo()
		}
		if b, ok := a.E.(*Binary); ok {
			return "not (" + b.Echo() + ")"
		}
	}
	var op = map[tok.Token]string{
		tok.Add: "+",
		tok.Sub: "-",
		tok.Not: "not ",
	}
	return op[a.Tok] + a.E.Echo()
}

func (a *Unary) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
}

type Binary struct {
	exprNodeT
	noPos
	Lhs     Expr
	Rhs     Expr
	Tok     tok.Token
	evalRaw bool
}

func (a *Binary) String() string {
	return "Binary(" + a.Tok.String() + " " + a.Lhs.String() + " " + a.Rhs.String() + ")"
}

var tokEcho = map[tok.Token]string{
	tok.Is:       " is ",
	tok.Isnt:     " isnt ",
	tok.Lt:       " < ",
	tok.Lte:      " <= ",
	tok.Gt:       " > ",
	tok.Gte:      " >= ",
	tok.Match:    " =~ ",
	tok.MatchNot: " !~ ",
	tok.Add:      " + ",
	tok.Sub:      " - ",
	tok.Cat:      " $ ",
	tok.Mul:      " * ",
	tok.Div:      " / ",
	tok.Mod:      " % ",
	tok.And:      " and ",
	tok.Or:       " or ",
}

func (a *Binary) Echo() string {
	return a.Lhs.Echo() + tokEcho[a.Tok] + a.Rhs.Echo()
}

func childExpr(fn func(Node) Node, pexpr *Expr) {
	if *pexpr != nil {
		*pexpr = fn(*pexpr).(Expr)
	}
}

func (a *Binary) Children(fn func(Node) Node) {
	childExpr(fn, &a.Lhs)
	childExpr(fn, &a.Rhs)
}

type Trinary struct {
	exprNodeT
	noPos
	Cond    Expr
	T       Expr
	F       Expr
	evalRaw bool
}

func (a *Trinary) String() string {
	return "Trinary(" + a.Cond.String() + " " + a.T.String() + " " + a.F.String() + ")"
}

func (a *Trinary) Echo() string {
	return a.Cond.Echo() + " ? " + a.T.Echo() + " : " + a.F.Echo()
}

func (a *Trinary) Children(fn func(Node) Node) {
	childExpr(fn, &a.Cond)
	childExpr(fn, &a.T)
	childExpr(fn, &a.F)
}

// Nary is used for associative binary operators e.g. add, multiply, and, or
type Nary struct {
	exprNodeT
	noPos
	Exprs   []Expr
	Tok     tok.Token
	evalRaw bool
}

func (a *Nary) String() string {
	s := "Nary(" + a.Tok.String()
	for _, e := range a.Exprs {
		s += " " + e.String()
	}
	return s + ")"
}

func (a *Nary) Echo() string {
	if len(a.Exprs) == 0 {
		return ""
	}
	s := a.Exprs[0].Echo()
	for _, e := range a.Exprs[1:] {
		s += tokEcho[a.Tok] + e.Echo()
	}
	return s
}

func (a *Nary) Children(fn func(Node) Node) {
	for i := range a.Exprs {
		childExpr(fn, &a.Exprs[i])
	}
}

type RangeTo struct {
	exprNodeT
	noPos
	E    Expr
	From Expr
	To   Expr
}

func (a *RangeTo) String() string {
	return "RangeTo(" + a.E.String() + " " + fmt.Sprint(a.From) + " " +
		fmt.Sprint(a.To) + ")"
}

func (a *RangeTo) Echo() string {
	s := a.E.String() + "["
	if a.From != nil {
		s += a.From.Echo()
	}
	s += ".."
	if a.To != nil {
		s += a.To.Echo()
	}
	return s + "]"
}

func (a *RangeTo) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
	childExpr(fn, &a.From)
	childExpr(fn, &a.To)
}

type RangeLen struct {
	exprNodeT
	noPos
	E    Expr
	From Expr
	Len  Expr
}

func (a *RangeLen) String() string {
	return "RangeLen(" + a.E.String() + " " + fmt.Sprint(a.From) + " " +
		fmt.Sprint(a.Len) + ")"
}

func (a *RangeLen) Echo() string {
	return a.E.Echo() + "[" + a.From.Echo() + "::" + a.Len.Echo() + "]"
}

func (a *RangeLen) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
	childExpr(fn, &a.From)
	childExpr(fn, &a.Len)
}

type Mem struct {
	exprNodeT
	noPos
	E      Expr
	M      Expr
	DotPos int32
}

func (a *Mem) String() string {
	return "Mem(" + a.E.String() + " " + a.M.String() + ")"
}

func (a *Mem) Echo() string {
	s := a.E.Echo()
	if c, ok := a.M.(*Constant); ok {
		if cs, ok := c.Val.(SuStr); ok && lexer.IsIdentifier(string(cs)) {
			return s + "." + string(cs)
		}
	}
	return s + "[" + a.M.Echo() + "]"
}

func (a *Mem) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
	childExpr(fn, &a.M)
}

type In struct {
	exprNodeT
	noPos
	E       Expr
	Exprs   []Expr
	evalRaw bool
}

func (a *In) String() string {
	s := "In(" + a.E.String() + " ["
	sep := ""
	for _, e := range a.Exprs {
		s += sep + e.String()
		sep = " "
	}
	return s + "])"
}

func (a *In) Echo() string {
	return a.E.Echo() + a.echo()
}

func (a *In) echo() string {
	s := " in ("
	sep := ""
	for _, e := range a.Exprs {
		s += sep + e.Echo()
		sep = ", "
	}
	return s + ")"
}

func (a *In) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
	for i := range a.Exprs {
		childExpr(fn, &a.Exprs[i])
	}
}

// InRange is added by Folder to bypass strict compare for same type ranges
type InRange struct {
	exprNodeT
	noPos
	E       Expr
	Org     Expr // *Constant
	End     Expr // *Constant
	OrgTok  tok.Token
	EndTok  tok.Token
	evalRaw bool
}

func (a *InRange) String() string {
	return "InRange(" + a.E.String() +
		" " + a.OrgTok.String() + " " + a.Org.String() +
		" " + a.EndTok.String() + " " + a.End.String() + ")"
}

func (a *InRange) Echo() string {
	return a.E.Echo() + tokEcho[a.OrgTok] + a.Org.Echo() +
		" and " + a.E.Echo() + tokEcho[a.EndTok] + a.End.Echo()
}

func (a *InRange) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
	childExpr(fn, &a.Org)
	childExpr(fn, &a.End)
}

type Call struct {
	exprNodeT
	noPos
	Fn      Expr
	Args    []Arg
	End     int32
	RawEval bool
}

func (a *Call) String() string {
	s := "Call(" + a.Fn.String()
	for _, arg := range a.Args {
		s += " " + arg.String()
	}
	return s + ")"
}

func (a *Call) Echo() string {
	s := a.Fn.Echo() + "("
	sep := ""
	for _, arg := range a.Args {
		s += sep + arg.Echo()
		sep = ", "
	}
	return s + ")"
}

func (a *Call) Children(fn func(Node) Node) {
	childExpr(fn, &a.Fn)
	for i := range a.Args {
		childExpr(fn, &a.Args[i].E)
	}
}

type Arg struct {
	SuAstNode
	Name Value // nil if not named
	E    Expr
	TwoPos
}

func (a *Arg) String() string {
	s := ""
	if a.Name != nil {
		if ks, ok := a.Name.(SuStr); ok && lexer.IsIdentifier(string(ks)) {
			s += string(ks) + ":"
		} else {
			s += a.Name.String() + ":"
		}
	}
	return s + a.E.String()
}

func (a *Arg) Echo() string {
	s := ""
	if a.Name != nil {
		if ks, ok := a.Name.(SuStr); ok && lexer.IsIdentifier(string(ks)) {
			s += string(ks) + ":"
		} else {
			s += a.Name.String() + ":"
		}
	}
	return s + a.E.Echo()
}

type Function struct {
	exprNodeT
	Final  map[string]uint8
	Params []Param
	Body   []Statement
	Base   Gnum
	TwoPos
	Pos1        int32
	Pos2        int32
	HasBlocks   bool
	IsNewMethod bool
}

func (a *Function) String() string {
	return a.str("Function")
}

func (a *Function) str(which string) string {
	s := which + "(" + params(a.Params)
	for _, stmt := range a.Body {
		if stmt != nil {
			s += "\n\t" + stmt.String()
		}
	}
	return s + ")"
}

func params(ps []Param) string {
	s := ""
	sep := ""
	for _, p := range ps {
		if sep == "" && p.String() == "this" {
			continue
		}
		s += sep + p.String()
		sep = ","
	}
	return s
}

func childStmt(fn func(Node) Node, pstmt *Statement) {
	if *pstmt != nil {
		stmt := fn(*pstmt)
		if stmt == nil {
			*pstmt = nil
		} else {
			*pstmt = stmt.(Statement)
		}
	}
}

func (a *Function) Children(fn func(Node) Node) {
	for i := range a.Body {
		childStmt(fn, &a.Body[i])
	}
}

func (a *Function) Position() int {
	return int(a.org)
}

type Params struct {
	SuAstNode
	Params []Param
}

func (a Params) String() string {
	return params(a.Params)
}

type Param struct {
	SuAstNode
	DefVal Value // may be nil
	Name   Ident // including prefix @ . _
	End    int32
	// Unused is set if the parameter was followed by /*unused*/
	Unused bool
}

func (a *Param) String() string {
	s := a.Name.Name
	if a.DefVal != nil {
		s += "=" + a.DefVal.String()
	}
	return s
}

func (a *Param) GetPos() int {
	return int(a.Name.Pos)
}

func (a *Param) GetEnd() int {
	return int(a.End)
}

type Block struct {
	Name string
	Function
	// CompileAsFunction is set and used by codegen
	CompileAsFunction bool
}

func (a *Block) String() string {
	s := "Block"
	if a.CompileAsFunction {
		s += "-func"
	}
	return a.Function.str(s)
}

func (a *Block) Children(fn func(Node) Node) {
	a.Function = *fn(&a.Function).(*Function)
}

type Statement interface {
	Node
	Position() int
	GetPos() int
	GetEnd() int
	stmtNode()
}
type stmtNodeT struct {
	astNodeT
	TwoPos
}

func (*stmtNodeT) stmtNode() {}
func (stmt *stmtNodeT) Position() int {
	return int(stmt.org)
}

type Compound struct {
	Body []Statement
	stmtNodeT
}

func (a *Compound) String() string {
	if len(a.Body) == 0 {
		return "{}"
	}
	if len(a.Body) == 1 {
		return a.Body[0].String()
	}
	s := "{\n"
	for _, stmt := range a.Body {
		if stmt != nil {
			s += stmt.String() + "\n"
		}
	}
	return s + "}"
}

func (a *Compound) Children(fn func(Node) Node) {
	for i := range a.Body {
		childStmt(fn, &a.Body[i])
	}
}

type If struct {
	Cond Expr
	Then Statement
	Else Statement // may be nil
	stmtNodeT
	ElseEnd int32
}

func (a *If) String() string {
	s := "If(" + a.Cond.String() + " "
	// if x.Then == nil {
	// 	s += "nil"
	// } else {
	s += a.Then.String()
	// }
	if a.Else != nil {
		s += "\nelse " + a.Else.String()
	}
	return s + ")"
}

func (a *If) Children(fn func(Node) Node) {
	childExpr(fn, &a.Cond)
	childStmt(fn, &a.Then)
	childStmt(fn, &a.Else)
}

type Return struct {
	Exprs []Expr
	stmtNodeT
	ReturnThrow bool
}

func (a *Return) String() string {
	s := "Return("
	if a.ReturnThrow {
		s = "ReturnThrow("
	}
	sep := ""
	for _, e := range a.Exprs {
		s += sep + e.String()
		sep = " "
	}
	return s + ")"
}

func (a *Return) Children(fn func(Node) Node) {
	for i := range a.Exprs {
		childExpr(fn, &a.Exprs[i])
	}
}

type Throw struct {
	E Expr
	stmtNodeT
}

func (a *Throw) String() string {
	return "Throw(" + a.E.String() + ")"
}

func (a *Throw) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
}

type TryCatch struct {
	Try         Statement
	Catch       Statement
	CatchFilter string
	CatchVar    Ident
	stmtNodeT
	CatchPos       int32
	CatchEnd       int32
	CatchVarUnused bool
}

func (a *TryCatch) String() string {
	s := "Try(" + a.Try.String()
	if a.Catch != nil {
		s += "\ncatch"
		if a.CatchVar.Name != "" {
			s += " (" + a.CatchVar.Name
			if a.CatchFilter != "" {
				s += ",'" + a.CatchFilter + "'"
			}
			s += ")"
		}
		s += " " + a.Catch.String()
	}
	return s + ")"
}

func (a *TryCatch) Children(fn func(Node) Node) {
	childStmt(fn, &a.Try)
	childStmt(fn, &a.Catch)
}

type Forever struct {
	Body Statement
	stmtNodeT
}

func (a *Forever) String() string {
	return "Forever(" + a.Body.String() + ")"
}

func (a *Forever) Children(fn func(Node) Node) {
	childStmt(fn, &a.Body)
}

type ForIn struct {
	E    Expr
	E2   Expr // used by for-range
	Body Statement
	Var  Ident // optional with for-range
	Var2 Ident // used with: for m,v in ob
	stmtNodeT
}

func (a *ForIn) String() string {
	s := "ForIn("
	if a.Var.Name != "" {
		s += a.Var.Name + " "
	}
	if a.Var2.Name != "" {
		s += a.Var2.Name + " "
	}
	s += a.E.String()
	if a.E2 != nil {
		s += " " + a.E2.String()
	}
	s += "\n" + a.Body.String() + ")"
	return s
}

func (a *ForIn) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
	childExpr(fn, &a.E2)
	childStmt(fn, &a.Body)
}

type For struct {
	Cond Expr
	Body Statement
	Init []Expr
	Inc  []Expr
	stmtNodeT
}

func (a *For) String() string {
	s := "For("
	sep := ""
	for _, e := range a.Init {
		s += sep + e.String()
		sep = ","
	}
	s += "; "
	if a.Cond != nil {
		s += a.Cond.String()
	}
	s += "; "
	sep = ""
	for _, e := range a.Inc {
		s += sep + e.String()
		sep = ","
	}
	return s + "\n" + a.Body.String() + ")"
}

func (a *For) Children(fn func(Node) Node) {
	for i := range a.Init {
		childExpr(fn, &a.Init[i])
	}
	childExpr(fn, &a.Cond)
	childStmt(fn, &a.Body)
	for i := range a.Inc {
		childExpr(fn, &a.Inc[i])
	}
}

type While struct {
	Cond Expr
	Body Statement
	stmtNodeT
}

func (a *While) String() string {
	return "While(" + a.Cond.String() + " " + a.Body.String() + ")"
}

func (a *While) Children(fn func(Node) Node) {
	childExpr(fn, &a.Cond)
	childStmt(fn, &a.Body)
}

type DoWhile struct {
	Body Statement
	Cond Expr
	stmtNodeT
}

func (a *DoWhile) String() string {
	return "DoWhile(" + a.Body.String() + " " + a.Cond.String() + ")"
}

func (a *DoWhile) Children(fn func(Node) Node) {
	childStmt(fn, &a.Body)
	childExpr(fn, &a.Cond)
}

type Break struct {
	stmtNodeT
}

func (*Break) String() string {
	return "Break"
}

type Continue struct {
	stmtNodeT
}

func (*Continue) String() string {
	return "Continue"
}

type ExprStmt struct {
	E Expr
	stmtNodeT
}

func (a *ExprStmt) String() string {
	return a.E.String() // NOTE: doesn't say "ExprStmt"
}

func (a *ExprStmt) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
}

type MultiAssign struct {
	stmtNodeT
	Lhs []Expr
	Rhs Expr
}

func (a *MultiAssign) String() string {
	s := "MultiAssign("
	for _, e := range a.Lhs {
		s += e.Echo() + " "
	}
	return s + a.Rhs.String() + ")"
}

func (a *MultiAssign) Children(fn func(Node) Node) {
	for i := range a.Lhs {
		childExpr(fn, &a.Lhs[i])
	}
	childExpr(fn, &a.Rhs)
}

type Switch struct {
	E       Expr
	Cases   []Case
	Default []Statement // may be nil
	stmtNodeT
	Pos1   int32
	Pos2   int32
	PosDef int32
}

type Case struct {
	SuAstNode
	Exprs []Expr
	Body  []Statement
	TwoPos
}

func (a *Switch) String() string {
	s := "Switch(" + a.E.String()
	for _, c := range a.Cases {
		s += "\n" + c.String()
	}
	if a.Default != nil {
		if len(a.Default) == 0 {
			s += "\n()"
		}
		for _, stmt := range a.Default {
			if stmt != nil {
				s += "\n" + stmt.String()
			}
		}
	}
	return s + ")"
}

func (a *Case) String() string {
	s := "Case("
	sep := ""
	for _, e := range a.Exprs {
		s += sep + e.String()
		sep = ","
	}
	for _, stmt := range a.Body {
		if stmt != nil {
			s += "\n" + stmt.String()
		}
	}
	s += ")"
	return s
}

func (a *Switch) Children(fn func(Node) Node) {
	childExpr(fn, &a.E)
	for i := range a.Cases {
		c := &a.Cases[i]
		for j := range c.Exprs {
			childExpr(fn, &c.Exprs[j])
		}
		for j := range c.Body {
			childStmt(fn, &c.Body[j])
		}
	}
	for i := range a.Default {
		childStmt(fn, &a.Default[i])
	}
}

type ExprPos struct {
	SuAstNode
	Expr
	Pos int32
	End int32
}
