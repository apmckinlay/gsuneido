// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ss

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNew(t *testing.T) {
	ss := New[string](3)
	assert.T(t).This(ss.Capacity()).Is(3)
	assert.T(t).This(ss.Count()).Is(0)
	assert.T(t).This(ss.Len()).Is(0)

	assert.T(t).This(func() { New[string](0) }).Panics("capacity")
}

func TestAddNoEviction(t *testing.T) {
	ss := New[string](4)
	ss.Add("a")
	ss.Add("b")
	ss.Add("a")

	assert.T(t).This(ss.Count()).Is(3)
	assert.T(t).This(ss.Len()).Is(2)

	c, e, ok := ss.Estimate("a")
	assert.T(t).That(ok)
	assert.T(t).This(c).Is(2)
	assert.T(t).This(e).Is(0)

	c, e, ok = ss.Estimate("b")
	assert.T(t).That(ok)
	assert.T(t).This(c).Is(1)
	assert.T(t).This(e).Is(0)

	_, _, ok = ss.Estimate("x")
	assert.T(t).That(!ok)
}

func TestEvictionBounds(t *testing.T) {
	ss := New[string](2)
	ss.Add("a")
	ss.Add("b")
	ss.Add("c")

	assert.T(t).This(ss.Len()).Is(2)
	assert.T(t).This(ss.Count()).Is(3)

	_, _, ok := ss.Estimate("a")
	assert.T(t).That(!ok)

	c, e, ok := ss.Estimate("c")
	assert.T(t).That(ok)
	assert.T(t).This(c).Is(2)
	assert.T(t).This(e).Is(1)
	assert.T(t).That(c-e <= 1)
}

func TestHeavyHitters(t *testing.T) {
	ss := New[string](8)
	for range 1000 {
		ss.Add("hot")
	}
	for i := range 200 {
		ss.Add(fmt.Sprintf("mid%03d", i%4))
	}
	for i := range 200 {
		ss.Add(fmt.Sprintf("cold%03d", i))
	}

	count, _, ok := ss.Estimate("hot")
	assert.T(t).That(ok)
	assert.T(t).That(count >= 1000)

	top := ss.Top()
	assert.T(t).That(len(top) > 0)
	assert.T(t).This(top[0].Value).Is("hot")

	for i := 1; i < len(top); i++ {
		assert.T(t).That(top[i-1].Count >= top[i].Count)
	}

	total := 1000 + 200 + 200
	assert.T(t).This(ss.Count()).Is(total)
}

