package value

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

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

const SMALL = 256

func Cat(x Value, y Value) Value {
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
