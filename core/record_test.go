// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBuilder(t *testing.T) {
	assert := assert.T(t).This
	var b RecordBuilder
	rec := b.Build()
	assert([]byte(rec)).Is([]byte{0})
	b.AddRaw("one")
	rec = b.Build()
	assert([]byte(rec)).Is([]byte{type8 << 6, 1, 7, 4, 'o', 'n', 'e'})
	assert(rec.GetRaw(0)).Is("one")

	b = RecordBuilder{}
	b.Add(SuInt(123))
	b.Add(SuStr("foobar"))

	rec = b.Build()
	assert(rec.mode()).Is(type8)
	assert(rec.Count()).Is(2)
	assert(rec.GetVal(0)).Is(SuInt(123))
	assert(rec.GetVal(1)).Is(SuStr("foobar"))

	s := strings.Repeat("helloworld", 30)
	b.AddRaw(s)
	rec = b.Build()
	assert(rec.mode()).Is(type16)
	assert(rec.GetRaw(2)).Is(s)

}

func TestLength(t *testing.T) {
	assert := assert.T(t).This
	assert(tblength(0, 0)).Is(1)
	assert(tblength(1, 1)).Is(5)
	assert(tblength(1, 200)).Is(204)
	assert(tblength(1, 248)).Is(252)

	assert(tblength(1, 252)).Is(258)
	assert(tblength(1, 300)).Is(306)

	assert(tblength(1, 0x10000)).Is(0x1000a)
}

func BenchmarkRecordBuilder(b *testing.B) {
	for range b.N {
		b := RecordBuilder{}
		for range 10 {
			b.Add(SuStr("hello"))
			b.AddRaw("world")
			b.Build()
		}
	}
}

func TestTruncate(t *testing.T) {
    b := RecordBuilder{}
	r := b.Build()
	assert.That(r == "\x00")
	r2 := r.Truncate(99)
	assert.That(r2 == r)
	r2 = r.Truncate(0)
	assert.That(r2 == r)
	
    b.Add(SuInt(123))
    b.Add(SuStr("hello world"))
	b.Add(SuBool(true))
	r = b.Build()
	r2 = r.Truncate(99)
	assert.That(r2 == r)
	r2 = r.Truncate(2)
	assert.This(r2.String()).Is(`{123, "hello world"}`)
	
	b.Add(SuStr(""))
	b.Add(SuStr(""))
	r = b.Build()
	r2 = r.Truncate(4)
	assert.This(r2.String()).Is(`{123, "hello world", true}`)
}