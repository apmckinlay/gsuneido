// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCheckCoTimeout(t *testing.T) {
	if testing.Short() {
		return
	}
	defer func(ma int) { MaxAge = ma }(MaxAge)
	MaxAge = 1
	ck := StartCheckCo(nil, nil, nil)
	tran := ck.StartTran()
	assert.T(t).False(tran.Failed())
	time.Sleep(2 * time.Second)
	assert.T(t).True(tran.Failed())
	close(ck.c)
}

func TestCheckCoRandom(*testing.T) {
	db, err := CreateDb(stor.HeapStor(8192))
	assert.That(err == nil)
	db.ck = StartCheckCo(db, mergeSink(), nil)
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
				randTran(db)
			}
		}()
	}
	wg.Wait()
	// fmt.Println("commit", nCommit, "conflict", nConflict)
}

func mergeSink() chan todo {
	c := make(chan todo)
	go func() {
		for range c {
		}
	}()
	return c
}

var nCommit, nConflict int64

func randTran(db *Database) {
	t := db.NewUpdateTran()
	nActions := rand.Intn(20)
	for i := 0; i < nActions; i++ {
		randAction(db.ck, t.ct)
	}
	if rand.Intn(2) == 1 {
		db.ck.Abort(t.ct, "")
	} else {
		if db.ck.Commit(t) {
			atomic.AddInt64(&nCommit, 1)
		} else {
			atomic.AddInt64(&nConflict, 1)
		}
	}

}

func randAction(ck Checker, t *CkTran) {
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
