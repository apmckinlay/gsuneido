package compile

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestFunction(t *testing.T) {
	result := ParseFunction("function () { 1 + 2 }")
	Assert(t).That(result.String(), Equals("(function () (STMTS (EXPR (+ 1 2))))"))
}
