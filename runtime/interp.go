package runtime

import (
	"runtime"
	"strings"

	op "github.com/apmckinlay/gsuneido/runtime/opcodes"
)

var blockBreak = &SuExcept{SuStr: SuStr("block:break")}
var blockContinue = &SuExcept{SuStr: SuStr("block:continue")}
var BlockReturn = &SuExcept{SuStr: SuStr("block return")}

// Start sets up a frame to Run a compiled Suneido function
// The stack must already be in the form required by the function (massaged)
func (t *Thread) Start(fn *SuFunc, this Value) Value {
	// reserve stack space for locals
	for expand := fn.Nlocals - fn.Nparams; expand > 0; expand-- {
		t.Push(nil)
	}
	t.frames[t.fp] = Frame{fn: fn, this: this,
		locals: t.stack[t.sp-int(fn.Nlocals) : t.sp]}
	return t.run()
}

// run is needed in addition to interp
// because we can only recover panic on the way out of a function
// so if the exception is caught we have to re-enter interp
// Called by Thread.Start and SuBlock.Call
func (t *Thread) run() Value {
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(strings.Repeat("    ", t.fp) + "run:", t.frames[t.fp].fn)
	sp := t.sp
	t.fp++
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
		for i := t.sp; i < t.spMax || t.stack[i] != nil; i++ {
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

	defer func() {
		if *catchJump == 0 && fr.fn.Id == 0 {
			return // this frame isn't catching
		}
		e := recover()
		if e == nil {
			return // not panic'ing, normal return
		}
		if e == BlockReturn {
			if t.frames[t.fp-1].fn.OuterId != fr.fn.Id {
				panic(e) // not our block, rethrow
			}
			return // normal return
		}
		if *catchJump == 0 {
			panic(e) // not catching
		}
		se, ok := e.(*SuExcept)
		if !ok {
			// first catch creates SuExcept with callstack
			var ss SuStr
			if re, ok := e.(runtime.Error); ok {
				// debug.PrintStack()
				ss = SuStr(re.Error())
			} else if s, ok := e.(string); ok {
				ss = SuStr(s)
			} else {
				ss = SuStr(ToStr(e.(Value)))
			}
			se = NewSuExcept(t, ss)
		}
		if catchMatch(string(se.SuStr), catchPat) {
			ret = se // tells run we're catching
		} else {
			panic(se)
		}
	}()

loop:
	for fr.ip < len(code) {
		// fmt.Println("stack:", t.sp, t.stack[ints.Max(0, t.sp-3):t.sp])
		// _, da := Disasm1(fr.fn, fr.ip)
		// fmt.Printf("%d: %d: %s\n", t.fp, fr.ip, da)
		oc = op.Opcode(code[fr.ip])
		fr.ip++
		switch oc {
		case op.Pop:
			t.Pop()
		case op.Dup:
			t.Push(t.Top())
		case op.Dup2:
			t.Dup2()
		case op.Dupx2:
			t.Dupx2()
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
			val := fr.locals[i]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[i])
			}
			t.Push(val)
		case op.Store:
			fr.locals[fetchUint8()] = t.Top()
		case op.Dyload:
			i := fetchUint8()
			if fr.locals[i] == nil {
				t.dyload(fr, i)
			}
			t.Push(fr.locals[i])
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
		case op.RangeTo:
			j := ToInt(t.Pop())
			i := Index(t.Pop())
			val := t.Pop()
			t.Push(val.RangeTo(i, j))
		case op.RangeLen:
			n := ToInt(t.Pop())
			i := Index(t.Pop())
			val := t.Pop()
			t.Push(val.RangeLen(i, n))
		case op.Is:
			t.sp--
			t.stack[t.sp-1] = Is(t.stack[t.sp-1], t.stack[t.sp])
		case op.Isnt:
			t.sp--
			t.stack[t.sp-1] = Isnt(t.stack[t.sp-1], t.stack[t.sp])
		case op.Match:
			t.sp--
			pat := t.RxCache.Get(ToStr(t.stack[t.sp]))
			s := t.stack[t.sp-1]
			t.stack[t.sp-1] = Match(s, pat)
		case op.MatchNot:
			t.sp--
			pat := t.RxCache.Get(ToStr(t.stack[t.sp]))
			s := t.stack[t.sp-1]
			t.stack[t.sp-1] = Match(s, pat).Not()
		case op.Lt:
			t.sp--
			t.stack[t.sp-1] = Lt(t.stack[t.sp-1], t.stack[t.sp])
		case op.Lte:
			t.sp--
			t.stack[t.sp-1] = Lte(t.stack[t.sp-1], t.stack[t.sp])
		case op.Gt:
			t.sp--
			t.stack[t.sp-1] = Gt(t.stack[t.sp-1], t.stack[t.sp])
		case op.Gte:
			t.sp--
			t.stack[t.sp-1] = Gte(t.stack[t.sp-1], t.stack[t.sp])
		case op.Add:
			t.sp--
			t.stack[t.sp-1] = Add(t.stack[t.sp-1], t.stack[t.sp])
		case op.Sub:
			t.sp--
			t.stack[t.sp-1] = Sub(t.stack[t.sp-1], t.stack[t.sp])
		case op.Cat:
			t.sp--
			t.stack[t.sp-1] = Cat(t, t.stack[t.sp-1], t.stack[t.sp])
		case op.Mul:
			t.sp--
			t.stack[t.sp-1] = Mul(t.stack[t.sp-1], t.stack[t.sp])
		case op.Div:
			t.sp--
			t.stack[t.sp-1] = Div(t.stack[t.sp-1], t.stack[t.sp])
		case op.Mod:
			t.sp--
			t.stack[t.sp-1] = Mod(t.stack[t.sp-1], t.stack[t.sp])
		case op.LeftShift:
			t.sp--
			t.stack[t.sp-1] = LeftShift(t.stack[t.sp-1], t.stack[t.sp])
		case op.RightShift:
			t.sp--
			t.stack[t.sp-1] = RightShift(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitOr:
			t.sp--
			t.stack[t.sp-1] = BitOr(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitAnd:
			t.sp--
			t.stack[t.sp-1] = BitAnd(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitXor:
			t.sp--
			t.stack[t.sp-1] = BitXor(t.stack[t.sp-1], t.stack[t.sp])
		case op.BitNot:
			t.stack[t.sp-1] = BitNot(t.stack[t.sp-1])
		case op.Not:
			t.stack[t.sp-1] = Not(t.stack[t.sp-1])
		case op.UnaryPlus:
			t.stack[t.sp-1] = UnaryPlus(t.stack[t.sp-1])
		case op.UnaryMinus:
			t.stack[t.sp-1] = UnaryMinus(t.stack[t.sp-1])
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
			x := t.Pop()
			iterable, ok := x.(interface{ Iter() Iter })
			if !ok {
				panic("can't iterate " + x.Type().String())
			}
			t.Push(SuIter{Iter: iterable.Iter()})
		case op.ForIn:
			brk := fetchInt16()
			local := fetchUint8()
			iter := t.Top()
			nextable := iter.(interface{ Next() Value })
			next := nextable.Next()
			if next != nil {
				fr.locals[local] = next
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
		case op.Block:
			fr.moveLocalsToHeap()
			fn := fr.fn.Values[fetchUint8()].(*SuFunc)
			block := &SuBlock{SuFunc: *fn, locals: fr.locals, this: fr.this}
			t.Push(block)
		case op.BlockBreak:
			panic(blockBreak)
		case op.BlockContinue:
			panic(blockContinue)
		case op.BlockReturnNil:
			t.Push(nil)
			fallthrough
		case op.BlockReturn:
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
		default:
			panic("invalid op code: " + oc.String()) // TODO fatal?
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
func (t *Thread) dyload(fr *Frame, idx int) {
	name := fr.fn.Names[idx]
	for i := t.fp - 2; i >= 0; i-- {
		fr2 := &t.frames[i]
		for j, s := range fr2.fn.Names {
			if s == name && fr2.locals[j] != nil {
				fr.locals[idx] = fr2.locals[j]
				return
			}
		}
	}
	panic("uninitialized variable: " + name)
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
