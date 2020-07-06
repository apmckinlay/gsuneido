// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCheckerRandom(*testing.T) {
	ck := NewChecker()
	nThreads := 8
	nTrans := 10000
	if testing.Short() {
		nThreads = 2
		nTrans = 1000
	}
	var wg sync.WaitGroup
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < nTrans; i++ {
				randTran(ck)
			}
		}()
	}
	wg.Wait()
	close(ck.c)
	fmt.Println("commit", nCommit, "conflict", nConflict)
	fmt.Println("maxTran", maxTran)
}

var nCommit, nConflict int64

func randTran(ck *Checker) {
	t := ck.StartTran()
	nActions := rand.Intn(20)
	for i := 0; i < nActions; i++ {
		randAction(ck, t)
	}
	if rand.Intn(2) == 1 {
		ck.Abort(t)
	} else {
		if ck.Commit(t) {
			atomic.AddInt64(&nCommit, 1)
		} else {
			atomic.AddInt64(&nConflict, 1)
		}
	}

}

func randAction(ck *Checker, t *CkTran) {
	nIndexes := 4
	table := randTable()
	if rand.Intn(3) == 1 {
		ck.Write(t, table, randKeys())
	} else {
		index := rand.Intn(nIndexes)
		from, to := randRange()
		ck.Read(t, table, index, from, to)
	}

}

func randTable() string {
	tables := []string{"one", "two", "three", "four", "five"}
	return tables[rand.Intn(len(tables))]
}

func randKeys() []string {
	nIndexes := 3
	keys := make([]string, nIndexes)
	for i := range keys {
		keys[i] = strconv.Itoa(rand.Intn(10000))
	}
	return keys
}

func randRange() (string, string) {
	from := rand.Intn(10000)
	to := from + rand.Intn(10)
	return strconv.Itoa(from), strconv.Itoa(to)
}
