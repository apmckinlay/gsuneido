package pack

import (
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDecoder(t *testing.T) {
	assert := assert.T(t).This

	// Test NewDecoder and basic operations
	d := NewDecoder("hello world")
	assert(d.Peek()).Is(byte('h'))
	assert(d.Get1()).Is(byte('h'))
	assert(d.Get(4)).Is("ello")
	assert(d.Remaining()).Is(true)
	assert(d.Remainder()).Is(" world")

	// Test with binary data
	binaryData := "\x00\x01\x02\x03\x04"
	d = NewDecoder(binaryData)
	assert(d.Get1()).Is(byte(0))
	assert(d.Get(2)).Is("\x01\x02")
	assert(d.Remaining()).Is(true)
	assert(d.Remainder()).Is("\x03\x04")

	// Test edge cases
	d = NewDecoder("")
	assert(d.Remaining()).Is(false)
	assert(d.Remainder()).Is("")

	// Test multiple operations
	d = NewDecoder("abcdefghij")
	assert(d.Get1()).Is(byte('a'))
	assert(d.Peek()).Is(byte('b'))
	assert(d.Get(3)).Is("bcd")
	assert(d.Remaining()).Is(true)
	assert(d.Get(2)).Is("ef")
	assert(d.Remainder()).Is("ghij")

	// Test integer decoding
	testVarUInt := func(n uint64) {
		t.Helper()
		e := NewEncoder(16)
		e.VarUint(n)
		d := NewDecoder(e.String())
		assert(d.VarUint()).Is(n)
	}
	testVarUInt(0)
	testVarUInt(1)
	testVarUInt(127)
	testVarUInt(128)
	testVarUInt(255)
	testVarUInt(256)
	testVarUInt(65535)
	testVarUInt(65536)
	testVarUInt(math.MaxUint32)
	testVarUInt(math.MaxUint64)

	testUint16 := func(n uint16) {
		t.Helper()
		e := NewEncoder(4)
		e.Uint16(n)
		d := NewDecoder(e.String())
		assert(d.Uint16()).Is(n)
	}
	testUint16(0)
	testUint16(1)
	testUint16(255)
	testUint16(256)
	testUint16(math.MaxUint16)

	testUint32 := func(n uint32) {
		t.Helper()
		e := NewEncoder(8)
		e.Uint32(n)
		d := NewDecoder(e.String())
		assert(d.Uint32()).Is(n)
	}
	testUint32(0)
	testUint32(1)
	testUint32(65535)
	testUint32(65536)
	testUint32(math.MaxUint32)
}
