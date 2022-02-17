// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"crypto/sha1"
	"io"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
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
	db, _ := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	db19.MakeSuTran = func(ut *db19.UpdateTran) *SuTran {
		return NewSuTran(nil, true)
	}
	qry.DoAdmin(db, "create users (user, passhash) key(user)", nil)
	ut := db.NewUpdateTran()
	qry.DoAction(nil, ut, "insert { user: 'fred', passhash: '123' } into users", nil)
	ut.Commit()

	dbms := NewDbmsLocal(db)
	GetDbms = func() IDbms { return dbms }
	nonce := Nonce()
	hash := sha1.New()
	io.WriteString(hash, nonce+passhash)
	s := user + "\x00" + string(hash.Sum(nil))
	assert.True(AuthUser(&Thread{}, s, nonce))
}
