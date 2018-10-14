// Package interp implements the virtual machine interpreter
package interp

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/interp/global"
	. "github.com/apmckinlay/gsuneido/interp/op"
)

// Call sets up a frame to Run a compiled Suneido function
// The stack must already be in the form required by the function (massaged)
func (t *Thread) Call(fn *SuFunc, self Value) Value {
	// expand stack if necessary for locals
	expand := fn.Nlocals - fn.Nparams
	for ; expand > 0; expand-- {
		t.Push(nil)
	}
	t.frames[t.fp] = Frame{fn: fn, locals: t.stack[t.sp-fn.Nlocals:], self: self}
	defer func(fp int) { t.fp = fp }(t.fp)
	t.fp++
	return t.Run()
}

var _ Context = (*Thread)(nil) // verify Thread satisfies Context

func (t *Thread) Run() Value {
	fr := &t.frames[t.fp-1]
	code := fr.fn.Code
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

	sp := t.sp
	for fr.ip < len(code) {
		// fmt.Println("stack:", t.stack[sp:t.sp])
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
		case TRUE:
			t.Push(True)
		case FALSE:
			t.Push(False)
		case ZERO:
			t.Push(Zero)
		case ONE:
			t.Push(One)
		case EMPTYSTR:
			t.Push(EmptyStr)
		case INT:
			t.Push(SuInt(fetchInt16()))
		case VALUE:
			t.Push(fr.fn.Values[fetchUint16()])
		case LOAD:
			t.Push(t.load(fr, fetchUint8()))
		case STORE:
			fr.locals[fetchUint8()] = t.Top()
		case DYLOAD:
			i := fetchUint8()
			if fr.locals[i] == nil {
				t.dyload(fr, i)
			}
			t.Push(fr.locals[i])
		case GLOBAL:
			gn := int(fetchUint16())
			val := global.Get(gn)
			if val == nil {
				panic("uninitialized global: " + global.Name(gn))
			}
			t.Push(val)
		case GET:
			m := t.Pop()
			ob := t.Pop()
			val := ob.Get(m)
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
			ob := t.Pop()
			t.Push(ob.RangeTo(i, j))
		case RANGELEN:
			n := Index(t.Pop())
			i := Index(t.Pop())
			ob := t.Pop()
			t.Push(ob.RangeLen(i, n))
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
		case RETURN:
			break
		case THROW:
			panic(t.Pop())
		case CALL:
			f := t.Pop()
			unnamed := code[fr.ip]
			fr.ip++
			named := int(code[fr.ip])
			fr.ip++
			var spec []byte
			if named > 0 {
				spec = code[fr.ip : fr.ip+named]
			}
			fr.ip += named
			argSpec := &ArgSpec{unnamed, spec, fr.fn.Strings}
			nargs := argSpec.Nargs()
			var result Value
			switch f := f.(type) {
			case Builtin:
				result = f(argSpec, t.stack[t.sp-nargs:]...)
			case Callable:
				result = f.Call(t, nil, t.args(f.Params(), argSpec)...)
			default:
				panic("can't call " + f.TypeName())
			}
			t.sp -= nargs
			t.Push(result)
		default:
			panic("invalid op code")
		}
	}
	if t.sp > sp {
		return t.Pop()
	}
	return nil
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

func (t *Thread) load(fr *Frame, i int) Value {
	val := fr.locals[i]
	if val == nil {
		panic("uninitialized variable: " + fr.fn.Strings[i])
	}
	return val
}

func (t *Thread) dyload(fr *Frame, idx int) {
	name := fr.fn.Strings[idx]
	for i := t.fp - 1; i >= 0; i-- {
		fr2 := t.frames[i]
		for j, s := range fr2.fn.Strings {
			if s == name {
				fr.locals[idx] = fr2.locals[j]
				return
			}
		}
	}
	panic("uninitialized variable: " + name)
}
