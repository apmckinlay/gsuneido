// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !portable && windows

package builtin

import (
	"testing"

	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
)

var result Value

func BenchmarkGetStr(b *testing.B) {
	const n = 100
	buf := heap.Alloc(n)
	var s SuStr
	for range b.N {
		s = SuStr(heap.GetStrN(buf, n))
	}
	result = s
}

func BenchmarkBufToStr(b *testing.B) {
	const n = 100
	buf := heap.Alloc(n)
	var s Value
	for range b.N {
		s = bufStrN(buf, n)
	}
	result = s
}
