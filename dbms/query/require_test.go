// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNoneReq(t *testing.T) {
	r := NoneReq(0.5)
	assert.T(t).This(r.use).Is(ReqNone)
	assert.T(t).That(r.cols == nil)
	assert.T(t).This(r.frac).Is(float32(0.5))
	assert.T(t).This(r.nseeks).Is(int32(0))
}

func TestOrderReq(t *testing.T) {
	r := OrderReq([]string{"a", "b"}, 0.3)
	assert.T(t).This(r.use).Is(ReqOrder)
	assert.T(t).This(r.cols).Is([]string{"a", "b"})
	assert.T(t).This(r.frac).Is(float32(0.3))
	assert.T(t).This(r.nseeks).Is(int32(0))

	r2 := OrderReq([]string{}, 0.7)
	assert.T(t).This(r2.use).Is(ReqNone)
	assert.T(t).That(r2.cols == nil)
}

func TestGroupReq(t *testing.T) {
	r := GroupReq([]string{"a", "b"}, 0.5, 10)
	assert.T(t).This(r.use).Is(ReqGroup)
	assert.T(t).This(r.cols).Is([]string{"a", "b"})
	assert.T(t).This(r.frac).Is(float32(0.5))
	assert.T(t).This(r.nseeks).Is(int32(10))

	r2 := GroupReq([]string{}, 0.5, 5)
	assert.T(t).This(r2.use).Is(ReqNone)
	assert.T(t).That(r2.cols == nil)
}

func TestUniqueReq(t *testing.T) {
	r := UniqueReq([]string{"id"}, 100)
	assert.T(t).This(r.use).Is(ReqUnique)
	assert.T(t).This(r.cols).Is([]string{"id"})
	assert.T(t).This(r.frac).Is(float32(0))
	assert.T(t).This(r.nseeks).Is(int32(100))
}

func TestSeekCount(t *testing.T) {
	assert.T(t).This(Require{nseeks: 7}.SeekCount(100)).Is(int32(7))
	assert.T(t).This(Require{nseeks: 7}.SeekCount(0)).Is(int32(7))

	assert.T(t).This(Require{frac: 0.5}.SeekCount(200)).Is(int32(100))
	assert.T(t).This(Require{frac: 0.25}.SeekCount(10)).Is(int32(2))
	assert.T(t).This(Require{frac: 0.01}.SeekCount(1)).Is(int32(0))
	assert.T(t).This(Require{frac: 0.5}.SeekCount(0)).Is(int32(0))
	assert.T(t).This(Require{frac: 1}.SeekCount(50)).Is(int32(50))
}

func TestSelectFrac(t *testing.T) {
	assert.T(t).This(Require{frac: 0.5}.SelectFrac(100)).Is(float32(0.5))
	assert.T(t).This(Require{frac: 1}.SelectFrac(100)).Is(float32(1))

	assert.T(t).This(Require{nseeks: 50}.SelectFrac(100)).Is(float32(0.5))
	assert.T(t).This(Require{nseeks: 200}.SelectFrac(100)).Is(float32(1))
	assert.T(t).This(Require{nseeks: 10}.SelectFrac(0)).Is(float32(1)) // nseeks/max(1,0) = 10/1 = 10, min(1,10) = 1

	assert.T(t).This(Require{}.SelectFrac(100)).Is(float32(1)) // frac=0, nseeks=0 -> 1
}

func TestRequireString(t *testing.T) {
	assert.T(t).This(NoneReq(0).String()).Is("ReqNone()")

	r := OrderReq([]string{"a", "b"}, 0.3)
	assert.T(t).This(r.String()).Is("ReqOrder(a,b) f0.3")

	r2 := GroupReq([]string{"x"}, 0.5, 10)
	assert.T(t).This(r2.String()).Is("ReqGroup(x) f0.5 s10")

	r3 := UniqueReq([]string{"id"}, 1)
	assert.T(t).This(r3.String()).Is("ReqUnique(id) s1")

	r4 := OrderReq([]string{"a"}, 0)
	assert.T(t).This(r4.String()).Is("ReqOrder(a)")

	r5 := OrderReq([]string{}, 0)
	assert.T(t).This(r5.String()).Is("ReqNone()")
}
