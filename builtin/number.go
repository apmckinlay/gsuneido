// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math"
	"strconv"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

var _ = builtin(Number, "(value)")

func Number(th *Thread, args []Value) Value {
	val := args[0]
	if s, ok := val.ToStr(); ok {
		s := strings.TrimSpace(s)
		if s == "" {
			return Zero
		}
<<<<<<< HEAD
	
		s = strings.ReplaceAll(s, ",", "")
<<<<<<< HEAD
  		s = strings.ReplaceAll(s, "_", "")
=======
		s = strings.ReplaceAll(s, "_", "")
>>>>>>> 3f31c5d975e671833591f03d8546cf9f34497776
		return numFromString(s)
=======
	s = strings.ReplaceAll(s, ",", "")
  	s = strings.ReplaceAll(s, "_", "")
	return numFromString(s)
>>>>>>> d153e5f20fa371083edf47945492408e69e22ebc
	}
	if _, ok := val.(SuDnum); ok {
		return val
	}
	if n, ok := SuIntToInt(val); ok {
		return IntVal(n)
	}
	if val == False {
		return Zero
	}
	panic("can't convert " + ErrType(val) + " to number")
}

func numFromString(s string) Value {
	defer func() {
		if e := recover(); e != nil {
			panic("can't convert string to number")
		}
	}()
	return NumFromString(s)
}

var minNarrow = dnum.FromInt(MinSuInt)
var maxNarrow = dnum.FromInt(MaxSuInt)

var _ = exportMethods(&NumMethods)

var _ = method(num_Chr, "()")

func num_Chr(this Value) Value {
	n := byte(ToInt(this))
	return SuStr(string([]byte{n}))
}

var _ = method(num_Int, "()")

func num_Int(this Value) Value {
	dn := ToDnum(this).Trunc()
	if dnum.Compare(dn, minNarrow) >= 0 && dnum.Compare(dn, maxNarrow) <= 0 {
		n, _ := dn.ToInt()
		return SuInt(n)
	}
	return SuDnum{Dnum: dn}
}

var _ = method(num_Format, "(mask)")

func num_Format(this, arg Value) Value {
	x := ToDnum(this)
	mask := ToStr(arg)
	return SuStr(x.Format(mask))
}

var _ = method(num_Frac, "()")

func num_Frac(this Value) Value {
	dn := ToDnum(this).Frac()
	if dn.IsZero() {
		return Zero
	}
	return SuDnum{Dnum: dn}
}

var _ = method(num_Hex, "()")

func num_Hex(this Value) Value {
	n := ToInt(this)
	return SuStr(strconv.FormatUint(uint64(uint32(n)), 16))
}

var _ = method(num_Round, "(number)")

func num_Round(this, arg Value) Value {
	x := ToDnum(this)
	r := ToInt(arg)
	return SuDnum{Dnum: x.Round(r, dnum.HalfUp)}
}

var _ = method(num_RoundUp, "(number)")

func num_RoundUp(this, arg Value) Value {
	x := ToDnum(this)
	r := ToInt(arg)
	return SuDnum{Dnum: x.Round(r, dnum.Up)}
}

var _ = method(num_RoundDown, "(number)")

func num_RoundDown(this, arg Value) Value {
	x := ToDnum(this)
	r := ToInt(arg)
	return SuDnum{Dnum: x.Round(r, dnum.Down)}
}

// float methods

var _ = method(num_Cos, "()")

func num_Cos(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Cos(f))
}

var _ = method(num_Sin, "()")

func num_Sin(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Sin(f))
}

var _ = method(num_Tan, "()")

func num_Tan(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Tan(f))
}

var _ = method(num_ACos, "()")

func num_ACos(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Acos(f))
}

var _ = method(num_ASin, "()")

func num_ASin(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Asin(f))
}

var _ = method(num_ATan, "()")

func num_ATan(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Atan(f))
}

var _ = method(num_Exp, "()")

func num_Exp(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Exp(f))
}

var _ = method(num_Log, "()")

func num_Log(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Log(f))
}

var _ = method(num_Log10, "()")

func num_Log10(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Log10(f))
}

var _ = method(num_Pow, "(number)")

func num_Pow(this, arg Value) Value {
	if p, ok := arg.ToInt(); ok && 0 <= p && p <= 10 {
		if p == 0 {
			return One
		}
		x := this
		for ; p > 1; p-- {
			x = OpMul(x, this)
		}
		return x
	}
	x := toFloat(this)
	y := toFloat(arg)
	return fromFloat(math.Pow(x, y))
}

var _ = method(num_Sqrt, "()")

func num_Sqrt(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Sqrt(f))
}

func toFloat(v Value) float64 {
	if i, ok := v.ToInt(); ok {
		return float64(i)
	}
	return ToDnum(v).ToFloat()
}

func fromFloat(f float64) Value {
	n := int64(f)
	if f == float64(n) {
		if MinSuInt <= n && n <= MaxSuInt {
			return SuInt(int(n))
		}
		return SuDnum{Dnum: dnum.FromInt(n)}
	}
	return SuDnum{Dnum: dnum.FromFloat(f)}
}

// Max and Min aren't specific to numbers,
// but that's normally what they're used for

var _ = builtin(Max, "(@args)")

func Max(_ *Thread, as *ArgSpec, args []Value) Value {
	if as.Nargs == 0 {
		panic("Max requires at least one value")
	}
	if as.Each == 0 {
		max := args[0]
		for i := 1; i < int(as.Nargs); i++ {
			if args[i].Compare(max) > 0 {
				max = args[i]
			}
		}
		return max
	}
	iterable, ok := args[0].(interface{ Iter() Iter })
	if !ok {
		panic("can't iterate " + args[0].Type().String())
	}
	it := iterable.Iter()
	max := it.Next()
	if as.Each == EACH1 && max != nil {
		max = it.Next()
	}
	if max == nil {
		panic("Max requires at least one value")
	}
	for v := it.Next(); v != nil; v = it.Next() {
		if v.Compare(max) > 0 {
			max = v
		}
	}
	return max
}

var _ = builtin(Min, "(@args)")

func Min(_ *Thread, as *ArgSpec, args []Value) Value {
	if as.Nargs == 0 {
		panic("Min requires at least one value")
	}
	if as.Each == 0 {
		min := args[0]
		for i := 1; i < int(as.Nargs); i++ {
			if args[i].Compare(min) < 0 {
				min = args[i]
			}
		}
		return min
	}
	iterable, ok := args[0].(interface{ Iter() Iter })
	if !ok {
		panic("can't iterate " + args[0].Type().String())
	}
	it := iterable.Iter()
	min := it.Next()
	if as.Each == EACH1 && min != nil {
		min = it.Next()
	}
	if min == nil {
		panic("Min requires at least one value")
	}
	for v := it.Next(); v != nil; v = it.Next() {
		if v.Compare(min) < 0 {
			min = v
		}
	}
	return min
}
