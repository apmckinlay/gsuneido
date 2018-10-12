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
	t.frames[t.fp] = Frame{fn: fn, bp: t.sp - fn.Nlocals, self: self}
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
			t.Push(SuInt(0))
		case ONE:
			t.Push(SuInt(1))
		case EMPTYSTR:
			t.Push(SuStr(""))
		case INT:
			t.Push(SuInt(fetchInt16()))
		case VALUE:
			t.Push(fr.fn.Values[fetchUint16()])
		case LOAD:
			t.Push(t.load(fr, fetchUint8()))
		case STORE:
			t.stack[fr.bp+fetchUint8()] = t.Top()
		case DYLOAD:
			i := fetchUint8()
			if t.stack[fr.bp+i] == nil {
				t.dyload(fr, i)
			}
			t.Push(t.stack[fr.bp+i])
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
			t.binop(Is)
		case ISNT:
			t.binop(Isnt)
		case MATCH:
			pat := t.rxcache.Get(t.Pop().ToStr())
			s := t.Pop()
			t.Push(Match(s, pat))
		case MATCHNOT:
			pat := t.rxcache.Get(t.Pop().ToStr())
			s := t.Pop()
			t.Push(Match(s, pat).Not())
		case LT:
			t.binop(Lt)
		case LTE:
			t.binop(Lte)
		case GT:
			t.binop(Gt)
		case GTE:
			t.binop(Gte)
		case ADD:
			t.binop(Add)
		case SUB:
			t.binop(Sub)
		case CAT:
			t.binop(Cat)
		case MUL:
			t.binop(Mul)
		case DIV:
			t.binop(Div)
		case MOD:
			t.binop(Mod)
		case LSHIFT:
			t.binop(Lshift)
		case RSHIFT:
			t.binop(Rshift)
		case BITOR:
			t.binop(Bitor)
		case BITAND:
			t.binop(Bitand)
		case BITXOR:
			t.binop(Bitxor)
		case BITNOT:
			t.unop(Bitnot)
		case NOT:
			t.unop(Not)
		case UPLUS:
			t.unop(Uplus)
		case UMINUS:
			t.unop(Uminus)
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
			spec := code[fr.ip : fr.ip+named]
			fr.ip += named
			switch f := f.(type) {
			case Callable:
				result := f.Call(t, nil,
					t.args(f.Params(), ArgSpec{unnamed, spec, fr.fn.Strings})...)
				t.sp -= int(unnamed) + named
				t.Push(result)
			default:
				panic("can't call " + f.TypeName())
			}
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
	val := t.stack[fr.bp+i]
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
				t.stack[fr.bp+idx] = t.stack[fr2.bp+j]
				return
			}
		}
	}
	panic("uninitialized variable: " + name)
}

func (t *Thread) unop(op func(Value) Value) {
	x := t.Pop()
	t.Push(op(x))
}

func (t *Thread) binop(op func(Value, Value) Value) {
	y := t.Pop()
	x := t.Pop()
	t.Push(op(x, y))
}
