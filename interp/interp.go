// Package interp implements the virtual machine interpreter
package interp

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/util/varint"
	. "github.com/apmckinlay/gsuneido/value"
)

func (t *Thread) Interp() Value {
	fr := &t.frames[len(t.frames)-1]
	code := fr.fn.Code
	for {
		fmt.Println("stack", t.stack)
		op := code[fr.ip]
		fr.ip++
		switch op {
		case PUSHINT:
			t.Push(IntVal(fetchInt(code, &fr.ip)))
		case PUSHVAL:
			t.Push(fr.fn.Values[fetchUint(code, &fr.ip)])
		case ADD:
			t.binop(Add)
		case SUB:
			t.binop(Sub)
		case CAT:
			t.binop(Cat)
		case RETURN:
			return t.Pop()
		}
	}
	return nil
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
