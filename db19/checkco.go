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

type ckWrite struct {
	t     *CkTran
	table string
	keys  []string
	ret   chan bool
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
	tables []string
	ret    chan bool
}

type ckEndExcl struct {
	tables []string
}

type ckPersist struct {
	ret chan *DbState
}

type ckTrans struct {
	ret chan []int
}

// var i = 0

func (ck *CheckCo) StartTran() *CkTran {
	// fmt.Print("\rStart ", i)
	// i++
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

func (ck *CheckCo) Write(t *CkTran, table string, keys []string) bool {
	if t.Failed() {
		return false
	}
	ret := make(chan bool, 1)
	ck.c <- &ckWrite{t: t, table: table, keys: keys, ret: ret}
	return <-ret
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
func (ck *CheckCo) AddExclusive(tables ...string) bool {
	ret := make(chan bool, 1)
	ck.c <- &ckAddExcl{ret: ret, tables: tables}
	return <-ret
}

func (ck *CheckCo) EndExclusive(tables ...string) {
	ck.c <- &ckEndExcl{tables: tables}
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
			// fmt.Println("checker chan", len(c), "merge chan", len(mergeChan))
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
	case *ckWrite:
		msg.ret <- ck.Write(msg.t, msg.table, msg.keys)
	case *ckAbort:
		ck.Abort(msg.t, msg.reason)
	case *ckCommit:
		result := ck.commit(msg.t)
		if result == nil {
			msg.ret <- false
			return
		}
		msg.t.commit()
		msg.ret <- true
		mergeChan <- todo{tables: result, meta: msg.t.meta}
	case *ckAddExcl:
		if !ck.AddExclusive(msg.tables...) {
			msg.ret <- false
		}
		// ensure pending merges are all complete
		ret := make(chan *DbState)
		mergeChan <- todo{ret: ret} // sync (meta == nil)
		<-ret
		msg.ret <- true
	case *ckEndExcl:
		ck.EndExclusive(msg.tables...)
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
	Write(t *CkTran, table string, keys []string) bool
	Abort(t *CkTran, reason string) bool
	Commit(t *UpdateTran) bool
	AddExclusive(tables ...string) bool
	EndExclusive(tables ...string)
	Persist() *DbState
	Stop()
	Transactions() []int
}

var _ Checker = (*Check)(nil)
var _ Checker = (*CheckCo)(nil)
