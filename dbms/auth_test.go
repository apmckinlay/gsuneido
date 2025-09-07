// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"crypto/sha1"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
	"golang.org/x/time/rate"
)

func TestToken(*testing.T) {
	assert.False(AuthToken("invalid"))
	tok1 := Token()
	tok2 := Token()
	assert.True(AuthToken(tok1))
	assert.False(AuthToken(tok1))
	assert.True(AuthToken(tok2))
	assert.False(AuthToken(tok2))
}

func TestAuthUser(*testing.T) {
	user := "fred"
	passhash := "123"
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	db19.MakeSuTran = func(ut *db19.UpdateTran) *SuTran {
		return NewSuTran(nil, true)
	}
	qry.DoAdmin(db, "create users (user, passhash) key(user)", nil)
	ut := db.NewUpdateTran()
	qry.DoAction(nil, ut, "insert { user: 'fred', passhash: '123' } into users")
	ut.Commit()

	dbms := NewDbmsLocal(db)
	GetDbms = func() IDbms { return dbms }
	nonce := Nonce()
	hash := sha1.Sum([]byte(nonce + passhash))
	s := user + "\x00" + string(hash[:])
	assert.True(AuthUser(&Thread{}, s, nonce))
}

func TestAuthRateLimit(t *testing.T) {
	originalLimiter := authLimiter
	defer func() { authLimiter = originalLimiter }()

	// Set a restrictive rate limiter for testing: 1 request per 100ms with burst of 1
	authLimiter = rate.NewLimiter(rate.Every(100*time.Millisecond), 1)

	token1 := Token()
	token2 := Token()

	// First attempt should succeed quickly (within burst limit)
	start := time.Now()
	assert.T(t).That(AuthToken(token1))
	duration1 := time.Since(start)

	// Second attempt should be rate limited and take at least 100ms
	start = time.Now()
	assert.T(t).That(AuthToken(token2))
	duration2 := time.Since(start)

	// Verify first call was fast (< 50ms) and second was delayed (>= 90ms)
	assert.T(t).That(duration1 < 50*time.Millisecond)
	assert.T(t).That(duration2 >= 90*time.Millisecond)
}
