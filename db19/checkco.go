// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta"
)

// CheckCo is the concurrent, channel based interface to Check
type CheckCo struct {
	c       chan interface{}
	allDone chan void
}

// message types

type ckStart struct {
	ret chan *CkTran
}

type ckRead struct {
	t        *CkTran
	table    string
	index    int
	from, to string
}

type ckOutput struct {
	t     *CkTran
	table string
	keys  []string
}

type ckDelete struct {
	t     *CkTran
	table string
	off   uint64
	keys  []string
}

type ckUpdate struct {
	t       *CkTran
	table   string
	oldoff  uint64
	oldkeys []string
	newkeys []string
}

type ckCommit struct {
	t   *UpdateTran
	ret chan bool
}

type ckResult struct {
}

type ckAbort struct {
	t      *CkTran
	reason string
}

type ckAddExcl struct {
	table string
	ret   chan struct{}
}

type ckEndExcl struct {
	table string
}

type ckPersist struct {
	ret chan *DbState
}

type ckTrans struct {
	ret chan []int
}

func (ck *CheckCo) StartTran() *CkTran {
	ret := make(chan *CkTran, 1)
	ck.c <- &ckStart{ret: ret}
	return <-ret
}

func (ck *CheckCo) Read(t *CkTran, table string, index int, from, to string) bool {
	if t.Failed() {
		return false
	}
	ck.c <- &ckRead{t: t, table: table, index: index, from: from, to: to}
	return true
}

func (ck *CheckCo) Output(t *CkTran, table string, keys []string) bool {
	if t.Failed() {
		return false
	}
	ck.c <- &ckOutput{t: t, table: table, keys: keys}
	return true
}

func (ck *CheckCo) Delete(t *CkTran, table string, off uint64, keys []string) bool {
	if t.Failed() {
		return false
	}
	ck.c <- &ckDelete{t: t, table: table, off: off, keys: keys}
	return true
}

func (ck *CheckCo) Update(t *CkTran, table string, oldoff uint64, oldkeys, newkeys []string) bool {
	if t.Failed() {
		return false
	}
	ck.c <- &ckUpdate{t: t, table: table,
		oldoff: oldoff, oldkeys: oldkeys, newkeys: newkeys}
	return true
}

func (ck *CheckCo) Commit(ut *UpdateTran) bool {
	if ut.ct.Failed() {
		return false
	}
	ret := make(chan bool, 1)
	ck.c <- &ckCommit{t: ut, ret: ret}
	return <-ret
}

func (ck *CheckCo) Abort(t *CkTran, reason string) bool {
	ck.c <- &ckAbort{t: t, reason: reason}
	return true
}

// AddExclusive also does sync (handled in dispatch)
func (ck *CheckCo) AddExclusive(table string) {
	ret := make(chan struct{}, 1)
	ck.c <- &ckAddExcl{ret: ret, table: table}
	<-ret
}

func (ck *CheckCo) EndExclusive(table string) {
	ck.c <- &ckEndExcl{table: table}
}

func (ck *CheckCo) Persist() *DbState {
	ret := make(chan *DbState, 1)
	ck.c <- &ckPersist{ret: ret}
	return <-ret
}

func (ck *CheckCo) Transactions() []int {
	ret := make(chan []int, 1)
	ck.c <- &ckTrans{ret: ret}
	return <-ret
}

//-------------------------------------------------------------------

func StartCheckCo(db *Database, mergeChan chan todo, allDone chan void) *CheckCo {
	ck := NewCheck(db)
	c := make(chan interface{}, 4)
	go checker(ck, c, mergeChan)
	return &CheckCo{c: c, allDone: allDone}
}

func (ck *CheckCo) Stop() {
	close(ck.c)
	<-ck.allDone // wait
}

func checker(ck *Check, c chan interface{}, mergeChan chan todo) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			log.Fatalln("FATAL ERROR in checker:", e)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg := <-c:
			if msg == nil { // channel closed
				if mergeChan != nil { // no channel when testing
					close(mergeChan)
				}
				return
			}
			ck.dispatch(msg, mergeChan)
		case <-ticker.C:
			ck.tick()
		}
	}
}

// dispatch runs in the checker goroutine
func (ck *Check) dispatch(msg interface{}, mergeChan chan todo) {
	switch msg := msg.(type) {
	case *ckStart:
		msg.ret <- ck.StartTran()
	case *ckRead:
		ck.Read(msg.t, msg.table, msg.index, msg.from, msg.to)
	case *ckOutput:
		ck.Output(msg.t, msg.table, msg.keys)
	case *ckDelete:
		ck.Delete(msg.t, msg.table, msg.off, msg.keys)
	case *ckUpdate:
		ck.Update(msg.t, msg.table, msg.oldoff, msg.oldkeys, msg.newkeys)
	case *ckAbort:
		ck.Abort(msg.t, msg.reason)
	case *ckCommit:
		tablesWritten := ck.commit(msg.t)
		if tablesWritten == nil {
			msg.ret <- false
			return
		}
		msg.t.commit()
		msg.ret <- true
		mergeChan <- todo{tables: tablesWritten, meta: msg.t.meta}
	case *ckAddExcl:
		ck.AddExclusive(msg.table)
		// ensure pending merges are all complete
		ret := make(chan *DbState)
		mergeChan <- todo{ret: ret} // sync (meta == nil)
		<-ret
		msg.ret <- struct{}{}
	case *ckEndExcl:
		ck.EndExclusive(msg.table)
	case *ckPersist:
		ret := make(chan *DbState)
		mergeChan <- todo{meta: persist, ret: ret}
		state := <-ret
		msg.ret <- state
	case *ckTrans:
		msg.ret <- ck.Transactions()
	default:
		panic("checker unknown message type")
	}
}

// persist is used to distinguish sync and persist
var persist = &meta.Meta{}

// Checker is the interface for Check and CheckCo
type Checker interface {
	StartTran() *CkTran
	Read(t *CkTran, table string, index int, from, to string) bool
	Output(t *CkTran, table string, keys []string) bool
	Delete(t *CkTran, table string, off uint64, keys []string) bool
	Update(t *CkTran, table string, oldoff uint64, oldkeys, newkeys []string) bool
	Abort(t *CkTran, reason string) bool
	Commit(t *UpdateTran) bool
	AddExclusive(table string)
	EndExclusive(table string)
	Persist() *DbState
	Stop()
	Transactions() []int
}

var _ Checker = (*Check)(nil)
var _ Checker = (*CheckCo)(nil)
