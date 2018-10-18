package compile

// TODO class

import (
	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// Constant compiles a Suneido constant (e.g. a library record)
// to a Suneido Value
func Constant(src string) Value {
	p := newParser(src)
	return p.constant()
}

func (p *parser) constant() Value {
	switch p.Token {
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
			panic("not implemented") // TODO parse classes
		default:
			s := p.Text
			p.next()
			return SuStr(s)
		}
	case STRING:
		return p.string()
	case NUMBER:
		return p.number()
	case ADD:
		p.next()
		return p.number()
	case SUB:
		p.next()
		return Uminus(p.number())
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
	case L_PAREN, L_CURLY, L_BRACKET:
		return p.object()
	}
	panic(p.error("invalid constant"))
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
	return p.memberList(close)
}

func (p *parser) memberList(closing Token) Value {
	ob := &SuObject{}
	for p.Token != closing {
		p.member(ob, closing)
		if p.Token == COMMA || p.Token == SEMICOLON {
			p.next()
		}
	}
	p.next()
	return ob
}

func (p *parser) member(ob *SuObject, closing Token) {
	mem := p.memberName()
	val := p.memberValue(mem, closing)
	if mem == nil {
		ob.Add(val)
	} else {
		ob.Put(mem, val)
	}
}

func (p *parser) memberName() Value {
	if p.isMemberName() {
		return p.constant()
		// does not absorb COLON
	}
	return nil
}

func (p *parser) isMemberName() bool {
	i := 0
	tok := p.Token
	if tok == ADD || tok == SUB || tok == HASH {
		tok = p.lxr.Ahead(0).Token
		i = 1
	}
	if (tok == IDENTIFIER || tok == STRING || tok == NUMBER) &&
		p.lxr.Ahead(i).Token == COLON {
		return true
	}
	return false
}

func (p *parser) memberValue(mem Value, closing Token) Value {
	if mem != nil {
		if p.Token == COLON {
			p.next()
		}
		if p.Token == COMMA || p.Token == closing {
			return True
		}
	}
	return p.constant()
}
