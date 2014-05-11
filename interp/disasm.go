package interp

import (
	"fmt"
	"io"

	"github.com/apmckinlay/gsuneido/util/verify"
	"github.com/apmckinlay/gsuneido/value"
)

var asm = []string{
	"return", "pop", "dup", "int", "value",
	"is", "isnt", "match", "matchnot", "lt", "lte", "gt", "gte",
	"add", "sub", "cat", "mul", "div", "mod",
	"lshift", "rshift", "bitor", "bitand", "bitxor",
	"bitnot", "not", "uplus", "uminus",
	"store", "load",
}

func init() {
	verify.That(asm[LOAD] == "load")
}

func Disasm(w io.Writer, fn *value.SuFunc) {
	for i := 0; i < len(fn.Code); {
		i = disasm1(w, fn, i)
	}
}

func disasm1(w io.Writer, fn *value.SuFunc, i int) int {
	op := fn.Code[i]
	i++
	fmt.Fprintf(w, "%d %s ", i, asm[op])
	switch op {
	case INT:
		n := fetchInt(fn.Code, &i)
		fmt.Fprintf(w, "%d", n)
	case VALUE:
		v := fn.Values[fetchUint(fn.Code, &i)]
		fmt.Fprintf(w, "%v", v)
	case STORE, LOAD:
		idx := fetchUint(fn.Code, &i)
		varname := fn.Strings[idx]
		fmt.Fprintf(w, "%s (%d)", varname, idx)
	}
	fmt.Fprintf(w, "\n")
	return i
}
