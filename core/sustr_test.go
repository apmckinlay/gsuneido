// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"math/rand"
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
			t.Errorf("%d\nwant %x\n got %x\n via %s", which, s, it.Text, got)
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
	test("\x19\x09", dq(`\x19\t`))
	which = 2
	test(`\`, dq(`\\`))

	for which = 0; which <= 2; which++ {
		// every 1 and 2 byte combination
		buf := make([]byte, 2)
		for i := range 256 {
			buf[0] = byte(i)
			test(string(buf[:1]), "")
			for j := range 256 {
				buf[1] = byte(j)
				test(string(buf), "")
			}
		}

		N := 10000
		if testing.Short() {
			N = 100
		}

		// random 8 byte strings
		buf = make([]byte, 8)
		for range N {
			for i := range buf {
				buf[i] = byte(rand.Intn(256))
			}
			test(string(buf), "")
		}

		// shuffles of all possible bytes
		buf = make([]byte, 256)
		for i := range 256 {
			buf[i] = byte(i)
		}
		for range N {
			rand.Shuffle(256, func(i, j int) { buf[i], buf[j] = buf[j], buf[i] })
			test(string(buf), "")
		}
	}
}

func TestSuStr1(t *testing.T) {
	for i := range 256 {
		assert.This(string(SuStr1s[i].(SuStr))[0]).Is(i)
	}
	test := func(s string) {
		assert.This(SuStr1(s)).Is(SuStr(s))
	}
	test("")
	test("x")
	test("foo")
}

var M Value

func BenchmarkSuStr(b *testing.B) {
	s := "x"
	for b.Loop() {
		M = SuStr(s)
	}
}

func BenchmarkSuStr1(b *testing.B) {
	s := "x"
	for b.Loop() {
		M = SuStr1(s)
	}
}
