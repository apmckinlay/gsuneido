// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package pack

import (
    "bytes"
    "math"
    "testing"

    "github.com/apmckinlay/gsuneido/util/assert"
)

func TestEncoder2(t *testing.T) {
    assert := assert.T(t).This

    // Test basic Put methods
    buf := &bytes.Buffer{}
    e := NewEncoder2(buf)
    e.Put1(12)
    e.Put2(34, 56)
    e.Put4(78, 90, 123, 234)
    e.Flush()

    result := buf.Bytes()
    expected := []byte{12, 34, 56, 78, 90, 123, 234}
    assert(result).Is(expected)

    // Test PutStr
    buf.Reset()
    e = NewEncoder2(buf)
    e.PutStr("hello")
    e.PutStr("world")
    e.Flush()

    assert(buf.String()).Is("helloworld")

    // Test Put with byte slice
    buf.Reset()
    e = NewEncoder2(buf)
    e.Put([]byte{1, 2, 3, 4})
    e.Flush()

    result = buf.Bytes()
    expected = []byte{1, 2, 3, 4}
    assert(result).Is(expected)

    // Test Uint16
    for _, n := range []uint16{0, 1, 1234, math.MaxUint16} {
        buf.Reset()
        e = NewEncoder2(buf)
        e.Uint16(n)
        e.Flush()

        result = buf.Bytes()
        assert(len(result)).Is(2)
        assert(uint16(result[0])<<8 | uint16(result[1])).Is(n)
    }

    // Test Uint32
    for _, n := range []uint32{0, 1, 12345678, math.MaxUint32} {
        buf.Reset()
        e = NewEncoder2(buf)
        e.Uint32(n)
        e.Flush()

        result = buf.Bytes()
        assert(len(result)).Is(4)
        assert(uint32(result[0])<<24 | uint32(result[1])<<16 | uint32(result[2])<<8 | uint32(result[3])).Is(n)
    }

    // Test VarUint
    testCases := []uint64{0, 1, 222, 22222, 12345678, math.MaxInt32, math.MaxUint64}
    for _, n := range testCases {
        buf.Reset()
        e = NewEncoder2(buf)
        e.VarUint(n)
        e.Flush()

        // Verify by reading back with a decoder
        d := NewDecoder(buf.String())
        assert(d.VarUint()).Is(n)
    }

    // Test method chaining
    buf.Reset()
    e = NewEncoder2(buf)
    e.Put1(1).Put2(2, 3).Put4(4, 5, 6, 7).PutStr("test").Uint16(1234).Uint32(5678).VarUint(9012)
    e.Flush()

    // Verify the buffer has content (not checking exact content since we already tested individual methods)
    assert(buf.Len()).Is(19)
}
