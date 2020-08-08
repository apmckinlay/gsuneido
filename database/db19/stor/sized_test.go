// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSized(t *testing.T) {
	st := HeapStor(1024)

	off, buf := st.AllocSized(10)
	copy(buf, "helloworld")

	buf = st.DataSized(off)
	Assert(t).That(string(buf), Is("helloworld"))
}
