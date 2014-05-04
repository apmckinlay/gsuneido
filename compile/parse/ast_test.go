package parse

import "fmt"
import "strings"
import "github.com/apmckinlay/gsuneido/compile/lex"

func ExampleAst_String() {
	fmt.Println(&AstNode{})
	fmt.Println(&AstNode{Token: lex.COMMENT, Value: "/* ... */"})
	fmt.Println(&AstNode{Token: lex.IDENTIFIER, Value: "foo"})
	a := AstNode{lex.ADD, "+", []AstNode{
		AstNode{Token: lex.IDENTIFIER, Value: "foo"},
		AstNode{Token: lex.NUMBER, Value: "123"},
	}}
	longid := AstNode{Token: lex.IDENTIFIER,
		Value: strings.Repeat("verylong", 10)}
	fmt.Println(&a)
	fmt.Println(&AstNode{lex.MUL, "*", []AstNode{
		a,
		AstNode{lex.DIV, "/", []AstNode{
			AstNode{Token: lex.NUMBER, Value: "123"}, longid}},
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
