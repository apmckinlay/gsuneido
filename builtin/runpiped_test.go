// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSplitCommand(t *testing.T) {
	test := func(input string, expected []string) {
		t.Helper()
		result := SplitCommand(input)
		assert.T(t).This(result).Is(expected)
	}
	test("", []string{})
	test("cmd", []string{"cmd"})
	test("cmd arg1", []string{"cmd", "arg1"})
	test("cmd arg1 arg2 arg3", []string{"cmd", "arg1", "arg2", "arg3"})

	// Whitespace
	test("   \t  ", []string{})
	test("   cmd arg", []string{"cmd", "arg"})
	test("cmd arg   ", []string{"cmd", "arg"})
	test("cmd    arg1   arg2", []string{"cmd", "arg1", "arg2"})
	test("\tcmd\t arg1  \targ2\t", []string{"cmd", "arg1", "arg2"})

	// Quoted
	test(`cmd "arg with spaces"`, []string{"cmd", "arg with spaces"})
	test(`cmd "first arg" "second arg"`, []string{"cmd", "first arg", "second arg"})
	test(`cmd arg1 "arg with spaces" arg2`, []string{"cmd", "arg1", "arg with spaces", "arg2"})
	test(`cmd "" arg`, []string{"cmd", "", "arg"})
	test(`"quoted cmd" arg`, []string{"quoted cmd", "arg"})
	test(`cmd "last arg"`, []string{"cmd", "last arg"})
	test(`"only quoted"`, []string{"only quoted"})
}
