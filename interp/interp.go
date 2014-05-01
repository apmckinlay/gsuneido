// Package interp implements the virtual machine interpreter
package interp

import (
	"fmt"
	. "github.com/apmckinlay/gsuneido/core/value"
)

func (t *Thread) Interp() Value {
	fr := &t.frames[len(t.frames)-1]
	code := fr.fn.code
	for {
		fmt.Println("stack", t.stack)
		op := code[fr.ip]
		fr.ip++
		switch op {
		case PUSHINT:
			t.Push(IntVal(fetchInt(code, &fr.ip)))
		case ADD:
			x := t.Pop()
			y := t.Pop()
			t.Push(Add(x, y))
		case RETURN:
			return t.Pop()
		}
	}
	return nil
}

func fetchInt(code []byte, ip *int) int {
	i := int(code[*ip])
	*ip++
	// TODO handle variable length ints
	return i
}
