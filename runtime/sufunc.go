// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// SuFunc is a compiled Suneido function, method, or block
type SuFunc struct {
	ParamSpec

	// Nlocals is the number of parameters and local variables
	Nlocals uint8

	// Code is the actual byte code
	Code string

	// ArgSpecs used by calls in the code
	ArgSpecs []ArgSpec

	// ClassName is used to privatize dot params
	ClassName string

	// Id is a unique identifier for a function defining blocks
	Id uint32
	// OuterId is the Id of the outer SuFunc
	// It is used by interp to handle block return
	OuterId uint32

	// SrcPos contains pairs of source and code position deltas
	SrcPos string
	// SrcBase is the starting point for the SrcPos source deltas
	SrcBase int

	// cover is used for coverage tracking. nil means no tracking.
	// If len(cover) < len(Code) then bool coverage else counts.
	cover []uint16
}

// Value interface (mostly handled by ParamSpec) --------------------

var _ Value = (*SuFunc)(nil)

func (f *SuFunc) Call(t *Thread, this Value, as *ArgSpec) Value {
	args := t.Args(&f.ParamSpec, as)
	for i, flag := range f.Flags {
		if flag&DotParam == DotParam {
			name := f.Names[i]
			if flag&PubParam == PubParam {
				name = str.Capitalize(name)
			} else { // privatize
				name = f.ClassName + "_" + name
			}
			this.Put(t, SuStr(name), args[i])
		}
	}
	return t.Start(f, this)
}

func (f *SuFunc) Type() types.Type {
	if f.OuterId != 0 {
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
	f.coverToOb(ob, counts)
	assert.That(counts == (len(f.cover) >= len(f.Code)))
	for _, v := range f.Values {
		if g, ok := v.(*SuFunc); ok {
			g.getCoverage(ob, counts) // RECURSE
		}
	}
}

func (f *SuFunc) coverToOb(ob *SuObject, counts bool) {
	for i, c := range f.cover {
		if c != 0 {
			if counts {
				srcpos := IntVal(f.CodeToSrcPos(i))
				val := IntVal(int(c))
				if x := ob.getIfPresent(srcpos); x != nil {
					val = OpAdd(x, val)
				}
				ob.set(srcpos, val)
			} else {
				for j := 0; j < 16; j++ {
					if c&(1<<j) != 0 {
						ob.Set(IntVal(f.CodeToSrcPos(i*16+j)), True)
					}
				}
			}
		}
	}
}
