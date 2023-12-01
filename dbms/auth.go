// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"crypto/rand"
	"crypto/sha1"
	"io"
	"sync"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

const nonceSize = 8
const tokenSize = 16

var tokens = make(map[string]bool)
var tokensLock sync.Mutex

func Nonce() string {
	buf := make([]byte, nonceSize)
	if _, err := rand.Read(buf); err != nil {
		panic("Nonce: " + err.Error())
	}
	return hacks.BStoS(buf)
}

// Token generates a random token.
// It is used by dbms.Token
func Token() string {
	buf := make([]byte, tokenSize)
	if _, err := rand.Read(buf); err != nil {
		panic("Token: " + err.Error())
	}
	s := hacks.BStoS(buf)
	tokensLock.Lock()
	defer tokensLock.Unlock()
	tokens[s] = true
	return s
}

// AuthToken verifies that the given token is valid.
// It is used by dbms.Auth
func AuthToken(s string) bool {
	tokensLock.Lock()
	defer tokensLock.Unlock()
	if tokens[s] {
		delete(tokens, s)
		return true
	}
	return false
}

func AuthUser(th *Thread, s, nonce string) bool {
	if nonce == "" {
		return false
	}
	user := str.BeforeFirst(s, "\x00")
	hash := sha1.New()
	passhash := getPassHash(th, user)
	io.WriteString(hash, nonce+passhash)
	t := user + "\x00" + string(hash.Sum(nil))
	return s == t
}

func getPassHash(th *Thread, user string) (result string) {
	defer func() {
		if e := recover(); e != nil {
			result = ""
		}
	}()
	dbms := th.Dbms()
	if u, ok := dbms.(*DbmsUnauth); ok {
		dbms = u.dbms
	}
	query := "users where user = " + SuStr(user).String() // handle quotes
	row, hdr, _ := dbms.Get(th, query, Next)
	if row == nil {
		return ""
	}
	hash := Unpack(row.GetRaw(hdr, "passhash"))
	return string(hash.(SuStr))
}
