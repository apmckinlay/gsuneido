// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"math"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	op "github.com/apmckinlay/gsuneido/core/opcodes"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

//TODO move these to frame.go

func (fr *Frame) fetchUint8() int {
	fr.ip++
	return int(fr.fn.Code[fr.ip-1])
}

func (fr *Frame) fetchInt16() int {
	fr.ip += 2
	return int(int16(uint16(fr.fn.Code[fr.ip-2])<<8 + uint16(fr.fn.Code[fr.ip-1])))
}

func (fr *Frame) fetchUint16() int {
	fr.ip += 2
	return int(uint16(fr.fn.Code[fr.ip-2])<<8 + uint16(fr.fn.Code[fr.ip-1]))
}

func (fr *Frame) jump() {
	fr.ip += fr.fetchInt16()
}

//-------------------------------------------------------------------

// invoke sets up a frame to Run a compiled Suneido function
// The stack must already be in the form required by the function (massaged)
// WARNING: invoke does not pop the stack, the caller is responsible for that
func (th *Thread) invoke2(fn *SuFunc, this Value) Value {
	// reserve stack space for locals
	for expand := fn.Nlocals - fn.Nparams; expand > 0; expand-- {
		th.Push(nil)
	}
	// return th.run2(Frame{fn: fn, this: this,
	// 	locals: locals{v: th.stack[th.sp-int(fn.Nlocals) : th.sp]}})
	fr := &th.frames[th.fp]
	fr.fn = fn
	fr.this = this
	fr.locals = locals{v: th.stack[th.sp-int(fn.Nlocals) : th.sp]}
	fr.ip = 0
	return th.run2()
}

