// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package opt

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// Int is an optional int stored in a single byte.
// The zero value is valid but not set.
// The actual value 0 is stored as math.MinInt
type Int struct {
	val int
}

func (i *Int) Set(v int) {
	assert.That(v != math.MinInt)
	if v == 0 {
		v = math.MinInt
	}
	i.val = v
}

func (i Int) IsSet() bool {
	return i.val != 0
}

func (i Int) NotSet() bool {
	return i.val == 0
}

func (i Int) Get() int {
	switch i.val {
	case 0:
		panic("opt.Int Get when not set")
	case math.MinInt:
		return 0
	}
	return i.val
}

func (i Int) GetOr(v int) int {
	switch i.val {
	case 0:
		return v
	case math.MinInt:
		return 0
	}
	return i.val
}
