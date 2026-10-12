// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/strp"
)

// common code for string types
// (SuStr, SuStrP, SuConcat, SuExcept)
// All string types are designed to convert to a Go string without allocation.

// IsEmptyStr is the correct way to check for an empty string.
// DO NOT use == EmptyStr
func IsEmptyStr(v Value) bool {
	s, ok := v.ToStr()
	return ok && len(s) == 0
}

// StrOrEmpty returns the string if m is a string, else "".
func StrOrEmpty(m Value) string {
	s, _ := m.ToStr()
	return s
}

// stringCall implements s(ob, ...) being treated as ob[s](...)
func stringCall(th *Thread, as *ArgSpec, method string) Value {
	base := th.sp - int(as.Nargs)
	args := th.stack[base : base+int(as.Nargs)]
	k, v := NewArgsIter(as, args)()
	if v == nil || k != nil {
		panic("string call requires 'this' argument")
	}
	fn := v.Lookup(th, method)
	if fn == nil {
		panic("method not found: " + ErrType(v) + "." + method)
	}
	return fn.Call(th, v, as.DropFirst())
}

func StrVal1(c byte) Value {
	return SuStrP1(c)
}

// StrVal converts a Go string to a Suneido Value.
func StrVal(s string) Value {
	if len(s) >= strp.MaxFromGo {
		return SuStr(s)
	}
	return NewSuStrP(s)
}

// StrCat concatenates Go strings to a Suneido Value.
func StrCat(parts ...string) Value {
	n := 0
	for _, s := range parts {
		n += len(s)
	}
	if n > strp.MaxLen {
		return SuStr(strings.Join(parts, ""))
	}
	return CatSuStrP(parts...)
}

// iterator

type stringIter struct {
	s string
	i int
	MayLock
}

func (si *stringIter) Next() Value {
	if si.Lock() {
		defer si.Unlock()
	}
	si.i++
	if si.i > len(si.s) {
		return nil
	}
	return StrVal1(si.s[si.i-1])
}

func (si *stringIter) Dup() Iter {
	return &stringIter{s: si.s}
}

func (si *stringIter) Infinite() bool {
	return false
}

func (si *stringIter) Instantiate() *SuObject {
	InstantiateMax(len(si.s))
	list := make([]Value, len(si.s))
	for i := range len(si.s) {
		list[i] = StrVal1(si.s[i])
	}
	return NewSuObject(list)
}
