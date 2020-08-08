// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package bits

import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func TestClp2(t *testing.T) {
	Assert(t).That(Clp2(0), Is(uint64(0)))
	Assert(t).That(Clp2(1), Is(uint64(1)))
	Assert(t).That(Clp2(2), Is(uint64(2)))
	Assert(t).That(Clp2(3), Is(uint64(4)))
	Assert(t).That(Clp2(123), Is(uint64(128)))
	Assert(t).That(Clp2(65536), Is(uint64(65536)))
}
