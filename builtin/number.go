// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math"
	"strconv"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

var _ = builtin(Number, "(value) :number")

func Number(th *Thread, args []Value) Value {
	val := args[0]
	if val.Type() == types.Number {
		return val
	}
	if s, ok := val.ToStr(); ok {
		s = strings.TrimSpace(s)
		if s == "" {
			return Zero
		}
		s = strings.ReplaceAll(s, ",", "")
		s = strings.ReplaceAll(s, "_", "")
		return numFromString(s)
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

var _ = exportMethods(&NumMethods, "num")

var _ = method(num_Binary, "() :string")

func num_Binary(this Value) Value {
	n := ToInt(this)
	return SuStr(strconv.FormatUint(uint64(n), 2))
}

var _ = method(num_Chr, "() :string")

func num_Chr(this Value) Value {
	return SuStr1s[ToInt(this)&0xff]
}

var _ = method(num_Int, "() :number")

func num_Int(this Value) Value {
	dn := ToDnum(this).Trunc()
	if dnum.Compare(dn, minNarrow) >= 0 && dnum.Compare(dn, maxNarrow) <= 0 {
		n, _ := dn.ToInt()
		return SuInt(n)
	}
	return SuDnum{Dnum: dn}
}

var _ = method(num_Format, "(mask :string) :string")

func num_Format(this, arg Value) Value {
	x := ToDnum(this)
	mask := ToStr(arg)
	return SuStr(x.Format(mask))
}

var _ = method(num_Frac, "() :number")

func num_Frac(this Value) Value {
	dn := ToDnum(this).Frac()
	if dn.IsZero() {
		return Zero
	}
	return SuDnum{Dnum: dn}
}

var _ = method(num_Hex, "() :string")

func num_Hex(this Value) Value {
	n := ToInt(this)
	return SuStr(strconv.FormatUint(uint64(n), 16))
}

var _ = method(num_Round, "(number) :number")

func num_Round(this, arg Value) Value {
	x := ToDnum(this)
	r := ToInt(arg)
	return SuDnum{Dnum: x.Round(r, dnum.HalfUp)}
}

var _ = method(num_RoundUp, "(number) :number")

func num_RoundUp(this, arg Value) Value {
	x := ToDnum(this)
	r := ToInt(arg)
	return SuDnum{Dnum: x.Round(r, dnum.Up)}
}

var _ = method(num_RoundDown, "(number) :number")

func num_RoundDown(this, arg Value) Value {
	x := ToDnum(this)
	r := ToInt(arg)
	return SuDnum{Dnum: x.Round(r, dnum.Down)}
}

// float methods

var _ = method(num_Cos, "() :number")

func num_Cos(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Cos(f))
}

var _ = method(num_Sin, "() :number")

func num_Sin(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Sin(f))
}

var _ = method(num_Tan, "() :number")

func num_Tan(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Tan(f))
}

var _ = method(num_ACos, "() :number")

func num_ACos(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Acos(f))
}

var _ = method(num_ASin, "() :number")

func num_ASin(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Asin(f))
}

var _ = method(num_ATan, "() :number")

func num_ATan(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Atan(f))
}

var _ = method(num_Exp, "() :number")

func num_Exp(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Exp(f))
}

var _ = method(num_Log, "() :number")

func num_Log(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Log(f))
}

var _ = method(num_Log2, "() :number")

func num_Log2(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Log2(f))
}

var _ = method(num_Log10, "() :number")

func num_Log10(this Value) Value {
	f := toFloat(this)
	return fromFloat(math.Log10(f))
}

var _ = method(num_Pow, "(number) :number")

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

var _ = method(num_Sqrt, "() :number")

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

var _ = builtin(Max, "(@args) :unknown")

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

var _ = builtin(Min, "(@args) :unknown")

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
