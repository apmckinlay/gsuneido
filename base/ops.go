package base

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var (
	Zero  Value = SuInt(0)
	One   Value = SuInt(1)
	True  Value = SuBool(true)
	False Value = SuBool(false)
	// EmptyStr defined in sustr.go
)

func Is(x Value, y Value) Value {
	return SuBool(x.Equal(y))
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
	if xi, xok := SmiToInt(x); xok {
		if yi, yok := SmiToInt(y); yok {
			return IntToValue(xi + yi)
		}
	}
	return SuDnum{dnum.Add(x.ToDnum(), y.ToDnum())}
}

func Sub(x Value, y Value) Value {
	if xi, xok := SmiToInt(x); xok {
		if yi, yok := SmiToInt(y); yok {
			return IntToValue(xi - yi)
		}
	}
	return SuDnum{dnum.Sub(x.ToDnum(), y.ToDnum())}
}

func Mul(x Value, y Value) Value {
	if xi, xok := SmiToInt(x); xok {
		if yi, yok := SmiToInt(y); yok {
			return IntToValue(xi * yi)
		}
	}
	return SuDnum{dnum.Mul(x.ToDnum(), y.ToDnum())}
}

func Div(x Value, y Value) Value {
	if xi, xok := SmiToInt(x); xok {
		if yi, yok := SmiToInt(y); yok {
			if xi%yi == 0 {
				return IntToValue(xi / yi)
			}
		}
	}
	return SuDnum{dnum.Div(x.ToDnum(), y.ToDnum())}
}

func Mod(x Value, y Value) Value {
	return IntToValue(x.ToInt() % y.ToInt())
}

func Lshift(x Value, y Value) Value {
	return IntToValue(int(uint(x.ToInt()) << uint(y.ToInt())))
}

func Rshift(x Value, y Value) Value {
	return IntToValue(int(uint(x.ToInt()) >> uint(y.ToInt())))
}

func Bitor(x Value, y Value) Value {
	return IntToValue(x.ToInt() | y.ToInt())
}

func Bitand(x Value, y Value) Value {
	return IntToValue(x.ToInt() & y.ToInt())
}

func Bitxor(x Value, y Value) Value {
	return IntToValue(x.ToInt() ^ y.ToInt())
}

func Bitnot(x Value) Value {
	return IntToValue(^x.ToInt())
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
	if _, ok := SmiToInt(x); ok {
		return x
	} else if _, ok := x.(SuDnum); ok {
		return x
	}
	return SuDnum{x.ToDnum()} // "" or false => 0, else throw
}

func Uminus(x Value) Value {
	if xi, ok := SmiToInt(x); ok {
		return IntToValue(-xi)
	}
	return SuDnum{x.ToDnum().Neg()}
}

// IntToValue returns an SuInt if it fits, else a SuDnum
func IntToValue(n int) Value {
	if math.MinInt16 < n && n < math.MaxInt16 {
		return SuInt(n)
	}
	return SuDnum{dnum.FromInt(int64(n))}
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
	}
	xs := x.ToStr()
	ys := y.ToStr()
	if len(xs)+len(ys) < SMALL {
		return SuStr(xs + ys)
	}
	return NewSuConcat().Add(xs).Add(ys)
}

func BitNot(x Value) Value {
	return IntToValue(^x.ToInt())
}

func Match(x Value, y regex.Pattern) SuBool {
	return SuBool(y.Matches(x.ToStr()))
}

// Index is basically the same as value.ToInt
// except it doesn't convert "" and false to 0
// and it has a different error message
// used by ranges and string[i]
func Index(v Value) int {
	if n, ok := SmiToInt(v); ok {
		return n
	}
	if dn, ok := v.(SuDnum); ok {
		if n, ok := dn.Dnum.ToInt(); ok {
			return n
		}
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
