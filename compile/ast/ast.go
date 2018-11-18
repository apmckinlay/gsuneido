// Package ast defines the node types
// used by the compiler to build syntax trees
package ast

import (
	"github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
)

type Node interface {
	astNode()
	String() string
}
type astNodeT struct{}

func (*astNodeT) astNode() {}

// Expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}
type exprNodeT struct{ astNodeT }

func (*exprNodeT) exprNode() {}

type Ident struct {
	exprNodeT
	Name string
}

func (a *Ident) String() string {
	return a.Name
}

type Constant struct {
	exprNodeT
	Val Value
}

func (a *Constant) String() string {
	return a.Val.String()
}

type Unary struct {
	exprNodeT
	Tok lexer.Token
	E   Expr
}

func (a *Unary) String() string {
	return "(" + a.Tok.String() + " " + a.E.String() + ")"
}

type Binary struct {
	exprNodeT
	Lhs Expr
	Tok lexer.Token
	Rhs Expr
}

func (a *Binary) String() string {
	return "(" + a.Tok.String() + " " + a.Lhs.String() + " " + a.Rhs.String() + ")"
}

type Trinary struct {
	exprNodeT
	Cond Expr
	T    Expr
	F    Expr
}

func (a *Trinary) String() string {
	return "(? " + a.Cond.String() + " " + a.T.String() + " " + a.F.String() + ")"
}

// Nary is used for associative binary operators e.g. add, multiply, and, or
type Nary struct {
	exprNodeT
	Tok   lexer.Token
	Exprs []Expr
}

func (a *Nary) String() string {
	s := "(" + a.Tok.String()
	for _, e := range a.Exprs {
		s += " " + e.String()
	}
	return s + ")"
}

type RangeTo struct {
	exprNodeT
	E    Expr
	From Expr
	To   Expr
}

func (a *RangeTo) String() string {
	s := a.E.String() + "["
	if a.From != nil {
		s += a.From.String()
	}
	s += ".."
	if a.To != nil {
		s += a.To.String()
	}
	return s + "]"
}

type RangeLen struct {
	exprNodeT
	E    Expr
	From Expr
	Len  Expr
}

func (a *RangeLen) String() string {
	s := a.E.String() + "["
	if a.From != nil {
		s += a.From.String()
	}
	s += "::"
	if a.Len != nil {
		s += a.Len.String()
	}
	return s + "]"
}

type Mem struct {
	exprNodeT
	E Expr
	M Expr
}

func (a *Mem) String() string {
	if c, ok := a.M.(*Constant); ok {
		if s, ok := c.Val.(SuStr); ok {
			return a.E.String() + "." + string(s)
		}
	}
	return a.E.String() + "[" + a.M.String() + "]"
}

type In struct {
	exprNodeT
	E     Expr
	Exprs []Expr
}

func (a *In) String() string {
	s := "(" + a.E.String() + " in"
	for _, e := range a.Exprs {
		s += " " + e.String()
	}
	return s + ")"
}

type Call struct {
	exprNodeT
	Fn   Expr
	Args []Arg
}

func (a *Call) String() string {
	s := "(call " + a.Fn.String()
	for _, arg := range a.Args {
		s += " " + arg.String()
	}
	return s + ")"
}

type Arg struct {
	Name Value // nil if not named
	E    Expr
}

func (a *Arg) String() string {
	s := ""
	if a.Name != nil {
		if ks, ok := a.Name.(SuStr); ok && IsIdentifier(string(ks)) {
			s += string(ks) + ": "
		} else {
			s += a.Name.String() + ": "
		}
	}
	return s + a.E.String()
}

type Function struct {
	exprNodeT
	Params []Param
	Body   []Statement
}

func (a *Function) String() string {
	s := "function("
	sep := ""
	for _, p := range a.Params {
		s += sep + p.String()
		sep = ","
	}
	s += ") {\n"
	for _, x := range a.Body {
		s += "\t" + x.String() + "\n"
	}
	return s + "}"
}

type Param struct {
	Name   string // including prefix @ . _
	DefVal Value  // may be nil
}

func (a *Param) String() string {
	s := a.Name
	if a.DefVal != nil {
		s += "=" + a.DefVal.String()
	}
	return s
}

type Block struct {
	Function
}

func (a *Block) String() string {
	s := "{"
	if len(a.Params) > 0 {
		s += "|"
		sep := ""
		for _, p := range a.Params {
			s += sep + p.String()
			sep = ","
		}
		s += "|"
	}
	return s + " }"
}

type Factory interface {
	Ident(name string) Expr
	Constant(val Value) Expr
	Unary(tok lexer.Token, expr Expr) Expr
	Binary(lhs Expr, tok lexer.Token, rhs Expr) Expr
	Mem(e Expr, m Expr) Expr
	Trinary(cond Expr, t Expr, f Expr) Expr
	Nary(tok lexer.Token, exprs []Expr) Expr
	In(e Expr, exprs []Expr) Expr
	Call(fn Expr, args []Arg) Expr
}

type Builder struct{}

