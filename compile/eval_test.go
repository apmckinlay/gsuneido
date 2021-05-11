// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestIsGlobal(t *testing.T) {
	assert := assert.T(t)
	assert.True(isGlobal("F"))
	assert.True(isGlobal("Foo"))
	assert.True(isGlobal("Foo_123_Bar"))
	assert.True(isGlobal("Foo!"))
	assert.True(isGlobal("Foo?"))

	assert.False(isGlobal(""))
	assert.False(isGlobal("f"))
	assert.False(isGlobal("foo"))
	assert.False(isGlobal("_foo"))
	assert.False(isGlobal("Foo!bar"))
	assert.False(isGlobal("Foo?bar"))
	assert.False(isGlobal("Foo.bar"))
}
