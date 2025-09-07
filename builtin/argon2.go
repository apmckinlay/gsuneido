// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"golang.org/x/crypto/argon2"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Argon2id, "(password, salt)")

func Argon2id(a, b Value) Value {
	password := ToStr(a)
	salt := ToStr(b)
	return SuStr(argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 1, 32))
}
