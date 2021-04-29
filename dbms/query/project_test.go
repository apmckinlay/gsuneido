// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPrefixed(t *testing.T) {
	p := Project{}
	oneval := []string{""}
	fixed := []Fixed{{col: "f1", values: oneval}, {col: "f2", values: oneval}}
	test := func(sidx, scols string) {
		t.Helper()
		idx := strings.Fields(sidx)
		cols := strings.Fields(scols)
		nu := countUnfixed(cols, fixed)
		assert.T(t).That(p.prefixed(fixed, idx, cols, nu))
		idx = append(idx, "x")
		assert.T(t).That(p.prefixed(fixed, idx, cols, nu))
		cols = append(cols, "y")
		assert.T(t).That(!p.prefixed(fixed, idx, cols, nu+1))
	}
	test("a", "a")
	test("a b", "a")
	test("a b", "b a")
	test("a f1", "f2 a")
	test("a f1 b f2", "a f1")
	test("a f1 b f2", "f1 b f2 a")
}
