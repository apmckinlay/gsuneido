// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build neworder

package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestConvertValue(t *testing.T) {
	test := func(val Value) {
		s := Pack(val.(Packable))
		assert.This(Unpack(s)).Is(val)
		buf := []byte(s)
		ConvertValue(xlat, buf, s)
		assert.This(Unpack(string(buf))).Is(val)
	}
	test(True)
	test(False)
	test(Zero)
	test(One)
	test(IntVal(123456789))
	test(EmptyStr)
	test(Now())
	r := NewSuRecord()
	test(r)
	r.Set(Zero, One)
	test(r)
	r.Set(One, True)
	test(r)
	ob := &SuObject{}
	test(ob)
	ob.Add(One)
	test(ob)
	ob.Add(False)
	test(ob)
	ob.Set(Zero, One)
	test(ob)
	ob.Set(One, True)
	test(ob)
	ob.Set(SuStr("rec"), r)
	test(ob)
}
