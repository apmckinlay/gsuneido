// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	"testing"

	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestKeywords(t *testing.T) {
	test := func(id string, expected tok.Token) {
		t.Helper()
		tok, val := keyword(id)
		Assert(t).That(tok, Is(expected))
		Assert(t).That(val, Is(id))
	}
	test("return", tok.Return)
	test("forever", tok.Forever)
	test("foo", tok.Identifier)
}

func TestLexer(t *testing.T) {
	first := func(src string, text string, token tok.Token) {
		t.Helper()
		Assert(t).That(NewLexer(src).Next(),
			Is(Item{text, 0, token}))
	}
	first("function", "function", tok.Function)
	first("foo", "foo", tok.Identifier)
	first("#foo", "foo", tok.String)
	first("#_foo?", "_foo?", tok.String)
	first("is", "is", tok.Is)
	first("is:", "is", tok.Identifier)
	first("0xff", "0xff", tok.Number)
	first("0xff.Chr()", "0xff", tok.Number)
	first("0x8002 //foo", "0x8002", tok.Number)
	first("'hello'", "hello", tok.String)
	first("'hello", "hello", tok.String)
	first("`hello`", "hello", tok.String)
	first("`hello", "hello", tok.String)
	first("'foo\\'bar'", "foo'bar", tok.String)
	first(`"\"foo\""`, `"foo"`, tok.String)
	first("\\", "\\", tok.Error)
	first("//foo\nbar", "//foo", tok.Comment) // not including newline
	first("'", "", tok.String)
	first("\"", "", tok.String)
	first("`", "", tok.String)

	check := func(source string, expected ...tok.Token) {
		t.Helper()
		lexer := NewLexer(source)
		for i := 0; i < len(expected); {
			item := lexer.Next()
			if item.Token == tok.Eof {
				Assert(t).That(i, Is(len(expected)-1).Comment("too few tokens"))
				break
			} else if item.Token == tok.Whitespace || item.Token == tok.Newline {
				continue
			}
			Assert(t).That(item.Token, Is(expected[i]).Comment(i, item))
			i++
		}
		Assert(t).That(lexer.Next().Token, Is(tok.Eof).Comment("didn't consume input"))
	}
	check("f()", tok.Identifier, tok.LParen, tok.RParen)
	check("4-1", tok.Number, tok.Sub, tok.Number)
	check("[1..]", tok.LBracket, tok.Number, tok.RangeTo, tok.RBracket)
	check("#20181112.End", tok.Hash, tok.Number, tok.Dot, tok.Identifier)
	check("0xff.Chr", tok.Number, tok.Dot, tok.Identifier)
	check("//foo\n0x8002 //bar", tok.Comment, tok.Number, tok.Comment)
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
		tok.And, tok.Break, tok.Case, tok.Catch, tok.Continue, tok.Class,
		tok.Default, tok.Do, tok.Else, tok.For, tok.Forever, tok.Function,
		tok.If, tok.Is, tok.Isnt, tok.Or, tok.Not, tok.New, tok.Switch,
		tok.Super, tok.Return, tok.Throw, tok.Try, tok.While,
		tok.True, tok.False, tok.Eq, tok.Match, tok.BitNot, tok.MatchNot,
		tok.LShiftEq, tok.LShift, tok.Isnt, tok.Lte, tok.Lt,
		tok.RShiftEq, tok.RShift, tok.Gte, tok.Gt, tok.BitOrEq, tok.BitOr,
		tok.BitAndEq, tok.BitAnd, tok.BitXorEq, tok.BitXor, tok.Dec,
		tok.SubEq, tok.Sub, tok.Inc, tok.AddEq, tok.Add, tok.DivEq, tok.Div,
		tok.MulEq, tok.Mul, tok.ModEq, tok.Mod, tok.CatEq, tok.Cat,
		tok.Identifier, tok.Identifier, tok.Identifier, tok.String, tok.String,
		tok.Number, tok.Number, tok.Identifier, tok.Dot, tok.Identifier,
		tok.Identifier, tok.Identifier, tok.Identifier, tok.Number, tok.Cat,
		tok.Number, tok.Add, tok.Number, tok.Identifier, tok.Eq, tok.Number,
		tok.Identifier, tok.AddEq, tok.Number, tok.Number, tok.Mod, tok.Number,
		tok.Comment, tok.Comment)
}

func TestAhead(t *testing.T) {
	lxr := NewLexer("a \n= /**/ 1 ")
	Assert(t).That(lxr.Ahead(0), Is(it(tok.Identifier, 0, "a")))
	Assert(t).That(lxr.Ahead(6), Is(it(tok.Number, 10, "1")))
	Assert(t).That(lxr.Ahead(2), Is(it(tok.Eq, 3, "=")))
	Assert(t).That(lxr.Ahead(8).Token, Is(tok.Eof))

	Assert(t).That(lxr.Next(), Is(it(tok.Identifier, 0, "a")))
	Assert(t).That(lxr.Next(), Is(it(tok.Newline, 1, " \n")))
	Assert(t).That(lxr.Next(), Is(it(tok.Eq, 3, "=")))
	Assert(t).That(lxr.Next(), Is(it(tok.Whitespace, 4, " ")))
	Assert(t).That(lxr.Next(), Is(it(tok.Comment, 5, "/**/")))
	Assert(t).That(lxr.Next(), Is(it(tok.Whitespace, 9, " ")))
	Assert(t).That(lxr.Next(), Is(it(tok.Number, 10, "1")))
	Assert(t).That(lxr.Next(), Is(it(tok.Whitespace, 11, " ")))
	Assert(t).That(lxr.Next().Token, Is(tok.Eof))
}

func TestAheadSkip(t *testing.T) {
	lxr := NewLexer(" a \n= /**/ 1 ")
	Assert(t).That(lxr.AheadSkip(0), Is(it(tok.Identifier, 1, "a")))
	Assert(t).That(lxr.AheadSkip(2), Is(it(tok.Number, 11, "1")))
	Assert(t).That(lxr.AheadSkip(1), Is(it(tok.Eq, 4, "=")))
	Assert(t).That(lxr.AheadSkip(3).Token, Is(tok.Eof))
}

func TestEscape(t *testing.T) {
	lexer := NewLexer(`"\x80"`)
	item := lexer.Next()
	s := item.Text
	Assert(t).That(len(s), Is(1))
	Assert(t).That(s[0], Is(byte(128)))
}
