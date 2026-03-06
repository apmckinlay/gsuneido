// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestRangeTo(t *testing.T) {
	test := func(s string, from int, to int, expected string) {
		t.Helper()

		a := SuStr(s).RangeTo(from, to)
		assert.T(t).This(a).Is(SuStr(expected))

		list := strToList(s)
		expectedList := strToList(expected)
		actualList := list.RangeTo(from, to)
		assert.T(t).This(actualList).Is(expectedList)
	}

	test("hello world", 0, 0, "")
	test("hello world", 0, 5, "hello")
	test("hello world", -99, 5, "hello")
	test("hello world", 0, 99, "hello world")
	test("hello world", 0, -99, "")
	test("hello world", 6, 0, "")
	test("hello world", 6, 5, "")
	test("hello world", 6, 11, "world")
	test("hello world", 6, 99, "world")
	test("hello world", 6, -5, "")
	test("hello world", -5, 99, "world")
	test("hello world", -99, 99, "hello world")
	test("hello world", -11, -6, "hello")
}

func TestRangeLen(t *testing.T) {
	test := func(s string, from int, n int, expected string) {
		t.Helper()

		a := SuStr(s).RangeLen(from, n)
		assert.T(t).This(a).Is(SuStr(expected))

		list := strToList(s)
		expectedList := strToList(expected)
		actualList := list.RangeLen(from, n)
		assert.T(t).This(actualList).Is(expectedList)
	}

	test("", 0, 9999, "")
	test("", 9, 9999, "")
	test("", -9, 9999, "")

	test("hello world", 0, 9999, "hello world")
	test("hello world", -99, 9999, "hello world")
	test("hello world", 6, 9999, "world")
	test("hello world", -5, 9999, "world")
	test("hello world", 99, 9999, "")

	test("hello world", 0, 0, "")
	test("hello world", 0, 4, "hell")
	test("hello world", 3, 2, "lo")
	test("hello world", 0, -5, "")
	test("hello world", 2, -2, "")
}

func strToList(s string) *SuObject {
	ob := SuObject{}
	for _, c := range s {
		ob.Add(SuStr(string(c)))
	}
	return &ob
}
