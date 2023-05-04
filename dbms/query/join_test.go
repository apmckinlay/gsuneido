// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestJoin_nrows(t *testing.T) {
	test := func(n1, p1, n2, p2, expected int) {
		t.Helper()
		j1n := Join{}
		j1n.joinType = one_n
		jn1 := Join{}
		jn1.joinType = n_one
		assert.T(t).Msg(n1, "/", p1, one_n, n2, "/", p2, "=>", expected).
			This(j1n.nrows(n1, p1, n2, p2)).Is(expected)
		assert.T(t).Msg(n1, "/", p1, n_one, n2, "/", p2, "=>", expected).
			This(jn1.nrows(n2, p2, n1, p1)).Is(expected)
	}
	test(0, 100, 2000, 2000, 0)
	test(1, 100, 2000, 2000, 20)
	test(2, 100, 2000, 2000, 40)
	test(90, 100, 2000, 2000, 1800)
	test(100, 100, 2000, 2000, 2000)

	test(100, 100, 0, 2000, 0)
	test(100, 100, 1, 2000, 1)
	test(100, 100, 10, 2000, 10)
	test(100, 100, 200, 2000, 200)
	test(100, 100, 1800, 2000, 1800)

	test(2, 100, 200, 2000, 40)
	test(2, 100, 10, 2000, 10)
}
