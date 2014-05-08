package parse

import "fmt"
import "strings"
import "github.com/apmckinlay/gsuneido/compile/lex"

func ExampleAst_String() {
	ast := func(tok lex.Token, val string, children ...AstNode) AstNode {
		return AstNode{lex.Item{Token: tok, Value: val}, children}
	}
	a := AstNode{}
	fmt.Println(&a)
	a = ast(lex.COMMENT, "/* ... */")
	fmt.Println(&a)
	a = ast(lex.IDENTIFIER, "foo")
	fmt.Println(&a)
	a = ast(lex.ADD, "+",
		ast(lex.IDENTIFIER, "foo"),
		ast(lex.NUMBER, "123"))
	fmt.Println(&a)
	longid := ast(lex.IDENTIFIER, strings.Repeat("verylong", 10))
	a = ast(lex.MUL, "*",
		a,
		ast(lex.DIV, "/",
			ast(lex.NUMBER, "123"), longid),
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
