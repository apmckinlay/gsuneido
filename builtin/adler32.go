// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"hash/adler32"

	. "github.com/apmckinlay/gsuneido/core"
)

// The built-in hashes are Adler32, Md5, Sha1, Sha256.
// The implementations are very similar.
// Modifications to any of them should probably be done to the others.

type suAdler32 struct {
	ValueBase[suAdler32]
	hash hash.Hash32
}

var _ = builtin(Adler32, "(@args)")

func Adler32(th *Thread, as *ArgSpec, args []Value) Value {
	h := suAdler32{hash: adler32.New()}
	iter := NewArgsIter(as, args)
	k, v := iter()
	if v == nil {
		return h
	}
	for ; k == nil && v != nil; k, v = iter() {
		io.WriteString(h.hash, ToStr(v))
	}
	return h.value()
}

var _ Value = suAdler32{}

func (h suAdler32) Equal(other any) bool {
	return h == other
}

func (suAdler32) Lookup(_ *Thread, method string) Callable {
	return adler32Methods[method]
}

var adler32Methods = methods()

var _ = method(adler32_Update, "(string)")

func adler32_Update(this, arg Value) Value {
	io.WriteString(this.(suAdler32).hash, ToStr(arg))
	return this
}

var _ = method(adler32_Value, "()")

func adler32_Value(this Value) Value {
	return this.(suAdler32).value()
}

func (h suAdler32) value() Value {
	return IntVal(int(int32(h.hash.Sum32())))
}
