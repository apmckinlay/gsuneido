// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	op "github.com/apmckinlay/gsuneido/runtime/opcodes"
)

// RunOnGoSide is injected
var RunOnGoSide = func() {}
var Interrupt = func() bool { return false }

var BlockBreak = BuiltinSuExcept("block:break")
var BlockContinue = BuiltinSuExcept("block:continue")
var BlockReturn = BuiltinSuExcept("block return")

// Invoke sets up a frame to Run a compiled Suneido function
// The stack must already be in the form required by the function (massaged)
func (t *Thread) Invoke(fn *SuFunc, this Value) Value {
	// reserve stack space for locals
	for expand := fn.Nlocals - fn.Nparams; expand > 0; expand-- {
		t.Push(nil)
	}
	return t.run(Frame{fn: fn, this: this,
		locals: Locals{v: t.stack[t.sp-int(fn.Nlocals) : t.sp]}})
}

// run is needed in addition to interp
// because we can only recover panic on the way out of a function
// so if the exception is caught we have to re-enter interp
// Called by Thread.Invoke and SuClosure.Call
func (t *Thread) run(frame Frame) Value {
	if t.fp >= len(t.frames) {
		panic("function call overflow")
	}
	if t.profile.enabled {
		t.profile.lock.Lock()
		t.profile.calls[frame.fn.Name]++
	}
	t.frames[t.fp] = frame
	t.fp++
	if t.profile.enabled {
		t.profile.lock.Unlock()
	}
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(strings.Repeat("    ", t.fp) + "run:", t.frames[t.fp].fn)
	sp := t.sp
	fp := t.fp
	if t.fp > t.fpMax {
		t.fpMax = t.fp // track high water mark
	} else if t.fpMax > t.fp {
		// clear the frames up to the high water mark
		for i := t.fp; i < t.fpMax; i++ {
			t.frames[i] = Frame{}
		}
		t.fpMax = t.fp // reset high water mark
	}
	if t.sp > t.spMax {
		t.spMax = t.sp
	} else if t.spMax > t.sp {
		// clear the value stack to high water mark
		// and following non-nil (expression temporaries)
		for i := t.sp; i < t.spMax || (i < maxStack && t.stack[i] != nil); i++ {
			t.stack[i] = nil
		}
		t.spMax = t.sp
	}

	catchJump := 0
	catchSp := -1
	for {
		result := t.interp(&catchJump, &catchSp)
		if result == nil {
			// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			t.fp = fp - 1
			if t.sp <= sp {
				return nil // implicit return from last statement had no value
			}
			return t.Top()
		}
		// try block threw
		t.sp = catchSp
		t.fp = fp
		fr := &t.frames[t.fp-1]
		fr.ip = catchJump
		catchJump = 0 // no longer catching
		catchSp = -1
		t.Push(result) // SuExcept
		// loop and re-enter interp
	}
}

// interp is the main interpreter loop
// It normally returns nil, with the return value (if any) on the stack
// Returns *SuExcept if there was an exception/panic
func (t *Thread) interp(catchJump, catchSp *int) (ret Value) {
	fr := &t.frames[t.fp-1]
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
			t.Push(result)
		default:
			// discard result
		}
	}

	profileOpCount := 0
	defer func() {
		if t.profile.enabled && profileOpCount > 0 {
			t.profile.ops[fr.fn.Name] += int32(profileOpCount)
		}
		// this is an optimization to avoid unnecessary recover/repanic
		if *catchJump == 0 && t.blockReturnFrame == nil {
			return // this frame isn't catching
		}
		e := recover()
		if e == nil {
			return // not panic'ing, normal return
		}
		if e == BlockReturn {
			if t.blockReturnFrame != fr {
				panic(e) // not our block, rethrow
			}
			return // normal return
		}
		if *catchJump == 0 {
			panic(e) // not catching
		}
		// return value (ret) tells run we're catching
		ret = OpCatch(t, e, catchPat)
	}()

