// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

func TestReplace(t *testing.T) {
	th := &Thread{}

	test := func(input, pattern, repl string, count int, expected string) {
		t.Helper()
		patValue := SuStr(pattern)
		replValue := SuStr(repl)
		result := replace(th, input, patValue, replValue, count)
		assert.T(t).This(result).Is(expected)

		// For identical replacements, verify we get back the same string instance
		if expected == input {
			assert.T(t).Msg("same").That(hacks.SameString(input, result))
		}
	}

	// No change with identical replacement
	test("hello world", "hello", "hello", 1, "hello world")

	// No change with identical replacement using regex
	test("hello world", "h.llo", "hello", 1, "hello world")

	// Partial change with some identical replacements
	test("hello hello world hello", "hello", "hi", 2, "hi hi world hello")

	// No change with count=0
	test("hello world", "hello", "hi", 0, "hello world")

	// Empty pattern and replacement
	test("hello world", "", "", 1, "hello world")

	// Actual change
	test("hello world", "hello", "hi", 1, "hi world")

	// zero width match
	test("world", "^x?", "hello", 1, "helloworld")

	// Test with a callable replacement that returns the same string
	t.Run("callable replacement - no change", func(t *testing.T) {
		input := "hello world"

		// Create a mock callable that returns the same string
		identityFn := builtinVal("", func(arg Value) Value {
			return arg // Return the input unchanged
		}, "(arg)")

		result := replace(th, input, SuStr("hello"), identityFn, 1)
		assert.T(t).This(result).Is(input)
	})

	// Test with a callable replacement that returns a different string
	t.Run("callable replacement - with change", func(t *testing.T) {
		input := "hello world"
		expected := "HI world"

		// Create a mock callable that returns a different string
		upperFn := builtinVal("", func(Value) Value {
			return SuStr("HI")
		}, "(arg)")

		result := replace(th, input, SuStr("hello"), upperFn, 1)
		assert.T(t).This(result).Is(expected)
	})

	// Test multiple replacements with some identical
	t.Run("mixed identical and different replacements", func(t *testing.T) {
		input := "abc abc def abc"

		// Create a mock callable that conditionally changes the string
		n := 0
		conditionalFn := builtinVal("", func(arg Value) Value {
			s := ToStr(arg)
			n++
			if s == "abc" && n == 3 {
				return SuStr("xyz") // Change only the 3rd occurrence
			}
			return arg // Return unchanged otherwise
		}, "(arg)")

		result := replace(th, input, SuStr("abc"), conditionalFn, 99)
		assert.T(t).This(result).Is("abc abc def xyz")
	})
}
