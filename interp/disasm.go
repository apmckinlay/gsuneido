package interp

import (
	"fmt"
	"io"
)

var asm = []string{
	"return", "pushint", "pushval", "add", "sub", "cat", "mul", "div", "mod",
	"storvar", "loadvar", "uplus", "uminus",
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
		case STORVAR, LOADVAR:
			idx := fetchUint(code, &i)
			varname := fn.Strings[idx]
			fmt.Fprintf(w, "%s (%d)", varname, idx)
		}
		fmt.Fprintf(w, "\n")
	}

}
