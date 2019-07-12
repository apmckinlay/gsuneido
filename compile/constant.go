package compile

import (
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/check"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

// Constant compiles an anonymous Suneido constant
func Constant(src string) Value {
	return NamedConstant("", "", src)
}

// NamedConstant compiles a Suneido constant with a name
// e.g. a library record
func NamedConstant(lib, name, src string) Value {
	p := NewParser(src)
	p.lib = lib
	p.name = name
	result := p.constant()
	if p.Token != tok.Eof {
		p.error("syntax error: did not parse all input")
	}
	return result
}

func Checked(src string) (Value, []string) {
	p := NewParser(src)
	var ck check.Check
	p.checker = func(ast *ast.Function) { ck.Check(ast) }
	v := p.constant()
	return v, ck.Results
}

func (p *parser) constant() Value {
	switch p.Token {
	case tok.String:
		return p.string()
	case tok.Add:
		p.next()
		fallthrough
	case tok.Number:
		return p.number()
	case tok.Sub:
		p.next()
		return UnaryMinus(p.number())
	case tok.LParen, tok.LCurly, tok.LBracket:
		return p.object()
	case tok.Hash:
		p.next()
		switch p.Token {
		case tok.Number:
			return p.date()
		case tok.LParen, tok.LCurly, tok.LBracket:
			return p.object()
		}
		panic("not implemented #" + p.Text)
	case tok.True:
		p.next()
		return True
	case tok.False:
		p.next()
		return False
	case tok.Function:
		return p.functionValue()
	case tok.Class:
		return p.class()
	default:
		if p.Token.IsIdent() {
			if okBase(p.Text) && p.lxr.AheadSkip(0).Token == tok.LCurly {
				return p.class()
			}
			if p.lxr.Ahead(0).Token != tok.Colon &&
				(p.Text == "struct" || p.Text == "dll" || p.Text == "callback") {
				panic("gSuneido does not implement " + p.Text)
			}
			s := p.Text
			p.next()
			return SuStr(s)
		}
	}
	panic(p.error("invalid constant, unexpected " + p.Token.String()))
}

func (p *parser) functionValue() Value {
	prevClassName := p.className
	p.className = "" // prevent privatization in standalone function
	ast := p.Function()
	p.className = prevClassName
	p.checker(ast)
	f := codegen(ast)
	f.Lib = p.lib
	f.Name = p.name
	return f
}

func (p *parser) string() Value {
	s := ""
	for {
		s += p.Text
		p.match(tok.String)
		if p.Token != tok.Cat || p.lxr.AheadSkip(0).Token != tok.String {
			break
		}
		p.next()
	}
	return SuStr(s)
}

func (p *parser) number() Value {
	s := p.Text
	p.match(tok.Number)
	return NumFromString(s)
}

func (p *parser) date() Value {
	s := p.Text
	p.match(tok.Number)
	date := DateFromLiteral(s)
	if date == NilDate {
		p.error("invalid date ", s)
	}
	return date
}

var closing = map[tok.Token]tok.Token{
	tok.LParen:   tok.RParen,
	tok.LCurly:   tok.RCurly,
	tok.LBracket: tok.RBracket,
}

const noBase = -1

func (p *parser) object() Value {
	close := closing[p.Token]
	p.next()
	var ob container
	if close == tok.RParen {
		ob = new(SuObject)
	} else {
		ob = NewSuRecord()
	}
	p.memberList(ob, close, noBase)
	if close == tok.RBracket && ob.(*SuRecord).ListSize() > 0 {
		suob := *ob.(*SuRecord).ToObject()
		ob = &suob
	}
	if p, ok := ob.(protectable); ok {
		p.SetReadOnly()
	}
	return ob.(Value)
}

// container allows using memberList etc. for both objects and classes
type container interface {
	Add(Value)
	HasKey(Value) bool
	Set(Value, Value)
}

type protectable interface {
	SetReadOnly()
}

func (p *parser) memberList(ob container, closing tok.Token, base Gnum) {
	for p.Token != closing {
		p.member(ob, closing, base)
		if p.Token == tok.Comma || p.Token == tok.Semicolon {
			p.next()
		}
	}
	p.next()
}

func (p *parser) member(ob container, closing tok.Token, base Gnum) {
	start := p.Token
	m := p.constant()
	inClass := base != noBase
	if inClass && start.IsIdent() && p.Token == tok.LParen { // method
		name := p.privatizeDef(m)
		prevName := p.name
		p.name += "." + name
		ast := p.method()
		ast.Base = base
		if name == "New" {
			ast.IsNewMethod = true
		}
		p.checker(ast)
		fn := codegen(ast)
		fn.Name = p.name
		p.name = prevName
		fn.ClassName = p.className
		p.putMem(ob, SuStr(name), fn)
	} else if p.matchIf(tok.Colon) {
		if inClass {
			m = SuStr(p.privatizeDef(m))
		}
		if p.Token == tok.Comma || p.Token == tok.Semicolon || p.Token == closing {
			p.putMem(ob, m, True)
		} else {
			prevName := p.name
			if s, ok := m.ToStr(); ok {
				p.name += "." + s
			}
			p.putMem(ob, m, p.constant())
			p.name = prevName
		}
	} else {
		ob.Add(m)
	}
}

func (p *parser) privatizeDef(m Value) string {
	ss, ok := m.(SuStr)
	if !ok {
		p.error("class member names must be strings")
	}
	name := string(ss)
	if strings.HasPrefix(name, "Getter_") &&
		len(name) > 7 && !ascii.IsUpper(name[7]) {
		p.error("invalid getter (" + name + ")")
	}
	if !ascii.IsLower(name[0]) {
		return name
	}
	if strings.HasPrefix(name, "getter_") {
		if len(name) <= 7 || !ascii.IsLower(name[7]) {
			p.error("invalid getter (" + name + ")")
		}
		return "Getter_" + p.className + name[6:]
	}
	return p.className + "_" + name
}

func (p *parser) putMem(ob container, m Value, v Value) {
	if ob.HasKey(m) {
		p.error("duplicate member name (" + m.String() + ")")
	} else {
		ob.Set(m, v)
	}
}

// classNum is used to generate names for anonymous classes
// should be referenced atomically
var classNum int32

// class parses a class definition
// like object, it builds a value rather than an ast
func (p *parser) class() Value {
	if p.Token == tok.Class {
		p.match(tok.Class)
		if p.Token == tok.Colon {
			p.match(tok.Colon)
		}
	}
	var base Gnum
	if p.Token == tok.Identifier {
		base = p.ckBase(p.Text)
		p.matchIdent()
	}
	p.match(tok.LCurly)
	prevClassName := p.className
	p.className = p.getClassName()
	mems := classcon{}
	p.memberList(mems, tok.RCurly, base)
	p.className = prevClassName
	return &SuClass{Base: base, MemBase: MemBase{Data: mems}, Lib: p.lib, Name: p.name}
}

func (p *parser) ckBase(name string) Gnum {
	if !okBase(name) {
		p.error("base class must be global defined in library, got: ", name)
	}
	if name[0] == '_' {
		return Global.Copy(name[1:])
	}
	return Global.Num(name)
}

func okBase(name string) bool {
	return ascii.IsUpper(name[0]) ||
		(name[0] == '_' && len(name) > 1 && ascii.IsUpper(name[1]))
}

func (p *parser) getClassName() string {
	last := p.name
	i := strings.LastIndexAny(last, " .")
	if i != -1 {
		last = last[i+1:]
	}
	if last == "" || last == "?" {
		cn := atomic.AddInt32(&classNum, 1)
		className := "Class" + strconv.Itoa(int(cn))
		return className
	}
	return last
}

type classcon map[string]Value

func (c classcon) Add(Value) {
	panic("class members must be named")
}
func (c classcon) HasKey(m Value) bool {
	if s, ok := m.(SuStr); ok {
		_, ok = c[string(s)]
		return ok
	}
	panic("class member names must be strings")
}
func (c classcon) Set(m, v Value) {
	c[string(m.(SuStr))] = v
}
