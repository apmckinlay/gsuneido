// Package interp implements the virtual machine interpreter
package interp

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/interp/globals"
	"github.com/apmckinlay/gsuneido/util/varint"
	. "github.com/apmckinlay/gsuneido/value"
)

func (t *Thread) Interp() Value {
	fr := &t.frames[len(t.frames)-1]
	code := fr.fn.Code
	sp := len(t.stack)
	for fr.ip < len(code) {
		fmt.Println("stack:", t.stack[sp:])
		_, da := Disasm1(fr.fn, fr.ip)
		fmt.Printf("%d: %s\n", fr.ip, da)
		op := code[fr.ip]
		fr.ip++
		switch op {
		case POP:
			t.Pop()
		case DUP:
			t.Push(t.Top())
		case DUP2:
			t.Dup2() // dup top two, used to dup member lvalues
		case TRUE:
			t.Push(True)
		case FALSE:
			t.Push(False)
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
			val := globals.Get(gn)
			if val == nil {
				panic("uninitialized global: " + globals.NumName(gn))
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
		case RETURN:
			break
		default:
			panic("invalid op code")
		}
	}
	if len(t.stack) > sp {
		return t.Pop()
	}
	return nil
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
