// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package verify

import (
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"testing"
)

func TestVerify(t *testing.T) {
	That(true) // does nothing
	Assert(t).That(func() { That(false) }, Panics("verify failed"))
}
