// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/md5"

	. "github.com/apmckinlay/gsuneido/runtime"
)

// The built-in hashes are Adler32, Md5, Sha1, Sha256.
// The implementations are very similar.
// Modifications to any of them should probably be done to the others.

type suMd5 struct {
	ValueBase[suMd5]
	hash hash.Hash
}

var _ = builtinRaw("Md5(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		h := suMd5{hash: md5.New()}
		iter := NewArgsIter(as, args)
		k, v := iter()
		if v == nil {
			return h
		}
		for ; k == nil && v != nil; k, v = iter() {
			io.WriteString(h.hash, ToStr(v))
		}
		return h.value()
	})

var _ Value = suMd5{}

func (h suMd5) Equal(other any) bool {
	return h == other
}

func (suMd5) Lookup(_ *Thread, method string) Callable {
	return md5Methods[method]
}

var md5Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(suMd5).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(suMd5).value()
	}),
}

func (h suMd5) value() Value {
	var buf [md5.Size]byte
	return SuStr(string(h.hash.Sum(buf[0:0])))
}
