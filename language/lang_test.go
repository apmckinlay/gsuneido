// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package language

import (
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSuFuncCall(t *testing.T) {
	fn := compile.Constant("function (a, b) { a - b }").(*SuFunc)
	th := NewThread()
	th.Push(SuInt(100))
	th.Push(SuInt(1))
	result := fn.Call(th, nil, &ArgSpec2)
	assert.T(t).This(result).Is(SuInt(99))
	Global.Add("F", fn)

	fn = compile.Constant("function () { F(100, 1) }").(*SuFunc)
	result = fn.Call(th, nil, &ArgSpec0)
	assert.T(t).This(result).Is(SuInt(99))

	fn = compile.Constant("function () { F(b: 1, a: 100) }").(*SuFunc)
	result = fn.Call(th, nil, &ArgSpec0)
	assert.T(t).This(result).Is(SuInt(99))
}

func BenchmarkInt(b *testing.B) {
	fn := compile.Constant("function (n) { for (i = 0; i < n; ++i){} }").(*SuFunc)
	th := NewThread()
	m := 1
	n := b.N
	for n > math.MaxInt16 {
		n /= 2
		m *= 2
	}
	for i := 0; i < m; i++ {
		th.Push(SuInt(n))
		th.Start(fn, nil)
	}
}
