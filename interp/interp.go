// Package interp implements the virtual machine interpreter
package interp

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/interp/global"
	. "github.com/apmckinlay/gsuneido/interp/op"
	"github.com/apmckinlay/gsuneido/util/varint"
)

type CallSpec struct {
	t  *Thread
	as ArgSpec
}

func (c CallSpec) CallSuFunc(f *SuFunc) Value {
	return c.t.Call(f, c.as)
}

var _ CallContext = CallSpec{}

// Call executes a SuFunc and returns the result.
// The arguments must be already on the stack as per the ArgSpec.
// On return, the arguments are removed from the stack.
func (t *Thread) Call(fn *SuFunc, as ArgSpec) Value {
	defer func(sp int) { t.stack = t.stack[:sp] }(len(t.stack) - as.Nargs())
	t.args(fn, as)
	base := len(t.stack) - fn.Nparams
	for i := fn.Nparams; i < fn.Nlocals; i++ {
		t.Push(nil)
	}
	frame := Frame{fn: fn, ip: 0, locals: t.stack[base:]}
	t.frames = append(t.frames, frame)
	defer func(fp int) { t.frames = t.frames[:fp] }(len(t.frames) - 1)
	return t.Run()
}

func (t *Thread) Run() Value {
	fr := &t.frames[len(t.frames)-1]
	code := fr.fn.Code
	sp := len(t.stack)
	for fr.ip < len(code) {
		// fmt.Println("stack:", t.stack[sp:])
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
			t.Push(SuInt(fetchInt(code, &fr.ip)))
		case VALUE:
			t.Push(fr.fn.Values[fetchUint(code, &fr.ip)])
		case LOAD:
			t.Push(t.load(fr, fetchUint(code, &fr.ip)))
		case STORE:
			fr.locals[fetchUint(code, &fr.ip)] = t.Top()
		case DYLOAD:
			i := fetchUint(code, &fr.ip)
			if fr.locals[i] == nil {
				t.dyload(fr, i)
			}
			t.Push(fr.locals[i])
		case GLOBAL:
			gn := int(fetchUint(code, &fr.ip))
			val := global.Get(gn)
			if val == nil {
				panic("uninitialized global: " + global.Name(gn))
			}
			t.Push(val)
		case GET:
			m := t.Pop()
			ob := t.Pop()
			t.Push(ob.Get(m))
		case PUT:
			val := t.Pop()
			m := t.Pop()
			ob := t.Pop()
			ob.Put(m, val)
			t.Push(val)
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
			jump(code, &fr.ip)
		case TJUMP:
			if t.popbool() {
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
			}
		case FJUMP:
			if !t.popbool() {
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
			}
		case AND:
			if !t.topbool() {
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
				t.Pop()
			}
		case OR:
			if t.topbool() {
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
				t.Pop()
			}
		case Q_MARK:
			if !t.popbool() {
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
			}
		case IN:
			y := t.Pop()
			x := t.Pop()
			if x.Equals(y) {
				t.Push(True)
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
				t.Push(x)
			}
		case EQJUMP:
			y := t.Pop()
			x := t.Pop()
			if x.Equals(y) {
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
				t.Push(x)
			}
		case NEJUMP:
			y := t.Pop()
			x := t.Pop()
			if !x.Equals(y) {
				t.Push(x)
				jump(code, &fr.ip)
			} else {
				fr.ip += 2
			}
		case RETURN:
			break
		case THROW:
			panic(t.Pop())
		case CALL:
			f := t.Pop()
			nargs := code[fr.ip]
			fr.ip++
			t.Push(f.Call(CallSpec{t, ArgSpec{Unnamed: nargs}}))
		case CALL_NAMED:
			//TODO
		default:
			panic("invalid op code")
		}
	}
	if len(t.stack) > sp {
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

func (t *Thread) load(fr *Frame, idx uint32) Value {
	val := fr.locals[idx]
	if val == nil {
		panic("uninitialized variable: " + fr.fn.Strings[idx])
	}
	return val
}

func (t *Thread) dyload(fr *Frame, idx uint32) {
	name := fr.fn.Strings[idx]
	for i := len(t.frames) - 1; i > 0; i-- {
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

func (t *Thread) unop(op func(Value) Value) {
	x := t.Pop()
	t.Push(op(x))
}

func (t *Thread) binop(op func(Value, Value) Value) {
	y := t.Pop()
	x := t.Pop()
	t.Push(op(x, y))
}

func fetchInt(code []byte, ip *int) (i int32) {
	i, *ip = varint.DecodeInt32(code, *ip)
	return i
}

func fetchUint(code []byte, ip *int) (i uint32) {
	i, *ip = varint.DecodeUint32(code, *ip)
	return i
}

func jump(code []byte, ip *int) {
	*ip += 2 + int(int16(uint16(code[*ip])<<8+uint16(code[*ip+1])))
}
