// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
)

// CheckCo is the concurrent, channel based interface to Check
type CheckCo struct {
	c       chan interface{}
	allDone chan void
}

// message types

type ckCreate struct {
	schema *schema.Schema
	ret    chan error
}

type ckDrop struct {
	table string
	ret   chan error
}

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

func (ck *CheckCo) Create(schema *schema.Schema) error {
	ret := make(chan error, 1)
	ck.c <- &ckCreate{schema: schema, ret: ret}
	return <-ret
}

func (ck *CheckCo) Drop(table string) error {
	ret := make(chan error, 1)
	ck.c <- &ckDrop{table: table, ret: ret}
	return <-ret
}

func (ck *CheckCo) StartTran() *CkTran {
	ret := make(chan *CkTran, 1)
	ck.c <- &ckStart{ret: ret}
	return <-ret
}

func (ck *CheckCo) Read(t *CkTran, table string, index int, from, to string) bool {
	if t.Aborted() {
		return false
	}
	ck.c <- &ckRead{t: t, table: table, index: index, from: from, to: to}
	return true
}

func (ck *CheckCo) Write(t *CkTran, table string, keys []string) bool {
	if t.Aborted() {
		return false
	}
	ck.c <- &ckWrite{t: t, table: table, keys: keys}
	return true
}

func (ck *CheckCo) Commit(ut *UpdateTran) bool {
	if ut.ct.Aborted() {
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

func (t *CkTran) Aborted() bool {
	return t.conflict.Load() != nil
}

//-------------------------------------------------------------------

func StartCheckCo(db *Database, mergeChan chan interface{}, resultChan chan error,
	allDone chan void) *CheckCo {
	ck := NewCheck(db)
	c := make(chan interface{}, 4)
	go checker(ck, c, mergeChan, resultChan)
	return &CheckCo{c: c, allDone: allDone}
}

func (ck *CheckCo) Stop() {
	close(ck.c)
	<-ck.allDone // wait
}

func checker(ck *Check, c chan interface{}, mergeChan chan interface{}, resultChan chan error) {
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
			ck.dispatch(msg, mergeChan, resultChan)
		case <-ticker.C:
			// fmt.Println("checker chan", len(c), "merge chan", len(mergeChan))
			ck.tick()
		}
	}
}

// dispatch runs in the checker goroutine
func (ck *Check) dispatch(msg interface{}, mergeChan chan interface{}, resultChan chan error) {
	switch msg := msg.(type) {
	case *ckCreate:
		mergeChan <- msg.schema
		err := <-resultChan
		msg.ret <- err
	case *ckDrop:
		mergeChan <- msg.table
		err := <-resultChan
		msg.ret <- err
	case *ckStart:
		msg.ret <- ck.StartTran()
	case *ckRead:
		ck.Read(msg.t, msg.table, msg.index, msg.from, msg.to)
	case *ckWrite:
		ck.Write(msg.t, msg.table, msg.keys)
	case *ckAbort:
		ck.Abort(msg.t, msg.reason)
	case *ckCommit:
		result := ck.commit(msg.t)
		if result == nil {
			msg.ret <- false
		}
		msg.t.commit()
		msg.ret <- true
		mergeChan <- result
	default:
		panic("checker unknown message type")
	}
}

// Checker is the interface for Check and CheckCo
type Checker interface {
	Create(schema *schema.Schema) error
	Drop(table string) error
	StartTran() *CkTran
	Read(t *CkTran, table string, index int, from, to string) bool
	Write(t *CkTran, table string, keys []string) bool
	Abort(t *CkTran, reason string) bool
	Commit(t *UpdateTran) bool
	Stop()
}

var _ Checker = (*Check)(nil)
var _ Checker = (*CheckCo)(nil)
