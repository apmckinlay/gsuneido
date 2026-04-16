// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSelsGet(t *testing.T) {
	sels := Sels{{"a", "1"}, {"b", "2"}}

	v, ok := sels.Get("a")
	assert.T(t).This(v).Is("1")
	assert.T(t).This(ok).Is(true)

	v, ok = sels.Get("b")
	assert.T(t).This(v).Is("2")
	assert.T(t).This(ok).Is(true)

	v, ok = sels.Get("c")
	assert.T(t).This(v).Is("")
	assert.T(t).This(ok).Is(false)
}

func TestSelsMustGet(t *testing.T) {
	sels := Sels{{"a", "1"}}
	assert.T(t).This(sels.MustGet("a")).Is("1")
	assert.T(t).This(func() { sels.MustGet("c") }).Panics("Sels.Get can't find c")
}

func TestSelsHasCol(t *testing.T) {
	sels := Sels{{"a", "1"}, {"b", "2"}}
	assert.T(t).This(sels.HasCol("a")).Is(true)
	assert.T(t).This(sels.HasCol("c")).Is(false)
}

func TestSelsFindCol(t *testing.T) {
	sels := Sels{{"a", "1"}, {"b", "2"}}
	assert.T(t).This(sels.FindCol("b")).Is(1)
	assert.T(t).This(sels.FindCol("c")).Is(-1)
}

func TestSelsColsAre(t *testing.T) {
	sels := Sels{{"a", "1"}, {"b", "2"}}

	assert.T(t).This(sels.ColsAre([]string{"a", "b"})).Is(true)
	assert.T(t).This(sels.ColsAre([]string{"b", "a"})).Is(true) // order independent
	assert.T(t).This(sels.ColsAre([]string{"a"})).Is(false)     // length mismatch
	assert.T(t).This(sels.ColsAre([]string{"a", "c"})).Is(false)
}

func TestSelsAll(t *testing.T) {
	sels := Sels{{"a", "1"}, {"b", "2"}}
	var cols, vals []string
	for c, v := range sels.All() {
		cols = append(cols, c)
		vals = append(vals, v)
	}
	assert.T(t).This(cols).Is([]string{"a", "b"})
	assert.T(t).This(vals).Is([]string{"1", "2"})
}
