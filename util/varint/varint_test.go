package varint

import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func TestVarint(t *testing.T) {
	Assert(t).That(EncodeUint32(0, []byte{}), Equals([]byte{0}))
	Assert(t).That(EncodeInt32(0, []byte{}), Equals([]byte{0}))
	Assert(t).That(EncodeUint32(45, []byte{}), Equals([]byte{45}))
	Assert(t).That(EncodeInt32(45, []byte{}), Equals([]byte{90}))
	Assert(t).That(EncodeInt32(-1, []byte{}), Equals([]byte{1}))
	Assert(t).That(EncodeUint32(256, []byte{}), Equals([]byte{128, 2}))

	for _, n := range []int32{0, 1, -1, 123, -123, 999999, -999999} {
		n2, _ := DecodeInt32(EncodeInt32(n, []byte{}), 0)
		Assert(t).That(n2, Equals(n).Comment("signed"))
		if n > 0 {
			n2, _ := DecodeUint32(EncodeUint32(uint32(n), []byte{}), 0)
			Assert(t).That(n2, Equals(uint32(n)).Comment("unsigned"))
		}
	}
}
