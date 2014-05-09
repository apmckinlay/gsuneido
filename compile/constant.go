package compile

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	v "github.com/apmckinlay/gsuneido/value"
)

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
			return v.True
		case FALSE:
			return v.False
		case FUNCTION:
			ast := p.function()
			return codegen(ast)
		default:
			panic("constant: not implemented: " + p.Value)
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
		default:
			panic("not implemented #" + p.Value)
		}
	}
	panic(p.error("invalid constant"))
}

func (p *parser) string() v.Value {
	s := ""
	for {
		s += p.Value
		p.match(STRING)
		if p.Token != CAT || p.lxr.Ahead(0).Token != STRING {
			break
		}
		p.matchSkipNewlines(CAT)
	}
	return v.SuStr(s)
}

func (p *parser) number() v.Value {
	s := p.Value
	p.match(NUMBER)

	n, err := strconv.ParseInt(s, 0, 32)
	if err == nil {
		return v.SuInt(n)
	}
	dn, err := dnum.Parse(s)
	if err == nil {
		return v.DnumToValue(dn)
	}
	panic(p.error("invalid number", s))
}

func (p *parser) date() v.Value {
	s := p.Value
	p.match(NUMBER)
	date := v.DateFromLiteral(s)
	if date == v.NilDate {
		p.error("invalid date", s)
	}
	return date
}
