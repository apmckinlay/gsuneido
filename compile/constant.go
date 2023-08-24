// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strconv"
	"strings"
	"sync/atomic"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

// Constant compiles an anonymous Suneido constant
func Constant(src string) Value {
	return NamedConstant("", "", src, nil)
}

// NamedConstant compiles a Suneido constant with a name
// e.g. a library record
func NamedConstant(lib, name, src string, prevDef Value) Value {
	p := NewParser(src)
	p.lib = lib
	p.name = name
	p.prevDef = prevDef
	result := p.constant()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return result
}

func Checked(th *Thread, src string) (Value, []string) {
	// can't do AST check after compile because that would miss nested functions
	p := CheckParser(src, th)
	v := p.constant()
	if p.Token != tok.Eof {
		p.Error("did not parse all input")
	}
	return v, p.CheckResults()
}

func (p *Parser) Const() (result Value) {
	defer func(org int32) {
		if r, ok := result.(iSetPos); ok {
			SetPos(r, org, p.EndPos)
		}
	}(p.Pos)
	return p.constant()
}

func (p *Parser) constant() Value {
	switch p.Token {
	case tok.String:
		return p.string()
	case tok.Symbol:
		s := SuStr(p.Text)
		p.Next()
		return s
	case tok.Add:
		p.Next()
		fallthrough
	case tok.Number:
		return p.number()
	case tok.Sub:
		p.Next()
		return OpUnaryMinus(p.number())
	case tok.LParen, tok.LCurly, tok.LBracket:
		return p.object()
	case tok.Hash:
		p.Next()
		switch p.Token {
		case tok.Number:
			return p.date()
		case tok.LParen, tok.LCurly, tok.LBracket:
			return p.object()
		}
		panic("not implemented #" + p.Text)
	case tok.True:
		p.Next()
		return True
	case tok.False:
		p.Next()
		return False
	case tok.Function:
		return p.functionValue()
	case tok.Class:
		return p.class()
	default:
		if p.Token.IsIdent() {
			if okBase(p.Text) && p.Lxr.AheadSkip(0).Token == tok.LCurly {
				return p.class()
			}
			if p.Lxr.Ahead(0).Token != tok.Colon &&
				(p.Text == "struct" || p.Text == "dll" || p.Text == "callback") {
				p.Error("gSuneido does not implement " + p.Text)
			}
			s := p.Text
			p.Next()
			return SuStr(s)
		}
	}
	panic(p.Error("invalid constant, unexpected " + p.Token.String()))
}

func (p *Parser) functionValue() Value {
	prevClassName := p.className
	p.className = "" // prevent privatization in standalone function
	ast := p.Function()
	p.className = prevClassName
	p.CheckFunc(ast)
	return p.codegen(p.lib, p.name, ast, p.prevDef)
}

// string handles compile time concatenation
func (p *Parser) string() Value {
	s := p.Text
	p.Match(tok.String)
	if !p.moreStr() {
		return SuStr(s) // normal case
	}
	strs := []string{s}
	for p.moreStr() {
		p.Match(tok.Cat)
		strs = append(strs, p.Text)
		p.Match(tok.String)
	}
	return p.mkConcat(strs)
}

func (p *Parser) moreStr() bool {
	return p.Token == tok.Cat && p.Lxr.AheadSkip(0).Token == tok.String
}

func (p *Parser) number() Value {
	s := p.Text
	p.Match(tok.Number)
	return NumFromString(s)
}

func (p *Parser) date() Value {
	s := p.Text
	p.Match(tok.Number)
	date := DateFromLiteral(s)
	if date == NilDate {
		p.Error("bad date literal ", s)
	}
	return date
}

var closing = map[tok.Token]tok.Token{
	tok.LParen:   tok.RParen,
	tok.LCurly:   tok.RCurly,
	tok.LBracket: tok.RBracket,
}

const noBase = -1

func (p *Parser) object() Value {
	close := closing[p.Token]
	p.Next()
	var ob container
	if close == tok.RParen {
		ob = p.mkObject()
	} else {
		ob = p.mkRecord()
	}
	p.memberList(ob, close, noBase)
	if close == tok.RBracket {
		ob = p.mkRecOrOb(ob)
	}
	if p, ok := ob.(protectable); ok {
		p.SetReadOnly()
	}
	return ob.(Value)
}

type protectable interface {
	SetReadOnly()
}

