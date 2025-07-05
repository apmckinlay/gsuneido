// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/tokens"
	op "github.com/apmckinlay/gsuneido/core/opcodes"
	"github.com/apmckinlay/gsuneido/util/str"
)

// DisasmOps returns the disassembled byte code for fn
func DisasmOps(fn *SuFunc) string {
	var sb strings.Builder
	Disasm(fn, func(_ *SuFunc, nest, i int, s string, _ int) {
		in := strings.Repeat("> ", nest)
		fmt.Fprintf(&sb, "%s%5d: %s\n", in, i, s)
	})
	return sb.String()
}

type level struct {
	sp  int
	spi int
	cp  int
}

// DisasmMixed returns a listing of the source and its disassembled byte code
func DisasmMixed(fn *SuFunc, src string) string {
	var sb strings.Builder
	src = str.BeforeLast(src, "}")
	sp := fn.SrcBase
	spi := 0
	cp := 0
	printSrc := func(in, s string) {
		pre := fmt.Sprintf("%s%5d: ", in, sp)
		for _, line := range strings.Split(s, "\n") {
			if line = strings.TrimSpace(line); line != "" {
				fmt.Fprintf(&sb, "%s%s\n", pre, line)
				pre = in + "       "
			}
		}
	}
	nestPrev := 0
	stack := []level{}
	Disasm(fn, func(fn *SuFunc, nest, i int, s string, srcLim int) {
		if nest > nestPrev {
			stack = append(stack, level{sp: sp, spi: spi, cp: cp})
			sp = fn.SrcBase
			spi = 0
			cp = 0
		} else if nest < nestPrev {
			top := stack[len(stack)-1]
			sp, spi, cp = top.sp, top.spi, top.cp
			stack = stack[:len(stack)-1]
		}
		nestPrev = nest
		in := strings.Repeat("> ", nest)
		for i >= cp {
			if spi < len(fn.SrcPos) {
				ds := int(fn.SrcPos[spi])
				printSrc(in, src[sp:min(sp+ds, srcLim)])
				cp += int(fn.SrcPos[spi+1])
				sp += ds
				spi += 2
			} else { // last
				s := str.Subn(src, 0, srcLim)
				if len(stack) > 0 {
					// limit by parent sp
					top := stack[len(stack)-1]
					s = str.Subn(s, 0, top.sp)
				}
				if sp < len(s) {
					printSrc(in, s[sp:])
				}
				cp = math.MaxInt
				if nest == 0 {
					sp = len(src)
				}
			}
		}
		fmt.Fprintf(&sb, "%s        %5d: %s\n", in, i, s)
	})
	return sb.String()
}

type outfn func(fn *SuFunc, nest, i int, s string, srcLim int)

// Disasm calls out for each disassembled byte code instruction in fn
func Disasm(fn *SuFunc, out outfn) {
	disasm(0, fn, out)
}

func disasm(nest int, fn *SuFunc, out outfn) {
	d := &dasm{fn: fn, out: out, nest: nest}
	for d.i < len(fn.Code) {
		nestedfn := d.next()
		if nestedfn != nil && nestedfn.SrcBase > 0 {
			disasm(d.nest+1, nestedfn, d.out) // recursive
		}
	}
}

func Disasm1(fn *SuFunc, ip int) string {
	var output string
	d := &dasm{fn: fn, i: ip, out: func(_ *SuFunc, _, i int, s string, _ int) {
		output = fmt.Sprintf("%5d: %s", i, s)
	}}
	d.next()
	return output
}

type dasm struct {
	fn   *SuFunc
	out  outfn
	i    int
	nest int
}

