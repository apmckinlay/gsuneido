// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type triggers struct {
	lock     sync.Mutex
	disabled map[string]int
}

// MakeSuTran is injected by dbms to avoid import cycle
var MakeSuTran func(tran *UpdateTran) *SuTran

func (t *triggers) DisableTrigger(table string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.disabled == nil {
		t.disabled = make(map[string]int)
	}
	t.disabled[table]++
}

func (t *triggers) EnableTrigger(table string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.disabled[table]--
	assert.That(t.disabled[table] >= 0)
}

func (t *triggers) enabled(table string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.disabled[table] == 0
}

func (t *triggers) CallTrigger(th *Thread, tran *UpdateTran, table string,
	oldrec, newrec Record) {
	sutran := MakeSuTran(tran)
	hdr := SimpleHeader(tran.GetSchema(table).Columns)
	t.call2(th, sutran, table,
		suRec(oldrec, hdr, sutran), suRec(newrec, hdr, sutran))
}

func suRec(rec Record, hdr *Header, tran *SuTran) Value {
	if rec == "" {
		return False
	}
	return SuRecordFromRow(Row{DbRec{Record: rec}}, hdr, "", tran)
}

func (t *triggers) call2(th *Thread, tran *SuTran, table string,
	oldrec, newrec Value) {
	if !t.enabled(table) {
		return
	}
	name := "Trigger_" + table
	fn := Global.FindName(th, name)
	if fn == nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			WrapPanic(e, name)
		}
	}()
	th.Call(fn, tran, oldrec, newrec)
}
