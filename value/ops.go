package value

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/regex"
)

func Is(x Value, y Value) Value {
	return SuBool(x.Equals(y))
}

func Isnt(x Value, y Value) Value {
	return SuBool(!x.Equals(y))
}

func Lt(x Value, y Value) Value {
	return SuBool(x.cmp(y) < 0)
}

func Lte(x Value, y Value) Value {
	return SuBool(x.cmp(y) <= 0)
}

func Gt(x Value, y Value) Value {
	return SuBool(x.cmp(y) > 0)
}

func Gte(x Value, y Value) Value {
	return SuBool(x.cmp(y) >= 0)
}

func Add(x Value, y Value) Value {
	if xi, xok := x.(SuInt); xok {
		if yi, yok := y.(SuInt); yok {
			return Int64ToValue(int64(xi) + int64(yi))
		}
	}
	return DnumToValue(dnum.Add(x.ToDnum(), y.ToDnum()))
}

func Sub(x Value, y Value) Value {
	if xi, xok := x.(SuInt); xok {
		if yi, yok := y.(SuInt); yok {
			return Int64ToValue(int64(xi) - int64(yi))
		}
	}
	return DnumToValue(dnum.Sub(x.ToDnum(), y.ToDnum()))
}

func Mul(x Value, y Value) Value {
	if xi, xok := x.(SuInt); xok {
		if yi, yok := y.(SuInt); yok {
			return Int64ToValue(int64(xi) * int64(yi))
		}
	}
	return DnumToValue(dnum.Mul(x.ToDnum(), y.ToDnum()))
}

func Div(x Value, y Value) Value {
	// TODO check if it's worth trying int division first
	// i.e. if x and y are ints and x % y == 0, then return x / y
	// could instrument existing suneido to see how common this is
	return DnumToValue(dnum.Div(x.ToDnum(), y.ToDnum()))
}

func Mod(x Value, y Value) Value {
	return Int64ToValue(int64(x.ToInt()) % int64(y.ToInt()))
}

func Lshift(x Value, y Value) Value {
	return Int64ToValue(int64(uint64(x.ToInt()) << uint64(y.ToInt())))
}

func Rshift(x Value, y Value) Value {
	return Int64ToValue(int64(uint64(x.ToInt()) >> uint64(y.ToInt())))
}

func Bitor(x Value, y Value) Value {
	return Int64ToValue(int64(uint64(x.ToInt()) | uint64(y.ToInt())))
}

func Bitand(x Value, y Value) Value {
	return Int64ToValue(int64(uint64(x.ToInt()) & uint64(y.ToInt())))
}

func Bitxor(x Value, y Value) Value {
	return Int64ToValue(int64(uint64(x.ToInt()) ^ uint64(y.ToInt())))
}

func Bitnot(x Value) Value {
	return Int64ToValue(^int64(x.ToInt()))
}

func Not(x Value) Value {
	if x == True {
		return False
	} else if x == False {
		return True
	}
	panic("not requires boolean")
}

func Uplus(x Value) Value {
	if _, ok := x.(SuInt); ok {
		return x
	} else if _, ok := x.(SuDnum); ok {
		return x
	}
	return DnumToValue(x.ToDnum())
}

func Uminus(x Value) Value {
	if xi, ok := x.(SuInt); ok {
		return Int64ToValue(-int64(xi))
	}
	return DnumToValue(x.ToDnum().Neg())
}

// Int64ToValue returns an SuInt if it fits, else a SuDnum
func Int64ToValue(n int64) Value {
	if math.MinInt32 < n && n < math.MaxInt32 {
		return SuInt(int32(n))
	} else {
		return SuDnum{dnum.FromInt64(n)}
	}
}

// DnumToValue returns an SuInt if it fits, else a SuDnum
func DnumToValue(dn dnum.Dnum) Value {
	if dn.IsInt() {
		if n, err := dn.Int32(); err == nil {
			return SuInt(n)
		}
	}
	return SuDnum{dn}
}

func Cat(x Value, y Value) Value {
	const SMALL = 256

	xc, xcok := x.(SuConcat)
	yc, ycok := y.(SuConcat)
	if xcok && ycok {
		return xc.AddSuConcat(yc)
	} else if xcok {
		return xc.Add(y.ToStr())
	} else if ycok {
		return NewSuConcat().Add(x.ToStr()).AddSuConcat(yc)
	} else {
		xs := x.ToStr()
		ys := y.ToStr()
		if len(xs)+len(ys) < SMALL {
			return SuStr(xs + ys)
		} else {
			return NewSuConcat().Add(xs).Add(ys)
		}
	}
}

func BitNot(x Value) Value {
	return SuInt(^x.ToInt())
}

func Cmp(x Value, y Value) int {
	xo := x.order()
	yo := y.order()
	if xo != yo {
		return cmpInt(int(xo), int(yo))
	}
	return x.cmp(y)
}

func cmpInt(x int, y int) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return +1
	default:
		return 0
	}
}

func Match(x Value, y regex.Pattern) SuBool {
	return SuBool(y.Matches(x.ToStr()))
}