func (d *dasm) next() *SuFunc {
	fetchUint8 := func() uint8 {
		d.i++
		return d.fn.Code[d.i-1]
	}
	fetchInt16 := func() int {
		d.i += 2
		return int(int16(uint16(d.fn.Code[d.i-2])<<8 + uint16(d.fn.Code[d.i-1])))
	}
	fetchUint16 := func() int {
		d.i += 2
		return int(uint16(d.fn.Code[d.i-2])<<8 + uint16(d.fn.Code[d.i-1]))
	}
	formatArgSpec := func(ai uint8) string {
		result := " "
		if int(ai) < len(StdArgSpecs) {
			result += StdArgSpecs[ai].String()[7:]
		} else {
			result += d.fn.ArgSpecs[int(ai)-len(StdArgSpecs)].String()[7:]
		}
		return result
	}

	ip := d.i
	oc := op.Opcode(d.fn.Code[ip])
	d.i++
	var nestedfn *SuFunc
	s := oc.String()
	switch oc {
	case op.Int:
		n := fetchInt16()
		s += fmt.Sprint(" ", n)
	case op.Value:
		v := d.fn.Values[fetchUint8()]
		s += fmt.Sprintf(" %v", v)
		if f, ok := v.(*SuFunc); ok {
			nestedfn = f
		}
	case op.Closure:
		f := d.fn.Values[fetchUint8()].(*SuFunc)
		nestedfn = f
	case op.Load, op.Store, op.Dyload:
		idx := fetchUint8()
		s += " " + d.fn.Names[idx]
	case op.LoadValue:
		localIdx := fetchUint8()
		valueIdx := fetchUint8()
		s += " " + d.fn.Names[localIdx] + fmt.Sprintf(" %v", d.fn.Values[valueIdx])
	case op.LoadLoad:
		idx1 := fetchUint8()
		idx2 := fetchUint8()
		s += " " + d.fn.Names[idx1] + " " + d.fn.Names[idx2]
	case op.ValueGet, op.ThisValue, op.GetValue:
		valueIdx := fetchUint8()
		s += fmt.Sprintf(" %v", d.fn.Values[valueIdx])
	case op.StorePop, op.ThisLoad, op.PopLoad:
		idx := fetchUint8()
		s += " " + d.fn.Names[idx]
	case op.LoadStore:
		idx := fetchUint8()
		s += " " + d.fn.Names[idx]
		fallthrough
	case op.GetPut:
		i := fetchUint8()
		op := tokens.Token(uint8(tokens.AddEq) + i>>1)
		s += " " + op.String()
		if i&1 == 1 {
			s += " retOrig"
		}
	case op.Global, op.Super:
		gn := fetchUint16()
		s += " " + Global.Name(gn)
	case op.Jump, op.JumpTrue, op.JumpFalse, op.And, op.Or, op.QMark, op.In, op.JumpIs,
		op.JumpIsnt, op.JumpLt, op.Catch:
		j := fetchInt16()
		s += fmt.Sprint(" ", d.i+j)
	case op.ForIn:
		idx := fetchUint8()
		j := fetchInt16()
		s += " " + d.fn.Names[idx] + fmt.Sprint(" ", d.i+j)
	case op.ForIn2:
		v1 := fetchUint8()
		v2 := fetchUint8()
		j := fetchInt16()
		s += " " + d.fn.Names[v1] + " " + d.fn.Names[v2] + fmt.Sprint(" ", d.i+j)
	case op.ForRange:
		j := fetchInt16()
		s += fmt.Sprint(" ", d.i+j)
	case op.ForRangeVar:
		idx := fetchUint8()
		j := fetchInt16()
		s += " " + d.fn.Names[idx] + fmt.Sprint(" ", d.i+j)
	case op.Try:
		j := fetchInt16()
		v := d.fn.Values[fetchUint8()]
		s += fmt.Sprintf(" %d %v", d.i+j-1, v)
	case op.GlobalCallFuncNoNil:
		gn := fetchUint16()
		s += " " + Global.Name(gn)
		fallthrough
	case op.CallFuncDiscard, op.CallFuncNoNil, op.CallFuncNilOk:
		s += formatArgSpec(fetchUint8())
	case op.ValueCallMethNoNil:
		valueIdx := fetchUint8()
		s += fmt.Sprintf(" %v", d.fn.Values[valueIdx])
		fallthrough
	case op.CallMethDiscard, op.CallMethNoNil, op.CallMethNilOk:
		s += formatArgSpec(fetchUint8())
	case op.InRange:
		orgTok := tokens.Token(fetchUint8())
		org := d.fn.Values[fetchUint8()]
		endTok := tokens.Token(fetchUint8())
		end := d.fn.Values[fetchUint8()]
		s += fmt.Sprint(" ", orgTok, " ", org, " ", endTok, " ", end)
	case op.ReturnMulti, op.PushReturn:
		n := fetchUint8()
		s += fmt.Sprint(" ", n)
	}
	srcLim := math.MaxInt
	if nestedfn != nil && nestedfn.SrcBase > 0 {
		srcLim = nestedfn.SrcBase
	}
	d.out(d.fn, d.nest, ip, s, srcLim)
	return nestedfn
}

func DisasmRaw(code string, fn func(i int)) {
	for i := 0; i < len(code); i++ {
		fn(i)
		switch op.Opcode(code[i]) {
		case op.Value, op.Closure, op.Load, op.Store, op.Dyload,
			op.GetPut, op.CallFuncDiscard, op.CallFuncNoNil, op.CallFuncNilOk,
			op.CallMethDiscard, op.CallMethNoNil, op.CallMethNilOk,
			op.ReturnMulti, op.PushReturn, op.ValueGet, op.ThisValue, op.GetValue,
			op.StorePop, op.ThisLoad, op.PopLoad:
			i++
		case op.Int, op.LoadStore, op.Global, op.Super,
			op.Jump, op.JumpTrue, op.JumpFalse, op.JumpIs, op.JumpIsnt,
			op.And, op.Or, op.QMark, op.In, op.Catch, op.LoadValue, op.LoadLoad,
			op.ValueCallMethNoNil:
			i += 2
		case op.ForIn, op.Try, op.GlobalCallFuncNoNil:
			i += 3
		case op.ForIn2:
			i += 4
		}
	}
}
