// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !gui

package dbms

import (
	"crypto/tls"
	"net"
	"slices"
	"strings"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/mux"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLogWithLimit(t *testing.T) {
	assert := assert.T(t)

	// Create a serverConn to test
	sc := &serverConn{}

	// Test 1: Normal logging under limit
	result := sc.limitLog("Hello World") // 11 bytes + 1 newline = 12
	assert.This(result).Is("Hello World")
	assert.This(sc.logSize.Load()).Is(12)

	// Test 2: Add more logs, still under limit
	result = sc.limitLog("Test message") // 12 bytes + 1 newline = 13, total 25
	assert.This(result).Is("Test message")
	assert.This(sc.logSize.Load()).Is(25)

	// Test 3: Add log that crosses the 10KB limit
	largeMsg := strings.Repeat("x", 10240-25-1) // 10214 + 1 newline = 10215, total = 10240
	result = sc.limitLog(largeMsg)
	assert.This(result).Is(largeMsg)
	assert.This(sc.logSize.Load()).Is(10240)

	// Test 4: Next log should trigger warning
	result = sc.limitLog("x") // 1 + 1 newline = 2, total 10242
	assert.This(result).Is("log size limit exceeded (10KB), ignoring further logs")
	assert.This(sc.logSize.Load()).Is(10242)

	// Test 5: Subsequent logs should be ignored (return empty string)
	oldSize := sc.logSize.Load()
	result = sc.limitLog("More logging") // 12 + 1 = 13
	assert.This(result).Is("")
	assert.This(sc.logSize.Load()).Is(oldSize + 13) // Size still increases even though ignored
}

func TestLogWithLimitBoundary(t *testing.T) {
	assert := assert.T(t)

	sc := &serverConn{}

	// Test exactly at the boundary
	result := sc.limitLog(strings.Repeat("a", 10239)) // 10239 + 1 newline = 10240 exactly
	assert.This(result).Is(strings.Repeat("a", 10239))
	assert.This(sc.logSize.Load()).Is(10240)

	// Next log should trigger the limit warning
	result = sc.limitLog("x") // 1 byte + 1 newline = 2, total 10242
	assert.This(result).Is("log size limit exceeded (10KB), ignoring further logs")
	assert.This(sc.logSize.Load()).Is(10242)

	// Subsequent logs should be ignored
	result = sc.limitLog("ignored")
	assert.This(result).Is("")
}

func TestLogWithLimitEmpty(t *testing.T) {
	assert := assert.T(t)

	sc := &serverConn{}

	// Test empty string - still counts 1 byte for the newline that log.Println adds
	result := sc.limitLog("")
	assert.This(result).Is("")
	assert.This(sc.logSize.Load()).Is(1)

	// Test that we can still log after empty
	result = sc.limitLog("hello")
	assert.This(result).Is("hello")
	assert.This(sc.logSize.Load()).Is(int32(7)) // 1 + (5 + 1) = 7
}

func TestRowToRecordSingleList(t *testing.T) {
	hdr := SimpleHeader([]string{"a", "-", "b"})
	rec := new(RecordBuilder).Add(One).Add(SuStr("gone")).Add(MinusOne).Build()
	row := NewRow(rec)

	r, fields := rowToRecord(row, hdr)
	// single-list headers preserve the original layout, including "-"
	assert.T(t).That(slices.Equal(fields, []string{"a", "-", "b"}))
	assert.T(t).That(slices.Equal(hdr.GetFields(), fields))
	assert.T(t).This(r.GetRaw(0)).Is(Pack(One))
	assert.T(t).This(r.GetRaw(2)).Is(Pack(MinusOne))
}

