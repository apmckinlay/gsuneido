package runtime

import (
	"fmt"
	"io"
	"strings"

	op "github.com/apmckinlay/gsuneido/runtime/opcodes"
	"github.com/apmckinlay/gsuneido/util/str"
)

func Disasm(w io.Writer, fn *SuFunc) {
	disasm(w, fn, 0)
}

func disasm(w io.Writer, fn *SuFunc, indent int) {
	var s string
	in := strings.Repeat("    ", indent)
	for i := 0; i < len(fn.Code); {
		j := i
		i, s = disasm1(fn, i, indent)
		fmt.Fprintf(w, "%s%d: %s\n", in, j, s)
	}
	// fmt.Fprintf(w, "%d:\n", len(fn.Code))
}

func Disasm1(fn *SuFunc, i int) (int, string) {
	return disasm1(fn, i, 0)
}

func disasm1(fn *SuFunc, i int, indent int) (int, string) {
	fetchUint8 := func() uint8 {
		i++
		return fn.Code[i-1]
	}
	fetchInt16 := func() int {
		i += 2
		return int(int16(uint16(fn.Code[i-2])<<8 + uint16(fn.Code[i-1])))
	}
	fetchUint16 := func() int {
		i += 2
		return int(uint16(fn.Code[i-2])<<8 + uint16(fn.Code[i-1]))
	}
	nested := func (fn *SuFunc) string {
		var sb strings.Builder
		sb.WriteString("\n")
		disasm(&sb, fn, indent + 1)
		return strings.TrimRight(sb.String(), "\n")
	}

	oc := op.Opcode(fn.Code[i])
	i++
	s := oc.String()
	switch oc {
	case op.Int:
		n := fetchInt16()
		s += fmt.Sprintf(" %d", n)
	case op.Value:
		v := fn.Values[fetchUint8()]
		s += fmt.Sprintf(" %v", v)
		if fn,ok := v.(*SuFunc); ok {
			s += nested(fn)
		}
	case op.Block:
		fn := fn.Values[fetchUint8()].(*SuFunc)
		s += nested(fn)
	case op.Load, op.Store, op.Dyload:
		idx := fetchUint8()
		s += " " + fn.Names[idx]
	case op.Global, op.Super:
		gn := fetchUint16()
		s += " " + Global.Name(gn)
	case op.Jump, op.JumpTrue, op.JumpFalse, op.And, op.Or, op.QMark, op.In, op.JumpIs,
		op.JumpIsnt, op.Catch:
		j := fetchInt16()
		s += fmt.Sprintf(" %d", i+j)
	case op.ForIn:
		j := fetchInt16()
		idx := fetchUint8()
		s += " " + fn.Names[idx] + fmt.Sprintf(" %d", i+j-1)
	case op.Try:
		j := fetchInt16()
		v := fn.Values[fetchUint8()]
		s += fmt.Sprintf(" %d %v", i+j-1, v)
	case op.CallFunc, op.CallMeth:
		ai := int(fetchUint8())
		s += " "
		if ai < len(StdArgSpecs) {
			s += StdArgSpecs[ai].String()[7:]
		} else {
			s += fn.ArgSpecs[ai-len(StdArgSpecs)].String()[7:]
		}
	}
	return i, s
}

func DisasmMixed(w io.Writer, fn *SuFunc, src string) {
	sp := fn.SrcBase
	printSrc := func (s string) {
		fmt.Fprintf(w, "%d: %s\n", sp,
			strings.TrimSpace(str.BeforeFirst(str.BeforeFirst(s, "}"), "//")))
	}
	src = strings.ReplaceAll(src, "\n", " ")
	cp := 0
	ip := 0
	var s string
	for i := 0; i < len(fn.SrcPos); i += 2 {
		ds := int(fn.SrcPos[i])
		printSrc(src[sp:sp+ds])
		cp += int(fn.SrcPos[i+1])
		for ip < cp {
			fmt.Fprintf(w, "\t%d: ", ip)
			ip, s = Disasm1(fn, ip)
			fmt.Fprintln(w, s)
		}
		sp += ds
	}
	printSrc(src[sp:])
	for ip < len(fn.Code) {
		fmt.Fprintf(w, "\t%d: ", ip)
		ip, s = Disasm1(fn, ip)
		fmt.Fprintln(w, s)
	}
}