// run is needed in addition to interp
// because we can only recover panic on the way out of a function
// so if the exception is caught we have to re-enter interp
// Called by Thread.Invoke (above) and SuClosure.Call
func (th *Thread) run2() Value {
	if th.fp >= len(th.frames) {
		panic("function call overflow")
	}
	fr := &th.frames[th.fp]
	th.fp++
	if th.profile.enabled {
		th.profile.calls[fr.fn]++
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(strings.Repeat("    ", t.fp) + "run:", t.frames[t.fp].fn)
	sp := th.sp
	fp := th.fp
	if th.fp > th.fpMax {
		th.fpMax = th.fp // track high water mark
	} else if th.fp < th.fpMax {
		// clear the frames up to the high water mark
		for i := th.fp; i < th.fpMax; i++ {
			th.frames[i] = Frame{}
		}
		th.fpMax = th.fp // reset high water mark
	}
	if th.sp > th.spMax {
		th.spMax = th.sp
	} else if th.sp < th.spMax {
		// clear the value stack to high water mark
		// and following non-nil (expression temporaries)
		for i := th.sp; i < th.spMax || (i < maxStack && th.stack[i] != nil); i++ {
			th.stack[i] = nil
		}
		th.spMax = th.sp
	}

	for {
		fr.catchJump = 0
		fr.catchSp = -1
		fr.catchPat = ""
		result := th.interp2()
		if result == nil {
			// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			th.fp = fp - 1
			if th.sp <= sp {
				return nil // implicit return from last statement had no value
			}
			return th.Top()
		}
		// try block threw
		th.sp = fr.catchSp
		th.fp = fp
		fr := &th.frames[th.fp-1]
		fr.ip = fr.catchJump
		th.Push(result) // SuExcept
		// loop and re-enter interp
	}
}

// interp is the main interpreter loop
// It normally returns nil, with the return value (if any) on the stack
// Returns *SuExcept if there was an exception/panic
func (th *Thread) interp2() (ret Value) {
	fr := &th.frames[th.fp-1]
	code := fr.fn.Code
	var fp int
	var profileBefore uint64
	if th.profile.enabled {
		fp = th.fp
		profileBefore = tsc.Read()
	}
	defer func() {
		if profileBefore > 0 {
			t := int64(tsc.Read() - profileBefore)
			th.profile.self[fr.fn] += t
			th.profile.total[fr.fn] += t
			if i := fp - 2; i >= 0 {
				fn := th.frames[i].fn
				th.profile.self[fn] -= t
			}
		}
		// this is an optimization to avoid unnecessary recover/repanic
		if fr.catchJump == 0 && th.blockReturnFrame == nil {
			return // this frame isn't catching
		}
		e := recover()
		if e == nil {
			return // not panic'ing, normal return
		}
		if e == BlockReturn {
			if th.blockReturnFrame != fr {
				panic(e) // not our block, rethrow
			}
			return // normal return
		}
		if fr.catchJump == 0 {
			panic(e) // not catching
		}
		// return value (ret) tells run we're catching
		ret = OpCatch(th, e, fr.catchPat)
	}()

	for fr.ip < len(code) {
		oc := code[fr.ip]
		fr.ip++
		interpFuncs[oc](th)
		
		if wingui { // const so should be compiled away
			if th == MainThread {
				opCount--
				if opCount <= 0 {
					opCount = opInterval // reset counter
					if Interrupt() {
						panic("interrupt")
					}
				}
			}
		}
	}
	return nil
}

func interpPop(th *Thread) {
	th.Pop()
}

func interpThis(th *Thread) {
	fr := &th.frames[th.fp-1]
	if fr.this == nil {
		panic("uninitialized: this")
	}
	th.Push(fr.this)
}

func interpTrue(th *Thread) {
	th.Push(True)
}

func interpFalse(th *Thread) {
	th.Push(False)
}

func interpZero(th *Thread) {
	th.Push(Zero)
}

func interpOne(th *Thread) {
	th.Push(One)
}

func interpMinusOne(th *Thread) {
	th.Push(MinusOne)
}

func interpMaxInt(th *Thread) {
	th.Push(MaxInt)
}

func interpEmptyStr(th *Thread) {
	th.Push(EmptyStr)
}

func interpInt(th *Thread) {
	fr := &th.frames[th.fp-1]
	th.Push(SuInt(fr.fetchInt16()))
}

func interpValue(th *Thread) {
	fr := &th.frames[th.fp-1]
	th.Push(fr.fn.Values[fr.fetchUint8()])
}

func interpLoad(th *Thread) {
	fr := &th.frames[th.fp-1]
	i := fr.fetchUint8()
	val := fr.locals.v[i]
	if val == nil {
		panic("uninitialized variable: " + fr.fn.Names[i])
	}
	th.Push(val)
}

func interpStore(th *Thread) {
	fr := &th.frames[th.fp-1]
	fr.locals.v[fr.fetchUint8()] = th.Top()
}

func interpLoadStore(th *Thread) {
	fr := &th.frames[th.fp-1]
	i := fr.fetchUint8()
	op := fr.fetchUint8()
	x := fr.locals.v[i]
	if x == nil {
		panic("uninitialized variable: " + fr.fn.Names[i])
	}
	y := th.stack[th.sp-1]
	var result Value
	switch op >> 1 {
	case 0:
		result = OpAdd(x, y)
	case 1:
		result = OpSub(x, y)
	case 2:
		result = OpCat(th, x, y)
	case 3:
		result = OpMul(x, y)
	case 4:
		result = OpDiv(x, y)
	case 5:
		result = OpMod(x, y)
	case 6:
		result = OpLeftShift(x, y)
	case 7:
		result = OpRightShift(x, y)
	case 8:
		result = OpBitOr(x, y)
	case 9:
		result = OpBitAnd(x, y)
	case 10:
		result = OpBitXor(x, y)
	}
	fr.locals.v[i] = result
	if op&1 == 0 {
		th.stack[th.sp-1] = result
	} else { // retOrig
		th.stack[th.sp-1] = x
	}
}

func interpDyload(th *Thread) {
	fr := &th.frames[th.fp-1]
	i := fr.fetchUint8()
	val := fr.locals.v[i]
	if val == nil {
		val = th.dyload(fr, i)
	}
	th.Push(val)
}

func interpGlobal(th *Thread) {
	fr := &th.frames[th.fp-1]
	gn := fr.fetchUint16()
	th.Push(Global.Get(th, gn))
}

func interpGet(th *Thread) {
	m := th.Pop()
	ob := th.Pop()
	val := ob.Get(th, m)
	if val == nil {
		if ss, ok := m.(SuStr); ok {
			val = ob.Lookup(th, string(ss))
			if val != nil {
				val = NewSuMethod(ob, val)
			}
		}
		if val == nil {
			MemberNotFound(m)
		}
	}
	th.Push(val)
}

func interpPut(th *Thread) {
	val := th.Pop()
	m := th.Pop()
	ob := th.Pop()
	ob.Put(th, m, val)
	th.Push(val)
}

func interpGetPut(th *Thread) {
	fr := &th.frames[th.fp-1]
	i := fr.fetchUint8()
	op := []func(x, y Value) Value{
		OpAdd, OpSub, th.Cat, OpMul, OpDiv, OpMod,
		OpLeftShift, OpRightShift, OpBitOr, OpBitAnd, OpBitXor}[i>>1]
	val := th.Pop()
	m := th.Pop()
	ob := th.Pop()
	th.Push(ob.GetPut(th, m, val, op, i&1 != 0))
}

func interpRangeTo(th *Thread) {
	j := ToInt(th.Pop())
	i := ToIndex(th.Pop())
	val := th.Pop()
	th.Push(val.RangeTo(i, j))
}

func interpRangeLen(th *Thread) {
	n := ToInt(th.Pop())
	i := ToIndex(th.Pop())
	val := th.Pop()
	th.Push(val.RangeLen(i, n))
}

func interpIs(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpIs(th.stack[th.sp-1], th.stack[th.sp])
}

func interpIsnt(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpIsnt(th.stack[th.sp-1], th.stack[th.sp])
}

func interpMatch(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpMatch(th, th.stack[th.sp-1], th.stack[th.sp])
}

func interpMatchNot(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpMatch(th, th.stack[th.sp-1], th.stack[th.sp]).Not()
}

func interpLt(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpLt(th.stack[th.sp-1], th.stack[th.sp])
}

func interpLte(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpLte(th.stack[th.sp-1], th.stack[th.sp])
}

func interpGt(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpGt(th.stack[th.sp-1], th.stack[th.sp])
}

func interpGte(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpGte(th.stack[th.sp-1], th.stack[th.sp])
}

func interpAdd(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpAdd(th.stack[th.sp-1], th.stack[th.sp])
}

func interpSub(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpSub(th.stack[th.sp-1], th.stack[th.sp])
}

func interpCat(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpCat(th, th.stack[th.sp-1], th.stack[th.sp])
}

func interpMul(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpMul(th.stack[th.sp-1], th.stack[th.sp])
}

func interpDiv(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpDiv(th.stack[th.sp-1], th.stack[th.sp])
}

func interpMod(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpMod(th.stack[th.sp-1], th.stack[th.sp])
}

func interpLeftShift(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpLeftShift(th.stack[th.sp-1], th.stack[th.sp])
}

func interpRightShift(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpRightShift(th.stack[th.sp-1], th.stack[th.sp])
}

func interpBitOr(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpBitOr(th.stack[th.sp-1], th.stack[th.sp])
}

func interpBitAnd(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpBitAnd(th.stack[th.sp-1], th.stack[th.sp])
}

func interpBitXor(th *Thread) {
	th.sp--
	th.stack[th.sp-1] = OpBitXor(th.stack[th.sp-1], th.stack[th.sp])
}

func interpBitNot(th *Thread) {
	th.stack[th.sp-1] = OpBitNot(th.stack[th.sp-1])
}

func interpNot(th *Thread) {
	th.stack[th.sp-1] = OpNot(th.stack[th.sp-1])
}

func interpUnaryPlus(th *Thread) {
	th.stack[th.sp-1] = OpUnaryPlus(th.stack[th.sp-1])
}

func interpUnaryMinus(th *Thread) {
	th.stack[th.sp-1] = OpUnaryMinus(th.stack[th.sp-1])
}
func interpInRange(th *Thread) {
	fr := &th.frames[th.fp-1]
	orgTok := tok.Token(fr.fetchUint8())
	org := fr.fn.Values[fr.fetchUint8()]
	endTok := tok.Token(fr.fetchUint8())
	end := fr.fn.Values[fr.fetchUint8()]
	th.stack[th.sp-1] = OpInRange(th.stack[th.sp-1], orgTok, org, endTok, end)
}
func interpBool(th *Thread) {
	th.topbool()
}

func interpJump(th *Thread) {
	fr := &th.frames[th.fp-1]
	fr.jump()
}

func interpJumpTrue(th *Thread) {
	fr := &th.frames[th.fp-1]
	if th.popbool() {
		fr.jump()
	} else {
		fr.ip += 2
	}
}

func interpJumpFalse(th *Thread) {
	fr := &th.frames[th.fp-1]
	if !th.popbool() {
		fr.jump()
	} else {
		fr.ip += 2
	}
}

func interpAnd(th *Thread) {
	fr := &th.frames[th.fp-1]
	if !th.topbool() {
		fr.jump()
	} else {
		th.Pop()
		fr.ip += 2
	}
}

func interpOr(th *Thread) {
	fr := &th.frames[th.fp-1]
	if th.topbool() {
		fr.jump()
	} else {
		th.Pop()
		fr.ip += 2
	}
}

func interpQMark(th *Thread) {
	fr := &th.frames[th.fp-1]
	if !th.popbool() {
		fr.jump()
	} else {
		fr.ip += 2
	}
}

func interpIn(th *Thread) {
	fr := &th.frames[th.fp-1]
	y := th.Pop()
	x := th.Pop()
	if x.Equal(y) {
		th.Push(True)
		fr.jump()
	} else {
		fr.ip += 2
		th.Push(x)
	}
}

func interpJumpIs(th *Thread) {
	fr := &th.frames[th.fp-1]
	y := th.Pop()
	x := th.Pop()
	if x.Equal(y) {
		fr.jump()
	} else {
		fr.ip += 2
		th.Push(x)
	}
}

func interpJumpIsnt(th *Thread) {
	fr := &th.frames[th.fp-1]
	y := th.Pop()
	x := th.Pop()
	if !x.Equal(y) {
		th.Push(x)
		fr.jump()
	} else {
		fr.ip += 2
	}
}

func interpJumpLt(th *Thread) {
	fr := &th.frames[th.fp-1]
	if strictCompare(th.stack[th.sp-1], th.stack[th.sp-2]) < 0 {
		fr.jump()
	} else {
		fr.ip += 2
	}
}

func interpIter(th *Thread) {
	th.stack[th.sp-1] = OpIter(th.stack[th.sp-1])
}

func interpForIn(th *Thread) {
	fr := &th.frames[th.fp-1]
	next := th.Top().(interface{ Next() Value }).Next()
	if next != nil {
		fr.locals.v[fr.fetchUint8()] = next
		fr.jump()
	} else {
		fr.ip += 3
	}
}

func interpIter2(th *Thread) {
	th.stack[th.sp-1] = OpIter2(th.stack[th.sp-1])
}

func interpForIn2(th *Thread) {
	fr := &th.frames[th.fp-1]
	m, v := th.Top().(SuIter2).iter2()
	if m != nil {
		fr.locals.v[fr.fetchUint8()] = m
		fr.locals.v[fr.fetchUint8()] = v
		fr.jump()
	} else {
		fr.ip += 4
	}
}

func interpForRange(th *Thread) {
	fr := &th.frames[th.fp-1]
	i := ToInt(th.Top())
	if i < ToInt(th.stack[th.sp-2]) {
		th.stack[th.sp-1] = SuInt(i + 1)
		fr.jump()
	} else {
		fr.ip += 2
	}
}

func interpForRangeVar(th *Thread) {
	fr := &th.frames[th.fp-1]
	i := ToInt(fr.locals.v[fr.fetchUint8()])
	if i < ToInt(th.Top()) {
		fr.locals.v[fr.ip-1] = SuInt(i + 1)
		fr.jump()
	} else {
		fr.ip += 2
	}
}

func interpReturnNil(th *Thread) {
	th.Push(nil)
	interpReturn(th)
}

func interpReturn(th *Thread) {
	th.frames[th.fp-1].ip = math.MaxInt
}

func interpReturnThrow(th *Thread) {
	th.ReturnThrow = true
	interpReturn(th)
}

func interpReturnMulti(th *Thread) {
	fr := &th.frames[th.fp-1]
	n := fr.fetchUint8()
	th.ReturnMulti = make([]Value, n)
	copy(th.ReturnMulti, th.stack[th.sp-n:th.sp])
	th.sp -= n - 1
}

func interpPushReturn(th *Thread) {
	fr := &th.frames[th.fp-1]
	n := fr.fetchUint8()
	rm := th.ReturnMulti
	th.ReturnMulti = th.ReturnMulti[:0]
	if n != len(rm) {
		panic("multiple return/assign mismatch")
	}
	for i := range n {
		th.Push(rm[i])
	}
	clear(th.ReturnMulti[:cap(th.ReturnMulti)])
}

func interpTry(th *Thread) {
	fr := &th.frames[th.fp-1]
	fr.catchJump = fr.ip + fr.fetchInt16()
	fr.catchSp = th.sp
	fr.catchPat = string(fr.fn.Values[fr.fetchUint8()].(SuStr))
}

func interpCatch(th *Thread) {
	fr := &th.frames[th.fp-1]
	fr.ip += fr.fetchInt16()
	fr.catchJump = 0 // no longer catching
}

func interpThrow(th *Thread) {
	panic(th.Pop())
}

func interpClosure(th *Thread) {
	fr := &th.frames[th.fp-1]
	fr.locals.moveToHeap()
	fn := fr.fn.Values[fr.fetchUint8()].(*SuFunc)
	parent := fr
	if fr.blockParent != nil {
		parent = fr.blockParent
	}
	block := &SuClosure{SuFunc: fn, locals: fr.locals.v, this: fr.this,
		parent: parent}
	th.Push(block)
}

func interpBlockBreak(th *Thread) {
	panic("block:break")
}

func interpBlockContinue(th *Thread) {
	panic("block:continue")
}

func interpBlockReturnNil(th *Thread) {
	th.Push(nil)
	panic("block return")
}

func interpBlockReturn(th *Thread) {
	panic("block return")
}

func interpCallFuncDiscard(th *Thread) {
	interpCallFunc(th, op.CallFuncDiscard)
}

func interpCallFuncNoNil(th *Thread) {
	interpCallFunc(th, op.CallFuncNoNil)
}

func interpCallFuncNilOk(th *Thread) {
	interpCallFunc(th, op.CallFuncNilOk)
}

func interpCallFunc(th *Thread, oc op.Opcode) {
	fr := &th.frames[th.fp-1]
	f := th.Pop()
	ai := fr.fetchUint8()
	var argSpec *ArgSpec
	if ai < len(StdArgSpecs) {
		argSpec = &StdArgSpecs[ai]
	} else {
		argSpec = &fr.fn.ArgSpecs[ai-len(StdArgSpecs)]
	}
	// fmt.Println(strings.Repeat("    ", t.fp+1), f)
	base := th.sp - int(argSpec.Nargs)
	result := f.Call(th, nil, argSpec)
	th.sp = base
	if th.ReturnThrow {
		// NOTE: this should be kept in sync with CallMeth & Finally
		th.ReturnThrow = false // default is to clear the flag
		if oc == op.CallFuncDiscard {
			if result != EmptyStr && result != True {
				if s, ok := result.ToStr(); ok {
					panic(s)
				}
				panic("return value not checked")
			}
		} else if fr.ip >= len(fr.fn.Code) ||
			op.Opcode(fr.fn.Code[fr.ip]) == op.Return {
			th.ReturnThrow = true // propagate if returning result
		}
	}
	pushResult(th, result, oc)
}

func pushResult(th *Thread, result Value, oc op.Opcode) {
	switch oc {
	case op.CallFuncNoNil, op.CallMethNoNil:
		if result == nil {
			panic("no return value")
		}
		fallthrough
	case op.CallFuncNilOk, op.CallMethNilOk:
		th.Push(result)
	default:
		// discard result
	}
}

func interpSuper(th *Thread) {
	fr := &th.frames[th.fp-1]
	th.super = fr.fetchUint16()
}

func interpCallMethDiscard(th *Thread) {
	interpCallMeth(th, op.CallMethDiscard)
}

func interpCallMethNoNil(th *Thread) {
	interpCallMeth(th, op.CallMethNoNil)
}

func interpCallMethNilOk(th *Thread) {
	interpCallMeth(th, op.CallMethNilOk)
}

func interpCallMeth(th *Thread, oc op.Opcode) {
	fr := &th.frames[th.fp-1]
	method := th.Pop()
	ai := fr.fetchUint8()
	var argSpec *ArgSpec
	if ai < len(StdArgSpecs) {
		argSpec = &StdArgSpecs[ai]
	} else {
		argSpec = &fr.fn.ArgSpecs[ai-len(StdArgSpecs)]
	}
	base := th.sp - int(argSpec.Nargs) - 1
	this := th.stack[base]
	if methstr, ok := method.ToStr(); ok {
		var f Value
		ob := this
		if th.super > 0 {
			if instance, ok := this.(*SuInstance); ok {
				for i, p := range instance.parents {
					if p.Base == th.super {
						c := instance.parents[i+1]
						f = c.lookup(th, methstr, instance.parents[i+1:])
						th.super = 0
						goto done
					}
				}
			}
			// else
			ob = Global.Get(th, th.super)
			th.super = 0
		}
		f = ob.Lookup(th, methstr)
	done:
		if f != nil {
			// fmt.Println(strings.Repeat("   ", t.fp+1), f)
			result := f.Call(th, this, argSpec)
			th.sp = base
			if th.ReturnThrow {
				// NOTE: this code should be kept in sync with CallFunc
				th.ReturnThrow = false // default is to clear the flag
				if oc == op.CallMethDiscard {
					if result != EmptyStr && result != True {
						if s, ok := result.ToStr(); ok {
							panic(s)
						}
						panic("return value not checked")
					}
				} else if fr.ip >= len(fr.fn.Code) ||
					op.Opcode(fr.fn.Code[fr.ip]) == op.Return {
					th.ReturnThrow = true // propagate if returning result
				}
			}
			pushResult(th, result, oc)
			return
		}
	}
	panic("method not found: " + ErrType(this) + "." + ToStrOrString(method))
}

func interpCover(th *Thread) {
	fr := &th.frames[th.fp-1]
	fn := fr.fn
	if fn.cover != nil {
		ip := fr.ip - 1
		if len(fn.cover) < len(fn.Code) {
			fn.cover[ip>>4] |= 1 << (ip & 15)
		} else { // count
			p := &fn.cover[ip]
			if x := *p + 1; x != 0 {
				*p = x
			}
		}
	}
}

type interpFunc func(*Thread)

var interpFuncs = []interpFunc{
	op.Pop:             interpPop,
	op.This:            interpThis,
	op.True:            interpTrue,
	op.False:           interpFalse,
	op.Zero:            interpZero,
	op.One:             interpOne,
	op.MinusOne:        interpMinusOne,
	op.MaxInt:          interpMaxInt,
	op.EmptyStr:        interpEmptyStr,
	op.Int:             interpInt,
	op.Value:           interpValue,
	op.Load:            interpLoad,
	op.Store:           interpStore,
	op.LoadStore:       interpLoadStore,
	op.Dyload:          interpDyload,
	op.Global:          interpGlobal,
	op.Get:             interpGet,
	op.Put:             interpPut,
	op.GetPut:          interpGetPut,
	op.RangeTo:         interpRangeTo,
	op.RangeLen:        interpRangeLen,
	op.Is:              interpIs,
	op.Isnt:            interpIsnt,
	op.Match:           interpMatch,
	op.MatchNot:        interpMatchNot,
	op.Lt:              interpLt,
	op.Lte:             interpLte,
	op.Gt:              interpGt,
	op.Gte:             interpGte,
	op.Add:             interpAdd,
	op.Sub:             interpSub,
	op.Cat:             interpCat,
	op.Mul:             interpMul,
	op.Div:             interpDiv,
	op.Mod:             interpMod,
	op.LeftShift:       interpLeftShift,
	op.RightShift:      interpRightShift,
	op.BitOr:           interpBitOr,
	op.BitAnd:          interpBitAnd,
	op.BitXor:          interpBitXor,
	op.BitNot:          interpBitNot,
	op.Not:             interpNot,
	op.UnaryPlus:       interpUnaryPlus,
	op.UnaryMinus:      interpUnaryMinus,
	op.InRange:         interpInRange,
	op.Bool:            interpBool,
	op.Jump:            interpJump,
	op.JumpTrue:        interpJumpTrue,
	op.JumpFalse:       interpJumpFalse,
	op.And:             interpAnd,
	op.Or:              interpOr,
	op.QMark:           interpQMark,
	op.In:              interpIn,
	op.JumpIs:          interpJumpIs,
	op.JumpIsnt:        interpJumpIsnt,
	op.JumpLt:          interpJumpLt,
	op.Iter:            interpIter,
	op.ForIn:           interpForIn,
	op.Iter2:           interpIter2,
	op.ForIn2:          interpForIn2,
	op.ForRange:        interpForRange,
	op.ForRangeVar:     interpForRangeVar,
	op.ReturnNil:       interpReturnNil,
	op.Return:          interpReturn,
	op.ReturnThrow:     interpReturnThrow,
	op.ReturnMulti:     interpReturnMulti,
	op.PushReturn:      interpPushReturn,
	op.Try:             interpTry,
	op.Catch:           interpCatch,
	op.Throw:           interpThrow,
	op.Closure:         interpClosure,
	op.BlockBreak:      interpBlockBreak,
	op.BlockContinue:   interpBlockContinue,
	op.BlockReturnNil:  interpBlockReturnNil,
	op.BlockReturn:     interpBlockReturn,
	op.CallFuncDiscard: interpCallFuncDiscard,
	op.CallFuncNoNil:   interpCallFuncNoNil,
	op.CallFuncNilOk:   interpCallFuncNilOk,
	op.Super:           interpSuper,
	op.CallMethDiscard: interpCallMethDiscard,
	op.CallMethNoNil:   interpCallMethNoNil,
	op.CallMethNilOk:   interpCallMethNilOk,
	op.Cover:           interpCover,
}