func (p *Parser) memberList(ob container, closing tok.Token, base Gnum) {
	for p.Token != closing {
		pos := p.Item.Pos
		k, v := p.member(closing, base)
		if p.Token == tok.Comma || p.Token == tok.Semicolon {
			p.Next()
		}
		if k == nil {
			p.set(ob, nil, v, pos, p.EndPos)
		} else {
			p.putMem(ob, k, v, pos)
		}
	}
	p.Next()
}

func (p *Parser) member(closing tok.Token, base Gnum) (k Value, v Value) {
	start := p.Token
	m := p.constant() // might be key or value
	inClass := base != noBase
	if inClass && start.IsIdent() && p.Token == tok.LParen { // method
		name := p.privatizeDef(m)
		prevName := p.name
		p.name += "." + name
		ast := p.function(true)
		ast.Base = base
		if name == "New" {
			ast.IsNewMethod = true
		}
		p.CheckFunc(ast)
		fn := p.codegen(p.lib, p.name, ast, p.prevDef)
		p.name = prevName
		if f, ok := fn.(*SuFunc); ok {
			f.ClassName = p.className
		}
		return SuStr(name), fn
	}
	if p.MatchIf(tok.Colon) { // named
		if inClass {
			m = SuStr(p.privatizeDef(m))
		}
		if p.Token == tok.Comma || p.Token == tok.Semicolon || p.Token == closing {
			return m, True
		}
		prevName := p.name
		if s, ok := m.ToStr(); ok {
			p.name += "." + s
		}
		c := p.constant()
		p.name = prevName
		return m, c
	}
	return nil, m
}

func (p *Parser) privatizeDef(m Value) string {
	ss, ok := m.(SuStr)
	if !ok {
		p.Error("class member names must be strings")
	}
	name := string(ss)
	if strings.HasPrefix(name, "Getter_") &&
		len(name) > 7 && !ascii.IsUpper(name[7]) {
		p.Error("invalid getter (" + name + ")")
	}
	if !ascii.IsLower(name[0]) {
		return name
	}
	if strings.HasPrefix(name, "getter_") &&
		(len(name) <= 7 || !ascii.IsLower(name[7])) {
		p.Error("invalid getter (" + name + ")")
	}
	return p.privatize(name, p.className)
}

// putMem checks for duplicate member and then calls p.set with endpos
func (p *Parser) putMem(ob container, m Value, v Value, pos int32) {
	if ob.HasKey(m) {
		p.ErrorAt(pos, "duplicate member name ("+m.String()+")")
	} else {
		p.set(ob, m, v, pos, p.EndPos)
	}
}

// classNum is used to generate names for anonymous classes
var classNum atomic.Int32

// class parses a class definition
// like object, it builds a value rather than an ast
func (p *Parser) class() (result Value) {
	if p.Token == tok.Class {
		p.Match(tok.Class)
		if p.Token == tok.Colon {
			p.Match(tok.Colon)
		}
	}
	var base Gnum
	baseName := "class"
	if p.Token.IsIdent() {
		baseName = p.Text
		base = p.ckBase(baseName)
		p.MatchIdent()
	}
	pos1 := p.EndPos
	p.Match(tok.LCurly)
	pos2 := p.EndPos
	prevClassName := p.className
	p.className = p.getClassName()
	mems := p.mkClass(baseName)
	p.memberList(mems, tok.RCurly, base)
	p.setPos(mems, pos1, pos2)
	p.className = prevClassName
	if cc, ok := mems.(classBuilder); ok {
		return &SuClass{Base: base, Lib: p.lib, Name: p.name,
			MemBase: MemBase{Data: cc}}
	}
	return mems.(Value)
}

func (p *Parser) ckBase(name string) Gnum {
	if !okBase(name) {
		p.Error("base class must be global defined in library, got: ", name)
	}
	if name[0] == '_' {
		if name == "_" || name[1:] != p.name {
			p.Error("invalid reference to " + name)
		}
		return Global.Overload(name, p.prevDef)
		// for _Name in expressions see codegen.go cgen.identifier
	}
	p.CheckGlobal(name, int(p.Pos))
	return Global.Num(name)
}

func okBase(name string) bool {
	return ascii.IsUpper(name[0]) ||
		(name[0] == '_' && len(name) > 1 && ascii.IsUpper(name[1]))
}

func (p *Parser) getClassName() string {
	last := p.name
	i := strings.LastIndexAny(last, " .")
	if i != -1 {
		last = last[i+1:]
	}
	if last == "" || last == "?" {
		cn := classNum.Add(1)
		className := "Class" + strconv.Itoa(int(cn))
		return className
	}
	return last
}
