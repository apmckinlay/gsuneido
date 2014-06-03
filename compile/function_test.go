package compile

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestFunction(t *testing.T) {
	result := ParseFunction("function () { a + b }")
	Assert(t).That(result.String(), Equals("(function () (STMTS (+ a b)))"))
}
