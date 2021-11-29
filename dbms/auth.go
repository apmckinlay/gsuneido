// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"crypto/rand"
	"crypto/sha1"
	"io"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

const nonceSize = 8
const tokenSize = 16

var tokens = make(map[string]bool)

func Nonce() string {
	buf := make([]byte, nonceSize)
	if _, err := rand.Read(buf); err != nil {
		panic("Nonce: " + err.Error())
	}
	return hacks.BStoS(buf)
}

func Token() string {
	buf := make([]byte, tokenSize)
	if _, err := rand.Read(buf); err != nil {
		panic("Token: " + err.Error())
	}
	s := hacks.BStoS(buf)
	tokens[s] = true
	return s
}

func AuthToken(s string) bool {
	if tokens[s] {
		delete(tokens, s)
		return true
	}
	return false
}

func AuthUser(s, nonce string) bool {
	if nonce == "" {
		return false
	}
	user := str.BeforeFirst(s, "\x00")
	hash := sha1.New()
	passhash := getPassHash(user)
	io.WriteString(hash, nonce+passhash)
	t := user+"\x00"+string(hash.Sum(nil))
	return s == t
}

func getPassHash(user string) (result string) {
	defer func() {
		if e := recover(); e != nil {
			result = ""
		}
	}()
	dbms := GetDbms()
	query := "users where user = '" + user + "'"
	row, hdr, _ := dbms.Get(query, Next)
	if row == nil {
		return ""
	}
	hash := Unpack(row.GetRaw(hdr, "passhash"))
	return string(hash.(SuStr))
}
