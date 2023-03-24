// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestMatcher(t *testing.T) {
	m := newMatcher("foo.bar")
	assert.True(m("foo.bar"))
	assert.False(m(""))
	assert.False(m("bar.foo"))

	m = newMatcher("foo.*")
	assert.True(m("foo."))
	assert.True(m("foo.bar"))
	assert.True(m("foo.baz"))
	assert.True(m("foo.bar.baz"))
	assert.False(m("foo"))
	assert.False(m(""))

	m = newMatcher("*.foo")
	assert.True(m(".foo"))
	assert.True(m("bar.foo"))
	assert.False(m(""))
	assert.False(m("foo."))
	assert.False(m("foo.bar"))
	assert.False(m("bar"))

	m = newMatcher("foo*.bar")
	assert.True(m("foo.bar"))
	assert.True(m("foobaz.bar"))
	assert.False(m(""))
	assert.False(m("foo."))
	assert.False(m(".bar"))
}
