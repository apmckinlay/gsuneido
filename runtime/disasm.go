package runtime

import (
	"fmt"
	"io"

	. "github.com/apmckinlay/gsuneido/runtime/op"
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
	"eqjump", "nejump", "throw", "rangeto", "rangelen", "this",
	"callfunc", "callfunc0", "callfunc1", "callfunc2", "callfunc3", "callfunc4",
	"callmeth", "callmeth0", "callmeth1", "callmeth2", "callmeth3",
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
	fetchUint8 := func() int {
		i++
		return int(fn.Code[i-1])
	}
	fetchInt16 := func() int {
		i += 2
		return int(int16(uint16(fn.Code[i-2])<<8 + uint16(fn.Code[i-1])))
	}
	fetchUint16 := func() int {
		i += 2
		return int(uint16(fn.Code[i-2])<<8 + uint16(fn.Code[i-1]))
	}

	op := fn.Code[i]
	i++
	if int(op) >= len(asm) {
		return i, fmt.Sprintf("bad op %d", op)
	}
	s := asm[op]
	switch op {
	case INT:
		n := fetchInt16()
		s += fmt.Sprintf(" %d", n)
	case VALUE:
		v := fn.Values[fetchUint16()]
		s += fmt.Sprintf(" %v", v)
	case LOAD, STORE, DYLOAD:
		idx := fetchUint8()
		s += " " + fn.Names[idx]
	case GLOBAL:
		idx := fetchUint16()
		s += " " + GlobalName(int(idx))
	case JUMP, TJUMP, FJUMP, AND, OR, Q_MARK, IN, EQJUMP, NEJUMP:
		j := fetchInt16()
		s += fmt.Sprintf(" %d", i+j)
	case CALLFUNC, CALLMETH:
		unnamed := fn.Code[i]
		i++
		named := int(fn.Code[i])
		i++
		spec := fn.Code[i : i+named]
		i += named
		s += ArgSpec{unnamed, spec, fn.Names}.String()[7:]
	}
	return i, s
}
