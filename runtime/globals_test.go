// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestGlobals(t *testing.T) {
	assert := assert.T(t).This
	foo := Global.Num("foo")
	assert(Global.Num("foo")).Is(foo)
	assert(Global.Add("bar", nil)).Is(foo + 1)
	assert(Global.Num("bar")).Is(foo + 1)
	Global.Add("baz", True)
	assert(func() { Global.Add("baz", False) }).Panics("duplicate")
	assert(Global.Name(foo)).Is("foo")
	assert(Global.Name(foo + 1)).Is("bar")
}
