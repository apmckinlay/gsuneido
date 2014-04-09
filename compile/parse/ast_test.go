package parse

import "fmt"
import "strings"
import "gsuneido/compile/lex"

func ExampleAst_String() {
	fmt.Println(&AstNode{})
	fmt.Println(&AstNode{token: lex.COMMENT, value: "/* ... */"})
	fmt.Println(&AstNode{token: lex.IDENTIFIER, value: "foo"})
	a := AstNode{lex.ADD, "+", []AstNode{
		AstNode{token: lex.IDENTIFIER, value: "foo"},
		AstNode{token: lex.NUMBER, value: "123"},
	}}
	longid := AstNode{token: lex.IDENTIFIER,
		value: strings.Repeat("verylong", 10)}
	fmt.Println(&a)
	fmt.Println(&AstNode{lex.MUL, "*", []AstNode{
		a,
		AstNode{lex.DIV, "/", []AstNode{
			AstNode{token: lex.NUMBER, value: "123"}, longid}},
		longid,
	}})
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
