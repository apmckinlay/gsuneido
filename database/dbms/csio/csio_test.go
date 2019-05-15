package csio

import (
	"bytes"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestInt(t *testing.T) {
	var buf bytes.Buffer
	rw := NewReadWrite(&buf)
	test := func(n int64) {
		rw.PutInt64(n)
		rw.Flush()
		Assert(t).That(rw.GetInt64(), Equals(n))
		buf.Reset()
	}
	test(0)
	test(23)
	test(-23)
	test(123456)
	test(-123456)
	test(0xffffffff)
}

func TestStr(t *testing.T) {
	var buf bytes.Buffer
	rw := NewReadWrite(&buf)
	test := func(s string) {
		rw.PutStr(s)
		rw.Flush()
		Assert(t).That(rw.GetStr(), Equals(s))
		buf.Reset()
	}
	test("")
	test("hello world")
	test("now is the time for all good men to come to the aid of their party")

	rw.PutInt(0xffffff)
	rw.Flush()
	Assert(t).That(func () { rw.GetStr() }, Panics("bad io size"))
}
