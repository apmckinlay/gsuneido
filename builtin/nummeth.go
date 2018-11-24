package builtin

import (
	"fmt"
	"math"
	"strconv"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

var minNarrow = dnum.FromInt(MinSuInt)
var maxNarrow = dnum.FromInt(MaxSuInt)

func init() {
	NumMethods = Methods{
		// TODO Format, Round, RoundDown, RoundUp

		"Chr": method0(func(self Value) Value {
			n := self.ToInt()
			return SuStr(string(rune(n)))
		}),
		"Int": method0(func(self Value) Value {
			dn := self.ToDnum().Int()
			if dnum.Compare(dn, minNarrow) >= 0 && dnum.Compare(dn, maxNarrow) <= 0 {
				n, _ := dn.ToInt()
				return SuInt(n)
			}
			return SuDnum{dn}
		}),
		"Frac": method0(func(self Value) Value {
			dn := self.ToDnum().Frac()
			if dn.IsZero() {
				return Zero
			}
			return SuDnum{dn}
		}),
		"Hex": method0(func(self Value) Value {
			n := self.ToInt()
			return SuStr(strconv.FormatInt(int64(n), 16))
		}),

		// float methods

		"Cos": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Cos(f))
		}),
		"Sin": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Sin(f))
		}),
		"Tan": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Tan(f))
		}),

		"ACos": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Acos(f))
		}),
		"ASin": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Asin(f))
		}),
		"ATan": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Atan(f))
		}),

		"Exp": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Exp(f))
		}),
		"Log": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Log(f))
		}),
		"Log10": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Log10(f))
		}),
		"Pow": method1("(number)", func(self, arg Value) Value {
			x := toFloat(self)
			y := toFloat(arg)
			return fromFloat(math.Pow(x, y))
		}),
		"Sqrt": method0(func(self Value) Value {
			f := toFloat(self)
			return fromFloat(math.Sqrt(f))
		}),
	}
}

func toFloat(v Value) float64 {
	if i, ok := SmiToInt(v); ok {
		fmt.Println("int toFloat")
		return float64(i)
	}
	return v.ToDnum().ToFloat()
}

func fromFloat(f float64) Value {
	n := int64(f)
	if f == float64(n) {
		if MinSuInt <= n && n <= MaxSuInt {
			fmt.Println("smi fromFloat")
			return SuInt(int(n))
		}
		fmt.Println("int fromFloat")
		return SuDnum{dnum.FromInt(n)}
	}
	return SuDnum{dnum.FromFloat(f)}
}
