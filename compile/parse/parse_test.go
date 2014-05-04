package parse

import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func TestParse(t *testing.T) {
	result := Parse("function () { 1 + 2 }").(AstNode)
	Assert(t).That(result.String(), Equals("(+ 1 2)"))
}
