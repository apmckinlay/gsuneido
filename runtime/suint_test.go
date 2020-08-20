// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"math"
	"reflect"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

func TestSuInt(t *testing.T) {
	assert := assert.T(t)
	test := func(n int) {
		assert.This(SuInt(n).toInt()).Is(int(n))
	}
	test(0)
	test(12345)
	test(-12345)
	test(math.MinInt16)
	test(math.MaxInt16)

	xtest := func(n int) {
		assert.This(func() { SuInt(n) }).Panics("index out of range")
	}
	xtest(123456)
	xtest(-123456)

	assert.False(reflect.DeepEqual(SuInt(2), SuInt(3)))
	assert.False(reflect.DeepEqual(SuInt(-2), SuInt(-3)))

	s10 := SuInt(10)
	d10 := SuDnum{Dnum: dnum.FromInt(10)}
	assert.True(s10.Equal(d10))
	assert.This(s10.Hash()).Is(d10.Hash())
}
