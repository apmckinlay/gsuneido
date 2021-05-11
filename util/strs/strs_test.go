// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package strs

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestEqual(t *testing.T) {
	test := func (x, y []string) {
		assert.T(t).That(Equal(x, y))
		assert.T(t).That(Equal(y, x))
	}
	xtest := func (x, y []string) {
		assert.T(t).That(!Equal(x, y))
		assert.T(t).That(!Equal(y, x))
	}
	empty := []string{}
	test(nil, nil)
	test(nil, empty)
	test(empty, nil)
	test(empty, empty)
	one := []string{"one"}
	test(one, one)
	xtest(one, nil)
	xtest(one, empty)
	four := []string{"one", "two", "three", "four"}
	test(four, four)
	xtest(four, nil)
	xtest(four, one)
}

func TestIndex(t *testing.T) {
	assert := assert.T(t).This
	assert(Index([]string{}, "five")).Is(-1)
	list := []string{"one", "two", "three", "two", "four"}
	assert(Index(list, "five")).Is(-1)
	assert(Index(list, "one")).Is(0)
	assert(Index(list, "two")).Is(1)
	assert(Index(list, "four")).Is(4)
}