func (Builder) Ident(name string) Expr {
	return &Ident{Name: name}
}
func (Builder) Constant(val Value) Expr {
	return &Constant{Val: val}
}
func (Builder) Unary(tok lexer.Token, expr Expr) Expr {
	return &Unary{Tok: tok, E: expr}
}
func (Builder) Binary(lhs Expr, tok lexer.Token, rhs Expr) Expr {
	return &Binary{Lhs: lhs, Tok: tok, Rhs: rhs}
}
func (Builder) Trinary(cond Expr, t Expr, f Expr) Expr {
	return &Trinary{Cond: cond, T: t, F: f}
}
func (Builder) Nary(tok lexer.Token, exprs []Expr) Expr {
	return &Nary{Tok: tok, Exprs: exprs}
}
func (Builder) Mem(e Expr, m Expr) Expr {
	return &Mem{E: e, M: m}
}
func (Builder) In(e Expr, exprs []Expr) Expr {
	return &In{E: e, Exprs: exprs}
}
func (Builder) Call(fn Expr, args []Arg) Expr {
	return &Call{Fn: fn, Args: args}
}

var _ Factory = Builder{}

// Statement nodes implement the Stmt interface.
type Statement interface {
	Node
	Position() int
	SetPos(pos int)
	stmtNode()
}
type stmtNodeT struct {
	astNodeT
	Pos int
}

func (*stmtNodeT) stmtNode() {}
func (stmt *stmtNodeT) Position() int {
	return stmt.Pos
}
func (stmt *stmtNodeT) SetPos(pos int) {
	stmt.Pos = pos
}

type Compound struct {
	stmtNodeT
	Body []Statement
}

func (x *Compound) String() string {
	if len(x.Body) == 0 {
		return "{}"
	}
	if len(x.Body) == 1 {
		return x.Body[0].String()
	}
	s := "{\n"
	for _, stmt := range x.Body {
		s += stmt.String() + "\n"
	}
	return s + "}"
}

type If struct {
	stmtNodeT
	Cond Expr
	Then Statement
	Else Statement // may be nil
}

func (x *If) String() string {
	s := "if " + x.Cond.String() + "\n" + x.Then.String()
	if x.Else != nil {
		s += "\nelse\n" + x.Else.String()
	}
	return s
}

type Return struct {
	stmtNodeT
	E Expr
}

func (x *Return) String() string {
	s := "return"
	if x.E != nil {
		s += " " + x.E.String()
	}
	return s
}

type Throw struct {
	stmtNodeT
	E Expr
}

func (x *Throw) String() string {
	return "throw " + x.E.String()
}

type TryCatch struct {
	stmtNodeT
	Try         Statement
	CatchVar    string
	CatchFilter string
	Catch       Statement
}

func (x *TryCatch) String() string {
	s := "try\n" + x.Try.String()
	if x.Catch != nil {
		s += "\ncatch"
		if x.CatchVar != "" {
			s += " (" + x.CatchVar
			if x.CatchFilter != "" {
				s += ",'" + x.CatchFilter + "'"
			}
			s += ")"
		}
		s += "\n" + x.Catch.String()
	}
	return s
}

type Forever struct {
	stmtNodeT
	Body Statement
}

func (x *Forever) String() string {
	return "forever\n" + x.Body.String()
}

type ForIn struct {
	stmtNodeT
	Var  string
	E    Expr
	Body Statement
}

func (x *ForIn) String() string {
	return "for " + x.Var + " in " + x.E.String() + "\n" + x.Body.String()
}

type For struct {
	stmtNodeT
	Init []Expr
	Cond Expr
	Inc  []Expr
	Body Statement
}

func (x *For) String() string {
	s := "for "
	sep := ""
	for _, e := range x.Init {
		s += sep + e.String()
		sep = ","
	}
	s += "; " + x.Cond.String() + "; "
	sep = ""
	for _, e := range x.Inc {
		s += sep + e.String()
		sep = ","
	}
	return s + "\n" + x.Body.String()
}

type While struct {
	stmtNodeT
	Cond Expr
	Body Statement
}

func (x *While) String() string {
	return "while " + x.Cond.String() + "\n" + x.Body.String()
}

type DoWhile struct {
	stmtNodeT
	Body Statement
	Cond Expr
}

func (x *DoWhile) String() string {
	return "do\n" + x.Body.String() + "\nwhile " + x.Cond.String()
}

type Break struct {
	stmtNodeT
}

func (*Break) String() string {
	return "break"
}

type Continue struct {
	stmtNodeT
}

func (*Continue) String() string {
	return "continue"
}

type Expression struct {
	stmtNodeT
	E Expr
}

func (x *Expression) String() string {
	return x.E.String()
}

type Switch struct {
	stmtNodeT
	E       Expr
	Cases   []Case
	Default []Statement // may be nil
}

type Case struct {
	Exprs []Expr
	Body  []Statement
}

func (x *Switch) String() string {
	s := "switch " + x.E.String() + " {"
	for _, c := range x.Cases {
		s += "\ncase "
		sep := ""
		for _, e := range c.Exprs {
			s += sep + e.String()
			sep = ","
		}
		for _, stmt := range c.Body {
			s += "\n\t" + stmt.String()
		}
	}
	if x.Default != nil {
		s += "\n" + "default:"
		for _, stmt := range x.Default {
			s += "\n\t" + stmt.String()
		}
	}
	return s + "\n}"
}
