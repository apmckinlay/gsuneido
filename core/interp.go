// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	op "github.com/apmckinlay/gsuneido/core/opcodes"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

var Interrupt func() bool // injected

const opInterval = 4001 // ???
var opCount int = opInterval

var BlockBreak = BuiltinSuExcept("block:break")
var BlockContinue = BuiltinSuExcept("block:continue")
var BlockReturn = BuiltinSuExcept("block return")

// invoke sets up a frame to run a compiled Suneido function
// The stack must already be in the form required by the function (massaged)
// WARNING: invoke does not pop the stack, the caller is responsible for that
func (th *Thread) invoke(fn *SuFunc, this Value) Value {
	// reserve stack space for locals
	for expand := fn.Nlocals - fn.Nparams; expand > 0; expand-- {
		th.Push(nil)
	}
	if th.fp >= len(th.frames) {
		panic("function call overflow")
	}
	fr := &th.frames[th.fp]
	fr.fn = fn
	fr.this = this
	fr.blockParent = nil
	fr.locals = locals{v: th.stack[th.sp-int(fn.Nlocals) : th.sp]}
	fr.ip = 0
	return th.run()
}

// run is needed in addition to interp
// because we can only recover panic on the way out of a function
// so if the exception is caught we have to re-enter interp
// Called by Thread.invoke (above) and SuClosure.Call
func (th *Thread) run() Value {
	fr := &th.frames[th.fp]
	fr.ip = 0
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
		result := th.interp()
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
		fr.ip = fr.catchJump
		th.Push(result) // SuExcept
		// loop and re-enter interp
	}
}

