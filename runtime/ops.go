// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var (
	Zero   Value = SuInt(0)
	One    Value = SuInt(1)
	MaxInt Value = SuDnum{Dnum: dnum.FromInt(math.MaxInt32)}
	Inf    Value = SuDnum{Dnum: dnum.PosInf}
	NegInf Value = SuDnum{Dnum: dnum.NegInf}
	True   Value = SuBool(true)
	False  Value = SuBool(false)
	// EmptyStr defined in sustr.go
)

func Is(x Value, y Value) Value {
	return SuBool(x == y || x.Equal(y))
}

func Isnt(x Value, y Value) Value {
	return SuBool(!x.Equal(y))
}

func Lt(x Value, y Value) Value {
	return SuBool(x.Compare(y) < 0)
}

func Lte(x Value, y Value) Value {
	return SuBool(x.Compare(y) <= 0)
}

func Gt(x Value, y Value) Value {
	return SuBool(x.Compare(y) > 0)
}

func Gte(x Value, y Value) Value {
	return SuBool(x.Compare(y) >= 0)
}

func Add(x Value, y Value) Value {
	if xi, xok := SuIntToInt(x); xok {
		if yi, yok := SuIntToInt(y); yok {
			return IntVal(xi + yi)
		}
	}
	return SuDnum{Dnum: dnum.Add(ToDnum(x), ToDnum(y))}
}

func Sub(x Value, y Value) Value {
	if xi, xok := SuIntToInt(x); xok {
		if yi, yok := SuIntToInt(y); yok {
			return IntVal(xi - yi)
		}
	}
	return SuDnum{Dnum: dnum.Sub(ToDnum(x), ToDnum(y))}
}

func Mul(x Value, y Value) Value {
	if xi, xok := SuIntToInt(x); xok {
		if yi, yok := SuIntToInt(y); yok {
			return IntVal(xi * yi)
		}
	}
	return SuDnum{Dnum: dnum.Mul(ToDnum(x), ToDnum(y))}
}

func Div(x Value, y Value) Value {
	if yi, yok := SuIntToInt(y); yok && yi != 0 {
		if xi, xok := SuIntToInt(x); xok {
			if xi%yi == 0 {
				return IntVal(xi / yi)
			}
		}
	}
	return SuDnum{Dnum: dnum.Div(ToDnum(x), ToDnum(y))}
}

func Mod(x Value, y Value) Value {
	return IntVal(ToInt(x) % ToInt(y))
}

func LeftShift(x Value, y Value) Value {
	result := int32(ToInt(x)) << ToInt(y)
	return IntVal(int(result))
}

func RightShift(x Value, y Value) Value {
	result := uint32(ToInt(x)) >> ToInt(y)
	return IntVal(int(result))
}

func BitOr(x Value, y Value) Value {
	return IntVal(ToInt(x) | ToInt(y))
}

func BitAnd(x Value, y Value) Value {
	return IntVal(ToInt(x) & ToInt(y))
}

func BitXor(x Value, y Value) Value {
	return IntVal(ToInt(x) ^ ToInt(y))
}

func BitNot(x Value) Value {
	return IntVal(^ToInt(x))
}

func Not(x Value) Value {
	if x == True {
		return False
	} else if x == False {
		return True
	}
	panic("not requires boolean")
}

func Bool(x Value) bool {
	switch x {
	case True:
		return true
	case False:
		return false
	default:
		panic("conditionals require true or false")
	}
}

func UnaryPlus(x Value) Value {
	if _, ok := x.(*smi); ok {
		return x
	}
	return SuDnum{Dnum: ToDnum(x)}
}

func UnaryMinus(x Value) Value {
	if xi, ok := SuIntToInt(x); ok {
		return IntVal(-xi)
	}
	return SuDnum{Dnum: ToDnum(x).Neg()}
}

func Cat(t *Thread, x, y Value) Value {
	if ssx, ok := x.(SuStr); ok {
		if ssy, ok := y.(SuStr); ok {
			return cat2(string(ssx), string(ssy))
		}
	}
	return cat3(t, x, y)
}

func cat2(xs, ys string) Value {
	const LARGE = 256

	if len(xs)+len(ys) < LARGE {
		return SuStr(xs + ys)
	}
	if len(xs) == 0 {
		return SuStr(ys)
	}
	if len(ys) == 0 {
		return SuStr(xs)
	}
	return NewSuConcat().Add(xs).Add(ys)
}

func cat3(t *Thread, x, y Value) Value {
	var result Value
	if xc, ok := x.(SuConcat); ok {
		result = xc.Add(catToStr(t, y))
	} else {
		result = cat2(catToStr(t, x), catToStr(t, y))
	}
	if xe, ok := x.(*SuExcept); ok {
		return &SuExcept{SuStr: SuStr(AsStr(result)), Callstack: xe.Callstack}
	}
	if ye, ok := y.(*SuExcept); ok {
		return &SuExcept{SuStr: SuStr(AsStr(result)), Callstack: ye.Callstack}
	}
	return result
}

func catToStr(t *Thread, v Value) string {
	if d, ok := v.(ToStringable); ok {
		return d.ToString(t)
	}
	return AsStr(v)
}

func Match(x Value, y regex.Pattern) SuBool {
	return SuBool(y.Matches(ToStr(x)))
}

// Index is used by ranges and string[i]
func Index(key Value) int {
	if n, ok := key.IfInt(); ok {
		return n
	}
	panic("indexes must be integers")
}

func prepFrom(from int, size int) int {
	if from < 0 {
		from += size
		if from < 0 {
			from = 0
		}
	}
	if from > size {
		from = size
	}
	return from
}

func prepTo(from int, to int, size int) int {
	if to < 0 {
		to += size
	}
	if to < from {
		to = from
	}
	if to > size {
		to = size
	}
	return to
}

func prepLen(len int, size int) int {
	if len < 0 {
		len = 0
	}
	if len > size {
		len = size
	}
	return len
}
