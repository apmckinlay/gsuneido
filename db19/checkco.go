// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"time"

	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/apmckinlay/gsuneido/util/queue"
)

// CheckCo is the concurrent, channel based interface to Check
type CheckCo struct {
	pq      *queue.PriorityQueue
	allDone chan void
}

// message types

type ckStart struct {
	ret chan *CkTran
}

type ckRead struct {
	t        *CkTran
	table    string
	from, to string
	index    int
}

type ckOutput struct {
	t     *CkTran
	table string
	keys  []string
}

type ckDelete struct {
	t     *CkTran
	table string
	keys  []string
	off   uint64
}

type ckUpdate struct {
	t       *CkTran
	table   string
	oldkeys []string
	newkeys []string
	oldoff  uint64
}

type ckCounts struct {
	t   *CkTran
	ret chan int
}

type ckCommit struct {
	t   *UpdateTran
	ret chan bool
}

type ckAbort struct {
	t      *CkTran
	reason string
}

type ckAddExcl struct {
	ret   chan bool
	table string
}

type ckEndExcl struct {
	ret   chan struct{}
	table string
}

type ckRunEndExcl struct {
	ckRunExcl
}

type ckRunExcl struct {
	fn    func()
	ret   chan any
	table string
}

type ckPersist struct {
	ret chan *DbState
}

type ckTrans struct {
	ret chan []int
}

type ckFinal struct {
	ret chan int
}

type ckTick struct{}

// Priority levels
const (
	stopPriority   = 0 // stop signal - lowest priority
	lowPriority    = 1 // start, admin, query operations
	mediumPriority = 2 // read, write operations
	highPriority   = 3 // commit, abort
)

func (ck *CheckCo) StartTran() *CkTran {
	ret := make(chan *CkTran, 1)
	ck.pq.Put(lowPriority, 0, &ckStart{ret: ret})
	return <-ret
}

func (ck *CheckCo) Read(t *CkTran, table string, index int, from, to string) bool {
	if t.Failed() {
		return false
	}
	ck.pq.Put(mediumPriority, t.start, &ckRead{t: t, table: table, index: index, from: from, to: to})
	return true
}

func (ck *CheckCo) Output(t *CkTran, table string, keys []string) bool {
	if t.Failed() {
		return false
	}
	ck.pq.Put(mediumPriority, t.start, &ckOutput{t: t, table: table, keys: keys})
	return true
}

func (ck *CheckCo) Delete(t *CkTran, table string, off uint64, keys []string) bool {
	if t.Failed() {
		return false
	}
	ck.pq.Put(mediumPriority, t.start, &ckDelete{t: t, table: table, off: off, keys: keys})
	return true
}

func (ck *CheckCo) Update(t *CkTran, table string, oldoff uint64, oldkeys, newkeys []string) bool {
	if t.Failed() {
		return false
	}
	ck.pq.Put(mediumPriority, t.start, &ckUpdate{t: t, table: table,
		oldoff: oldoff, oldkeys: oldkeys, newkeys: newkeys})
	return true
}

func (ck *CheckCo) ReadCount(t *CkTran) int {
	if t.Failed() {
		return -1
	}
	ret := make(chan int)
	ck.pq.Put(lowPriority, t.start, &ckCounts{t: t, ret: ret})
	return <-ret
}

func (ck *CheckCo) Commit(ut *UpdateTran) bool {
	if ut.ct.Failed() {
		return false
	}
	ret := make(chan bool, 1)
	ck.pq.Put(highPriority, ut.ct.start, &ckCommit{t: ut, ret: ret})
	return <-ret
}

func (ck *CheckCo) Abort(t *CkTran, reason string) bool {
	ck.pq.Put(highPriority, t.start, &ckAbort{t: t, reason: reason})
	return true
}

// AddExclusive is used by load table and add index
func (ck *CheckCo) AddExclusive(table string) bool {
	ret := make(chan bool, 1)
	ck.pq.Put(mediumPriority, 0, &ckAddExcl{table: table, ret: ret})
	return <-ret
}

func (ck *CheckCo) EndExclusive(table string) {
	ret := make(chan struct{}, 1)
	ck.pq.Put(highPriority, 0, &ckEndExcl{table: table, ret: ret})
	<-ret
}

func (ck *CheckCo) RunEndExclusive(table string, fn func()) any {
	ret := make(chan any, 1)
	ck.pq.Put(mediumPriority, 0, &ckRunEndExcl{ckRunExcl{table: table, fn: fn, ret: ret}})
	return <-ret
}

func (ck *CheckCo) RunExclusive(table string, fn func()) any {
	ret := make(chan any, 1)
	ck.pq.Put(mediumPriority, 0, &ckRunExcl{table: table, fn: fn, ret: ret})
	return <-ret
}

