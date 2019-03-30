package pack

import (
	"math"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestPack(t *testing.T) {
	e := NewEncoder(32)
	e.Put4(12, 34, 56, 78)
	d := NewDecoder(e.String())
	Assert(t).That(d.Get1(), Equals(12))
	Assert(t).That(d.Get1(), Equals(34))
	Assert(t).That(d.Get1(), Equals(56))
	Assert(t).That(d.Get1(), Equals(78))

	e = NewEncoder(32)
	e.PutStr("helloworld")
	d = NewDecoder(e.String())
	Assert(t).That(d.Get(5), Equals("hello"))
	Assert(t).That(d.Get(5), Equals("world"))

	for _, n := range []uint16{0, 1, 1234, math.MaxUint16} {
		e = NewEncoder(32)
		e.Uint16(n)
		d = NewDecoder(e.String())
		Assert(t).That(d.Uint16(), Equals(n))
	}
	for _, n := range []uint32{0, 1, 12345678, math.MaxUint32} {
		e = NewEncoder(32)
		e.Uint32(n)
		d = NewDecoder(e.String())
		Assert(t).That(d.Uint32(), Equals(n))
	}
	for _, n := range []int{0, 1, -1, 12345678, -12345678, math.MaxInt32, math.MinInt32} {
		e = NewEncoder(32)
		e.Int32(n)
		d = NewDecoder(e.String())
		Assert(t).That(d.Int32(), Equals(n))
	}
}
