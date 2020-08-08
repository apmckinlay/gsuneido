// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ints

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestFill(t *testing.T) {
	Fill([]int{}, 0)
	data := make([]int, 4)
	Fill(data, 123)
	for _, x := range data {
		Assert(t).That(x, Is(123))
	}
}

func TestIndex(t *testing.T) {
	data := []int{}
	Assert(t).That(Index(data, 0), Is(-1))
	data = []int{123}
	Assert(t).That(Index(data, 0), Is(-1))
	Assert(t).That(Index(data, 123), Is(0))
	data = []int{3, 7, 4, 9, 2, 4, 1}
	Assert(t).That(Index(data, 0), Is(-1))
	Assert(t).That(Index(data, 3), Is(0))
	Assert(t).That(Index(data, 4), Is(2))
	Assert(t).That(Index(data, 1), Is(6))
}

func TestCompare(t *testing.T) {
	Assert(t).That(Compare(0, 0), Is(0))
	Assert(t).That(Compare(123, 0), Is(+1))
	Assert(t).That(Compare(123, 456), Is(-1))
}

func TestMin(t *testing.T) {
	Assert(t).That(Min(0, 0), Is(0))
	Assert(t).That(Min(1, 2), Is(1))
	Assert(t).That(Min(2, 1), Is(1))
}
