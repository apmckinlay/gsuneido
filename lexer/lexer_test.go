package lexer

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestKeywords(t *testing.T) {
	test := func(id string, expected Token) {
		t.Helper()
		tok, val := Keyword(id)
		Assert(t).That(tok, Equals(expected))
		Assert(t).That(val, Equals(id))
	}
	test("return", RETURN)
	test("forever", FOREVER)
	test("foo", IDENTIFIER)
}

func TestLexer(t *testing.T) {
	first := func(src string, text string, tok Token) {
		t.Helper()
		Assert(t).That(NewLexer(src).Next(),
			Equals(Item{text, 0, tok}))
	}
	first("function", "function", FUNCTION)
	first("foo", "foo", IDENTIFIER)
	first("#foo", "foo", STRING)
	first("#_foo?", "_foo?", STRING)
	first("is", "is", IS)
	first("is:", "is", IDENTIFIER)
	first("0xff", "0xff", NUMBER)
	first("0xff.Chr()", "0xff", NUMBER)
	first("0x8002 //foo", "0x8002", NUMBER)
	first("'hello'", "hello", STRING)
	first("'hello", "hello", STRING)
	first("`hello`", "hello", STRING)
	first("`hello", "hello", STRING)
	first("'foo\\'bar'", "foo'bar", STRING)
	first(`"\"foo\""`, `"foo"`, STRING)
	first("\\", "\\", ERROR)
	first("//foo\nbar", "//foo", COMMENT) // not including newline

	check := func(source string, expected ...Token) {
		t.Helper()
		lexer := NewLexer(source)
		for i := 0; i < len(expected); {
			item := lexer.Next()
			if item.Token == EOF {
				Assert(t).That(i, Equals(len(expected)-1).Comment("too few tokens"))
				break
			} else if item.Token == WHITESPACE || item.Token == NEWLINE {
				continue
			}
			Assert(t).That(item.Token, Equals(expected[i]).Comment(i, item))
			i++
		}
		Assert(t).That(lexer.Next().Token, Equals(EOF).Comment("didn't consume input"))
	}
	check("f()", IDENTIFIER, L_PAREN, R_PAREN)
	check("4-1", NUMBER, SUB, NUMBER)
	check("[1..]", L_BRACKET, NUMBER, RANGETO, R_BRACKET)
	check("#20181112.End", HASH, NUMBER, DOT, IDENTIFIER)
	check("0xff.Chr", NUMBER, DOT, IDENTIFIER)
	check("//foo\n0x8002 //bar", COMMENT, NUMBER, COMMENT)
	check(`and break
		case catch continue class default do
		else for forever function if is isnt or not
		new switch super return throw try while
		true false
		= =~ ~ !~ <<= << <> <= <
		>>= >> >= > |= | &= &
		^= ^ -- -= - ++ += + /= /
		*= * %= % $= $ name _name name123 'single'
		"double" 123 123name .name  Name Name123 name? 1$2 +1 num=1
		num+=1 1%2 /*comments*/ //comments`,
		AND, BREAK, CASE, CATCH, CONTINUE, CLASS, DEFAULT, DO,
		ELSE, FOR, FOREVER, FUNCTION, IF, IS, ISNT, OR, NOT,
		NEW, SWITCH, SUPER, RETURN, THROW, TRY, WHILE,
		TRUE, FALSE,
		EQ, MATCH, BITNOT, MATCHNOT, LSHIFTEQ, LSHIFT,
		ISNT, LTE, LT, RSHIFTEQ, RSHIFT, GTE, GT, BITOREQ, BITOR,
		BITANDEQ, BITAND, BITXOREQ, BITXOR, DEC, SUBEQ, SUB, INC,
		ADDEQ, ADD, DIVEQ, DIV, MULEQ, MUL, MODEQ, MOD, CATEQ, CAT,
		IDENTIFIER, IDENTIFIER, IDENTIFIER, STRING, STRING, NUMBER,
		NUMBER, IDENTIFIER, DOT, IDENTIFIER, IDENTIFIER, IDENTIFIER,
		IDENTIFIER, NUMBER, CAT, NUMBER, ADD, NUMBER, IDENTIFIER,
		EQ, NUMBER, IDENTIFIER, ADDEQ, NUMBER, NUMBER, MOD, NUMBER,
		COMMENT, COMMENT)
}

func TestAhead(t *testing.T) {
	lxr := NewLexer("a \n= /**/ 1 ")
	Assert(t).That(lxr.Ahead(0), Equals(it(IDENTIFIER, 0, "a")))
	Assert(t).That(lxr.Ahead(6), Equals(it(NUMBER, 10, "1")))
	Assert(t).That(lxr.Ahead(2), Equals(it(EQ, 3, "=")))
	Assert(t).That(lxr.Ahead(8).Token, Equals(EOF))

	Assert(t).That(lxr.Next(), Equals(it(IDENTIFIER, 0, "a")))
	Assert(t).That(lxr.Next(), Equals(it(NEWLINE, 1, " \n")))
	Assert(t).That(lxr.Next(), Equals(it(EQ, 3, "=")))
	Assert(t).That(lxr.Next(), Equals(it(WHITESPACE, 4, " ")))
	Assert(t).That(lxr.Next(), Equals(it(COMMENT, 5, "/**/")))
	Assert(t).That(lxr.Next(), Equals(it(WHITESPACE, 9, " ")))
	Assert(t).That(lxr.Next(), Equals(it(NUMBER, 10, "1")))
	Assert(t).That(lxr.Next(), Equals(it(WHITESPACE, 11, " ")))
	Assert(t).That(lxr.Next().Token, Equals(EOF))
}

func TestAheadSkip(t *testing.T) {
	lxr := NewLexer(" a \n= /**/ 1 ")
	Assert(t).That(lxr.AheadSkip(0), Equals(it(IDENTIFIER, 1, "a")))
	Assert(t).That(lxr.AheadSkip(2), Equals(it(NUMBER, 11, "1")))
	Assert(t).That(lxr.AheadSkip(1), Equals(it(EQ, 4, "=")))
	Assert(t).That(lxr.AheadSkip(3).Token, Equals(EOF))
}

func TestEscape(t *testing.T) {
	lexer := NewLexer(`"\x80"`)
	item := lexer.Next()
	s := item.Text
	Assert(t).That(len(s), Equals(1))
	Assert(t).That(s[0], Equals(byte(128)))
}
