// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package cache

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCache(t *testing.T) {
	n := 0
	c := New(func(s string) string { n++; return s + s })
	assert.T(t).This(c.Get("foo")).Is("foofoo")
	assert.T(t).This(n).Is(1)
	assert.T(t).This(c.Get("bar")).Is("barbar")
	assert.T(t).This(n).Is(2)
	assert.T(t).This(c.Get("foo")).Is("foofoo")
	assert.T(t).This(c.Get("bar")).Is("barbar")
	assert.T(t).This(n).Is(2)
	for i := range 100 {
		s := strconv.Itoa(i)
		assert.T(t).This(c.Get(s)).Is(s + s)
	}
	assert.T(t).This(n).Is(102)
	assert.T(t).This(c.Get("95")).Is("9595")
	assert.T(t).This(c.Get("99")).Is("9999")
	assert.T(t).This(n).Is(102)
}
