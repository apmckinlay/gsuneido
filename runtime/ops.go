// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"encoding/base64"
	"errors"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var (
	Zero     Value = SuInt(0)
	One      Value = SuInt(1)
	MinusOne Value = SuInt(-1)
	MaxInt   Value = SuDnum{Dnum: dnum.FromInt(math.MaxInt32)}
	Inf      Value = SuDnum{Dnum: dnum.PosInf}
	NegInf   Value = SuDnum{Dnum: dnum.NegInf}
	True     Value = SuBool(true)
	False    Value = SuBool(false)
	// EmptyStr defined in sustr.go
)

func OpIs(x Value, y Value) Value {
	return SuBool(x.Equal(y))
}

func OpIsnt(x Value, y Value) Value {
	return SuBool(!x.Equal(y))
}

func OpLt(x Value, y Value) Value {
	return SuBool(x.Compare(y) < 0)
}

func OpLte(x Value, y Value) Value {
	return SuBool(x.Compare(y) <= 0)
}

func OpGt(x Value, y Value) Value {
	return SuBool(x.Compare(y) > 0)
}

func OpGte(x Value, y Value) Value {
	return SuBool(x.Compare(y) >= 0)
}

func OpAdd(x Value, y Value) Value {
	if xi, xok := SuIntToInt(x); xok {
		if yi, yok := SuIntToInt(y); yok {
			return IntVal(xi + yi)
		}
	}
	return SuDnum{Dnum: dnum.Add(ToDnum(x), ToDnum(y))}
}

func OpSub(x Value, y Value) Value {
	if xi, xok := SuIntToInt(x); xok {
		if yi, yok := SuIntToInt(y); yok {
			return IntVal(xi - yi)
		}
	}
	return SuDnum{Dnum: dnum.Sub(ToDnum(x), ToDnum(y))}
}

func OpMul(x Value, y Value) Value {
	if xi, xok := SuIntToInt(x); xok {
		if yi, yok := SuIntToInt(y); yok {
			return IntVal(xi * yi)
		}
	}
	return SuDnum{Dnum: dnum.Mul(ToDnum(x), ToDnum(y))}
}

func OpDiv(x Value, y Value) Value {
	if yi, yok := SuIntToInt(y); yok && yi != 0 {
		if xi, xok := SuIntToInt(x); xok {
			if xi%yi == 0 {
				return IntVal(xi / yi)
			}
		}
	}
	return SuDnum{Dnum: dnum.Div(ToDnum(x), ToDnum(y))}
}

func OpMod(x Value, y Value) Value {
	return IntVal(ToInt(x) % ToInt(y))
}

func OpLeftShift(x Value, y Value) Value {
	result := int32(ToInt(x)) << ToInt(y)
	return IntVal(int(result))
}

func OpRightShift(x Value, y Value) Value {
	result := uint32(ToInt(x)) >> ToInt(y)
	return IntVal(int(result))
}

func OpBitOr(x Value, y Value) Value {
	return IntVal(ToInt(x) | ToInt(y))
}

func OpBitAnd(x Value, y Value) Value {
	return IntVal(ToInt(x) & ToInt(y))
}

func OpBitXor(x Value, y Value) Value {
	return IntVal(ToInt(x) ^ ToInt(y))
}

func OpBitNot(x Value) Value {
	return IntVal(^ToInt(x))
}

func OpNot(x Value) Value {
	if x == True {
		return False
	} else if x == False {
		return True
	}
	panic("not requires boolean")
}

func OpBool(x Value) bool {
	switch x {
	case True:
		return true
	case False:
		return false
	default:
		panic("conditionals require true or false")
	}
}

func OpUnaryPlus(x Value) Value {
	if _, ok := x.(*smi); ok {
		return x
	}
	return SuDnum{Dnum: ToDnum(x)}
}

func OpUnaryMinus(x Value) Value {
	if xi, ok := SuIntToInt(x); ok {
		return IntVal(-xi)
	}
	return SuDnum{Dnum: ToDnum(x).Neg()}
}

func OpCat(t *Thread, x, y Value) Value {
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

func OpMatch(t *Thread, x Value, y Value) SuBool {
	var pat regex.Pattern
	if r, ok := y.(SuRegex); ok {
		pat = r.Pat
	} else if t != nil {
		pat = t.RxCache.Get(ToStr(y))
	} else {
		pat = regex.Compile(ToStr(y))
	}
	return SuBool(pat.Matches(ToStr(x)))
}

// ToIndex is used by ranges and string[i]
func ToIndex(key Value) int {
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

func OpIter(x Value) SuIter {
	iterable, ok := x.(interface{ Iter() Iter })
	if !ok {
		panic("can't iterate " + x.Type().String())
	}
	return SuIter{Iter: iterable.Iter()}
}

func OpCatch(t *Thread, e any, catchPat string) *SuExcept {
	se := ToSuExcept(t, e)
	if catchMatch(string(se.SuStr), catchPat) {
		return se
	}
	panic(se) // propagate panic if not caught
}

// ToSuExcept converts to SuExcept, and also logs runtime and assert errors
func ToSuExcept(t *Thread, e any) *SuExcept {
	se, ok := e.(*SuExcept)
	if !ok {
		// first catch creates SuExcept with callstack
		var ss SuStr
		switch e := e.(type) {
		case error:
			var perr runtime.Error
			if errors.As(e, &perr) {
				log.Println("ERROR", e)
				dbg.PrintStack()
				printSuStack(t, e)
			}
			ss = SuStr(e.Error())
		case string:
			if isLocalAssertFail(e) {
				// assert has already logged error and Go call stack
				printSuStack(t, e)
			}
			ss = SuStr(e)
		default:
			ss = SuStr(ToStr(e.(Value)))
		}
		se = NewSuExcept(t, ss)
	}
	return se
}

func isLocalAssertFail(s string) bool {
	return strings.HasPrefix(s, "assert failed") &&
		!strings.HasSuffix(s, "(from server)")
}

func printSuStack(th *Thread, e any) {
	if se, ok := e.(*SuExcept); ok {
		PrintStack(se.Callstack)
	} else {
		PrintStack(th.Callstack())
	}
}

func LogUncaught(th *Thread, where string, e any) {
	log.Println("ERROR", th.Name, "uncaught in", where+":", e)
	if isRuntimeError(e) {
		dbg.PrintStack()
	}
	printSuStack(th, e)
}

func isRuntimeError(e any) bool {
	switch e := e.(type) {
	case runtime.Error:
		return true
	case error:
		var perr runtime.Error
		return errors.As(e, &perr)
	}
	return false
}

// catchMatch matches an exception string with a catch pattern
func catchMatch(e, pat string) bool {
	for {
		p := pat
		i := strings.IndexByte(p, '|')
		if i >= 0 {
			pat = pat[i+1:]
			p = p[:i]
		}
		if strings.HasPrefix(p, "*") {
			if strings.Contains(e, p[1:]) {
				return true
			}
		} else if strings.HasPrefix(e, p) {
			return true
		}
		if i < 0 {
			break
		}
	}
	return false
}

func Unpack64(s string) Value {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic("Unpack64 bad data")
	}
	return Unpack(hacks.BStoS(data))
}
