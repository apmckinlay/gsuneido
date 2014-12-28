package compile

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestKeywords(t *testing.T) {
	Assert(t).That(Keyword("forever"), Equals(FOREVER))
}

func TestLexer(t *testing.T) {
	first := func(src string, text string, id, kw Token) {
		Assert(t).That(NewLexer(src).Next(),
			Equals(Item{text, 0, id, kw}))
	}
	first("function", "function", IDENTIFIER, FUNCTION)
	first("foo", "foo", IDENTIFIER, NIL)
	first("is", "is", IDENTIFIER, IS)
	first("is:", "is", IDENTIFIER, NIL)
	first("0xff", "0xff", NUMBER, NIL)
	first("'hello'", "hello", STRING, NIL)
	first("'foo\\'bar'", "foo'bar", STRING, STRING)
	first("\\", "\\", ERROR, NIL)
	first("//foo\nbar", "//foo", COMMENT, NIL) // not including newline

	check := func(source string, expected ...Token) {
		lexer := NewLexer(source)
		for i := 0; i < len(expected); {
			item := lexer.Next()
			if item.Token == EOF {
				Assert(t).That(i, Equals(len(expected)-1).Comment("too few tokens"))
				break
			} else if item.Token == WHITESPACE || item.Token == NEWLINE {
				continue
			}
			Assert(t).That(item.KeyTok(), Equals(expected[i]).Comment(i, item))
			i++
		}
	}
	check("f()", IDENTIFIER, L_PAREN, R_PAREN)
	check("4-1", NUMBER, SUB, NUMBER)
	check("[1..]", L_BRACKET, NUMBER, RANGETO, R_BRACKET)
	check(`and break
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

func TestAhead(t *testing.T) {
	lxr := NewLexer("a=1")
	Assert(t).That(lxr.Ahead(0), Equals(it(IDENTIFIER, 0, "a")))
	Assert(t).That(lxr.Ahead(2), Equals(it(NUMBER, 2, "1")))
	Assert(t).That(lxr.Ahead(1), Equals(it(EQ, 1, "=")))
	Assert(t).That(lxr.Ahead(3).Token, Equals(EOF))

	Assert(t).That(lxr.Next(), Equals(it(IDENTIFIER, 0, "a")))
	Assert(t).That(lxr.Next(), Equals(it(EQ, 1, "=")))
	Assert(t).That(lxr.Next(), Equals(it(NUMBER, 2, "1")))
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
