package compile

import v "github.com/apmckinlay/gsuneido/value"

// Constant compiles a Suneido constant (e.g. a library record)
// to a Suneido Value
func Constant(src string) v.Value {
	p := newParser(src)
	return p.constant()
}

func (p *parser) constant() v.Value {
	switch p.Token {
	case IDENTIFIER:
		switch p.Keyword {
		case TRUE:
			p.next()
			return v.True
		case FALSE:
			p.next()
			return v.False
		case FUNCTION:
			ast := p.function()
			return codegen(ast)
		default:
			s := p.Text
			p.next()
			return v.SuStr(s)
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
		return v.Uminus(p.number())
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

func (p *parser) string() v.Value {
	s := ""
	for {
		s += p.Text
		p.match(STRING)
		if p.Token != CAT || p.lxr.Ahead(0).Token != STRING {
			break
		}
		p.nextSkipNL()
	}
	return v.SuStr(s)
}

func (p *parser) number() v.Value {
	s := p.Text
	p.match(NUMBER)
	val, err := v.NumFromString(s)
	if err != nil {
		panic(p.error("invalid number", s))
	}
	return val
}

func (p *parser) date() v.Value {
	s := p.Text
	p.match(NUMBER)
	date := v.DateFromLiteral(s)
	if date == v.NilDate {
		p.error("invalid date", s)
	}
	return date
}

func (p *parser) object() v.Value {
	closing := p.Token.closing()
	p.next()
	return p.memberList(closing)
}

func (p *parser) memberList(closing Token) v.Value {
	ob := &v.SuObject{}
	for p.Token != closing {
		p.member(ob, closing)
		if p.Token == COMMA || p.Token == SEMICOLON {
			p.next()
		}
	}
	p.next()
	return ob
}

func (p *parser) member(ob *v.SuObject, closing Token) {
	mem := p.memberName()
	val := p.memberValue(mem, closing)
	if mem == nil {
		ob.Add(val)
	} else {
		ob.Put(mem, val)
	}
}

func (p *parser) memberName() v.Value {
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

func (p *parser) memberValue(mem v.Value, closing Token) v.Value {
	if mem != nil {
		if p.Token == COLON {
			p.next()
		}
		if p.Token == COMMA || p.Token == closing {
			return v.True
		}
	}
	return p.constant()
}
