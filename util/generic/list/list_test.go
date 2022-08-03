// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package list

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"testing"
)

type Num int

func (n Num) Equal(m any) bool {
	return n == m
}

func TestList(t *testing.T) {
	x := List[Num]{}
	x.Push(123)
	x.Push(456)
	assert.T(t).That(x.Has(123))
	assert.T(t).That(x.Has(456))
	x.Pop()
	assert.T(t).That(x.Has(123))
	assert.T(t).That(!x.Has(456))
	x.Remove(123)
	assert.T(t).That(!x.Has(123))
	assert.T(t).That(!x.Has(456))
}
