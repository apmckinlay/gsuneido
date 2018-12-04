package builtin

import (
	"math"
	"strconv"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

var minNarrow = dnum.FromInt(MinSuInt)
var maxNarrow = dnum.FromInt(MaxSuInt)

func init() {
	NumMethods = Methods{
		"Chr": method0(func(this Value) Value {
			n := this.ToInt()
			return SuStr(string(rune(n)))
		}),
		"Int": method0(func(this Value) Value {
			dn := this.ToDnum().Int()
			if dnum.Compare(dn, minNarrow) >= 0 && dnum.Compare(dn, maxNarrow) <= 0 {
				n, _ := dn.ToInt()
				return SuInt(n)
			}
			return SuDnum{dn}
		}),
		"Format": method1("(mask)", func(this, arg Value) Value {
			x := this.ToDnum()
			mask := arg.ToStr()
			return SuStr(x.Format(mask))
		}),
		"Frac": method0(func(this Value) Value {
			dn := this.ToDnum().Frac()
			if dn.IsZero() {
				return Zero
			}
			return SuDnum{dn}
		}),
		"Hex": method0(func(this Value) Value {
			n := this.ToInt()
			return SuStr(strconv.FormatUint(uint64(uint32(n)), 16))
		}),

		"Round": method1("(number)", func(this, arg Value) Value {
			x := this.ToDnum()
			r := arg.ToInt()
			return SuDnum{x.Round(r, dnum.HalfUp)}
		}),
		"RoundUp": method1("(number)", func(this, arg Value) Value {
			x := this.ToDnum()
			r := arg.ToInt()
			return SuDnum{x.Round(r, dnum.Up)}
		}),
		"RoundDown": method1("(number)", func(this, arg Value) Value {
			x := this.ToDnum()
			r := arg.ToInt()
			return SuDnum{x.Round(r, dnum.Down)}
		}),

		// float methods

		"Cos": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Cos(f))
		}),
		"Sin": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Sin(f))
		}),
		"Tan": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Tan(f))
		}),

		"ACos": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Acos(f))
		}),
		"ASin": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Asin(f))
		}),
		"ATan": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Atan(f))
		}),

		"Exp": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Exp(f))
		}),
		"Log": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Log(f))
		}),
		"Log10": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Log10(f))
		}),
		"Pow": method1("(number)", func(this, arg Value) Value {
			x := toFloat(this)
			y := toFloat(arg)
			return fromFloat(math.Pow(x, y))
		}),
		"Sqrt": method0(func(this Value) Value {
			f := toFloat(this)
			return fromFloat(math.Sqrt(f))
		}),
	}
}

func toFloat(v Value) float64 {
	if i, ok := SmiToInt(v); ok {
		return float64(i)
	}
	return v.ToDnum().ToFloat()
}

func fromFloat(f float64) Value {
	n := int64(f)
	if f == float64(n) {
		if MinSuInt <= n && n <= MaxSuInt {
			return SuInt(int(n))
		}
		return SuDnum{dnum.FromInt(n)}
	}
	return SuDnum{dnum.FromFloat(f)}
}