func TestRowToRecordMultiList(t *testing.T) {
	hdr := NewHeader([][]string{{"a", "-"}, {"-", "b"}}, []string{"a", "b"})
	rec1 := new(RecordBuilder).Add(One).Add(SuStr("gone1")).Build()
	rec2 := new(RecordBuilder).Add(SuStr("gone2")).Add(MinusOne).Build()
	row := JoinRows(NewRow(rec1), NewRow(rec2))

	r, fields := rowToRecord(row, hdr)
	// Preserve the first list's deleted slots for Record.Update serialization.
	assert.T(t).That(slices.Equal(fields, []string{"a", "-", "b"}))
	assert.T(t).That(slices.Equal(hdr.GetFields(), fields))
	assert.T(t).This(r.Count()).Is(3)
	assert.T(t).This(r.GetRaw(0)).Is(Pack(One))
	assert.T(t).This(r.GetRaw(1)).Is("")
	assert.T(t).This(r.GetRaw(2)).Is(Pack(MinusOne))
}

func TestRowToRecordExtendedUpdate(t *testing.T) {
	base := []string{"id", "-", "type"}
	hdr := NewHeader([][]string{base, {"status"}}, []string{"id", "type", "status"})
	row := NewRow(new(RecordBuilder).Add(One).Build()).
		With(new(RecordBuilder).Add(One).Build())
	rec, fields := rowToRecord(row, hdr)
	clientHdr := SimpleHeader(fields)
	r := SuRecordFromRow(NewRow(rec), clientHdr, "", nil, 0)
	r.Put(nil, SuStr("type"), SuStr("Change"))
	updated := r.ToRecord(nil, clientHdr)
	assert.T(t).This(updated.GetRaw(1)).Is("")
	assert.T(t).This(updated.GetRaw(2)).Is(Pack(SuStr("Change")))
	assert.T(t).This(updated.GetRaw(3)).Is(Pack(One))
}

func TestNewServerConnUnauthorized(t *testing.T) {
	assert := assert.T(t)

	options.BuiltDate = "Dec 29 2020 12:34"
	db := db19.CreateDb(stor.HeapStor(8192))
	dbmsLocal := NewDbmsLocal(db)
	p1, p2 := net.Pipe()
	workers = mux.NewWorkers(doRequest)
	go newServerConn(dbmsLocal, p1, serverTLSConfig())
	errmsg := checkHello(p2)
	assert.This(errmsg).Is("")
	p2.Write(hello())
	tlsConn := tls.Client(p2, clientTLSConfig())
	if err := tlsConn.Handshake(); err != nil {
		panic(err)
	}
	c := NewDbmsClient(tlsConn)
	ses := c.NewSession()

	// Without Auth, the connection is unauthorized and Get should be rejected
	assert.This(func() {
		args := SuObjectOf(SuStr("tables sort table"))
		ses.Get(nil, args, Next)
	}).Panics("not authorized")

	time.Sleep(25 * time.Millisecond)
}

func TestUnauthCmds(t *testing.T) {
	assert := assert.T(t)

	sc := &serverConn{dbms: &DbmsUnauth{}}
	ss := &serverSession{sc: sc, thread: &Thread{}}

	// Commands that route through ss.sc.dbms should panic with "not authorized"
	ss.SetBuf([]byte{0}) // empty string (varint size 0)
	assert.This(func() { cmdAdmin(ss) }).Panics(notauth)

	ss.SetBuf([]byte{0}) // bool false
	assert.This(func() { cmdCheck(ss) }).Panics(notauth)

	ss.SetBuf([]byte{0}) // empty string (varint size 0)
	assert.This(func() { cmdCursor(ss) }).Panics(notauth)

	assert.This(func() { cmdFinal(ss) }).Panics(notauth)

	assert.This(func() { cmdInfo(ss) }).Panics(notauth)

	ss.SetBuf([]byte{0}) // empty string
	assert.This(func() { cmdKill(ss) }).Panics(notauth)

	// Encode a non-empty string for Log (varint 1 + "x" = 0x02 0x78)
	ss.SetBuf([]byte{0x02, 'x'})
	assert.This(func() { cmdLog(ss) }).Panics(notauth)

	assert.This(func() { cmdSize(ss) }).Panics(notauth)

	assert.This(func() { cmdTimestamp(ss) }).Panics(notauth)

	ss.SetBuf([]byte{0}) // bool false
	assert.This(func() { cmdTransaction(ss) }).Panics(notauth)

	assert.This(func() { cmdTransactions(ss) }).Panics(notauth)
}
