// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile/lexer"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestEscape(t *testing.T) {
	which := 0
	test := func(s, want string) {
		t.Helper()
		got := escapeStr(s, which)
		if want != "" && got != want {
			t.Errorf("%d %q want %q got %q", which, s, want, got)
		}
		lxr := lexer.NewLexer(got)
		it := lxr.Next()
		assert.This(it.Token == tokens.String)
		if it.Text != s {
			t.Errorf("%d want %q got %q", which, s, it.Text)
		}
	}
	sq := func(s string) string {
		return string(singleQuote) + s + string(singleQuote)
	}
	dq := func(s string) string {
		return string(doubleQuote) + s + string(doubleQuote)
	}
	bq := func(s string) string {
		return string(backQuote) + s + string(backQuote)
	}
	test("", `""`)
	test("x", sq("x")) // prefer single quotes for single characters
	test("'", dq("'"))
	test(`\`, bq(`\`))
	test("foo", dq("foo"))
	test("foo'bar", dq("foo'bar"))
	test(`foo"bar`, sq(`foo"bar`))
	test(`foo'"bar`, bq(`foo'"bar`))
	test("foo'\"\tbar", dq(`foo'\"\tbar`))
	test(`\/`, bq(`\/`))
	test("\x00\x01\x05\x0c\x0f\xf0\xff", dq(`\x00\x01\x05\x0c\x0f\xf0\xff`))
	test("\n", sq(`\n`))
	which = 2
	test(`\`, dq(`\\`))
}
