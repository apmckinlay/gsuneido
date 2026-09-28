// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"context"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/atomics"
	"golang.org/x/time/rate"
)

/*
Auth Flow
---------
1. client connects to server
2. client connection starts out with no permissions
3. the server rejects all requests except for loading code
4. clients gets name and password from user using login dialog
5. client calls Database.Auth
6. auth request goes to server
7. server rate limits Auth calls, max of 3 tries, then disconnect
8. server calls application defined Auth function
9. app Auth grants itself read permission to user/password table
10. app Auth verifies user and password, return false on fail
11. app Auth grants itself read permission to permission table
12. app Auth sets up permissions (for tables and ServerEval whitelist)
13. for developers, the app Auth can grant permission to alter permissions later
14. if app Auth fails (e.g. throws exception) permissions are cleared (defer)
15. once app Auth returns true, permissions are moved to serverConn
*/

func Unauth(dbms *DbmsLocal) IDbms {
	return &DbmsUnauth{dbms: dbms}
}

// StandaloneDbms is the current dbms for standalone mode.
// It starts as an unauth wrapper and is replaced by DbmsLocal on auth.
// Needs atomic because Auth runs on a different thread than GetDbms callers.
var StandaloneDbms atomics.Intfc[IDbms]

// DbmsUnauth is a wrapper for DbmsLocal for unauthorized client connections.
// Only allows LibGet, Libraries, SessionId, and Use
type DbmsUnauth struct {
	dbms *DbmsLocal
}

var _ IDbms = (*DbmsUnauth)(nil)

const notauth = "not authorized"

func (du *DbmsUnauth) Admin(string, *Sviews) {
	panic(notauth)
}

func (du *DbmsUnauth) Auth(th *Thread, data Value) (result bool) {
	// This is only used by standalone mode.
	// Give the thread the unwrapped dbms since the app Auth may query the db.
	// Restore it if the app Auth fails or throws.
	prev := th.SetDbms(du.dbms)
	defer func() {
		if !result {
			th.SetDbms(prev)
		}
	}()
	result, _ = auth(th, data)
	if result {
		StandaloneDbms.Store(du.dbms)
	}
	return result
}

// authLimiter limits the rate of authentication attempts
var authLimiter = rate.NewLimiter(rate.Limit(4), 1)
var authContext = context.Background()

func auth(th *Thread, data Value) (result bool, perms *Perms) {
	authLimiter.Wait(authContext)
	authFn := Global.FindName(th, "Auth")
	if authFn == nil {
		return false, nil
	}
	// the app Auth sets up permissions with Perm.ServerEval
	th.SetPerms(&Perms{})
	defer func() {
		th.SetPerms(nil) // clear perms from thread, caller takes ownership
	}()
	result = ToBool(th.CallEach(authFn, data))
	if result {
		perms = th.Perms()
	}
	return result, perms
}

func (du *DbmsUnauth) Check(bool) string {
	panic(notauth)
}

func (du *DbmsUnauth) Close() {
	du.dbms.Close()
}

func (du *DbmsUnauth) Connections() Value {
	panic(notauth)
}

func (du *DbmsUnauth) Cursor(string, *Sviews) ICursor {
	panic(notauth)
}

func (du *DbmsUnauth) Cursors() int {
	panic(notauth)
}

func (du *DbmsUnauth) DisableTrigger(string) {
	panic(notauth)
}

func (du *DbmsUnauth) EnableTrigger(string) {
	panic(notauth)
}

func (du *DbmsUnauth) Dump(string) string {
	panic(notauth)
}

func (du *DbmsUnauth) Exec(*Thread, Value) Value {
	panic(notauth)
}

func (du *DbmsUnauth) Final() int {
	panic(notauth)
}

func (du *DbmsUnauth) Get(*Thread, Value, Dir) (Row, *Header, string) {
	panic(notauth)
}

func (du *DbmsUnauth) Info() Value {
	panic(notauth)
}

func (du *DbmsUnauth) Kill(string) int {
	panic(notauth)
}

func (du *DbmsUnauth) LibGet(name string) []string {
	return du.dbms.LibGet(name)
}

func (du *DbmsUnauth) Libraries() []string {
	return du.dbms.Libraries()
}

func (du *DbmsUnauth) Load(string) int {
	panic(notauth)
}

func (du *DbmsUnauth) Log(s string) {
	panic(notauth)
}

func (du *DbmsUnauth) Schema(string) string {
	panic(notauth)
}

func (du *DbmsUnauth) SessionId(th *Thread, id string) string {
	return du.dbms.SessionId(th, id)
}

func (du *DbmsUnauth) Size() uint64 {
	panic(notauth)
}

func (du *DbmsUnauth) Timestamp() SuDate {
	panic(notauth)
}

func (du *DbmsUnauth) Transaction(bool) ITran {
	panic(notauth)
}

func (du *DbmsUnauth) Transactions() *SuObject {
	panic(notauth)
}

func (du *DbmsUnauth) Unuse(lib string) bool {
	return du.dbms.Unuse(lib)
}

func (du *DbmsUnauth) Use(lib string) bool {
	return du.dbms.Use(lib)
}
