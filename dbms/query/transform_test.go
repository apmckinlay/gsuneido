// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTransform(t *testing.T) {
	test := func(from, to string) {
		t.Helper()
		if to == "" {
			to = from
		}
		q := ParseQuery(from).Transform()
		assert.T(t).This(q.String()).Is(to)
	}
	test("table", "")
	test("table1 join table2", "")
	test("table rename a to b, c to d", "")
	test("table rename a to b rename c to d rename e to f",
		"table rename a to b, c to d, e to f")
	test("table rename a to b rename b to c rename c to d",
		"table rename a to d")
	test("table extend a = 1, b = 2", "")
	test("table extend a = 1 extend b = 2",
		"table extend a = 1, b = 2")
}
