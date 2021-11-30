// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strconv"
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

func TestCacheSet(*testing.T) {
	var cs cacheSet
	itoa := strconv.Itoa
	add := func(i int) { cs.Add(itoa(i)) }
	for i := 0; i < 2*cacheSetSize; i++ {
		add(i)
	}
	has := func(i int) bool { return cs.Has(itoa(i)) }
	for i := 0; i < cacheSetSize; i++ {
		assert.This(has(i)).Is(i >= cacheSetSize)
	}
}
