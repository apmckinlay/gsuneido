// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/options"
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
	return SuBool(strictCompare(x, y) < 0)
}

func OpLte(x Value, y Value) Value {
	return SuBool(strictCompare(x, y) <= 0)
}

func OpGt(x Value, y Value) Value {
	return SuBool(strictCompare(x, y) > 0)
}

func OpGte(x Value, y Value) Value {
	return SuBool(strictCompare(x, y) >= 0)
}

func strictCompare(x Value, y Value) int {
	cmp := x.Compare(y)
	if (cmp&3) == 2 && options.StrictCompare && x != False && y != False {
		panic(fmt.Sprint("StrictCompare: ", x, " <=> ", y))
	}
	return cmp
}

func OpInRange(x Value, orgOp tok.Token, org Value, endOp tok.Token, end Value) Value {
	if (orgOp == tok.Gt && !(x.Compare(org) > 0)) || !(x.Compare(org) >= 0) {
		return False
	}
	if (endOp == tok.Lt && !(x.Compare(end) < 0)) || !(x.Compare(end) <= 0) {
		return False
	}
	return True
}

func OpAdd(x Value, y Value) Value {
	if xi, xok := SuIntToInt(x); xok {
		if yi, yok := SuIntToInt(y); yok {
			return IntVal(xi + yi)
		}
	}
	return SuDnum{Dnum: dnum.Add(ToDnum(x), ToDnum(y))}
}

func OpAdd1(x Value) Value {
	if n, ok := SuIntToInt(x); ok {
		return IntVal(n + 1)
	}
	return SuDnum{Dnum: dnum.Add(ToDnum(x), dnum.One)}
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
	result := ToInt(x) << ToInt(y)
	return IntVal(result)
}

func OpRightShift(x Value, y Value) Value {
	result := uint(ToInt(x)) >> ToInt(y)
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
	switch x {
	case True:
		return False
	case False:
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
	if x.Type() == types.Number {
		return x
	}
	if x == EmptyStr || x == False {
		return Zero
	}
	panic("can't convert " + ErrType(x) + " to number")
}

func OpUnaryMinus(x Value) Value {
	if xi, ok := SuIntToInt(x); ok {
		return IntVal(-xi)
	}
	if x == EmptyStr || x == False {
		return Zero
	}
	return SuDnum{Dnum: ToDnum(x).Neg()}
}

func OpCat(th *Thread, x, y Value) Value {
	if ssx, ok := x.(SuStr); ok {
		if ssy, ok := y.(SuStr); ok {
			return cat2(string(ssx), string(ssy))
		}
	}
	return cat3(th, x, y)
}

func cat2(xs, ys string) Value {
	const LARGE = 256 // ??? exact value not critical

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

func cat3(th *Thread, x, y Value) Value {
	var result Value
	if xc, ok := x.(SuConcat); ok {
		result = xc.Add(catToStr(th, y))
	} else {
		result = cat2(catToStr(th, x), catToStr(th, y))
	}
	if xe, ok := x.(*SuExcept); ok {
		return &SuExcept{SuStr: SuStr(AsStr(result)), Callstack: xe.Callstack}
	}
	if ye, ok := y.(*SuExcept); ok {
		return &SuExcept{SuStr: SuStr(AsStr(result)), Callstack: ye.Callstack}
	}
	return result
}

func catToStr(th *Thread, v Value) string {
	if d, ok := v.(ToStringable); ok {
		return d.ToString(th)
	}
	return AsStr(v)
}

func OpCatN(th *Thread, count int) Value {
	values := th.stack[th.sp-count : th.sp]
	totalLen := 0
	var firstExcept *SuExcept
	for i, v := range values {
		if firstExcept == nil {
			if xe, ok := v.(*SuExcept); ok {
				firstExcept = xe
			}
		}
		if ss, ok := v.(SuStr); ok {
			totalLen += len(ss)
		} else {
			ss = SuStr(catToStr(th, v))
			values[i] = ss
			totalLen += len(ss)
		}
	}
	result := make([]byte, totalLen)
	pos := 0
	for _, v := range values {
		pos += copy(result[pos:], string(v.(SuStr)))
	}
	th.sp -= count // pop
	resultStr := SuStr(hacks.BStoS(result))
	if firstExcept != nil {
		return &SuExcept{SuStr: resultStr, Callstack: firstExcept.Callstack}
	}
	return resultStr
}

func OpMatch(th *Thread, x Value, y Value) SuBool {
	var pat regex.Pattern
	if th != nil {
		pat = th.Regex(y)
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

type iter2 = func() (Value, Value)
type iter2able interface{ Iter2(bool, bool) iter2 }

var _ iter2able = (*SuObject)(nil)

func OpIter2(x Value) SuIter2 {
	iterable, ok := x.(iter2able)
	if !ok {
		panic("can't iterate " + x.Type().String())
	}
	return SuIter2{iter2: iterable.Iter2(true, true)}
}

func OpCatch(th *Thread, e any, catchPat string) *SuExcept {
	se := ToSuExcept(th, e)
	if catchMatch(string(se.SuStr), catchPat) {
		return se
	}
	panic(se) // propagate panic if not caught
}

// ToSuExcept converts to SuExcept, and also logs runtime and assert errors
func ToSuExcept(th *Thread, e any) *SuExcept {
	se, ok := e.(*SuExcept)
	if !ok {
		// first catch creates SuExcept with callstack
		var ss SuStr
		switch e := e.(type) {
		case error:
			var perr runtime.Error
			if errors.As(e, &perr) {
				log.Println("ERROR:", e)
				dbg.PrintStack()
				printSuStack(th, e)
			}
			ss = SuStr(e.Error())
		case string:
			logStringError(th, e)
			ss = SuStr(e)
		default:
			ss = SuStr(ToStr(e.(Value)))
		}
		se = NewSuExcept(th, ss)
	}
	return se
}

func printSuStack(th *Thread, e any) {
	if th.Name != "" {
		fmt.Fprintln(os.Stderr, th.Name)
	}
	if se, ok := e.(*SuExcept); ok {
		PrintStack(se.Callstack)
	} else {
		th.PrintStack()
	}
}

// LogInternalError logs the error and the call stacks, if an InternalError.
// It is used by dbmsserver
func LogInternalError(th *Thread, from string, e any) {
	if isRuntimeError(e) {
		log.Println("ERROR:", from, e)
		dbg.PrintStack()
		printSuStack(th, e)
	} else if s, ok := e.(string); ok {
		logStringError(th, s)
	}
}

func logStringError(th *Thread, e string) {
	if strings.HasPrefix(e, "ASSERT FAILED") &&
		!strings.HasSuffix(e, "(from server)") {
		// assert has already logged error and Go call stack
		printSuStack(th, e)
	}
}

func LogUncaught(th *Thread, where string, e any) {
	log.Println("ERROR:", "uncaught in", where+":", e)
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
