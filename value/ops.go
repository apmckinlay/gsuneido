package value

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

func Add(x Value, y Value) Value {
	if xi, xok := x.(IntVal); xok {
		if yi, yok := y.(IntVal); yok {
			return Int64ToValue(int64(xi) + int64(yi))
		}
	}
	return DnumToValue(dnum.Add(x.ToDnum(), y.ToDnum()))
}

func Sub(x Value, y Value) Value {
	if xi, xok := x.(IntVal); xok {
		if yi, yok := y.(IntVal); yok {
			return Int64ToValue(int64(xi) - int64(yi))
		}
	}
	return DnumToValue(dnum.Sub(x.ToDnum(), y.ToDnum()))
}

// Int64ToValue returns an IntVal if it fits, else a DnumVal
func Int64ToValue(n int64) Value {
	if math.MinInt32 < n && n < math.MaxInt32 {
		return IntVal(int32(n))
	} else {
		return DnumVal(dnum.FromInt64(n))
	}
}

// DnumToValue returns an IntVal if it fits, else a DnumVal
func DnumToValue(dn dnum.Dnum) Value {
	if n, err := dn.Int32(); err == nil {
		return IntVal(n)
	} else {
		return DnumVal(dn)
	}
}

const SMALL = 256

func Cat(x Value, y Value) Value {
	xc, xcok := x.(CatVal)
	yc, ycok := y.(CatVal)
	if xcok && ycok {
		return xc.AddCatVal(yc)
	} else if xcok {
		return xc.Add(y.ToStr())
	} else if ycok {
		return NewCatVal().Add(x.ToStr()).AddCatVal(yc)
	} else {
		xs := x.ToStr()
		ys := y.ToStr()
		if len(xs)+len(ys) < SMALL {
			return StrVal(xs + ys)
		} else {
			return NewCatVal().Add(xs).Add(ys)
		}
	}
}

func BitNot(x Value) Value {
	return IntVal(^x.ToInt())
}
