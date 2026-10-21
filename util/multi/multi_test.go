// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package multi

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestNew(t *testing.T) {
	assert := assert.T(t).This
	a := 1
	ps := New(&a)
	assert(ps.Len()).Is(1)
	assert(ps.Get(0)).Is(&a)
	assert(func() { ps.Get(1) }).Panics("Ptrs.Get: index out of range")
	assert(func() { ps.Get(-1) }).Panics("Ptrs.Get: index out of range")
}

func TestWith(t *testing.T) {
	all := func(ps Ptrs[int]) []*int {
		var res []*int
		for i := range ps.Len() {
			res = append(res, ps.Get(i))
		}
		return res
	}
	assert := assert.T(t).This
	var empty Ptrs[int]
	a, b, c := 1, 2, 3
	one := With(empty, &a)
	assert(one.Len()).Is(1)
	assert(one.Get(0)).Is(&a)
	two := With(one, &b)
	assert(two.Len()).Is(2)
	assert(all(two)).Is([]*int{&a, &b})
	three := With(two, &c)
	assert(three.Len()).Is(3)
	assert(all(three)).Is([]*int{&a, &b, &c})
	assert(empty.Len()).Is(0)
	assert(all(one)).Is([]*int{&a})
	assert(all(two)).Is([]*int{&a, &b})
}

func TestJoin(t *testing.T) {
	assert := assert.T(t).This
	a, b, c := 1, 2, 3
	ab := Join(New(&a), New(&b))
	assert(ab.Len()).Is(2)
	assert(ab.Get(0)).Is(&a)
	assert(ab.Get(1)).Is(&b)
	abc := Join(ab, New(&c))
	assert(abc.Len()).Is(3)
	assert(abc.Get(0)).Is(&a)
	assert(abc.Get(1)).Is(&b)
	assert(abc.Get(2)).Is(&c)
	assert(func() { abc.Get(3) }).Panics("Ptrs.Get: index out of range")
}

func TestJoinEmpty(t *testing.T) {
	assert := assert.T(t).This
	var empty Ptrs[int]
	a := 1
	one := New(&a)
	assert(Join(empty, one).Get(0)).Is(&a)
	assert(Join(one, empty).Get(0)).Is(&a)
	assert(func() { Join(empty, empty).Get(0) }).Panics("Ptrs.Get: index out of range")
}
