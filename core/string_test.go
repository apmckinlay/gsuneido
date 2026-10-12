// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/strp"
)

// stringTypes lists every string Value type with a constructor.
// SuExcept embeds SuStr so it shares its Equal/Compare/Hash.
var stringTypes = []struct {
	name string
	new  func(string) Value
}{
	{"SuStr", func(s string) Value { return SuStr(s) }},
	{"SuStrP", func(s string) Value { return NewSuStrP(s) }},
	{"SuConcat", func(s string) Value { return NewSuConcat().Add(s) }},
	{"SuExcept", func(s string) Value { return &SuExcept{SuStr: SuStr(s)} }},
}

var stringContents = func() []string {
	strs := []string{"", "a", "ab", "abc",
		"hello\x00world", "the quick brown fox jumps over the lazy dog"}
	strs = append(strs, strings.Repeat("x", strp.MaxLen))
	return strs
}()

func TestStringEqual(t *testing.T) {
	a := assert.T(t)
	// equal contents are equal in both directions, across all type pairs
	for _, ca := range stringTypes {
		for _, cb := range stringTypes {
			for _, s := range stringContents {
				x, y := ca.new(s), cb.new(s)
				a.Msg(ca.name, cb.name, s).This(x.Equal(y)).Is(true)
				a.Msg(ca.name, cb.name, s).This(y.Equal(x)).Is(true)
			}
		}
	}
	// different contents are not equal, in both directions
	for _, ca := range stringTypes {
		for _, cb := range stringTypes {
			for _, s := range stringContents {
				var s2 string
				if len(s) > 1 {
					s2 = s[:len(s)-1]
				} else {
					s2 = s + "x"
				}
				x, y := ca.new(s), cb.new(s2)
				a.Msg(ca.name, cb.name, s).This(x.Equal(y)).Is(false)
				a.Msg(ca.name, cb.name, s).This(y.Equal(x)).Is(false)
			}
		}
	}
	// non-string values are never equal
	for _, ca := range stringTypes {
		x := ca.new("abc")
		a.Msg(ca.name).This(x.Equal(SuInt16(5))).Is(false)
		a.Msg(ca.name).This(x.Equal(nil)).Is(false)
		a.Msg(ca.name).This(x.Equal("abc")).Is(false) // plain string is not a Value
	}
}

func TestStringCompare(t *testing.T) {
	a := assert.T(t)
	for _, ca := range stringTypes {
		for _, cb := range stringTypes {
			a.This(ca.new("abc").Compare(cb.new("abc"))).Is(0)
			a.This(ca.new("abc").Compare(cb.new("abd"))).Is(-1)
			a.This(ca.new("abd").Compare(cb.new("abc"))).Is(1)
			a.This(ca.new("").Compare(cb.new("a"))).Is(-1)
			// strings order after numbers and before objects
			a.Msg(ca.name, cb.name).That(ca.new("abc").Compare(SuInt16(5)) > 0)
			a.Msg(ca.name, cb.name).That(ca.new("abc").Compare(&SuObject{}) < 0)
		}
	}
}

func TestStringHash(t *testing.T) {
	a := assert.T(t)
	contents := append(append([]string{}, stringContents...),
		strings.Repeat("x", 100), strings.Repeat("abc", 80))
	for _, s := range contents {
		x0 := stringTypes[0].new(s)
		want := x0.Hash()
		a.Msg(stringTypes[0].name, len(s)).This(x0.Hash2()).Is(want)
		for _, st := range stringTypes[1:] {
			x := st.new(s)
			a.Msg(st.name, len(s)).This(x.Hash()).Is(want)
			a.Msg(st.name, len(s)).This(x.Hash2()).Is(want)
		}
	}
}

func TestInternStrVal(t *testing.T) {
	// interned names share storage instead of copying into SuStrP
	a := InternStrVal("hello").(*SuStrP)
	b := InternStrVal("hello").(*SuStrP)
	assert.That(a == b)

	// too long to intern falls back to SuStr
	long := InternStrVal(strings.Repeat("x", strp.MaxLen+1))
	_, ok := long.(SuStr)
	assert.That(ok)
}
