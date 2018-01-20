package interp

import (
	"fmt"
	"io"

	. "github.com/apmckinlay/gsuneido/interp/op"
	"github.com/apmckinlay/gsuneido/util/verify"
)

var asm = []string{
	"return", "pop", "dup", "dup2", "dupx2", "int", "value",
	"is", "isnt", "match", "matchnot", "lt", "lte", "gt", "gte",
	"add", "sub", "cat", "mul", "div", "mod",
	"lshift", "rshift", "bitor", "bitand", "bitxor",
	"bitnot", "not", "uplus", "uminus",
	"load", "store", "dyload", "get", "put", "global",
	"true", "false", "zero", "one", "emptystr",
	"or", "and", "bool", "qmark", "in", "jump", "tjump", "fjump",
	"eqjump", "nejump", "throw", "call",
}

func init() {
	verify.That(asm[FALSE] == "false")
}

func Disasm(w io.Writer, fn *SuFunc) {
	var s string
	for i := 0; i < len(fn.Code); {
		j := i
		i, s = Disasm1(fn, i)
		fmt.Fprintf(w, "%d: %s\n", j, s)
	}
	fmt.Fprintf(w, "%d:\n", len(fn.Code))
}

func Disasm1(fn *SuFunc, i int) (int, string) {
	op := fn.Code[i]
	i++
	if int(op) >= len(asm) {
		return i, fmt.Sprintf("bad op %d", op)
	}
	s := asm[op]
	switch op {
	case INT:
		n := fetchInt(fn.Code, &i)
		s += fmt.Sprintf(" %d", n)
	case VALUE:
		v := fn.Values[fetchUint(fn.Code, &i)]
		s += fmt.Sprintf(" %v", v)
	case LOAD, STORE, DYLOAD:
		idx := fetchUint(fn.Code, &i)
		s += " " + fn.Strings[idx]
	case GLOBAL:
		idx := fetchUint(fn.Code, &i)
		s += " " + NumNameG(int(idx))
	case JUMP, TJUMP, FJUMP, AND, OR, Q_MARK, IN, EQJUMP, NEJUMP:
		ip := i
		i += 2
		jump(fn.Code, &ip)
		s += fmt.Sprintf(" %d", ip)
	}
	return i, s
}
