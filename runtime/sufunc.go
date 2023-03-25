// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/opcodes"
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// SuFunc is a compiled Suneido function, method, or block
type SuFunc struct {

	// Code is the actual byte code
	Code string

	// ClassName is used to privatize dot params
	ClassName string

	// SrcPos contains pairs of source and code position deltas
	SrcPos string

	// ArgSpecs used by calls in the code
	ArgSpecs []ArgSpec

	// cover is used for coverage tracking. nil means no tracking.
	// If len(cover) < len(Code) then bool coverage else counts.
	cover []uint16

	ParamSpec

	// SrcBase is the starting point for the SrcPos source deltas
	SrcBase int

	// Nlocals is the number of parameters and local variables
	Nlocals uint8

	IsBlock bool
}

// Value interface (mostly handled by ParamSpec) --------------------

var _ Value = (*SuFunc)(nil)

func (f *SuFunc) Call(th *Thread, this Value, as *ArgSpec) Value {
	args := th.Args(&f.ParamSpec, as)
	for i, flag := range f.Flags {
		if flag&DotParam == DotParam {
			name := f.Names[i]
			if flag&PubParam == PubParam {
				name = str.Capitalize(name)
			} else { // privatize
				name = f.ClassName + "_" + name
			}
			this.Put(th, SuStr(name), args[i])
		}
	}
	return th.invoke(f, this)
}

func (f *SuFunc) Type() types.Type {
	if f.IsBlock {
		return types.Block
	}
	return types.Function
}

// SuFuncMethods is initialized by the builtin package
var SuFuncMethods Methods

func (*SuFunc) Lookup(_ *Thread, method string) Callable {
	if m, ok := ParamsMethods[method]; ok {
		return m
	}
	return SuFuncMethods[method]
}

func (f *SuFunc) String() string {
	s := ""
	if f.Name != "" && f.Name != "?" {
		s = f.Name + " "
	}
	s += "/* " + str.Opt(f.Lib, " ")
	if f.ClassName != "" {
		s += "method"
	} else if f.IsBlock {
		s += "block"
	} else {
		s += "function"
	}
	s += " */"
	return s
}

func (f *SuFunc) CodeToSrcPos(ip int) int {
	sp := f.SrcBase
	cp := 0
	for i := 0; i < len(f.SrcPos); i += 2 {
		prev := sp
		sp += int(f.SrcPos[i])
		cp += int(f.SrcPos[i+1])
		if cp > ip {
			return prev
		}
	}
	return sp
}

// coverage ---------------------------------------------------------

func (f *SuFunc) StartCoverage(count bool) {
	if len(f.Code) == 0 {
		return
	}
	f.startCoverage(count)
	for _, v := range f.Values {
		if g, ok := v.(*SuFunc); ok {
			g.StartCoverage(count) // RECURSE
		}
		if c, ok := v.(*SuClass); ok {
			c.StartCoverage(count)
		}
	}
}

func (f *SuFunc) startCoverage(count bool) {
	n := len(f.Code)
	if !count {
		n /= 16
	}
	f.cover = make([]uint16, n+1)
}

func (f *SuFunc) StopCoverage() *SuObject {
	ob := &SuObject{}
	f.getCoverage(ob, len(f.cover) >= len(f.Code))
	f.cover = nil
	return ob
}

func (f *SuFunc) getCoverage(ob *SuObject, counts bool) {
	if len(f.Code) == 0 || len(f.cover) == 0 {
		return
	}
	assert.That(counts == (len(f.cover) >= len(f.Code)))
	f.coverToOb(ob, counts)
	for _, v := range f.Values {
		if g, ok := v.(*SuFunc); ok {
			g.getCoverage(ob, counts) // RECURSE
		}
	}
}

func (f *SuFunc) coverToOb(ob *SuObject, counts bool) {
	DisasmRaw(f.Code, func(i int) {
		if f.Code[i] != byte(opcodes.Cover) {
			return
		}
		var v Value
		if counts {
			v = IntVal(int(f.cover[i]))
		} else {
			if f.cover[i>>4]&(1<<(i&15)) == 0 {
				v = False
			} else {
				v = True
			}
		}
		ob.Set(IntVal(f.CodeToSrcPos(i)), v)
	})
}
