// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	op "github.com/apmckinlay/gsuneido/core/opcodes"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestInterp(t *testing.T) {
	test := func(expected Value, code ...byte) {
		fn := &SuFunc{Code: string(code)}
		var th Thread
		result := th.Call(fn)
		assert.T(t).This(result).Is(expected)
	}
	test(SuInt(8), byte(op.Int), 0, 3, byte(op.Int), 0, 5, byte(op.Add), byte(op.Return))
}

func TestCatchMatch(t *testing.T) {
	match := func(e, pat string) {
		assert.T(t).True(catchMatch(e, pat))
	}
	match("", "")
	match("foo", "")
	match("foo", "|")
	match("foo", "*")
	match("foobar", "foo")
	match("foobar", "*bar")
	match("foobar", "*ba")
	match("foobar", "*foo")
	match("foobar", "*foobar")

	match("foobar", "foo|def")
	match("foobar", "abc|foo")
	match("foobar", "abc|foo|def")

	match("foobar", "abc|*bar")

	nomatch := func(e, pat string) {
		assert.T(t).False(catchMatch(e, pat))
	}
	nomatch("", "foo")
	nomatch("foo", "bar")
	nomatch("foo", "*bar")
	nomatch("foobar", "bar")
	nomatch("foobar", "far|boo|x")
}

// compare to BenchmarkInterp in execute_test.go

func BenchmarkJit(b *testing.B) {
	th := &Thread{}
	for b.Loop() {
		th.Reset()
		result := jitfn(th)
		if !result.Equal(SuInt(4950)) {
			panic("wrong result")
		}
	}
}

var hundred = SuInt(100)

func jitfn(th *Thread) Value {
	th.sp += 2
	th.stack[0] = Zero // sum
	th.stack[1] = Zero // i
	for {
		th.stack[0] = OpAdd(th.stack[0], th.stack[1]) // sum += i
		th.stack[1] = OpAdd(th.stack[1], One)         // ++i
		if OpLt(th.stack[1], hundred) != True {
			break
		}
	}
	return th.stack[0] // return sum
}

func BenchmarkTranspile(b *testing.B) {
	for b.Loop() {
		result := transpilefn()
		if !result.Equal(SuInt(4950)) {
			panic("wrong result")
		}
	}
}

func transpilefn() Value {
	sum := Zero
	i := Zero
	for {
		sum = OpAdd(sum, i) // sum += i
		i = OpAdd(i, One)   // ++i
		if OpLt(i, hundred) != True {
			break
		}
	}
	return sum
}

func BenchmarkSpecialize(b *testing.B) {
	for b.Loop() {
		result := specialized()
		if !result.Equal(SuInt(4950)) {
			panic("wrong result")
		}
	}
}

func specialized() Value {
	sum := 0
	i := 0
	for {
		sum += i
		i++
		if i >= 100 {
			break
		}
	}
	return SuInt(sum)
}

var r Value

func BenchmarkLoadStore(b *testing.B) {
	x := One
	y := One
	for n := range b.N {
		switch n % 11 {
		case 0:
			r = OpAdd(x, y)
		case 1:
			r = OpSub(x, y)
		case 2:
			r = OpCat(nil, x, y)
		case 3:
			r = OpMul(x, y)
		case 4:
			r = OpDiv(x, y)
		case 5:
			r = OpMod(x, y)
		case 6:
			r = OpLeftShift(x, y)
		case 7:
			r = OpRightShift(x, y)
		case 8:
			r = OpBitOr(x, y)
		case 9:
			r = OpBitAnd(x, y)
		case 10:
			r = OpBitXor(x, y)
		}
	}
}

func BenchmarkLoadStore2(b *testing.B) {
	x := One
	y := One
	var th *Thread
	for n := range b.N {
		op := []func(x, y Value) Value{
			OpAdd, OpSub, th.Cat, OpMul, OpDiv, OpMod,
			OpLeftShift, OpRightShift, OpBitOr, OpBitAnd, OpBitXor}[n%11]
		r = op(x, y)
	}
}

func BenchmarkLoadStore3(b *testing.B) {
	x := One
	y := MinusOne
	for b.Loop() {
		r = OpSub(x, y)
	}
}

func BenchmarkInterp(b *testing.B) {
	code := []byte{byte(op.Int), 0, 3, byte(op.Int), 0, 5, byte(op.Add),
		byte(op.Return)}
	fn := &SuFunc{Code: string(code)}
	var th Thread
	for b.Loop() {
		th.invoke(fn, nil)
		th.sp = 0
	}
}

func BenchmarkLoadValue(b *testing.B) {
	// LoadValue optimized path - loads local var then constant value
	code := []byte{
		byte(op.LoadValue), 0, 0, // LoadValue local_idx=0, value_idx=0
		byte(op.Return),
	}
	fn := &SuFunc{
		Code: string(code),
		ParamSpec: ParamSpec{
			Values:  []Value{SuStr("test")}, // Value at index 0
			Names:   []string{"x"},          // Local var name at index 0
			Nparams: 1,                      // Treat as parameter so it gets initialized
		},
		Nlocals: 1,
	}
	var th Thread
	th.Push(SuStr("local_value")) // Push parameter value outside loop
	for b.Loop() {
		th.invoke(fn, nil)
		th.sp = 1 // Reset to just the parameter
	}
}

func BenchmarkLoadValueUnoptimized(b *testing.B) {
	// Separate Load + Value ops - what LoadValue optimizes
	code := []byte{
		byte(op.Load), 0, // Load local var at index 0
		byte(op.Value), 0, // Value constant at index 0
		byte(op.Return),
	}
	fn := &SuFunc{
		Code: string(code),
		ParamSpec: ParamSpec{
			Values:  []Value{SuStr("test")}, // Value at index 0
			Names:   []string{"x"},          // Local var name at index 0
			Nparams: 1,                      // Treat as parameter so it gets initialized
		},
		Nlocals: 1,
	}
	var th Thread
	th.Push(SuStr("local_value")) // Push parameter value outside loop
	for b.Loop() {
		th.invoke(fn, nil)
		th.sp = 1 // Reset to just the parameter
	}
}

func BenchmarkCatChained(b *testing.B) {
	code := []byte{
		byte(op.Value), 0,
		byte(op.Value), 1,
		byte(op.Cat),
		byte(op.Value), 2,
		byte(op.Cat),
		byte(op.Return),
	}
	fn := &SuFunc{
		Code: string(code),
		ParamSpec: ParamSpec{
			Values: []Value{
				SuStr("hello"),
				SuStr(" "),
				SuStr("world"),
			},
		},
	}
	var th Thread
	for b.Loop() {
		th.invoke(fn, nil)
		th.sp = 0
	}
}

func BenchmarkCatN(b *testing.B) {
	code := []byte{
		byte(op.Value), 0,
		byte(op.Value), 1,
		byte(op.Value), 2,
		byte(op.CatN), 3,
		byte(op.Return),
	}
	fn := &SuFunc{
		Code: string(code),
		ParamSpec: ParamSpec{
			Values: []Value{
				SuStr("hello"),
				SuStr(" "),
				SuStr("world"),
			},
		},
	}
	var th Thread
	for b.Loop() {
		th.invoke(fn, nil)
		th.sp = 0
	}
}
