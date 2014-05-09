package compile

import "fmt"

func ExampleToken_String() {
	fmt.Println(NIL, EOF, WHITESPACE, COMMENT, NEWLINE)
	// Output:
	// NIL EOF WHITESPACE COMMENT NEWLINE
}
