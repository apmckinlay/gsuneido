// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestGrouped(t *testing.T) {
	oneval := []string{""}
	fixed := []Fixed{{col: "f1", values: oneval}, {col: "f2", values: oneval}}
	test := func(sidx, scols string) {
		t.Helper()
		idx := strings.Fields(sidx)
		cols := strings.Fields(scols)
		nu := countUnfixed(cols, fixed)
		assert.T(t).That(grouped(idx, cols, nu, fixed))
		idx = append(idx, "x")
		assert.T(t).That(grouped(idx, cols, nu, fixed))
		cols = append(cols, "y")
		assert.T(t).That(!grouped(idx, cols, nu+1, fixed))
	}
	test("a", "a")
	test("a b", "a")
	test("a b", "b a")
	test("a f1", "f2 a")
	test("a f1 b f2", "a f1")
	test("a f1 b f2", "f1 b f2 a")
}

func TestRemove(t *testing.T) {
	tbl := &Table{}
	cols := []string{"a", "a_deps", "b", "b_deps", "c", "c_deps",
		"d", "d_deps", "x_lower!"}
	tbl.header = SimpleHeader(cols)
	proj := NewRemove(tbl, []string{"a", "c"})
	assert.T(t).This(proj.columns).Is([]string{"b", "b_deps", "d", "d_deps"})
}
