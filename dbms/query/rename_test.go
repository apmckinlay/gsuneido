// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestRename_check(t *testing.T) {
	test := func(input string) {
		ParseQuery("table rename "+input, testTran{}, nil)
	}
	test("a to x")

	xtest := func(input, expected string) {
		defer func() {
			if r := recover(); r != nil {
				assert.T(t).This(fmt.Sprint(r)).Is(expected)
			}
		}()
		test(input)
		t.Error("expected " + expected)
	}
	xtest("a to b", "rename: column already exists: b")
	xtest("a to x, b to x", "rename: column already exists: x")
	xtest("x to y", "rename: nonexistent column: x")
	xtest("a to x, a to y", "rename: nonexistent column: a")
}

func TestRename_fwd(t *testing.T) {
	test := func(rename, slist, expected string) {
		r := ParseQuery("table rename "+rename, testTran{}, nil).(*Rename)
		list := strings.Split(slist, ",")
		result := r.renameFwd(list)
		assert.T(t).This(result).Is(strings.Split(expected, ","))
	}
	test("a to x", "x,y,z", "x,y,z")
	test("a to x", "a,b,c", "x,b,c")
	test("b to x, x to y, y to z", "a,b,c", "a,z,c")
}

func TestRename_rev(t *testing.T) {
	test := func(rename, slist, expected string) {
		r := ParseQuery("table rename "+rename, testTran{}, nil).(*Rename)
		list := strings.Split(slist, ",")
		result := r.renameRev(list)
		assert.T(t).This(result).Is(strings.Split(expected, ","))
	}
	test("a to x", "x,y,z", "a,y,z")
	test("a to x", "a,b,c", "a,b,c")
	test("b to x, x to y, y to z", "a,z,c", "a,b,c")
}
