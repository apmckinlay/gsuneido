package compile

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	v "github.com/apmckinlay/gsuneido/value"
)

func TestConstant(t *testing.T) {
	Assert(t).That(Constant("true"), Equals(v.True))
	Assert(t).That(Constant("false"), Equals(v.False))
	Assert(t).That(Constant("0"), Equals(v.SuInt(0)))
	Assert(t).That(Constant("-123"), Equals(v.SuInt(-123)))
	Assert(t).That(Constant("+456"), Equals(v.SuInt(456)))
	Assert(t).That(Constant("0xff"), Equals(v.SuInt(255)))
	Assert(t).That(Constant("0377"), Equals(v.SuInt(255)))
	Assert(t).That(Constant("'hi wo'"), Equals(v.SuStr("hi wo")))
	Assert(t).That(Constant("#20140425").String(), Equals("#20140425"))
	Assert(t).That(Constant("function () {}").TypeName(), Equals("Function"))
}
