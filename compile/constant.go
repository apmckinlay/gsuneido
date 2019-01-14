package compile

import (
	"strconv"
	"strings"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

func NamedConstant(name, src string) Value {
	p := newParser(src)
	p.className = name
	result := p.constant()
	if p.Token != EOF {
		p.error("syntax error: did not parse all input")
	}
	if named, ok := result.(Named); ok {
		named.SetName(name)
	}
	return result
}

// Constant compiles a Suneido constant (e.g. a library record)
// to a Suneido Value
func Constant(src string) Value {
	return NamedConstant("", src)
}

func (p *parser) constant() Value {
	switch p.Token {
	case STRING:
		return p.string()
	case ADD:
		p.next()
		fallthrough
	case NUMBER:
		return p.number()
	case SUB:
		p.next()
		return Uminus(p.number())
	case L_PAREN, L_CURLY, L_BRACKET:
		return p.object()
	case HASH:
		p.next()
		switch p.Token {
		case NUMBER:
			return p.date()
		case L_PAREN, L_CURLY, L_BRACKET:
			return p.object()
		}
		panic("not implemented #" + p.Text)
	case TRUE:
		p.next()
		return True
	case FALSE:
		p.next()
		return False
	case FUNCTION:
		ast := p.function()
		return codegen(ast)
	case CLASS:
		return p.class()
	default:
		if IsIdent[p.Token] {
			if okBase(p.Text) && p.lxr.AheadSkip(0).Token == L_CURLY {
				return p.class()
			}
			s := p.Text
			p.next()
			return SuStr(s)
		}
	}
	panic(p.error("invalid constant, unexpected " + p.Token.String()))
}

func (p *parser) string() Value {
	s := ""
	for {
		s += p.Text
		p.match(STRING)
		if p.Token != CAT || p.lxr.AheadSkip(0).Token != STRING {
			break
		}
		p.next()
	}
	return SuStr(s)
}

func (p *parser) number() Value {
	s := p.Text
	p.match(NUMBER)
	return NumFromString(s)
}

func (p *parser) date() Value {
	s := p.Text
	p.match(NUMBER)
	date := DateFromLiteral(s)
	if date == NilDate {
		p.error("invalid date ", s)
	}
	return date
}

var closing = map[Token]Token{
	L_PAREN:   R_PAREN,
	L_CURLY:   R_CURLY,
	L_BRACKET: R_BRACKET,
}

const noBase = -1

func (p *parser) object() Value {
	close := closing[p.Token]
	p.next()
	var ob container
	if close == R_PAREN {
		ob = new(SuObject)
	} else {
		ob = NewSuRecord()
	}
	p.memberList(ob, close, noBase)
	if p,ok := ob.(protectable); ok {
	p.SetReadOnly()
	}
	return ob.(Value)
}

type container interface {
	Add(Value)
	Has(Value) bool
	Put(Value, Value)
}

type protectable interface {
	SetReadOnly()
}

func (p *parser) memberList(ob container, closing Token, base Global) {
	for p.Token != closing {
		p.member(ob, closing, base)
		if p.Token == COMMA || p.Token == SEMICOLON {
			p.next()
		}
	}
	p.next()
}

func (p *parser) member(ob container, closing Token, base Global) {
	start := p.Token
	m := p.constant()
	inClass := base != noBase
	if inClass && IsIdent[start] && p.Token == L_PAREN { // method
		ast := p.method()
		ast.Base = base
		name := p.privatizeDef(m)
		if name == "New" {
			ast.IsNewMethod = true
		}
		fn := codegen(ast)
		fn.ClassName = p.className
		fn.Name = "." + name
		p.putMem(ob, SuStr(name), fn)
	} else if p.matchIf(COLON) {
		if inClass {
			m = SuStr(p.privatizeDef(m))
		}
		if p.Token == COMMA || p.Token == SEMICOLON || p.Token == closing {
			p.putMem(ob, m, True)
		} else {
			p.putMem(ob, m, p.constant())
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
	if ob.Has(m) {
		p.error("duplicate member name (" + m.String() + ")")
	} else {
		ob.Put(m, v)
	}
}

// classNum is used to generate names for anonymous classes
// should be referenced atomically
var classNum int32

// class parses a class definition
// like object, it builds a value rather than an ast
func (p *parser) class() Value {
	if p.Token == CLASS {
		p.match(CLASS)
		if p.Token == COLON {
			p.match(COLON)
		}
	}
	var base Global
	if p.Token == IDENTIFIER {
		base = p.ckBase(p.Text)
		p.matchIdent()
	}
	p.match(L_CURLY)
	mems := classcon{}
	classNamePrev := p.className
	if p.className == "" {
		cn := atomic.AddInt32(&classNum, 1)
		p.className = "Class" + strconv.Itoa(int(cn))
	}
	p.memberList(mems, R_CURLY, base)
	p.className = classNamePrev
	return &SuClass{Base: base, MemBase: MemBase{Data: mems}}
}

func (p *parser) ckBase(name string) Global {
	if !okBase(name) {
		p.error("base class must be global defined in library, got: ", name)
	}
	return GlobalNum(name)
}

func okBase(name string) bool {
	return ascii.IsUpper(name[0]) ||
		(name[0] == '_' && len(name) > 1 && ascii.IsUpper(name[1]))
}

type classcon map[string]Value

func (c classcon) Add(Value) {
	panic("class members must be named")
}
func (c classcon) Has(m Value) bool {
	if s, ok := m.(SuStr); ok {
		_, ok = c[string(s)]
		return ok
	}
	panic("class member names must be strings")
}
func (c classcon) Put(m, v Value) {
	c[string(m.(SuStr))] = v
}
