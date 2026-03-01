// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSuFuncString(t *testing.T) {
	sf := SuFunc{}
	sf.Flags = make([]Flag, 8)
	assert.T(t).This(sf.Params()).Is("()")
	sf.Nparams = 3
	sf.Names = []string{"a", "b", "c"}
	assert.T(t).This(sf.Params()).Is("(a,b,c)")
	sf.Names = []string{"a", "b", "c"}
	sf.Ndefaults = 1
	sf.Values = []Value{SuInt(123)}
	assert.T(t).This(sf.Params()).Is("(a,b,c=123)")
}

func TestVarName(t *testing.T) {
	assert := assert.T(t).This
	// Test without shared variables
	ps := &ParamSpec{
		Names:   []string{"a", "b", "c"},
		Nstack: 3, // no shared vars
	}
	assert(ps.VarName(0)).Is("a")
	assert(ps.VarName(1)).Is("b")
	assert(ps.VarName(2)).Is("c")

	// Test with shared variables
	// Names = [local0, local1, shared0, shared1]
	// Nshared = 2 (locals are indices 0-1)
	ps = &ParamSpec{
		Names:   []string{"local0", "local1", "shared0", "shared1"},
		Nstack: 2,
	}
	assert(ps.VarName(0)).Is("local0")
	assert(ps.VarName(1)).Is("local1")
	assert(ps.VarName(192)).Is("shared0")
	assert(ps.VarName(193)).Is("shared1")
}

func TestNshared(t *testing.T) {
	assert := assert.T(t).This
	// No shared variables
	ps := &ParamSpec{
		Names:   []string{"a", "b", "c"},
		Nstack: 3,
	}
	assert(len(ps.Names[ps.Nstack:])).Is(0)

	// With shared variables
	ps = &ParamSpec{
		Names:   []string{"local", "shared1", "shared2"},
		Nstack: 1,
	}
	assert(ps.Nstack).Is(uint8(1))
	assert(ps.Names[ps.Nstack:]).Is([]string{"shared1", "shared2"})
}