loop:
	for fr.ip < len(code) {
		profileOpCount++
		// fmt.Println("stack:", t.sp, t.stack[ints.Max(0, t.sp-3):t.sp])
		// _, da := Disasm1(fr.fn, fr.ip)
		// fmt.Printf("%d: %d: %s\n", t.fp, fr.ip, da)
		if t.UIThread {
			if t.OpCount == 0 {
				RunOnGoSide()
				if Interrupt() {
					panic("interrupt")
				}
				t.OpCount = 1009 // otherwise it won't trigger again
			}
			t.OpCount--
		}
		oc = op.Opcode(code[fr.ip])
		fr.ip++
		switch oc {
		case op.Pop:
			t.Pop()
		case op.Dup:
			t.Push(t.Top())
		case op.Swap:
			t.Swap()
		case op.This:
			if fr.this == nil {
				panic("uninitialized: this")
			}
			t.Push(fr.this)
		case op.True:
			t.Push(True)
		case op.False:
			t.Push(False)
		case op.Zero:
			t.Push(Zero)
		case op.One:
			t.Push(One)
		case op.MinusOne:
			t.Push(MinusOne)
		case op.MaxInt:
			t.Push(MaxInt)
		case op.EmptyStr:
			t.Push(EmptyStr)
		case op.Int:
			t.Push(SuInt(fetchInt16()))
		case op.Value:
			t.Push(fr.fn.Values[fetchUint8()])
		case op.Load:
			i := fetchUint8()
			val := fr.locals.v[i]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[i])
			}
			t.Push(val)
		case op.Store:
			fr.locals.v[fetchUint8()] = t.Top()
		case op.LoadStore:
			i := fetchUint8()
			op := fetchUint8()
			x := fr.locals.v[i]
			y := t.stack[t.sp-1]
			var result Value
			switch op >> 1 {
			case 0:
				result = OpAdd(x, y)
			case 1:
				result = OpSub(x, y)
			case 2:
				result = OpCat(t, x, y)
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
				t.stack[t.sp-1] = result
			} else { // retOrig
				t.stack[t.sp-1] = x
			}
		case op.Dyload:
			i := fetchUint8()
			val := fr.locals.v[i]
			if val == nil {
				val = t.dyload(fr, i)
			}
			t.Push(val)
		case op.Global:
			gn := fetchUint16()
			t.Push(Global.Get(t, gn))
		case op.Get:
			m := t.Pop()
			ob := t.Pop()
			val := ob.Get(t, m)
			if val == nil {
				panic("uninitialized member: " + m.String())
			}
			t.Push(val)
		case op.Put:
			val := t.Pop()
			m := t.Pop()
			ob := t.Pop()
			ob.Put(t, m, val)
			t.Push(val)
		case op.GetPut:
			i := fetchUint8()
			op := []func(x, y Value) Value{
				OpAdd, OpSub, t.Cat, OpMul, OpDiv, OpMod,
				OpLeftShift, OpRightShift, OpBitOr, OpBitAnd, OpBitXor}[i>>1]
			val := t.Pop()
			m := t.Pop()
			ob := t.Pop()
			t.Push(ob.GetPut(t, m, val, op, i&1 != 0))
		case op.RangeTo:
			j := ToInt(t.Pop())
			i := ToIndex(t.Pop())
			val := t.Pop()
			t.Push(val.RangeTo(i, j))
		case op.RangeLen:
			n := ToInt(t.Pop())
			i := ToIndex(t.Pop())
			val := t.Pop()
			t.Push(val.RangeLen(i, n))
		case op.Is:
			t.sp--
			t.stack[t.sp-1] = OpIs(t.stack[t.sp-1], t.stack[t.sp])
		case op.Isnt:
			t.sp--
			t.stack[t.sp-1] = OpIsnt(t.stack[t.sp-1], t.stack[t.sp])
		case op.Match:
			t.sp--
			t.stack[t.sp-1] = OpMatch(t, t.stack[t.sp-1], t.stack[t.sp])
		case op.MatchNot:
			t.sp--
			t.stack[t.sp-1] = OpMatch(t, t.stack[t.sp-1], t.stack[t.sp]).Not()
		case op.Lt:
			t.sp--
			t.stack[t.sp-1] = OpLt(t.stack[t.sp-1], t.stack[t.sp])
		case op.Lte:
			t.sp--
			t.stack[t.sp-1] = OpLte(t.stack[t.sp-1], t.stack[t.sp])
		case op.Gt:
			t.sp--
			t.stack[t.sp-1] = OpGt(t.stack[t.sp-1], t.stack[t.sp])
		case op.Gte:
			t.sp--
			t.stack[t.sp-1] = OpGte(t.stack[t.sp-1], t.stack[t.sp])
		case op.Add:
			t.sp--
			t.stack[t.sp-1] = OpAdd(t.stack[t.sp-1], t.stack[t.sp])
		case op.Sub:
			t.sp--
			t.stack[t.sp-1] = OpSub(t.stack[t.sp-1], t.stack[t.sp])
		case op.Cat:
			t.sp--
			t.stack[t.sp-1] = OpCat(t, t.stack[t.sp-1], t.stack[t.sp])
		case op.Mul:
			t.sp--
			t.stack[t.sp-1] = OpMul(t.stack[t.sp-1], t.stack[t.sp])
		case op.Div:
			t.sp--
			t.stack[t.sp-1] = OpDiv(t.stack[t.sp-1], t.stack[t.sp])
		case op.Mod:
			t.sp--
			t.stack[t.sp-1] = OpMod(t.stack[t.sp-1], t.stack[t.sp])
		case op.LeftShift:
			t.sp--
			t.stack[t.sp-1] = OpLeftShift(t.stack[t.sp-1], t.stack[t.sp])
		case op.RightShift:
			t.sp--
			t.stack[t.sp-1] = OpRightShift(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitOr:
			t.sp--
			t.stack[t.sp-1] = OpBitOr(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitAnd:
			t.sp--
			t.stack[t.sp-1] = OpBitAnd(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitXor:
			t.sp--
			t.stack[t.sp-1] = OpBitXor(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitNot:
			t.stack[t.sp-1] = OpBitNot(t.stack[t.sp-1])
		case op.Not:
			t.stack[t.sp-1] = OpNot(t.stack[t.sp-1])
		case op.UnaryPlus:
			t.stack[t.sp-1] = OpUnaryPlus(t.stack[t.sp-1])
		case op.UnaryMinus:
			t.stack[t.sp-1] = OpUnaryMinus(t.stack[t.sp-1])
		case op.Bool:
			t.topbool()
		case op.Jump:
			jump()
		case op.JumpTrue:
			if t.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case op.JumpFalse:
			if !t.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case op.And:
			if !t.topbool() {
				jump()
			} else {
				fr.ip += 2
				t.Pop()
			}
		case op.Or:
			if t.topbool() {
				jump()
			} else {
				fr.ip += 2
				t.Pop()
			}
		case op.QMark:
			if !t.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case op.In:
			y := t.Pop()
			x := t.Pop()
			if x.Equal(y) {
				t.Push(True)
				jump()
			} else {
				fr.ip += 2
				t.Push(x)
			}
		case op.JumpIs:
			y := t.Pop()
			x := t.Pop()
			if x.Equal(y) {
				jump()
			} else {
				fr.ip += 2
				t.Push(x)
			}
		case op.JumpIsnt:
			y := t.Pop()
			x := t.Pop()
			if !x.Equal(y) {
				t.Push(x)
				jump()
			} else {
				fr.ip += 2
			}
		case op.Iter:
			t.stack[t.sp-1] = OpIter(t.stack[t.sp-1])
		case op.ForIn:
			brk := fetchInt16()
			local := fetchUint8()
			iter := t.Top()
			nextable := iter.(interface{ Next() Value })
			next := nextable.Next()
			if next != nil {
				fr.locals.v[local] = next
			} else {
				fr.ip += brk - 1 // break
			}
		case op.ReturnNil:
			t.Push(nil)
			fallthrough
		case op.Return:
			break loop
		case op.Try:
			*catchJump = fr.ip + fetchInt16()
			*catchSp = t.sp
			catchPat = string(fr.fn.Values[fetchUint8()].(SuStr))
		case op.Catch:
			fr.ip += fetchInt16()
			*catchJump = 0 // no longer catching
		case op.Throw:
			panic(t.Pop())
		case op.Closure:
			fr.locals.moveToHeap()
			fn := fr.fn.Values[fetchUint8()].(*SuFunc)
			parent := fr
			if fr.blockParent != nil {
				parent = fr.blockParent
			}
			block := &SuClosure{SuFunc: *fn, locals: fr.locals.v, this: fr.this,
				parent: parent}
			t.Push(block)
		case op.BlockBreak:
			panic(BlockBreak)
		case op.BlockContinue:
			panic(BlockContinue)
		case op.BlockReturnNil:
			t.Push(nil)
			fallthrough
		case op.BlockReturn:
			t.blockReturnFrame = fr.blockParent
			panic(BlockReturn)
		case op.CallFuncDiscard, op.CallFuncNoNil, op.CallFuncNilOk:
			f := t.Pop()
			ai := fetchUint8()
			var argSpec *ArgSpec
			if ai < len(StdArgSpecs) {
				argSpec = &StdArgSpecs[ai]
			} else {
				argSpec = &fr.fn.ArgSpecs[ai-len(StdArgSpecs)]
			}
			// fmt.Println(strings.Repeat("    ", t.fp+1), f)
			base := t.sp - int(argSpec.Nargs)
			result := f.Call(t, nil, argSpec)
			t.sp = base
			pushResult(result)
		case op.Super:
			super = fetchUint16()
		case op.CallMethDiscard, op.CallMethNoNil, op.CallMethNilOk:
			method := t.Pop()
			ai := fetchUint8()
			var argSpec *ArgSpec
			if ai < len(StdArgSpecs) {
				argSpec = &StdArgSpecs[ai]
			} else {
				argSpec = &fr.fn.ArgSpecs[ai-len(StdArgSpecs)]
			}
			base := t.sp - int(argSpec.Nargs) - 1
			this := t.stack[base]
			if methstr, ok := method.ToStr(); ok {
				ob := this
				if super > 0 {
					ob = Global.Get(t, super)
					super = 0
				}
				if f := ob.Lookup(t, string(methstr)); f != nil {
					// fmt.Println(strings.Repeat("   ", t.fp+1), f)
					result := f.Call(t, this, argSpec)
					t.sp = base
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
func (t *Thread) topbool() bool {
	switch t.Top() {
	case True:
		return true
	case False:
		return false
	default:
		panic("conditionals require true or false")
	}
}

// popbool pops the top of the stack and returns it as bool, panicing if not True or False
func (t *Thread) popbool() bool {
	switch t.Pop() {
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
func (t *Thread) dyload(fr *Frame, idx int) Value {
	name := fr.fn.Names[idx]
	for i := t.fp - 2; i >= 0; i-- {
		fr2 := &t.frames[i]
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
