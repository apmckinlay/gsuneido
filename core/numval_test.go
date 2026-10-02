// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

func TestNumValHash(t *testing.T) {
	assert := assert.T(t)
	vals := []int{0, 1, -1, 100, -100, 12345, -12345,
		math.MinInt16, math.MaxInt16}
	for _, n := range vals {
		si := SuInt(n)
		si64 := SuInt64{n: n}
		dn := SuDnum{Dnum: dnum.FromInt(n)}
		assert.True(si.Equal(si64))
		assert.True(si.Equal(dn))
		assert.True(si64.Equal(dn))
		hash := si.Hash()
		assert.This(si64.Hash()).Is(hash)
		assert.This(dn.Hash()).Is(hash)
	}
	// values outside SuInt range but within Dnum's 16-digit precision
	vals = []int{1 << 20, -(1 << 20), math.MaxInt32, math.MinInt32,
		1e12, -1e12, 9999999999999999, -9999999999999999}
	for _, n := range vals {
		si64 := SuInt64{n: n}
		dn := SuDnum{Dnum: dnum.FromInt(n)}
		assert.True(si64.Equal(dn))
		hash := si64.Hash()
		assert.This(dn.Hash()).Is(hash)
	}
}
