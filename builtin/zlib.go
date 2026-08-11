// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"compress/zlib"
	"io"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type suZlib struct {
	staticClass[suZlib]
}

func init() {
	Global.Builtin("Zlib", &suZlib{})
}

func (*suZlib) String() string {
	return "Zlib /* builtin class */"
}

func (z *suZlib) Equal(other any) bool {
	return z == other
}

func (*suZlib) Lookup(_ *Thread, method string) Value {
	return zlibMethods[method]
}

var zlibMethods = methods("zlib")

var _ = staticMethod(zlib_Compress, "(string :string) :string")

func zlib_Compress(arg Value) Value {
	s := ToStr(arg)
	var sb strings.Builder
	w := zlib.NewWriter(&sb)
	n, err := io.WriteString(w, s)
	if err != nil {
		panic("Zlib.Compress: " + err.Error())
	}
	assert.That(n == len(s))
	err = w.Close()
	if err != nil {
		panic("Zlib.Compress: " + err.Error())
	}
	return SuStr(sb.String())
}

var _ = staticMethod(zlib_Uncompress, "(string :string) :string")

func zlib_Uncompress(arg Value) Value {
	data := ToStr(arg)
	r, err := zlib.NewReader(strings.NewReader(data))
	if err != nil {
		panic("Zlib.Uncompress: " + err.Error())
	}
	var sb strings.Builder
	n, err := io.Copy(&sb, r)
	if err != nil {
		panic("Zlib.Uncompress: " + err.Error())
	}
	r.Close()
	assert.That(int(n) == len(sb.String()))
	return SuStr(sb.String())
}

var _ = staticMethod(zlib_Members, "() :object")

func zlib_Members() Value {
	return methodList(zlibMethods)
}
