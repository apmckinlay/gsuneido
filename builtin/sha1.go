// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/sha1"

	. "github.com/apmckinlay/gsuneido/core"
)

type suSha1 struct {
	ValueBase[suSha1]
	hash hash.Hash
}

// The built-in hashes are Adler32, Md5, Sha1, Sha256.
// The implementations are very similar.
// Modifications to any of them should probably be done to the others.

var _ = builtin(Sha1, "(@args)")

func Sha1(th *Thread, as *ArgSpec, args []Value) Value {
	h := suSha1{hash: sha1.New()}
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

var _ Value = suSha1{}

func (h suSha1) Equal(other any) bool {
	return h == other
}

func (suSha1) Lookup(_ *Thread, method string) Value {
	return sha1Methods[method]
}

var sha1Methods = methods()

var _ = method(Sha1_Update, "(string)")

func Sha1_Update(this, arg Value) Value {
	io.WriteString(this.(suSha1).hash, ToStr(arg))
	return this
}

var _ = method(Sha1_Value, "()")

func Sha1_Value(this Value) Value {
	return this.(suSha1).value()
}

func (h suSha1) value() Value {
	var buf [sha1.Size]byte
	return SuStr(string(h.hash.Sum(buf[0:0])))
}
