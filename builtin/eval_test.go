package builtin

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestIsGlobal(t *testing.T) {
	Assert(t).True(isGlobal("F"))
	Assert(t).True(isGlobal("Foo"))
	Assert(t).True(isGlobal("Foo_123_Bar"))
	Assert(t).True(isGlobal("Foo!"))
	Assert(t).True(isGlobal("Foo?"))

	Assert(t).False(isGlobal(""))
	Assert(t).False(isGlobal("f"))
	Assert(t).False(isGlobal("foo"))
	Assert(t).False(isGlobal("_foo"))
	Assert(t).False(isGlobal("Foo!bar"))
	Assert(t).False(isGlobal("Foo?bar"))
	Assert(t).False(isGlobal("Foo.bar"))
}
