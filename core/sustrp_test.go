// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/strp"
)

var sink any // package-level, forces escape

var lengths = []int{1, 2, 4, 8, 15, 23, 31, 64, 255}

func TestSuStrP(t *testing.T) {
	assert := assert.T(t).This

	sp := NewSuStrP("")
	assert(sp.Len()).Is(0)
	assert(sp.toStr()).Is("")

	sp = NewSuStrP("hello")
	assert(sp.Len()).Is(5)
	assert(sp.toStr()).Is("hello")

	// maximum length
	s := strings.Repeat("x", strp.MaxLen)
	sp = NewSuStrP(s)
	assert(sp.Len()).Is(strp.MaxLen)
	assert(sp.toStr()).Is(s)

	// single characters, all possible bytes
	for i := range 256 {
		c := byte(i)
		sp = SuStrP1(c)
		assert(sp.Len()).Is(1)
		assert(sp.toStr()).Is(string([]byte{c}))
	}

	// concatenation
	sp = CatSuStrP("foo", "bar")
	assert(sp.Len()).Is(6)
	assert(sp.toStr()).Is("foobar")
	sp = CatSuStrP("", "")
	assert(sp.Len()).Is(0)
	assert(sp.toStr()).Is("")

	// too long panics
	long := strings.Repeat("x", strp.MaxLen+1)
	assert(func() { NewSuStrP(long) }).Panics("string too large for SuStrP")
	assert(func() { CatSuStrP(long, "") }).Panics("string too large for SuStrP")
	assert(func() {
		CatSuStrP(strings.Repeat("x", 200), strings.Repeat("y", 100))
	}).Panics("string too large for SuStrP")
}

func TestSuStrPValue(t *testing.T) {
	a := assert.T(t)
	assert := a.This

	sp := NewSuStrP("abc")
	assert(sp.AsStr()).Is("abc", true)
	assert(sp.ToStr()).Is("abc", true)
	assert(sp.ToInt()).Is(0, false)
	assert(sp.IfInt()).Is(0, false)
	assert(sp.ToDnum()).Is(dnum.Zero, false)
	assert(sp.ToContainer()).Is(nil, false)
	assert(sp.Type()).Is(types.String)

	empty := NewSuStrP("")
	assert(empty.ToInt()).Is(0, true)
	assert(empty.ToDnum()).Is(dnum.Zero, true)

	// quoting/escaping
	assert(NewSuStrP("").String()).Is(`""`)
	assert(NewSuStrP("foo").String()).Is(`"foo"`)
	assert(NewSuStrP("x").String()).Is(`'x'`)
	assert(NewSuStrP("foo").Display(nil)).Is(`"foo"`)

	// hash is consistent with an equal SuStr
	a.That(sp.Hash() == SuStr("abc").Hash())
	a.That(sp.Hash2() == sp.Hash())
}

func TestSuStrPGet(t *testing.T) {
	assert := assert.T(t).This
	sp := NewSuStrP("abc")
	assert(ToStr(sp.Get(nil, SuInt16(0)))).Is("a")
	assert(ToStr(sp.Get(nil, SuInt16(2)))).Is("c")
	assert(ToStr(sp.Get(nil, SuInt16(-1)))).Is("c")
	assert(ToStr(sp.Get(nil, SuInt16(-3)))).Is("a")
	assert(ToStr(sp.Get(nil, SuInt16(3)))).Is("")
	assert(ToStr(sp.Get(nil, SuInt16(100)))).Is("")
	assert(ToStr(sp.Get(nil, SuInt16(-4)))).Is("")
	assert(sp.Get(nil, StrVal("x"))).Is(nil)
}

