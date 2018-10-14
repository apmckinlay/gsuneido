package compile

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestAstString(t *testing.T) {
	test := func(ast *Ast, expected string) {
		Assert(t).That(ast.String(), Equals(expected))
	}
	a := func(tok Token, txt string, children ...*Ast) *Ast {
		return &Ast{Item: Item{Token: tok, Text: txt}, Children: children}
	}
	test(&Ast{}, "()")
	test(&Ast{Item: Item{Token: EOF}}, "EOF")
	test(a(COMMENT, "/* ... */"), "COMMENT")
	test(a(IDENTIFIER, "foo"), "foo")
	test(&Ast{value: SuInt(123)}, "123")
	test(&Ast{Item: Item{Text: "foo"}}, "foo")
	test(&Ast{Item: Item{Text: "num"}, value: SuInt(123)}, "(num 123)")
}

func TestFold(t *testing.T) {
	Assert(t).That(fold(Item{Token: IDENTIFIER, Keyword: TRUE}, nil, nil),
		Equals(&Ast{value: True}))
}
