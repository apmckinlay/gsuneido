// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	"testing"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestKeywords(t *testing.T) {
	test := func(id string, expected tok.Token) {
		t.Helper()
		tok, val := keyword(id)
		assert.T(t).This(tok).Is(expected)
		assert.T(t).This(val).Is(id)
	}
	test("return", tok.Return)
	test("forever", tok.Forever)
	test("foo", tok.Identifier)
}

func TestLexer(t *testing.T) {
	assert := assert.T(t).This
	first := func(src string, text string, token tok.Token) {
		t.Helper()
		assert(NewLexer(src).Next()).Is(Item{Text: text, Pos: 0, Token: token})
	}
	first("function", "function", tok.Function)
	first("foo", "foo", tok.Identifier)
	first("#foo", "foo", tok.Symbol)
	first("#_foo?", "_foo?", tok.Symbol)
	first("is", "is", tok.Is)
	first("is:", "is", tok.Identifier)
	first("0xff", "0xff", tok.Number)
	first("0xff.Chr()", "0xff", tok.Number)
	first("0x8002 //foo", "0x8002", tok.Number)
	first("1000", "1000", tok.Number)
	first("1_2", "12", tok.Number)
	first("1_000", "1000", tok.Number)
	first("1_000_000", "1000000", tok.Number)
	first("123_456_789", "123456789", tok.Number)
	first("1_2.3", "12.3", tok.Number)
	first("1_2.3_4", "12.34", tok.Number)
	first("12.34_5", "12.345", tok.Number)
	first("9_", "9", tok.Number)
	first("9._2", "9.2", tok.Number)
	first("9.3_", "9.3", tok.Number)
	first("'hello'", "hello", tok.String)
	first("'hello", "hello", tok.String)
	first("`hello`", "hello", tok.String)
	first("`hello", "hello", tok.String)
	first("'foo\\'bar'", "foo'bar", tok.String)
	first(`"\"foo\""`, `"foo"`, tok.String)
	first("\\", "\\", tok.Error)
	first("//foo\r\nbar", "//foo", tok.Comment) // not including \r\n
	first("'", "", tok.String)
	first("\"", "", tok.String)
	first("`", "", tok.String)

	check := func(source string, expected ...tok.Token) {
		t.Helper()
		lexer := NewLexer(source)
		for i := 0; i < len(expected); {
			item := lexer.Next()
			if item.Token == tok.Eof {
				assert(i).Msg("too few tokens").Is(len(expected) - 1)
				break
			} else if item.Token == tok.Whitespace || item.Token == tok.Newline {
				continue
			}
			assert(item.Token).Msg(i, item).Is(expected[i])
			i++
		}
		assert(lexer.Next().Token).Msg("didn't consume input").Is(tok.Eof)
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
	assert := assert.T(t).This
	lxr := NewLexer("a \n= /**/ 1 ")
	assert(lxr.Ahead(0)).Is(it(tok.Identifier, 0, "a"))
	assert(lxr.Ahead(6)).Is(it(tok.Number, 10, "1"))
	assert(lxr.Ahead(2)).Is(it(tok.Eq, 3, "="))
	assert(lxr.Ahead(8).Token).Is(tok.Eof)

	assert(lxr.Next()).Is(it(tok.Identifier, 0, "a"))
	assert(lxr.Next()).Is(it(tok.Newline, 1, " \n"))
	assert(lxr.Next()).Is(it(tok.Eq, 3, "="))
	assert(lxr.Next()).Is(it(tok.Whitespace, 4, " "))
	assert(lxr.Next()).Is(it(tok.Comment, 5, "/**/"))
	assert(lxr.Next()).Is(it(tok.Whitespace, 9, " "))
	assert(lxr.Next()).Is(it(tok.Number, 10, "1"))
	assert(lxr.Next()).Is(it(tok.Whitespace, 11, " "))
	assert(lxr.Next().Token).Is(tok.Eof)
}

func TestAheadSkip(t *testing.T) {
	assert := assert.T(t).This
	lxr := NewLexer(" a \n= /**/ 1 ")
	assert(lxr.AheadSkip(0)).Is(it(tok.Identifier, 1, "a"))
	assert(lxr.AheadSkip(2)).Is(it(tok.Number, 11, "1"))
	assert(lxr.AheadSkip(1)).Is(it(tok.Eq, 4, "="))
	assert(lxr.AheadSkip(3).Token).Is(tok.Eof)
}

func TestEscape(t *testing.T) {
	lexer := NewLexer(`"\x80"`)
	item := lexer.Next()
	s := item.Text
	assert.T(t).This(len(s)).Is(1)
	assert.T(t).This(s[0]).Is(byte(128))
}

func TestQueryKeyword(t *testing.T) {
	for k, t := range queryKeywords {
		t2, k2 := queryKeyword(k)
		assert.This(t2).Is(t)
		assert.This(k2).Is(k)
	}
}

var queryKeywords = map[string]tok.Token{
	"alter":     tok.Alter,
	"and":       tok.And,
	"average":   tok.Average,
	"by":        tok.By,
	"cascade":   tok.Cascade,
	"class":     tok.Class,
	"count":     tok.Count,
	"create":    tok.Create,
	"delete":    tok.Delete,
	"destroy":   tok.Drop,
	"drop":      tok.Drop,
	"ensure":    tok.Ensure,
	"extend":    tok.Extend,
	"false":     tok.False,
	"function":  tok.Function,
	"history":   tok.History,
	"in":        tok.In,
	"index":     tok.Index,
	"insert":    tok.Insert,
	"intersect": tok.Intersect,
	"into":      tok.Into,
	"is":        tok.Is,
	"isnt":      tok.Isnt,
	"join":      tok.Join,
	"key":       tok.Key,
	"leftjoin":  tok.Leftjoin,
	"list":      tok.List,
	"max":       tok.Max,
	"min":       tok.Min,
	"minus":     tok.Minus,
	"not":       tok.Not,
	"or":        tok.Or,
	"project":   tok.Project,
	"remove":    tok.Remove,
	"rename":    tok.Rename,
	"reverse":   tok.Reverse,
	"set":       tok.Set,
	"sort":      tok.Sort,
	"summarize": tok.Summarize,
	"sview":     tok.Sview,
	"tempindex": tok.TempIndex,
	"times":     tok.Times,
	"to":        tok.To,
	"total":     tok.Total,
	"true":      tok.True,
	"union":     tok.Union,
	"unique":    tok.Unique,
	"update":    tok.Update,
	"view":      tok.View,
	"where":     tok.Where,
}
