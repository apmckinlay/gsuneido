// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestQueryParser(t *testing.T) {
	test := func(qs string) {
		sc := ParseRequest(qs).(*Schema)
		assert.T(t).This("create " + sc.String()).Is(qs)
	}
	test("create mytable (one,two,three) key()")
	test("create mytable (one,two,three) key(one)")
	test("create mytable (one,two,three) key(one,two)")
	test("create mytable (one,two,three) index(one,two)")
	test("create mytable (one,two,three) index unique(one,two)")

	test("create mytable (one,two,three) index(two) in other")
	test("create mytable (one,two,three) index(two) in other cascade")
	test("create mytable (one,two,three) index(two) in other cascade update")
	test("create mytable (one,two,three) index(two) in other (six)")
	test("create mytable (one,two,three) index(two) in other (six) cascade")
	test("create mytable (one,two,three) index(two) in other (six) cascade update")

	test("create mytable (one,Two,Three) key(one)")
	test("create mytable (one,two,two_lower!) key(two_lower!)")

	xtest := func(qs, err string) {
		fn := func() { ParseRequest(qs) }
		assert.T(t).This(fn).Panics(err)
	}
	xtest("create mytable () key(foo)", "invalid index column: foo")
	xtest("create mytable (one,two,three) index(bar)", "invalid index column: bar")
}
