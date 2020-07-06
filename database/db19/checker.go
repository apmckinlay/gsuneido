// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
)

type Checker struct {
	c chan interface{}
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
	t   *CkTran
	ret chan bool
}

type ckResult struct {
}

type ckAbort struct {
	t *CkTran
}

func (ck *Checker) StartTran() *CkTran {
	ret := make(chan *CkTran, 1)
	ck.c <- &ckStart{ret: ret}
	return <-ret
}

func (ck *Checker) Read(t *CkTran, table string, index int, from, to string) bool {
	if t.Aborted() {
		return false
	}
	ck.c <- &ckRead{t: t, table: table, index: index, from: from, to: to}
	return true
}

func (ck *Checker) Write(t *CkTran, table string, keys []string) bool {
	if t.Aborted() {
		return false
	}
	ck.c <- &ckWrite{t: t, table: table, keys: keys}
	return true
}

func (ck *Checker) Commit(t *CkTran) bool {
	if t.Aborted() {
		return false
	}
	ret := make(chan bool, 1)
	ck.c <- &ckCommit{t: t, ret: ret}
	return <-ret
}

func (ck *Checker) Abort(t *CkTran) {
	ck.c <- &ckAbort{t: t}
}

func (t *CkTran) Aborted() bool {
	return t.conflict.Load() != nil
}

//-------------------------------------------------------------------

func NewChecker() *Checker {
	c := make(chan interface{}, 4)
	go checker(c)
	return &Checker{c: c}
}

func checker(c chan interface{}) {
	ck := NewCheck()
	for msg := range c {
		switch msg := msg.(type) {
		case *ckStart:
			if len(ck.trans) > maxTran {
				maxTran = len(ck.trans)
			}
			msg.ret <- ck.StartTran()
		case *ckRead:
			ck.Read(msg.t.start, msg.table, msg.index, msg.from, msg.to)
		case *ckWrite:
			ck.Write(msg.t.start, msg.table, msg.keys)
		case *ckAbort:
			ck.Abort(msg.t.start)
		case *ckCommit:
			msg.ret <- ck.Commit(msg.t.start)
		default:
			panic("checker unknown message type")
		}
	}
	fmt.Println("trans", len(ck.trans))
}

var maxTran int
