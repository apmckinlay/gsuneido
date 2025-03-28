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
	for range b.N {
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
	for range b.N {
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
	for range b.N {
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
	for range b.N {
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

func BenchmarkInterp2(b *testing.B) {
	code := []byte{byte(op.Int), 0, 3, byte(op.Int), 0, 5, byte(op.Add),
		byte(op.Return)}
	fn := &SuFunc{Code: string(code)}
	var th Thread
	for b.Loop() {
		th.invoke2(fn, nil)
        th.sp = 0
	}
}

func BenchmarkJit2(b *testing.B) {
	var th Thread
	for b.Loop() {
		block(&th)
	}
}

//go:noinline
func block(th *Thread) {
	push(th, SuInt(3))
	push(th, SuInt(5))
	add(th)
	th.Pop()
}

//go:noinline
func push(th *Thread, v Value) {
	th.Push(v)
}

//go:noinline
func add(th *Thread) {
	th.Push(OpAdd(th.Pop(), th.Pop()))
}

func BenchmarkInterp3(b *testing.B) {
	push3 := func(th *Thread) {
		th.Push(SuInt(3))
	}
	push5 := func(th *Thread) {
		th.Push(SuInt(5))
	}
	add := func(th *Thread) {
		th.Push(OpAdd(th.Pop(), th.Pop()))
	}
	pop := func(th *Thread) {
		th.Pop()
	}
	fns := [83]func(*Thread){10: push3, 20: push5, 30: add, 40: pop}
	code := []byte{10, 20, 30, 40}
	var th Thread
	for b.Loop() {
		for _, i := range code {
			fns[i](&th)
		}
	}
}

func BenchmarkInterp4(b *testing.B) {
	code := []byte{10, 20, 30, 40}
    var th Thread
    for b.Loop() {
        for _, i := range code {
            switch i {
            case 0:
                th.Push(SuInt(3))
            case 1:
                th.Push(SuInt(5))
            case 2:
                th.Push(OpAdd(th.Pop(), th.Pop()))
            case 3:
                th.Pop()
            case 4:
                th.Push(SuInt(4))
            case 5:
                th.Push(SuInt(6))
            case 6:
                th.Push(OpSub(th.Pop(), th.Pop()))
            case 7:
                th.Push(OpMul(th.Pop(), th.Pop()))
            case 8:
                th.Push(OpDiv(th.Pop(), th.Pop()))
            case 9:
                th.Push(OpMod(th.Pop(), th.Pop()))
            case 10:
                th.Push(SuInt(3))
            case 11:
                th.Push(OpLte(th.Pop(), th.Pop()))
            case 12:
                th.Push(OpGt(th.Pop(), th.Pop()))
            case 13:
                th.Push(OpGte(th.Pop(), th.Pop()))
            case 14:
                th.Push(OpIs(th.Pop(), th.Pop()))
            case 15:
                th.Push(OpIsnt(th.Pop(), th.Pop()))
            case 16:
                th.Push(OpBitAnd(th.Pop(), th.Pop()))
            case 17:
                th.Push(OpBitOr(th.Pop(), th.Pop()))
            case 18:
                th.Push(OpNot(th.Pop()))
            case 19:
                th.Push(OpBitAnd(th.Pop(), th.Pop()))
            case 20:
                th.Push(SuInt(5))
            case 21:
                th.Push(OpBitXor(th.Pop(), th.Pop()))
            case 22:
                th.Push(OpBitNot(th.Pop()))
            case 23:
                th.Push(OpLeftShift(th.Pop(), th.Pop()))
            case 24:
                th.Push(OpRightShift(th.Pop(), th.Pop()))
            case 25:
                th.Push(SuInt(25))
            case 26:
                th.Push(SuInt(26))
            case 27:
                th.Push(SuInt(27))
            case 28:
                th.Push(SuInt(28))
            case 29:
                th.Push(SuInt(29))
            case 30:
                th.Push(OpAdd(th.Pop(), th.Pop()))
            case 31:
                th.Push(SuInt(31))
            case 32:
                th.Push(SuInt(32))
            case 33:
                th.Push(SuInt(33))
            case 34:
                th.Push(SuInt(34))
            case 35:
                th.Push(SuInt(35))
            case 36:
                th.Push(SuInt(36))
            case 37:
                th.Push(SuInt(37))
            case 38:
                th.Push(SuInt(38))
            case 39:
                th.Push(SuInt(39))
            case 40:
                th.Pop()
            case 41:
                th.Push(SuInt(41))
            case 42:
                th.Push(SuInt(42))
            case 43:
                th.Push(SuInt(43))
            case 44:
                th.Push(SuInt(44))
            case 45:
                th.Push(SuInt(45))
            case 46:
                th.Push(SuInt(46))
            case 47:
                th.Push(SuInt(47))
            case 48:
                th.Push(SuInt(48))
            case 49:
                th.Push(SuInt(49))
            case 50:
                th.Push(SuInt(50))
            case 51:
                th.Push(SuInt(51))
            case 52:
                th.Push(SuInt(52))
            case 53:
                th.Push(SuInt(53))
            case 54:
                th.Push(SuInt(54))
            case 55:
                th.Push(SuInt(55))
            case 56:
                th.Push(SuInt(56))
            case 57:
                th.Push(SuInt(57))
            case 58:
                th.Push(SuInt(58))
            case 59:
                th.Push(SuInt(59))
            case 60:
                th.Push(SuInt(60))
            case 61:
                th.Push(SuInt(61))
            case 62:
                th.Push(SuInt(62))
            case 63:
                th.Push(SuInt(63))
            case 64:
                th.Push(SuInt(64))
            case 65:
                th.Push(SuInt(65))
            case 66:
                th.Push(SuInt(66))
            case 67:
                th.Push(SuInt(67))
            case 68:
                th.Push(SuInt(68))
            case 69:
                th.Push(SuInt(69))
            case 70:
                th.Push(SuInt(70))
            case 71:
                th.Push(SuInt(71))
            case 72:
                th.Push(SuInt(72))
            case 73:
                th.Push(SuInt(73))
            case 74:
                th.Push(SuInt(74))
            case 75:
                th.Push(SuInt(75))
            case 76:
                th.Push(SuInt(76))
            case 77:
                th.Push(SuInt(77))
            case 78:
                th.Push(SuInt(78))
            case 79:
                th.Push(SuInt(79))
            case 80:
                th.Push(SuInt(80))
            case 81:
                th.Push(SuInt(81))
            case 82:
                th.Push(SuInt(82))
            case 83:
                th.Push(SuInt(83))
            }
        }
    }
}
