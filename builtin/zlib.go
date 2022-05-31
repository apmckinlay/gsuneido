// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"compress/zlib"
	"io"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type suZlib struct {
	SuBuiltin
}

func init() {
	Global.Builtin("Zlib", &suZlib{})
}

func (*suZlib) String() string {
	return "Zlib /* builtin class */"
}

func (z *suZlib) Lookup(_ *Thread, method string) Callable {
	return zlibMethods[method]
}

var zlibMethods = Methods{
	"Compress": method1("(string)", func(_, arg Value) Value {
		s := ToStr(arg)
		var b strings.Builder
		w := zlib.NewWriter(&b)
		n, err := io.WriteString(w, s)
		if err != nil {
			panic("Zlib.Compress: " + err.Error())
		}
		assert.That(n == len(s))
		err = w.Close()
		if err != nil {
			panic("Zlib.Compress: " + err.Error())
		}
		return SuStr(b.String())
	}),
	"Uncompress": method1("(string)", func(_, arg Value) Value {
		data := ToStr(arg)
		r, err := zlib.NewReader(strings.NewReader(data))
		if err != nil {
			panic("Zlib.Uncompress: " + err.Error())
		}
		var b strings.Builder
		n, err := io.Copy(&b, r)
		if err != nil {
			panic("Zlib.Uncompress: " + err.Error())
		}
		r.Close()
		assert.That(int(n) == len(b.String()))
		return SuStr(b.String())
	}),
}

func (z *suZlib) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Zlib")
}
