// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package csio

import (
	"bytes"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestInt(t *testing.T) {
	var buf bytes.Buffer
	rw := NewReadWrite(&buf)
	test := func(n int64) {
		rw.PutInt64(n)
		rw.Flush()
		assert.T(t).This(rw.GetInt64()).Is(n)
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
		assert.T(t).This(rw.GetStr()).Is(s)
		buf.Reset()
	}
	test("")
	test("hello world")
	test("now is the time for all good men to come to the aid of their party")
}
