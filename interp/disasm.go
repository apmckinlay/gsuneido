package interp

import (
	"fmt"
	"io"
)

var asm = []string{
	"return", "pushint", "pushval", "add", "sub", "cat", "mul", "div", "mod",
}

func Disasm(w io.Writer, fn *Function) {
	code := fn.Code
	for i := 0; i < len(code); {
		op := code[i]
		i++
		fmt.Fprintf(w, "%d %s ", i, asm[op])
		switch op {
		case PUSHINT:
			n := fetchInt(code, &i)
			fmt.Fprintf(w, "%d", n)
		case PUSHVAL:
			v := fn.Values[fetchUint(code, &i)]
			fmt.Fprintf(w, "%v", v)
		}
		fmt.Fprintf(w, "\n")
	}

}