func (ck *CheckCo) Persist() *DbState {
	ret := make(chan *DbState, 1)
	ck.pq.Put(lowPriority, 0, &ckPersist{ret: ret})
	return <-ret
}

func (ck *CheckCo) Transactions() []int {
	ret := make(chan []int, 1)
	ck.pq.Put(lowPriority, 0, &ckTrans{ret: ret})
	return <-ret
}

func (ck *CheckCo) Final() int {
	ret := make(chan int, 1)
	ck.pq.Put(lowPriority, 0, &ckFinal{ret: ret})
	return <-ret
}

//-------------------------------------------------------------------

const ckchanSize = 4 // ???

func StartCheckCo(db *Database, mergeChan chan todo, allDone chan void) *CheckCo {
	ck := NewCheck(db)
	pq := queue.NewPriorityQueue()
	stopTicker := make(chan struct{})
	go tickGenerator(pq, stopTicker)
	go checker(ck, pq, mergeChan, stopTicker)
	return &CheckCo{pq: pq, allDone: allDone}
}

func (ck *CheckCo) Stop() {
	exit.Progress("  checker stopping")
	// send nil rather than closing
	// so other threads don't get "send on closed channel"
	ck.pq.Put(stopPriority, 0, nil)
	<-ck.allDone // wait until closed by concur.go merger
	exit.Progress("  checker stopped")
}

func tickGenerator(pq *queue.PriorityQueue, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case <-time.After(time.Second):
			pq.Put(lowPriority, 0, &ckTick{})
		}
	}
}

func checker(ck *Check, pq *queue.PriorityQueue, mergeChan chan todo, stopTicker chan struct{}) {
	defer func() {
		if e := recover(); e != nil {
			dbg.PrintStack()
			log.Fatalln("FATAL: in checker:", e)
		}
		close(stopTicker)
	}()
	for {
		msg := pq.Get()
		if msg == nil { // Stop sends nil
			if mergeChan != nil { // no channel when testing
				close(mergeChan)
			}
			return
		}
		ck.dispatch(msg, mergeChan)
	}
}

// dispatch runs in the checker goroutine
func (ck *Check) dispatch(msg any, mergeChan chan todo) {
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
	case *ckCounts:
		rc := ck.ReadCount(msg.t)
		msg.ret <- rc
	case *ckAbort:
		ck.Abort(msg.t, msg.reason)
	case *ckCommit:
		tablesWritten := ck.commit(msg.t)
		if tablesWritten == nil {
			msg.ret <- false
		} else if len(tablesWritten) == 0 {
			msg.ret <- true
		} else {
			msg.t.commit()
			msg.ret <- true
			mergeChan <- todo{tables: tablesWritten, meta: msg.t.meta}
		}
	case *ckAddExcl:
		msg.ret <- ck.AddExclusive(msg.table)
	case *ckEndExcl:
		ck.EndExclusive(msg.table)
		msg.ret <- struct{}{}
	case *ckRunEndExcl:
		defer ck.EndExclusive(msg.table)
		ck.run(&msg.ckRunExcl, mergeChan)
	case *ckRunExcl:
		if !ck.AddExclusive(msg.table) {
			msg.ret <- "already exclusive: " + msg.table
			return
		}
		defer ck.EndExclusive(msg.table)
		ck.run(msg, mergeChan)
	case *ckPersist:
		ret := make(chan any)
		mergeChan <- todo{ret: ret}
		state := <-ret
		msg.ret <- state.(*DbState)
	case *ckTrans:
		msg.ret <- ck.Transactions()
	case *ckFinal:
		msg.ret <- ck.Final()
	case *ckTick:
		ck.tick()
	default:
		panic("checker unknown message type")
	}
}

func (ck *Check) run(msg *ckRunExcl, mergeChan chan todo) {
	ret := make(chan any)
	td := todo{fn: msg.fn, ret: ret}
	mergeChan <- td
	err := <-ret
	msg.ret <- err
}

// Checker is the interface for Check and CheckCo
type Checker interface {
	StartTran() *CkTran
	Read(t *CkTran, table string, index int, from, to string) bool
	Output(t *CkTran, table string, keys []string) bool
	Delete(t *CkTran, table string, off uint64, keys []string) bool
	Update(t *CkTran, table string, oldoff uint64, oldkeys, newkeys []string) bool
	ReadCount(t *CkTran) int
	Abort(t *CkTran, reason string) bool
	Commit(t *UpdateTran) bool
	Persist() *DbState
	Stop()
	Transactions() []int
	Final() int
	AddExclusive(table string) bool
	EndExclusive(table string)
	RunEndExclusive(table string, fn func()) any
	RunExclusive(table string, fn func()) any
}

var _ Checker = (*Check)(nil)
var _ Checker = (*CheckCo)(nil)