// interp is the main interpreter loop
// It normally returns nil, with the return value (if any) on the stack
// Returns *SuExcept if there was an exception/panic
func (th *Thread) interp() (ret Value) {
	fr := &th.frames[th.fp-1]
	code := fr.fn.Code
	super := 0
	catchPat := ""
	var oc op.Opcode

	fetchUint8 := func() int {
		fr.ip++
		return int(code[fr.ip-1])
	}
	fetchInt16 := func() int {
		fr.ip += 2
		return int(int16(uint16(code[fr.ip-2])<<8 + uint16(code[fr.ip-1])))
	}
	fetchUint16 := func() int {
		fr.ip += 2
		return int(uint16(code[fr.ip-2])<<8 + uint16(code[fr.ip-1]))
	}
	jump := func() {
		fr.ip += fetchInt16()
	}
	pushResult := func(result Value) {
		switch oc {
		case op.CallFuncNoNil, op.CallMethNoNil, op.ValueCallMethNoNil, op.GlobalCallFuncNoNil:
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
		ret = OpCatch(th, e, catchPat)
	}()

loop:
	for fr.ip < len(code) {
		// fmt.Println("stack:", th.sp, th.stack[:th.sp])
		// fmt.Println(Disasm1(fr.fn, fr.ip))
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
		oc = op.Opcode(code[fr.ip])
		fr.ip++
		switch oc {
		case op.Pop:
			th.Pop()
		case op.This:
			if fr.this == nil {
				panic("uninitialized: this")
			}
			th.Push(fr.this)
		case op.True:
			th.Push(True)
		case op.False:
			th.Push(False)
		case op.Zero:
			th.Push(Zero)
		case op.One:
			th.Push(One)
		case op.MinusOne:
			th.Push(MinusOne)
		case op.MaxInt:
			th.Push(MaxInt)
		case op.EmptyStr:
			th.Push(EmptyStr)
		case op.Int:
			th.Push(SuInt(fetchInt16()))
		case op.Value:
			th.Push(fr.fn.Values[fetchUint8()])
		case op.LoadLoad:
			i := fetchUint8()
			val := fr.locals.v[i]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[i])
			}
			th.Push(val)
			fallthrough
		case op.Load:
			i := fetchUint8()
			val := fr.locals.v[i]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[i])
			}
			th.Push(val)
		case op.Store:
			fr.locals.v[fetchUint8()] = th.Top()
		case op.LoadStore:
			i := fetchUint8()
			op := fetchUint8()
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
		case op.Dyload:
			i := fetchUint8()
			val := fr.locals.v[i]
			if val == nil {
				val = th.dyload(fr, i)
			}
			th.Push(val)
		case op.LoadValue:
			// Load
			localIdx := fetchUint8()
			val := fr.locals.v[localIdx]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[localIdx])
			}
			th.Push(val)
			// Value
			valueIdx := fetchUint8()
			th.Push(fr.fn.Values[valueIdx])
		case op.ThisValue:
			// This
			if fr.this == nil {
				panic("uninitialized: this")
			}
			th.Push(fr.this)
			// Value
			valueIdx := fetchUint8()
			th.Push(fr.fn.Values[valueIdx])
		case op.StorePop:
			i := fetchUint8()
			fr.locals.v[i] = th.Pop()
		case op.ThisLoad:
			// This
			if fr.this == nil {
				panic("uninitialized: this")
			}
			th.Push(fr.this)
			// Load
			i := fetchUint8()
			val := fr.locals.v[i]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[i])
			}
			th.Push(val)
		case op.GetValue:
			// Get
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
			// Value
			valueIdx := fetchUint8()
			th.Push(fr.fn.Values[valueIdx])
		case op.PopLoad:
			th.Pop()
			// Load
			i := fetchUint8()
			val := fr.locals.v[i]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[i])
			}
			th.Push(val)
		case op.Global:
			gn := fetchUint16()
			th.Push(Global.Get(th, gn))
		case op.ValueGet:
			valueIdx := fetchUint8()
			th.Push(fr.fn.Values[valueIdx])
			fallthrough
		case op.Get:
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
		case op.Put:
			val := th.Pop()
			m := th.Pop()
			ob := th.Pop()
			ob.Put(th, m, val)
			th.Push(val)
		case op.GetPut:
			i := fetchUint8()
			op := []func(x, y Value) Value{
				OpAdd, OpSub, th.Cat, OpMul, OpDiv, OpMod,
				OpLeftShift, OpRightShift, OpBitOr, OpBitAnd, OpBitXor}[i>>1]
			val := th.Pop()
			m := th.Pop()
			ob := th.Pop()
			th.Push(ob.GetPut(th, m, val, op, i&1 != 0))
		case op.RangeTo:
			j := ToInt(th.Pop())
			i := ToIndex(th.Pop())
			val := th.Pop()
			th.Push(val.RangeTo(i, j))
		case op.RangeLen:
			n := ToInt(th.Pop())
			i := ToIndex(th.Pop())
			val := th.Pop()
			th.Push(val.RangeLen(i, n))
		case op.Is:
			th.sp--
			th.stack[th.sp-1] = OpIs(th.stack[th.sp-1], th.stack[th.sp])
		case op.Isnt:
			th.sp--
			th.stack[th.sp-1] = OpIsnt(th.stack[th.sp-1], th.stack[th.sp])
		case op.Match:
			th.sp--
			th.stack[th.sp-1] = OpMatch(th, th.stack[th.sp-1], th.stack[th.sp])
		case op.MatchNot:
			th.sp--
			th.stack[th.sp-1] = OpMatch(th, th.stack[th.sp-1], th.stack[th.sp]).Not()
		case op.Lt:
			th.sp--
			th.stack[th.sp-1] = OpLt(th.stack[th.sp-1], th.stack[th.sp])
		case op.Lte:
			th.sp--
			th.stack[th.sp-1] = OpLte(th.stack[th.sp-1], th.stack[th.sp])
		case op.Gt:
			th.sp--
			th.stack[th.sp-1] = OpGt(th.stack[th.sp-1], th.stack[th.sp])
		case op.Gte:
			th.sp--
			th.stack[th.sp-1] = OpGte(th.stack[th.sp-1], th.stack[th.sp])
		case op.Add:
			th.sp--
			th.stack[th.sp-1] = OpAdd(th.stack[th.sp-1], th.stack[th.sp])
		case op.Sub:
			th.sp--
			th.stack[th.sp-1] = OpSub(th.stack[th.sp-1], th.stack[th.sp])
		case op.Cat:
			th.sp--
			th.stack[th.sp-1] = OpCat(th, th.stack[th.sp-1], th.stack[th.sp])
		case op.CatN:
			count := fetchUint8()
			result := OpCatN(th, count)
			th.Push(result)
		case op.Mul:
			th.sp--
			th.stack[th.sp-1] = OpMul(th.stack[th.sp-1], th.stack[th.sp])
		case op.Div:
			th.sp--
			th.stack[th.sp-1] = OpDiv(th.stack[th.sp-1], th.stack[th.sp])
		case op.Mod:
			th.sp--
			th.stack[th.sp-1] = OpMod(th.stack[th.sp-1], th.stack[th.sp])
		case op.InRange:
			orgTok := tok.Token(fetchUint8())
			org := fr.fn.Values[fetchUint8()]
			endTok := tok.Token(fetchUint8())
			end := fr.fn.Values[fetchUint8()]
			th.stack[th.sp-1] = OpInRange(th.stack[th.sp-1], orgTok, org, endTok, end)
		case op.LeftShift:
			th.sp--
			th.stack[th.sp-1] = OpLeftShift(th.stack[th.sp-1], th.stack[th.sp])
		case op.RightShift:
			th.sp--
			th.stack[th.sp-1] = OpRightShift(th.stack[th.sp-1], th.stack[th.sp])
		case op.BitOr:
			th.sp--
			th.stack[th.sp-1] = OpBitOr(th.stack[th.sp-1], th.stack[th.sp])
		case op.BitAnd:
			th.sp--
			th.stack[th.sp-1] = OpBitAnd(th.stack[th.sp-1], th.stack[th.sp])
		case op.BitXor:
			th.sp--
			th.stack[th.sp-1] = OpBitXor(th.stack[th.sp-1], th.stack[th.sp])
		case op.BitNot:
			th.stack[th.sp-1] = OpBitNot(th.stack[th.sp-1])
		case op.Not:
			th.stack[th.sp-1] = OpNot(th.stack[th.sp-1])
		case op.UnaryPlus:
			th.stack[th.sp-1] = OpUnaryPlus(th.stack[th.sp-1])
		case op.UnaryMinus:
			th.stack[th.sp-1] = OpUnaryMinus(th.stack[th.sp-1])
		case op.Bool:
			th.topbool()
		case op.Jump:
			jump()
		case op.JumpTrue:
			if th.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case op.JumpFalse:
			if !th.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case op.And:
			if !th.topbool() {
				jump()
			} else {
				fr.ip += 2
				th.Pop()
			}
		case op.Or:
			if th.topbool() {
				jump()
			} else {
				fr.ip += 2
				th.Pop()
			}
		case op.QMark:
			if !th.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case op.In:
			y := th.Pop()
			x := th.Pop()
			if x.Equal(y) {
				th.Push(True)
				jump()
			} else {
				fr.ip += 2
				th.Push(x)
			}
		case op.JumpIs:
			y := th.Pop()
			x := th.Pop()
			if x.Equal(y) {
				jump()
			} else {
				fr.ip += 2
				th.Push(x)
			}
		case op.JumpIsnt:
			y := th.Pop()
			x := th.Pop()
			if !x.Equal(y) {
				th.Push(x)
				jump()
			} else {
				fr.ip += 2
			}
		case op.JumpLt:
			if strictCompare(th.stack[th.sp-1], th.stack[th.sp-2]) < 0 {
				jump()
			} else {
				fr.ip += 2
			}
		case op.Iter:
			th.stack[th.sp-1] = OpIter(th.stack[th.sp-1])
		case op.ForIn:
			next := th.Top().(interface{ Next() Value }).Next()
			if next != nil {
				fr.locals.v[fetchUint8()] = next
				jump()
			} else {
				fr.ip += 3
			}
		case op.Iter2:
			th.stack[th.sp-1] = OpIter2(th.stack[th.sp-1])
		case op.ForIn2:
			m, v := th.Top().(SuIter2).iter2()
			if m != nil {
				fr.locals.v[fetchUint8()] = m
				fr.locals.v[fetchUint8()] = v
				jump()
			} else {
				fr.ip += 4
			}
		case op.ForRange:
			th.stack[th.sp-1] = OpAdd1(th.stack[th.sp-1])
			if strictCompare(th.stack[th.sp-1], th.stack[th.sp-2]) < 0 {
				jump()
			} else {
				fr.ip += 2
			}
		case op.ForRangeVar:
			th.stack[th.sp-1] = OpAdd1(th.stack[th.sp-1])
			fr.locals.v[fetchUint8()] = th.stack[th.sp-1]
			if strictCompare(th.stack[th.sp-1], th.stack[th.sp-2]) < 0 {
				jump()
			} else {
				fr.ip += 2
			}
		case op.ReturnNil:
			th.Push(nil)
			break loop
		case op.Return:
			break loop
		case op.ReturnThrow:
			th.ReturnThrow = true
			break loop
		case op.ReturnMulti:
			th.ReturnMulti = th.ReturnMulti[:0]
			n := fetchUint8()
			for range n {
				th.ReturnMulti = append(th.ReturnMulti, th.Pop())
			}
			break loop
		case op.PushReturn:
			th.Pop() // discard the normal nil return value
			n := fetchUint8()
			rm := th.ReturnMulti
			th.ReturnMulti = th.ReturnMulti[:0]
			if n != len(rm) {
				panic("multiple return/assign mismatch")
			}
			for i := range n {
				th.Push(rm[i])
			}
			clear(th.ReturnMulti[:cap(th.ReturnMulti)])
		case op.Try:
			fr.catchJump = fr.ip + fetchInt16()
			fr.catchSp = th.sp
			catchPat = string(fr.fn.Values[fetchUint8()].(SuStr))
		case op.Catch:
			fr.ip += fetchInt16()
			fr.catchJump = 0 // no longer catching
		case op.Throw:
			panic(th.Pop())
		case op.Closure:
			fr.locals.moveToHeap()
			fn := fr.fn.Values[fetchUint8()].(*SuFunc)
			parent := fr
			if fr.blockParent != nil {
				parent = fr.blockParent
			}
			block := &SuClosure{SuFunc: fn, locals: fr.locals.v, this: fr.this,
				parent: parent, thread: th}
			th.Push(block)
		case op.BlockBreak:
			panic(BlockBreak)
		case op.BlockContinue:
			panic(BlockContinue)
		case op.BlockReturnNil:
			th.Push(nil)
			fallthrough
		case op.BlockReturn:
			th.blockReturnFrame = fr.blockParent
			panic(BlockReturn)
		case op.GlobalCallFuncNoNil:
			gn := fetchUint16()
			th.Push(Global.Get(th, gn))
			fallthrough
		case op.CallFuncDiscard, op.CallFuncNoNil, op.CallFuncNilOk:
			f := th.Pop()
			ai := fetchUint8()
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
				} else if fr.ip >= len(code) ||
					op.Opcode(code[fr.ip]) == op.Return {
					th.ReturnThrow = true // propagate if returning result
				}
			}
			pushResult(result)
		case op.Super:
			super = fetchUint16()
		case op.ValueCallMethNoNil:
			th.Push(fr.fn.Values[fetchUint8()])
			fallthrough
		case op.CallMethDiscard, op.CallMethNoNil, op.CallMethNilOk:
			method := th.Pop()
			ai := fetchUint8()
			var argSpec *ArgSpec
			if ai < len(StdArgSpecs) {
				argSpec = &StdArgSpecs[ai]
			} else {
				argSpec = &fr.fn.ArgSpecs[ai-len(StdArgSpecs)]
			}
			base := th.sp - int(argSpec.Nargs) - 1
			this := th.stack[base]
			if methstr, ok := method.ToStr(); ok {
				if class, ok := this.(*SuClass); ok && class.Base > 0 {
					this = getParents(th, class)
				}
				var f Value
				ob := this
				if super > 0 {
					if p, ok := this.(interface{ Parents() []*SuClass }); ok {
						parents := p.Parents()
						for i, p := range parents {
							if p.Base == super {
								c := parents[i+1]
								f = c.lookup(th, methstr, parents[i+1:])
								super = 0
								goto done
							}
						}
					}
					// else
					ob = Global.Get(th, super)
					super = 0
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
						} else if fr.ip >= len(code) ||
							op.Opcode(code[fr.ip]) == op.Return {
							th.ReturnThrow = true // propagate if returning result
						}
					}
					pushResult(result)
					break
				}
			}
			panic("method not found: " + ErrType(this) + "." + ToStrOrString(method))
		case op.Cover:
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
		// avoid range check on switch at the cost of a larger jump table
		case 0, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255:
			Fatal("invalid op code:", int(oc))
		}
	}
	return nil
}

// topbool return the top of the stack as bool, panicing if not True or False
func (th *Thread) topbool() bool {
	switch th.Top() {
	case True:
		return true
	case False:
		return false
	default:
		panic("conditionals require true or false")
	}
}

// popbool pops the top of the stack and returns it as bool, panicing if not True or False
func (th *Thread) popbool() bool {
	switch th.Pop() {
	case True:
		return true
	case False:
		return false
	default:
		panic("conditionals require true or false")
	}
}

// dyload pushes a dynamic variable onto the stack
// It looks up the frame stack to find it, and copies it locally
func (th *Thread) dyload(fr *Frame, idx int) Value {
	name := fr.fn.Names[idx]
	for i := th.fp - 2; i >= 0; i-- {
		fr2 := &th.frames[i]
		for j, s := range fr2.fn.Names {
			if s == name {
				if x := fr2.locals.v[j]; x != nil {
					fr.locals.v[idx] = x
					return x
				}
			}
		}
	}
	panic("uninitialized variable: " + name)
}
