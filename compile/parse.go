package compile

import "fmt"

func newParser(src string) *parser {
	lxr := NewLexer(src)
	return &parser{lxr: lxr, Item: lxr.Next()}
}

type parser struct {
	lxr *Lexer
	Item
	nest int
	bld  builder // used by expression
}

func (p *parser) match(tok Token) {
	if tok != p.Token && tok != p.Keyword {
		p.error("expecting", tok)
	}
	p.next()
}

func (p *parser) matchIf(tok Token) bool {
	if tok != p.Token && tok != p.Keyword {
		return false
	}
	p.next()
	return true
}

func (p *parser) matchSkipNewlines(tok Token) {
	if tok != p.Token && tok != p.Keyword {
		p.error("expecting", tok)
	}
	for {
		p.next()
		if p.Token != NEWLINE {
			break
		}
	}
}

func (p *parser) evalMatch(result T, tok Token) T {
	p.match(tok)
	return result
}

func (p *parser) evalNext(result T) T {
	p.next()
	return result
}

// next advances to the next non-white token, tracking nesting
// NOTE: it does NOT skip newlines
func (p *parser) next() {
	for {
		p.Item = p.lxr.Next()
		switch p.Token {
		case COMMENT, WHITESPACE:
			continue
		case L_CURLY, L_PAREN, L_BRACKET:
			p.nest++
		case R_CURLY, R_PAREN, R_BRACKET:
			p.nest--
		}
		break
	}
	if p.Token == STRING && p.Keyword != STRING {
		// make a copy of strings that are slices of the source
		p.Value = " " + p.Value
		p.Value = p.Value[1:]
	}
	//fmt.Println("item:", p.Item)
}

// returns string so it can be called inside panic
// so compiler knows we don't return
func (p *parser) error(args ...interface{}) string {
	panic("syntax error" + fmt.Sprint(args...))
}
