// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hll

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestEmpty(t *testing.T) {
	h := New()
	assert.T(t).This(h.Count()).Is(uint64(0))
}

func TestDuplicate(t *testing.T) {
	h := New()
	for range 1000 {
		h.Add("same")
	}
	c := h.Count()
	assert.T(t).That(c >= 1 && c <= 2)
}

func TestCardinality(t *testing.T) {
	test := func(n int, tol float64) {
		t.Helper()
		h := New()
		for i := range n {
			h.Add(fmt.Sprintf("v%08d", i))
		}
		got := float64(h.Count())
		want := float64(n)
		err := abs(got-want) / want
		assert.T(t).Msg("n", n, "got", got, "err", err).That(err <= tol)
	}

	test(1_000, 0.05)
	test(10_000, 0.03)
	test(100_000, 0.03)
}

func TestWithDuplicates(t *testing.T) {
	h := New()
	for i := range 50_000 {
		s := fmt.Sprintf("v%06d", i)
		h.Add(s)
		h.Add(s)
		h.Add(s)
	}
	got := float64(h.Count())
	want := float64(50_000)
	err := abs(got-want) / want
	assert.T(t).Msg("got", got, "err", err).That(err <= 0.03)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
