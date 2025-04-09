// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import (
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPack(t *testing.T) {
	assert := assert.T(t).This
	e := NewEncoder(32)
	e.Put4(12, 34, 56, 78)
	d := MakeDecoder(e.String())
	assert(d.Get1()).Is(12)
	assert(d.Get1()).Is(34)
	assert(d.Get1()).Is(56)
	assert(d.Get1()).Is(78)

	e = NewEncoder(32)
	e.PutStr("helloworld")
	d = MakeDecoder(e.String())
	assert(d.Get(5)).Is("hello")
	assert(d.Get(5)).Is("world")

	for _, n := range []uint16{0, 1, 1234, math.MaxUint16} {
		e = NewEncoder(32)
		e.Uint16(n)
		d = MakeDecoder(e.String())
		assert(d.Uint16()).Is(n)
	}
	for _, n := range []uint32{0, 1, 12345678, math.MaxUint32} {
		e = NewEncoder(32)
		e.Uint32(n)
		d = MakeDecoder(e.String())
		assert(d.Uint32()).Is(n)
	}
	for _, n := range []uint64{0, 1, 222, 22222, 12345678, math.MaxInt32} {
		e = NewEncoder(32)
		e.VarUint(n)
		d = MakeDecoder(e.String())
		assert(d.VarUint()).Is(n)
	}
}
