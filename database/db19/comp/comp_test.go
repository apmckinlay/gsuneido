// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package comp

import (
	"bytes"
	"math/rand"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestEncode(t *testing.T) {
	Assert(t).That(Encode([][]byte{}), Equals(""))
	Assert(t).That(Encode([][]byte{{'a', 'b'}}), Equals("ab"))
	Assert(t).That(Encode([][]byte{{'a'}, {'b'}}), Equals("a\x00\x00b"))
	Assert(t).That(Encode([][]byte{{'a', 0, 'b'}}), Equals("a\x00\x01b"))
	Assert(t).That(Encode([][]byte{{0, 'a', 'b'}}), Equals("\x00\x01ab"))
	Assert(t).That(Encode([][]byte{{'a', 0, 0, 'b'}}), Equals("a\x00\x01\x00\x01b"))
	Assert(t).That(Encode([][]byte{{'a', 0, 1, 'b'}}), Equals("a\x00\x01\x01b"))
	Assert(t).That(Encode([][]byte{{'a', 'b', 0}}), Equals("ab\x00\x01"))
	Assert(t).That(Encode([][]byte{{'a', 'b', 0, 0}}), Equals("ab\x00\x01\x00\x01"))
}

const m = 3

func TestEncodeRandom(t *testing.T) {
	var n = 1000000
	if testing.Short() {
		n = 10000
	}
	for i := 0; i < n; i++ {
		x := gen()
		y := gen()
		xenc := Encode(x)
		yenc := Encode(y)
		Assert(t).That(xenc < yenc, Equals(lt(x, y)))
	}
}

func gen() [][]byte {
	x := make([][]byte, m)
	for i := 0; i < m; i++ {
		n := rand.Intn(6) + 1
		x[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			x[i][j] = byte(rand.Intn(4)) // 25% zeros
		}
	}
	return x
}

func lt(x [][]byte, y [][]byte) bool {
	for i := 0; i < len(x) && i < len(y); i++ {
		if cmp := bytes.Compare(x[i], y[i]); cmp != 0 {
			return cmp < 0
		}
	}
	return len(x) < len(y)
}
