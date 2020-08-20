// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ints

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFill(t *testing.T) {
	Fill([]int{}, 0)
	data := make([]int, 4)
	Fill(data, 123)
	for _, x := range data {
		assert.T(t).This(x).Is(123)
	}
}

func TestIndex(t *testing.T) {
	assert := assert.T(t).This
	data := []int{}
	assert(Index(data, 0)).Is(-1)
	data = []int{123}
	assert(Index(data, 0)).Is(-1)
	assert(Index(data, 123)).Is(0)
	data = []int{3, 7, 4, 9, 2, 4, 1}
	assert(Index(data, 0)).Is(-1)
	assert(Index(data, 3)).Is(0)
	assert(Index(data, 4)).Is(2)
	assert(Index(data, 1)).Is(6)
}

func TestCompare(t *testing.T) {
	assert.T(t).This(Compare(0, 0)).Is(0)
	assert.T(t).This(Compare(123, 0)).Is(+1)
	assert.T(t).This(Compare(123, 456)).Is(-1)
}

func TestMin(t *testing.T) {
	assert.T(t).This(Min(0, 0)).Is(0)
	assert.T(t).This(Min(1, 2)).Is(1)
	assert.T(t).This(Min(2, 1)).Is(1)
}
