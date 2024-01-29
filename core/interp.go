// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	op "github.com/apmckinlay/gsuneido/core/opcodes"
)

// RunOnGoSide is injected
var RunOnGoSide = func() {}
var Interrupt = func() bool { return false }

var BlockBreak = BuiltinSuExcept("block:break")
var BlockContinue = BuiltinSuExcept("block:continue")
var BlockReturn = BuiltinSuExcept("block return")

// invoke sets up a frame to Run a compiled Suneido function
// The stack must already be in the form required by the function (massaged)
// WARNING: invoke does not pop the stack, the caller is responsible for that
func (th *Thread) invoke(fn *SuFunc, this Value) Value {
	// reserve stack space for locals
	for expand := fn.Nlocals - fn.Nparams; expand > 0; expand-- {
		th.Push(nil)
	}
	return th.run(Frame{fn: fn, this: this,
		locals: locals{v: th.stack[th.sp-int(fn.Nlocals) : th.sp]}})
}

// run is needed in addition to interp
// because we can only recover panic on the way out of a function
// so if the exception is caught we have to re-enter interp
// Called by Thread.Invoke and SuClosure.Call
func (th *Thread) run(frame Frame) Value {
	if th.fp >= len(th.frames) {
		panic("function call overflow")
	}
	if th.profile.enabled {
		th.profile.lock.Lock()
		th.profile.calls[frame.fn.Name]++
	}
	th.frames[th.fp] = frame
	th.fp++
	if th.profile.enabled {
		th.profile.lock.Unlock()
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(strings.Repeat("    ", t.fp) + "run:", t.frames[t.fp].fn)
	sp := th.sp
	fp := th.fp
	if th.fp > th.fpMax {
		th.fpMax = th.fp // track high water mark
	} else if th.fpMax > th.fp {
		// clear the frames up to the high water mark
		for i := th.fp; i < th.fpMax; i++ {
			th.frames[i] = Frame{}
		}
		th.fpMax = th.fp // reset high water mark
	}
	if th.sp > th.spMax {
		th.spMax = th.sp
	} else if th.spMax > th.sp {
		// clear the value stack to high water mark
		// and following non-nil (expression temporaries)
		for i := th.sp; i < th.spMax || (i < maxStack && th.stack[i] != nil); i++ {
			th.stack[i] = nil
		}
		th.spMax = th.sp
	}

	catchJump := 0
	catchSp := -1
	for {
		result := th.interp(&catchJump, &catchSp)
		if result == nil {
			// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			th.fp = fp - 1
			if th.sp <= sp {
				return nil // implicit return from last statement had no value
			}
			return th.Top()
		}
		// try block threw
		th.sp = catchSp
		th.fp = fp
		fr := &th.frames[th.fp-1]
		fr.ip = catchJump
		catchJump = 0 // no longer catching
		catchSp = -1
		th.Push(result) // SuExcept
		// loop and re-enter interp
	}
}

// interp is the main interpreter loop
// It normally returns nil, with the return value (if any) on the stack
// Returns *SuExcept if there was an exception/panic
func (th *Thread) interp(catchJump, catchSp *int) (ret Value) {
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

	profileOpCount := 0
	defer func() {
		if th.profile.enabled && profileOpCount > 0 {
			th.profile.ops[fr.fn.Name] += int32(profileOpCount)
		}
		// this is an optimization to avoid unnecessary recover/repanic
		if *catchJump == 0 && th.blockReturnFrame == nil {
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
		if *catchJump == 0 {
			panic(e) // not catching
		}
		// return value (ret) tells run we're catching
		ret = OpCatch(th, e, catchPat)
	}()

loop:
	for fr.ip < len(code) {
		profileOpCount++
		// fmt.Println("stack:", t.sp, t.stack[max(0, t.sp-3):t.sp])
		// _, da := Disasm1(fr.fn, fr.ip)
		// fmt.Printf("%d: %d: %s\n", t.fp, fr.ip, da)
		if th.UIThread {
			if th.OpCount == 0 {
				RunOnGoSide()
				if Interrupt() {
					panic("interrupt")
				}
				th.OpCount = 1009 // otherwise it won't trigger again
			}
			th.OpCount--
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
		case op.Global:
			gn := fetchUint16()
			th.Push(Global.Get(th, gn))
		case op.Get:
			m := th.Pop()
			ob := th.Pop()
			val := ob.Get(th, m)
			if val == nil {
				panic("uninitialized member: " + m.String())
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
			if OpLt(th.stack[th.sp-1], th.stack[th.sp-2]) == True {
				jump()
			} else {
				fr.ip += 2
			}
		case op.Iter:
			th.stack[th.sp-1] = OpIter(th.stack[th.sp-1])
		case op.ForIn:
			brk := fetchInt16()
			local := fetchUint8()
			iter := th.Top()
			nextable := iter.(interface{ Next() Value })
			next := nextable.Next()
			if next != nil {
				fr.locals.v[local] = next
			} else {
				fr.ip += brk - 1 // break
			}
		case op.ReturnNil:
			th.Push(nil)
			break loop
		case op.Return:
			break loop
		case op.ReturnThrow:
			th.ReturnThrow = true
			break loop
		case op.Try:
			*catchJump = fr.ip + fetchInt16()
			*catchSp = th.sp
			catchPat = string(fr.fn.Values[fetchUint8()].(SuStr))
		case op.Catch:
			fr.ip += fetchInt16()
			*catchJump = 0 // no longer catching
		case op.Throw:
			panic(th.Pop())
		case op.Closure:
			fr.locals.moveToHeap()
			fn := fr.fn.Values[fetchUint8()].(*SuFunc)
			parent := fr
			if fr.blockParent != nil {
				parent = fr.blockParent
			}
			block := &SuClosure{SuFunc: *fn, locals: fr.locals.v, this: fr.this,
				parent: parent}
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
				// NOTE: this code should be kept in sync with CallMeth
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
				var f Callable
				ob := this
				if super > 0 {
					if instance, ok := this.(*SuInstance); ok {
						for i, p := range instance.parents {
							if p.Base == super {
								c := instance.parents[i+1]
								f = c.lookup(th, methstr, instance.parents[i+1:])
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
		default:
			Fatal("invalid op code: " + oc.String())
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
