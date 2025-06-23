// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTokens(t *testing.T) {
	tzr := NewInput("Now is (the) 1 time for 123foo bar456 abc999def 123zzz456 " +
		"US T4 T4A W-5 401k TD1X 123 343-8887 P&L")
	expected := []string{"time", "123foo", "123", "foo", "bar456", "bar", "456",
		"abc999def", "abc", "999", "def", "123zzz456", "123", "zzz", "456", "us",
		"t4", "t4a", "w-5", "401k", "401", "td1x", "td", "123",
		"343-8887", "343", "8887", "p&l"}
	e := 0
	for tok := tzr.Next(); tok != ""; tok = tzr.Next() {
		assert.This(tok).Is(expected[e])
		assert.This(NewInput(tok).Next()).Is(tok)
		e++
	}
	assert.That(e == len(expected))

	assert.This(NewInput("T4A").Next()).Is("t4a")

	in := NewInput("foo " + strings.Repeat("a-5", 99) + " bar")
	assert.This(in.Next()).Is("foo")
	assert.This(in.Next()).Is("bar")
	assert.This(in.Next()).Is("")
}
