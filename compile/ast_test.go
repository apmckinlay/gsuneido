package compile

import (
	"fmt"
	"strings"
)

func ExampleAst_String() {
	ast := func(tok Token, txt string, children ...Ast) Ast {
		return Ast{Item: Item{Token: tok, Text: txt}, Children: children}
	}
	a := Ast{}
	fmt.Println(&a)
	a = ast(COMMENT, "/* ... */")
	fmt.Println(&a)
	a = ast(IDENTIFIER, "foo")
	fmt.Println(&a)
	a = ast(ADD, "+",
		ast(IDENTIFIER, "foo"),
		ast(NUMBER, "123"))
	fmt.Println(&a)
	longid := ast(IDENTIFIER, strings.Repeat("verylong", 10))
	a = ast(MUL, "*",
		a,
		ast(DIV, "/",
			ast(NUMBER, "123"), longid),
		longid,
	)
	fmt.Println(&a)
	// Output:
	// NIL
	// COMMENT
	// foo
	// (+ foo 123)
	// (*
	//     (+ foo 123)
	//     (/
	//         123
	//         verylongverylongverylongverylongverylongverylongverylongverylongverylongverylong)
	//     verylongverylongverylongverylongverylongverylongverylongverylongverylongverylong)
}
