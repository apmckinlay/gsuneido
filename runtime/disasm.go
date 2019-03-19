package runtime

import (
	"fmt"
	"io"

	op "github.com/apmckinlay/gsuneido/runtime/opcodes"
)

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
	case op.Block:
		fetchUint8()
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
