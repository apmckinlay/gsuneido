// Package interp implements the virtual machine interpreter
package interp

import (
	"fmt"
	"os"

	"github.com/apmckinlay/gsuneido/util/varint"
	. "github.com/apmckinlay/gsuneido/value"
)

func (t *Thread) Interp() Value {
	fr := &t.frames[len(t.frames)-1]
	code := fr.fn.Code
	sp := len(t.stack)
	for {
		fmt.Println("stack:", t.stack[sp:])
		disasm1(os.Stdout, fr.fn, fr.ip)
		op := code[fr.ip]
		fr.ip++
		switch op {
		case POP:
			t.Pop()
		case DUP:
			t.Push(t.Top())
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
			if len(t.stack) > sp {
				return t.Pop()
			}
			return nil
		default:
			panic("invalid op code")
		}
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
