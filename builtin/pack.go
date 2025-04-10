// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"compress/zlib"
	"fmt"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(packSize, "(value)")

func packSize(arg Value) Value {
	return IntVal(PackSize(arg))
}

var _ = builtin(pack, "(value, zip=false)")

func pack(arg, zip Value) Value {
	if !ToBool(zip) {
		return SuStr(PackValue(arg))
	}
	var dst strings.Builder
	w := zlib.NewWriter(&dst)
	if err := PackTo(arg, w); err != nil {
		panic("Pack: " + err.Error())
	}
	if err := w.Close(); err != nil {
		panic("Pack: " + err.Error())
	}
	return SuStr(dst.String())
}

var _ = builtin(unpack, "(string, zip=false)")

func unpack(arg, zip Value) Value {
	defer func() {
		if e := recover(); e != nil {
			panic("Unpack: " + fmt.Sprint(e))
		}
	}()
	if !ToBool(zip) {
		return Unpack(ToStr(arg))
	}
	src := strings.NewReader(ToStr(arg))
	r, err := zlib.NewReader(src)
	if err != nil {
		panic(err)
	}
	return UnpackFrom(r)
}