func TestSuStrPRange(t *testing.T) {
	assert := assert.T(t).This
	sp := NewSuStrP("abcdef")
	assert(ToStr(sp.RangeTo(1, 3))).Is("bc")
	assert(ToStr(sp.RangeTo(-2, -1))).Is("e")
	assert(ToStr(sp.RangeTo(2, 10))).Is("cdef")
	assert(ToStr(sp.RangeTo(10, 20))).Is("")
	assert(ToStr(sp.RangeLen(1, 2))).Is("bc")
	assert(ToStr(sp.RangeLen(4, 10))).Is("ef")
	assert(ToStr(sp.RangeLen(2, -1))).Is("")
	assert(ToStr(sp.RangeLen(1, 1))).Is("b")
}

func TestSuStrPPack(t *testing.T) {
	a := assert.T(t)
	assert := a.This

	sp := NewSuStrP("hello")
	assert(sp.PackSize(nil)).Is(6)
	assert(sp.PackSize2(nil, nil)).Is(6)
	packed := Pack(sp)
	assert(len(packed)).Is(6)
	assert(Unpack(packed)).Is(SuStr("hello"))

	empty := NewSuStrP("")
	assert(empty.PackSize(nil)).Is(0)
	assert(empty.PackSize2(nil, nil)).Is(0)
	assert(Pack(empty)).Is("")
	a.This(Unpack(Pack(empty))).Is(EmptyStr)
}

func TestSuStrPIter(t *testing.T) {
	assert := assert.T(t).This
	sp := NewSuStrP("abc")
	it := sp.Iter()
	assert(it.Infinite()).Is(false)
	assert(ToStr(it.Next())).Is("a")
	assert(ToStr(it.Next())).Is("b")
	assert(ToStr(it.Next())).Is("c")
	assert(it.Next()).Is(nil)

	assert(ToStr(it.Dup().Next())).Is("a")

	ob := sp.Iter().Instantiate()
	assert(ToStr(ob.ListGet(0))).Is("a")
	assert(ToStr(ob.ListGet(1))).Is("b")
	assert(ToStr(ob.ListGet(2))).Is("c")

	assert(NewSuStrP("").Iter().Next()).Is(nil)
}

func TestSuStrPNoPut(t *testing.T) {
	a := assert.T(t)
	sp := NewSuStrP("abc")
	a.This(func() { sp.Put(nil, Zero, EmptyStr) }).
		Panics("string does not support put")
	a.This(func() { sp.GetPut(nil, Zero, EmptyStr, nil, false) }).
		Panics("string does not support update")
}

func TestSuStrPSetConcurrent(t *testing.T) {
	assert := assert.T(t).This
	sp := NewSuStrP("abc")
	sp.SetConcurrent()
	assert(sp.toStr()).Is("abc")
	assert(sp.Len()).Is(3)
}

func BenchmarkConcat(b *testing.B) {
	for _, n := range lengths {
		x := strings.Repeat("a", n/2)
		y := strings.Repeat("b", n-n/2)
		b.Run(fmt.Sprintf("old/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sink = SuStr(x + y) // data alloc + header alloc
			}
		})
		b.Run(fmt.Sprintf("new/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sink = (*SuStrP)(strp.Concat(x, y))
			}
		})
	}
}

var source = "hello"

func BenchmarkStrValue(b *testing.B) {
	b.Run("SuStr", func(b *testing.B) {
		b.ReportAllocs()
		var v Value
		for i := 0; i < b.N; i++ {
			v = SuStr(source)
		}
		sink = v
	})
	b.Run("SuStrP", func(b *testing.B) {
		b.ReportAllocs()
		var v Value
		for i := 0; i < b.N; i++ {
			v = NewSuStrP(source)
		}
		sink = v
	})
}

func BenchmarkSubstr(b *testing.B) {
	rec := strings.Repeat("0123456789", 100) // stands in for a record
	for _, n := range lengths {
		b.Run(fmt.Sprintf("old/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sink = SuStr(rec[10 : 10+n]) // header alloc only
			}
		})
		b.Run(fmt.Sprintf("new/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sink = (*SuStrP)(strp.From(rec[10 : 10+n]))
			}
		})
	}
}
