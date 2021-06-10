// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestParseAdmin(t *testing.T) {
	test := func(cmd string) {
		t.Helper()
		admin := ParseAdmin(cmd)
		assert.T(t).This(admin.String()).Is(cmd)
	}
	test("drop mytable")

	test("rename mytable to newtable")

	test("create mytable (one,two,three) key()")
	test("create mytable (one,two,three) key(one)")
	test("create mytable (one,two,three) key(one,two)")

	test("ensure mytable index(one,two)")
	test("ensure mytable (one,two,three) index(one,two)")
	test("ensure mytable (one,two,three) index unique(one,two)")

	test("ensure mytable (one,two,three) index(two) in other")
	test("ensure mytable (one,two,three) index(two) in other cascade")
	test("ensure mytable (one,two,three) index(two) in other cascade update")
	test("ensure mytable (one,two,three) index(two) in other(six)")
	test("ensure mytable (one,two,three) index(two) in other(six) cascade")
	test("ensure mytable (one,two,three) index(two) in other(six) cascade update")

	test("create mytable (one,Two,Three) key(one)")
	test("create mytable (one,two,two_lower!) key(two_lower!)")

	test("alter mytable drop (one,two,three) index(two)")
	test("alter mytable create (one,two,three) index(two)")
	test("alter mytable rename one to two, three to four")

	test("view tc = tables join columns")

	xtest := func(cmd, err string) {
		t.Helper()
		fn := func() { ParseAdmin(cmd) }
		assert.T(t).This(fn).Panics(err)
	}
	xtest("create mytable () key(foo)", "invalid index column: foo")
	xtest("create mytable (one,two,three) index(one)", "key required")
	xtest("create mytable (one,two,three) key(bar)", "invalid index column: bar")
	xtest("create mytable (one,two,three_lower!) key(one)",
		"_lower! nonexistent column: three")
}
