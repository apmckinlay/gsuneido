// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"sync"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/time/rate"
)

const nonceSize = 8
const tokenSize = 16

var tokens = make(map[string]bool)
var tokensLock sync.Mutex

// authLimiter limits the rate of authentication attempts
var authLimiter = rate.NewLimiter(rate.Limit(4), 1) // ???
var authContext = context.Background()

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
	tokens[s] = false // not old
	return s
}

// AuthToken verifies that the given token is valid.
// It is used by dbms.Auth
func AuthToken(s string) bool {
	authLimiter.Wait(authContext)
	tokensLock.Lock()
	defer tokensLock.Unlock()
	if _, ok := tokens[s]; ok {
		delete(tokens, s)
		return true
	}
	return false
}

func AuthUser(th *Thread, s, nonce string) bool {
	authLimiter.Wait(authContext)
	if nonce == "" {
		return false
	}
	user := str.BeforeFirst(s, "\x00")
	passhash := getPassHash(th, user)
	hash := sha1.Sum([]byte(nonce + passhash)) //TODO upgrade sha1
	t := user + "\x00" + string(hash[:])
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
	args := SuObjectOf(SuStr("users"))
	args.Set(SuStr("user"), SuStr(user))
	row, hdr, _ := dbms.Get(th, args, Only)
	if row == nil {
		return ""
	}
	hash := Unpack(row.GetRaw(hdr, "passhash"))
	return string(hash.(SuStr))
}
