// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/sha256"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type suSha256 struct {
	ValueBase[suSha256]
	hash hash.Hash
}

// The built-in hashes are Adler32, Md5, Sha1, Sha256.
// The implementations are very similar.
// Modifications to any of them should probably be done to the others.

var _ = builtinRaw("Sha256(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		h := suSha256{hash: sha256.New()}
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

var _ Value = suSha256{}

func (h suSha256) Equal(other interface{}) bool {
	h2, ok := other.(suSha256)
	return ok && h == h2
}

func (suSha256) Lookup(_ *Thread, method string) Callable {
	return sha256Methods[method]
}

var sha256Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(suSha256).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(suSha256).value()
	}),
}

func (h suSha256) value() Value {
	var buf [sha256.Size]byte
	return SuStr(string(h.hash.Sum(buf[0:0])))
}
