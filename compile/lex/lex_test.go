package lex

import "testing"

import . "hamcrest"

func TestKeywords(t *testing.T) {
	Assert(t).That(Keyword("forever"), Equals(FOREVER))
}

func TestIsInfix(t *testing.T) {
	assert := Assert(t)
	assert.That(ADD.IsInfix(), Equals(true))
	assert.That(FOREVER.IsInfix(), Equals(false))
}

func TestLexer(t *testing.T) {
	assert := Assert(t)
	assert.That(first("forever"),
		Equals(Item{IDENTIFIER, FOREVER, "forever"}))
	assert.That(first("foo"),
		Equals(Item{IDENTIFIER, NIL, "foo"}))
	assert.That(first("is"),
		Equals(Item{IDENTIFIER, IS, "is"}))
	assert.That(first("is:"),
		Equals(Item{IDENTIFIER, NIL, "is"}))
	assert.That(first("\\"),
		Equals(Item{ERROR, NIL, "\\"}))

	check(assert, "f()", IDENTIFIER, L_PAREN, R_PAREN)

	check(assert, `and break 
		case catch continue class callback default dll do
		else for forever function if is isnt or not
		new switch struct super return throw try while
		true false 
		== = =~ ~ != !~ ! <<= << <> <= <
		>>= >> >= > || |= | && &= &
		^= ^ -- -= - ++ += + /= /
		*= * %= % $= $ name _name name123 'string'
		"string" 123 123name .name  Name Name123 name? 1$2 +1 num=1
		num+=1 1%2 /*comments*/ //comments`,
		AND, BREAK, CASE, CATCH, CONTINUE, CLASS, CALLBACK, DEFAULT, DLL, DO,
		ELSE, FOR, FOREVER, FUNCTION, IF, IS, ISNT, OR, NOT,
		NEW, SWITCH, STRUCT, SUPER, RETURN, THROW, TRY, WHILE,
		TRUE, FALSE,
		IS, EQ, MATCH, BITNOT, ISNT, MATCHNOT, NOT, LSHIFTEQ, LSHIFT,
		ISNT, LTE, LT, RSHIFTEQ, RSHIFT, GTE, GT, OR, BITOREQ, BITOR,
		AND, BITANDEQ, BITAND, BITXOREQ, BITXOR, DEC, SUBEQ, SUB, INC,
		ADDEQ, ADD, DIVEQ, DIV, MULEQ, MUL, MODEQ, MOD, CATEQ, CAT,
		IDENTIFIER, IDENTIFIER, IDENTIFIER, STRING, STRING, NUMBER,
		NUMBER, IDENTIFIER, DOT, IDENTIFIER, IDENTIFIER, IDENTIFIER,
		IDENTIFIER, NUMBER, CAT, NUMBER, ADD, NUMBER, IDENTIFIER,
		EQ, NUMBER, IDENTIFIER, ADDEQ, NUMBER, NUMBER, MOD, NUMBER,
		COMMENT, COMMENT)
}

func first(src string) Item {
	return Lexer(src).Next()
}

func check(assert Asserter, source string, expected ...Token) {
	lexer := Lexer(source)
	for i := 0; i < len(expected); {
		item := lexer.Next()
		if item.token == EOF {
			break
		} else if item.token == WHITESPACE || item.token == NEWLINE {
			continue
		}
		tok := item.token
		if item.keyword != NIL {
			tok = item.keyword
		} else {
			tok = item.token
		}
		assert.That(tok, Equals(expected[i]).Comment(i, item))
		i++
	}
}
