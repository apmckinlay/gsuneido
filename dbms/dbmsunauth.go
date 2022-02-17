// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import . "github.com/apmckinlay/gsuneido/runtime"

// DbmsUnauth is a wrapper for DbmsLocal for unauthorized client connections.
// Only allows Auth, LibGet, Libraries, Nonce, SessionId, and Use
type DbmsUnauth struct {
	dbms *DbmsLocal
}

var _ IDbms = (*DbmsUnauth)(nil)

const notauth = "not authorized"

func (du *DbmsUnauth) Admin(string) {
	panic(notauth)
}

func (du *DbmsUnauth) Auth(data string) bool {
	return du.dbms.Auth(data)
}

func (du *DbmsUnauth) Check() string {
	panic(notauth)
}

func (du *DbmsUnauth) Close() {
	du.dbms.Close()
}

func (du *DbmsUnauth) Connections() Value {
	panic(notauth)
}

func (du *DbmsUnauth) Cursor(string) ICursor {
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

func (du *DbmsUnauth) Get(*Thread, string, Dir) (Row, *Header, string) {
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
	du.dbms.Log(s);
}

func (du *DbmsUnauth) Nonce() string {
	return du.dbms.Nonce()
}

func (du *DbmsUnauth) Run(*Thread, string) Value {
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

func (du *DbmsUnauth) Token() string {
	panic(notauth)
}

func (du *DbmsUnauth) Transaction(bool) ITran {
	panic(notauth)
}

func (du *DbmsUnauth) Transactions() *SuObject {
	panic(notauth)
}

func (du *DbmsUnauth) Unuse(string) bool {
	panic(notauth)
}

func (du *DbmsUnauth) Use(lib string) bool {
	return du.dbms.Use(lib)
}