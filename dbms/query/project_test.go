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

func TestProjectIndexes(t *testing.T) {
	test := func(idxs string, cols string, expected string) {
		t.Helper()
		result := projectIndexes(sToIdxs(idxs), strings.Fields(cols))
		assert.T(t).This(idxsToS(result)).Is(expected)
	}
	test("a, b+c, b+c+x, d+e+f", "a b c d", "a, b+c, d")
	test("a, b+c, d+e+f", "c e f", "")
	test("a, b+c, d+e+f", "", "")
}

// sToIdxs splits strings like: "a+b, c, d+e+f"
func sToIdxs(s string) [][]string {
	var idxs [][]string
	for _, ix := range strings.Split(s, ", ") {
		idxs = append(idxs, strings.Split(ix, "+"))
	}
	return idxs
}

// idxsToS converts [][]string to a string like: "a+b, c, d+e+f"
func idxsToS(idxs [][]string) string {
	tmp := make([]string, len(idxs))
	for i, ix := range idxs {
		tmp[i] = strings.Join(ix, "+")
	}
	return strings.Join(tmp, ", ")
}

func TestProjectKeys(t *testing.T) {
	test := func(idxs string, cols string, expected string) {
		t.Helper()
		result := projectKeys(sToIdxs(idxs), strings.Fields(cols))
		assert.T(t).This(idxsToS(result)).Is(expected)
	}
	test("a, b+c, b+c+x, d+e+f", "a b c d", "a, b+c")
	test("a, b+c, d+e+f", "c e f", "c+e+f") // fallback to all columns
}

func TestHasKey(t *testing.T) {
	var fixed []Fixed
	test := func(cols string, keys string, expected bool) {
		t.Helper()
		result := hasKey(strings.Fields(cols), sToIdxs(keys), fixed)
		assert.T(t).This(result).Is(expected)
	}
	// no fixed
	test("", "a+b, c", false)
	test("a b", "a", true)
	test("a b", "b", true)
	test("a b", "x, a+b, y", true)

	fixed = []Fixed{{col: "b", values: []string{"1"}}}
	test("", "a", false)
	test("", "b", true)
	test("a c", "a+b+c", true)
}
