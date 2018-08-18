package dnum

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestDiv128(t *testing.T) {
	Assert(t).That(div128(1, 4), Equals(uint64(2500000000000000)))
	Assert(t).That(div128(1, 3), Equals(uint64(3333333333333333)))
	Assert(t).That(div128(2, 3), Equals(uint64(6666666666666666)))
	Assert(t).That(div128(1, 11), Equals(uint64(909090909090909)))
	Assert(t).That(div128(11, 13), Equals(uint64(8461538461538461)))
	Assert(t).That(div128(1234567890123456, 9876543210987654),
		Equals(uint64(1249999988609374)))
}
