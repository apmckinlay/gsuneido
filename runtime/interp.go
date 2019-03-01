package runtime

import (
	"runtime"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime/op"
)

var blockBreak = &SuExcept{SuStr: SuStr("block:break")}
var blockContinue = &SuExcept{SuStr: SuStr("block:continue")}
var blockReturn = &SuExcept{SuStr: SuStr("block return")}

// Call sets up a frame to Run a compiled Suneido function
// The stack must already be in the form required by the function (massaged)
func (t *Thread) Call(fn *SuFunc) Value {
	// reserve stack space for locals
	for expand := fn.Nlocals - fn.Nparams; expand > 0; expand-- {
		t.Push(nil)
	}
	t.frames[t.fp] = Frame{fn: fn, this: t.this,
		locals: t.stack[t.sp-int(fn.Nlocals) : t.sp]}
	return t.run()
}

func CallMethod(t *Thread, this Value, f Value, as *ArgSpec) Value {
	t.this = this
	return f.Call(t, as)
}

// run is needed in addition to interp
// because we can only recover panic on the way out of a function
// so if the exception is caught we have to re-enter interp
// Called by Thread.Call and SuBlock.Call
func (t *Thread) run() Value {
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
	for i := 0; i < 4; i++ {
		result := t.interp(&catchJump)
		if result == nil {
			t.fp = fp - 1
			if t.sp > sp {
				return t.Pop()
			}
			return nil
		}
		// try block threw
		t.sp = sp
		t.fp = fp
		fr := &t.frames[t.fp-1]
		fr.ip = catchJump
		catchJump = 0  // no longer catching
		t.Push(result) // SuExcept
		// loop and re-enter interp
	}
	panic("Run too many loops")
}

