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
	t *CkTran
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

func (ck *CheckCo) Abort(t *CkTran) bool {
	ck.c <- &ckAbort{t: t}
	return true
}

func (t *CkTran) Aborted() bool {
	return t.conflict.Load() != nil
}

//-------------------------------------------------------------------

func StartCheckCo(commitChan chan *UpdateTran, allDone chan void) *CheckCo {
	c := make(chan interface{}, 4)
	go checker(c, commitChan)
	return &CheckCo{c: c, allDone: allDone}
}

func (ck *CheckCo) Stop() {
	close(ck.c)
	<-ck.allDone // wait
}

func checker(c chan interface{}, commitChan chan *UpdateTran) {
	ck := NewCheck(commitChan)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg := <-c:
			if msg == nil { // channel closed
				if commitChan != nil {
					close(commitChan)
				}
				return
			}
			ck.dispatch(msg)
		case <-ticker.C:
			ck.tick()
		}
	}
}

func (ck *Check) dispatch(msg interface{}) {
	switch msg := msg.(type) {
	case *ckStart:
		msg.ret <- ck.StartTran()
	case *ckRead:
		ck.Read(msg.t, msg.table, msg.index, msg.from, msg.to)
	case *ckWrite:
		ck.Write(msg.t, msg.table, msg.keys)
	case *ckAbort:
		ck.Abort(msg.t)
	case *ckCommit:
		msg.ret <- ck.Commit(msg.t)
	default:
		panic("checker unknown message type")
	}
}

type Checker interface {
	StartTran() *CkTran
	Read(t *CkTran, table string, index int, from, to string) bool
	Write(t *CkTran, table string, keys []string) bool
	Abort(t *CkTran) bool
	Commit(t *UpdateTran) bool
	Stop()
}

var _ Checker = (*Check)(nil)
var _ Checker = (*CheckCo)(nil)
