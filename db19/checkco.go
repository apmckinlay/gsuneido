// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"time"
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

func StartCheckCo(mergeChan chan merge, allDone chan void) *CheckCo {
	c := make(chan interface{}, 4)
	go checker(c, mergeChan)
	return &CheckCo{c: c, allDone: allDone}
}

func (ck *CheckCo) Stop() {
	close(ck.c)
	<-ck.allDone // wait
}

func checker(c chan interface{}, mergeChan chan merge) {
	ck := NewCheck()
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
func (ck *Check) dispatch(msg interface{}, mergeChan chan merge) {
	switch msg := msg.(type) {
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
		// checking complete so we can send result and let client code continue
		msg.ret <- result != nil
		if result != nil && mergeChan != nil { // no channel when testing
			// passed checking so we can asynchronously commit it
			// (can't fail at this point)
			// Since we haven't returned, no other activity will happen
			// until we finish the commit. i.e. serialized
			msg.t.commit()
			mergeChan <- result
			//TODO commit and send merge in another goroutine ?
		}
	default:
		panic("checker unknown message type")
	}
}

// Checker is the interface for Check and CheckCo
type Checker interface {
	StartTran() *CkTran
	Read(t *CkTran, table string, index int, from, to string) bool
	Write(t *CkTran, table string, keys []string) bool
	Abort(t *CkTran, reason string) bool
	Commit(t *UpdateTran) bool
	Stop()
}

var _ Checker = (*Check)(nil)
var _ Checker = (*CheckCo)(nil)
