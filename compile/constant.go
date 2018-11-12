package compile

// TODO class

import (
	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ascii"
)

// Constant compiles a Suneido constant (e.g. a library record)
// to a Suneido Value
func Constant(src string) Value {
	p := newParser(src)
	return p.constant()
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
		default:
			panic("not implemented #" + p.Text)
		}
	case IDENTIFIER:
		switch p.Keyword {
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
		if p.Token != CAT || p.lxr.Ahead(0).Token != STRING {
			break
		}
		p.nextSkipNL()
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
		p.error("invalid date", s)
	}
	return date
}

var closing = map[Token]Token{
	L_PAREN:   R_PAREN,
	L_CURLY:   R_CURLY,
	L_BRACKET: R_BRACKET,
}

func (p *parser) object() Value {
	close := closing[p.Token]
	p.next()
	var ob container
	if close == R_PAREN {
		ob = new(SuObject)
	} else {
		ob = new(SuRecord)
	}
	p.memberList(ob, close, false)
	return ob.(Value)
}

type container interface {
	Add(Value)
	Has(Value) bool
	Put(Value, Value)
}

func (p *parser) memberList(ob container, closing Token, inClass bool) {
	for p.Token != closing {
		p.member(ob, closing, inClass)
		if p.Token == COMMA || p.Token == SEMICOLON {
			p.next()
		}
	}
	p.next()
}

func (p *parser) member(ob container, closing Token, inClass bool) {
	start := p.Token
	m := p.constant()
	if inClass && start == IDENTIFIER && p.Token == L_PAREN {
		fn := codegen(p.functionWithoutKeyword(true))
		p.putMem(ob, m, fn)
	} else if p.matchIf(COLON) {
		if p.Token == COMMA || p.Token == SEMICOLON || p.Token == closing {
			p.putMem(ob, m, True)
		} else {
			p.putMem(ob, m, p.constant())
		}
	} else {
		ob.Add(m)
	}
}

func (p *parser) putMem(ob container, m Value, v Value) {
	if ob.Has(m) {
		p.error("duplicate member name (" + m.String() + ")")
	} else {
		ob.Put(m, v)
	}
}

// class parses a class definition
// like object, it builds a value rather than an ast
func (p *parser) class() Value {
	if p.Keyword == CLASS {
		p.matchSkipNL(CLASS)
		if p.Token == COLON {
			p.matchSkipNL(COLON)
		}
	}
	var base string
	if p.Token == IDENTIFIER {
		base = p.ckBase(p.Text)
		p.matchSkipNL(IDENTIFIER)
	}
	p.matchSkipNL(L_CURLY)
	mems := classcon{}
	p.memberList(mems, R_CURLY, true)
	return &SuClass{Base: base, Data: mems}
}

func (p *parser) ckBase(name string) string {
	if !okBase(name) {
		p.error("base class must be global defined in library, got: ", name)
	}
	return name
}

func okBase(name string) bool {
	return ascii.IsUpper(name[0]) || (name[0] == '_' && ascii.IsUpper(name[1]))
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