// interp is the main interpreter loop
// It normally returns nil, with the return value (if any) on the stack
// Returns *SuExcept if there was an exception/panic
func (t *Thread) interp(catchJump *int) (ret Value) {
	fr := &t.frames[t.fp-1]
	code := fr.fn.Code
	super := 0
	catchPat := ""

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
	defer func() {
		if *catchJump == 0 && fr.fn.Id == 0 {
			return // this frame isn't catching
		}
		e := recover()
		if e == nil {
			return // not panic'ing, normal return
		}
		if e == blockReturn {
			if t.frames[t.fp-1].fn.OuterId != fr.fn.Id {
				panic(e) // not our block, rethrow
			}
			return // normal return
		}
		if *catchJump == 0 {
			return // not catching
		}
		se, ok := e.(*SuExcept)
		if !ok {
			// first catch creates SuExcept with callstack
			var ss SuStr
			if re, ok := e.(runtime.Error); ok {
				ss = SuStr(re.Error())
			} else if s, ok := e.(string); ok {
				ss = SuStr(s)
			} else {
				ss = e.(SuStr)
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
		// fmt.Println("stack:", t.stack[spBase:t.sp])
		// _, da := Disasm1(fr.fn, fr.ip)
		// fmt.Printf("%d: %s\n", fr.ip, da)
		op := code[fr.ip]
		fr.ip++
		switch op {
		case POP:
			t.Pop()
		case DUP:
			t.Push(t.Top())
		case DUP2:
			t.Dup2() // dup top two, used to dup member lvalues
		case DUPX2:
			t.Dupx2() // dup top under next two, used for post inc/dec
		case THIS:
			t.Push(fr.this)
		case TRUE:
			t.Push(True)
		case FALSE:
			t.Push(False)
		case ZERO:
			t.Push(Zero)
		case ONE:
			t.Push(One)
		case MAXINT:
			t.Push(MaxInt)
		case EMPTYSTR:
			t.Push(EmptyStr)
		case INT:
			t.Push(SuInt(fetchInt16()))
		case VALUE:
			t.Push(fr.fn.Values[fetchUint8()])
		case LOAD:
			i := fetchUint8()
			val := fr.locals[i]
			if val == nil {
				panic("uninitialized variable: " + fr.fn.Names[i])
			}
			t.Push(val)
		case STORE:
			fr.locals[fetchUint8()] = t.Top()
		case DYLOAD:
			i := fetchUint8()
			if fr.locals[i] == nil {
				t.dyload(fr, i)
			}
			t.Push(fr.locals[i])
		case GLOBAL:
			gn := Global(fetchUint16())
			val := GetGlobal(gn)
			if val == nil {
				panic("uninitialized global: " + GlobalName(gn))
			}
			t.Push(val)
		case GET:
			m := t.Pop()
			ob := t.Pop()
			val := ob.Get(t, m)
			if val == nil {
				panic("uninitialized member: " + m.String())
			}
			t.Push(val)
		case PUT:
			val := t.Pop()
			m := t.Pop()
			ob := t.Pop()
			ob.Put(m, val)
			t.Push(val)
		case RANGETO:
			j := Index(t.Pop())
			i := Index(t.Pop())
			val := t.Pop()
			t.Push(val.RangeTo(i, j))
		case RANGELEN:
			n := Index(t.Pop())
			i := Index(t.Pop())
			val := t.Pop()
			t.Push(val.RangeLen(i, n))
		case IS:
			t.sp--
			t.stack[t.sp-1] = Is(t.stack[t.sp-1], t.stack[t.sp])
		case ISNT:
			t.sp--
			t.stack[t.sp-1] = Isnt(t.stack[t.sp-1], t.stack[t.sp])
		case MATCH:
			t.sp--
			pat := t.rxcache.Get(t.stack[t.sp].ToStr())
			s := t.stack[t.sp-1]
			t.stack[t.sp-1] = Match(s, pat)
		case MATCHNOT:
			t.sp--
			pat := t.rxcache.Get(t.stack[t.sp].ToStr())
			s := t.stack[t.sp-1]
			t.stack[t.sp-1] = Match(s, pat).Not()
		case LT:
			t.sp--
			t.stack[t.sp-1] = Lt(t.stack[t.sp-1], t.stack[t.sp])
		case LTE:
			t.sp--
			t.stack[t.sp-1] = Lte(t.stack[t.sp-1], t.stack[t.sp])
		case GT:
			t.sp--
			t.stack[t.sp-1] = Gt(t.stack[t.sp-1], t.stack[t.sp])
		case GTE:
			t.sp--
			t.stack[t.sp-1] = Gte(t.stack[t.sp-1], t.stack[t.sp])
		case ADD:
			t.sp--
			t.stack[t.sp-1] = Add(t.stack[t.sp-1], t.stack[t.sp])
		case SUB:
			t.sp--
			t.stack[t.sp-1] = Sub(t.stack[t.sp-1], t.stack[t.sp])
		case CAT:
			t.sp--
			t.stack[t.sp-1] = Cat(t.stack[t.sp-1], t.stack[t.sp])
		case MUL:
			t.sp--
			t.stack[t.sp-1] = Mul(t.stack[t.sp-1], t.stack[t.sp])
		case DIV:
			t.sp--
			t.stack[t.sp-1] = Div(t.stack[t.sp-1], t.stack[t.sp])
		case MOD:
			t.sp--
			t.stack[t.sp-1] = Mod(t.stack[t.sp-1], t.stack[t.sp])
		case LSHIFT:
			t.sp--
			t.stack[t.sp-1] = Lshift(t.stack[t.sp-1], t.stack[t.sp])
		case RSHIFT:
			t.sp--
			t.stack[t.sp-1] = Rshift(t.stack[t.sp-1], t.stack[t.sp])
		case BITOR:
			t.sp--
			t.stack[t.sp-1] = Bitor(t.stack[t.sp-1], t.stack[t.sp])
		case BITAND:
			t.sp--
			t.stack[t.sp-1] = Bitand(t.stack[t.sp-1], t.stack[t.sp])
		case BITXOR:
			t.sp--
			t.stack[t.sp-1] = Bitxor(t.stack[t.sp-1], t.stack[t.sp])
		case BITNOT:
			t.stack[t.sp-1] = Bitnot(t.stack[t.sp-1])
		case NOT:
			t.stack[t.sp-1] = Not(t.stack[t.sp-1])
		case UPLUS:
			t.stack[t.sp-1] = Uplus(t.stack[t.sp-1])
		case UMINUS:
			t.stack[t.sp-1] = Uminus(t.stack[t.sp-1])
		case BOOL:
			t.topbool()
		case JUMP:
			jump()
		case TJUMP:
			if t.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case FJUMP:
			if !t.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case AND:
			if !t.topbool() {
				jump()
			} else {
				fr.ip += 2
				t.Pop()
			}
		case OR:
			if t.topbool() {
				jump()
			} else {
				fr.ip += 2
				t.Pop()
			}
		case Q_MARK:
			if !t.popbool() {
				jump()
			} else {
				fr.ip += 2
			}
		case IN:
			y := t.Pop()
			x := t.Pop()
			if x.Equal(y) {
				t.Push(True)
				jump()
			} else {
				fr.ip += 2
				t.Push(x)
			}
		case EQJUMP:
			y := t.Pop()
			x := t.Pop()
			if x.Equal(y) {
				jump()
			} else {
				fr.ip += 2
				t.Push(x)
			}
		case NEJUMP:
			y := t.Pop()
			x := t.Pop()
			if !x.Equal(y) {
				t.Push(x)
				jump()
			} else {
				fr.ip += 2
			}
		case ITER:
			t.Push(t.callMethod("Iter", ArgSpec0))
		case FORIN:
			brk := fetchInt16()
			local := fetchUint8()
			iter := t.Top()
			t.Push(iter) // since call will pop it
			next := t.callMethod("Next", ArgSpec0)
			fr.locals[local] = next
			if next.Equal(iter) {
				fr.ip += brk - 1 // jump
			}
		case RETURN:
			break loop
		case TRY:
			*catchJump = fr.ip + fetchInt16()
			catchPat = string(fr.fn.Values[fetchUint8()].(SuStr))
		case CATCH:
			fr.ip += fetchInt16()
			*catchJump = 0
		case THROW:
			panic(t.Pop())
		case BLOCK:
			fr.moveLocalsToHeap()
			fn := fr.fn.Values[fetchUint8()].(*SuFunc)
			block := &SuBlock{SuFunc: *fn, locals: fr.locals, this: fr.this}
			t.Push(block)
		case BLOCK_BREAK:
			panic(blockBreak)
		case BLOCK_CONTINUE:
			panic(blockContinue)
		case BLOCK_RETURN_NULL:
			t.Push(nil)
			fallthrough
		case BLOCK_RETURN:
			panic(blockReturn)
		case CALLFUNC:
			f := t.Pop()
			ai := fetchUint8()
			var argSpec *ArgSpec
			if ai < len(StdArgSpecs) {
				argSpec = StdArgSpecs[ai]
			} else {
				argSpec = fr.fn.ArgSpecs[ai-len(StdArgSpecs)]
			}
			base := t.sp - int(argSpec.Nargs)
			result := f.Call(t, argSpec)
			t.sp = base
			t.Push(result)
		case SUPER:
			super = Global(fetchUint16())
		case CALLMETH:
			method := t.Pop()
			ai := fetchUint8()
			var argSpec *ArgSpec
			if ai < len(StdArgSpecs) {
				argSpec = StdArgSpecs[ai]
			} else {
				argSpec = fr.fn.ArgSpecs[ai-len(StdArgSpecs)]
			}
			base := t.sp - int(argSpec.Nargs) - 1
			this := t.stack[base]
			if methstr, ok := method.(SuStr); ok {
				ob := this
				if super > 0 {
					ob = GetGlobal(super)
					super = 0
				}
				if f := ob.Lookup(string(methstr)); f != nil {
					t.this = this
					result := f.Call(t, argSpec)
					t.sp = base
					t.Push(result)
					break
				}
			}
			panic("method not found " + this.TypeName() + "." + method.ToStr())
		default:
			panic("invalid op code: " + asm[op]) // TODO fatal?
		}
	}
	t.this = nil
	return nil
}

// callMethod is used by ITER and FORIN
func (t *Thread) callMethod(method string, argSpec *ArgSpec) Value {
	base := t.sp - int(argSpec.Nargs) - 1
	ob := t.stack[base]
	f := ob.Lookup(method)
	if f == nil {
		panic("method not found " + ob.TypeName() + "." + method)
	}
	t.this = ob
	result := f.Call(t, argSpec)
	t.sp = base
	return result
}

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

func (t *Thread) dyload(fr *Frame, idx int) {
	name := fr.fn.Names[idx]
	for i := t.fp - 1; i >= 0; i-- {
		fr2 := &t.frames[i]
		for j, s := range fr2.fn.Names {
			if s == name {
				fr.locals[idx] = fr2.locals[j]
				return
			}
		}
	}
	panic("uninitialized variable: " + name)
}

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
