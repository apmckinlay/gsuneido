// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestUnorderedReq(t *testing.T) {
	r := UnorderedReq(0.5)
	assert.T(t).This(r.Use()).Is(ReqUnordered)
	assert.T(t).That(r.cols == nil)
	assert.T(t).This(r.frac).Is(float32(0.5))
	assert.T(t).This(r.nlookups).Is(int32(0))
}

func TestOrderedReq(t *testing.T) {
	r := OrderedReq([]string{"a", "b"}, 0.3)
	assert.T(t).This(r.Use()).Is(ReqOrdered)
	assert.T(t).This(r.cols).Is([]string{"a", "b"})
	assert.T(t).This(r.frac).Is(float32(0.3))
	assert.T(t).This(r.nlookups).Is(int32(0))

	r2 := OrderedReq([]string{}, 0.7)
	assert.T(t).This(r2.Use()).Is(ReqUnordered)
	assert.T(t).That(r2.cols == nil)
}

func TestGroupedReq(t *testing.T) {
	r := GroupedReq([]string{"a", "b"}, 0.5, 10)
	assert.T(t).This(r.Use()).Is(ReqGrouped)
	assert.T(t).This(r.cols).Is([]string{"a", "b"})
	assert.T(t).This(r.frac).Is(float32(0.5))
	assert.T(t).This(r.nlookups).Is(int32(10))

	r2 := GroupedReq([]string{}, 0.5, 5)
	assert.T(t).This(r2.Use()).Is(ReqUnordered)
	assert.T(t).That(r2.cols == nil)

	r3 := GroupedReq([]string{"x"}, 0.2, 0)
	assert.T(t).This(r3.Use()).Is(ReqGrouped)
	assert.T(t).This(r3.nlookups).Is(int32(1))

	r4 := GroupedReq([]string{"x"}, 0.2, -5)
	assert.T(t).This(r4.Use()).Is(ReqGrouped)
	assert.T(t).This(r4.nlookups).Is(int32(1))
}

func TestLookupReq(t *testing.T) {
	r := LookupReq([]string{"id"}, 100)
	assert.T(t).This(r.Use()).Is(ReqLookup)
	assert.T(t).This(r.cols).Is([]string{"id"})
	assert.T(t).This(r.frac).Is(float32(0))
	assert.T(t).This(r.nlookups).Is(int32(100))
}

func TestRequireUse(t *testing.T) {
	assert.T(t).This(Require{cols: nil, frac: 0, nlookups: 0}.Use()).Is(ReqUnordered)
	assert.T(t).This(Require{cols: []string{"x"}, frac: 0.5, nlookups: 0}.Use()).Is(ReqOrdered)
	assert.T(t).This(Require{cols: []string{"x"}, frac: 0.5, nlookups: 1}.Use()).Is(ReqGrouped)
	assert.T(t).This(Require{cols: []string{"x"}, frac: 0, nlookups: 1}.Use()).Is(ReqLookup)
}

func TestLookupCount(t *testing.T) {
	assert.T(t).This(Require{nlookups: 7}.LookupCount(100)).Is(int32(7))
	assert.T(t).This(Require{nlookups: 7}.LookupCount(0)).Is(int32(7))

	assert.T(t).This(Require{frac: 0.5}.LookupCount(200)).Is(int32(100))
	assert.T(t).This(Require{frac: 0.25}.LookupCount(10)).Is(int32(2))
	assert.T(t).This(Require{frac: 0.01}.LookupCount(1)).Is(int32(0))
	assert.T(t).This(Require{frac: 0.5}.LookupCount(0)).Is(int32(0))
	assert.T(t).This(Require{frac: 1}.LookupCount(50)).Is(int32(50))
}

func TestSelectFrac(t *testing.T) {
	assert.T(t).This(Require{frac: 0.5}.SelectFrac(100)).Is(float32(0.5))
	assert.T(t).This(Require{frac: 1}.SelectFrac(100)).Is(float32(1))

	assert.T(t).This(Require{nlookups: 50}.SelectFrac(100)).Is(float32(0.5))
	assert.T(t).This(Require{nlookups: 200}.SelectFrac(100)).Is(float32(1))
	assert.T(t).This(Require{nlookups: 10}.SelectFrac(0)).Is(float32(1)) // nlookups/max(1,0) = 10/1 = 10, min(1,10) = 1

	assert.T(t).This(Require{}.SelectFrac(100)).Is(float32(1)) // frac=0, nlookups=0 -> 1
}

func TestRequireString(t *testing.T) {
	assert.T(t).This(UnorderedReq(0).String()).Is("ReqUnordered()")

	r := OrderedReq([]string{"a", "b"}, 0.3)
	assert.T(t).This(r.String()).Is("ReqOrdered(a,b) f0.3 n0")

	r2 := GroupedReq([]string{"x"}, 0.5, 10)
	assert.T(t).This(r2.String()).Is("ReqGrouped(x) f0.5 n10")

	r3 := LookupReq([]string{"id"}, 1)
	assert.T(t).This(r3.String()).Is("ReqLookup(id) f0 n1")

	r4 := OrderedReq([]string{"a"}, 0)
	assert.T(t).This(r4.String()).Is("ReqOrdered(a)")

	r5 := OrderedReq([]string{}, 0)
	assert.T(t).This(r5.String()).Is("ReqUnordered()")
}

func TestNewRequire(t *testing.T) {
	r := NewRequire([]string{"a", "b"}, 0.5, 10)
	assert.T(t).This(r.cols).Is([]string{"a", "b"})
	assert.T(t).This(r.frac).Is(float32(0.5))
	assert.T(t).This(r.nlookups).Is(int32(10))
}
