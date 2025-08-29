// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"golang.org/x/time/rate"
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
	ck.pq.Put(stopPriority, 0, nil)
}

func TestCheckCoRandom(*testing.T) {
	db := CreateDb(stor.HeapStor(8192))
	db.ck = StartCheckCo(db, mergeSink(), nil)
	nThreads := 6
	nTrans := 10000
	if testing.Short() {
		nThreads = 2
		nTrans = 1000
	}
	var wg sync.WaitGroup
	for range nThreads {
		wg.Go(func() {
			for range nTrans {
				randTran(db)
			}
		})
	}
	wg.Wait()
	f := float32(nConflict.Load()) / float32(nCommit.Load())
	fmt.Println("commit", nCommit.Load(), "conflict", nConflict.Load(), "=", f)
	assert.That(f < .1)
}

func mergeSink() chan todo {
	c := make(chan todo)
	go func() {
		for range c {
		}
	}()
	return c
}

var nCommit, nConflict atomic.Int32

func randTran(db *Database) {
	t := db.NewUpdateTran()
	nActions := rand.Intn(20)
	for range nActions {
		randAction(db.ck, t.ct)
	}
	if rand.Intn(2) == 1 {
		db.ck.Abort(t.ct, "")
	} else {
		if db.ck.Commit(t) {
			nCommit.Add(1)
		} else {
			nConflict.Add(1)
		}
	}

}

func randAction(ck Checker, t *CkTran) {
	nIndexes := 4
	table := randTable()
	if rand.Intn(3) == 1 {
		// Always do a read before update to make the test realistic
		index := rand.Intn(nIndexes)
		oldKeys := randKeys()
		// Make the read range match the key that will be updated
		from := oldKeys[index]
		to := from + "\x00" // range that includes the key
		ck.Read(t, table, index, from, to)
		ck.Update(t, table, uint64(1+rand.Intn(2000)), oldKeys, randKeys())
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
	nIndexes := 4
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

func TestRateForOutstanding(t *testing.T) {
	// Below threshold => Inf
	if r := rateForOutstanding(0); r != rate.Inf {
		t.Fatalf("expected Inf for 0, got %v", r)
	}
	if r := rateForOutstanding(199); r != rate.Inf {
		t.Fatalf("expected Inf for 199, got %v", r)
	}

	// At 200 => 100/sec
	if r := rateForOutstanding(200); !almostEqual(float64(r), 100.0) {
		t.Fatalf("expected 100 for 200, got %v", r)
	}

	// Midpoint ~ (200+499)/2 => around ~50.5/sec
	mid := (200 + 499) / 2 // 349
	expectedMid := 100.0 - float64(mid-200)*99.0/299.0
	if r := rateForOutstanding(mid); !almostEqual(float64(r), expectedMid) {
		t.Fatalf("expected %v for %v, got %v", expectedMid, mid, r)
	}

	// At 499 => ~1/sec
	if r := rateForOutstanding(499); !almostEqual(float64(r), 1.0) {
		t.Fatalf("expected ~1 for 499, got %v", r)
	}

	// At/above 500 => 1/sec
	if r := rateForOutstanding(500); !almostEqual(float64(r), 1.0) {
		t.Fatalf("expected 1 for 500, got %v", r)
	}
	if r := rateForOutstanding(1000); !almostEqual(float64(r), 1.0) {
		t.Fatalf("expected 1 for 1000, got %v", r)
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}
