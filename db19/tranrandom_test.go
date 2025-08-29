// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"math/rand/v2"
	"runtime/metrics"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sample"
)

//TODO add a second index

//TODO add 2 phase locking (2PL) style verification (data doesn't catch enough)

func TestRandomConcurrentTransactions(t *testing.T) {
	const numTransactions = 500_000
	const numRecords = 9000
	const numThreads = 16

	if testing.Short() {
		t.Skip("skipping random concurrent transaction test in short mode")
	}

	db := createRandomTestDb()
	defer db.Close()

	testData := initializeTestData(db, numRecords)

	start := time.Now()
	var wg sync.WaitGroup
	for range numThreads {
		wg.Go(func() {
			runRandomTransactions(db, testData, numTransactions)
		})
	}
	wg.Wait()

	commits := testData.commitCount.Load()
	aborts := testData.abortCount.Load()
	ntrans := commits + aborts
	abortPercent := float64(aborts) * 100.0 / float64(ntrans)
	fmt.Println("memory alloc per transactions", memory()/ntrans)
	fmt.Printf("transactions %d, aborts %.2f%% (%d)\n",
		ntrans, abortPercent, aborts)
	perSec := float64(ntrans) / time.Since(start).Seconds()
	fmt.Println("transactions per second", trace.Number(int(perSec)))

	verifyData(t, db, testData)
	db.MustCheck()
}

var nreadonly, nabort1of int

func TestRandomSingleThreadTransactions(t *testing.T) {
	const numRecords = 10

	if testing.Short() {
		t.Skip("skipping random single thread transaction test in short mode")
	}

	db := createRandomTestDb()
	defer db.Close()

	testData := initializeTestData(db, numRecords)

	//	runRandomTransactions(db, testData, 10_000)
	t1 := db.NewUpdateTran()
	t2 := db.NewUpdateTran()
	key := 5
	r1 := t1.Lookup("test_table", 0, intKey(key))
	off1 := r1.Off
	r2 := t2.Lookup("test_table", 0, intKey(key))
	off2 := r2.Off
	newRecord := int2rec(key, 123)
	t1.Update(nil, "test_table", off1, newRecord)
	t2.Update(nil, "test_table", off2, newRecord)

	verifyData(t, db, testData)
	db.MustCheck()
}

type RandomTestData struct {
	mirrorData  []int
	mutex       sync.RWMutex // guards mirrorData
	numRecords  int
	commitCount atomic.Uint64
	abortCount  atomic.Uint64
	ntrans      atomic.Uint64
}

func createRandomTestDb() *Database {
	db := CreateDb(stor.HeapStor(8192))
	StartConcur(db, 50*time.Millisecond)

	db.Create(&schema.Schema{
		Table:   "test_table",
		Columns: []string{"key", "count"},
		Indexes: []schema.Index{{Mode: 'k', Columns: []string{"key"}}},
	})

	return db
}

func initializeTestData(db *Database, numRecords int) *RandomTestData {
	testData := &RandomTestData{
		mirrorData: make([]int, numRecords),
		numRecords: numRecords,
	}
	ut := db.NewUpdateTran()
	for i := range numRecords {
		ut.Output(nil, "test_table", int2rec(i, 0))
		testData.mirrorData[i] = 0
	}
	ut.Complete()
	return testData
}

func runRandomTransactions(db *Database, testData *RandomTestData, count uint64) {
	for testData.ntrans.Add(1) <= count {
		if performRandomTransaction(db, testData) {
			testData.commitCount.Add(1)
		} else {
			testData.abortCount.Add(1)
		}
	}
}

func performRandomTransaction(db *Database, testData *RandomTestData) (success bool) {
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && strings.Contains(s, "abort") {
				success = false
			} else {
				panic(r)
			}
		}
	}()

	ut := db.NewUpdateTran()
	if ut == nil {
		return false
	}

	n := 2 + rand.IntN(7) // 1 to 7
	keys := make([]int, n)
	offsets := make([]uint64, n)
	changes := make([]int, n)
	sum := 0
	var change int
	for i, key := range sample.Take(n, testData.numRecords) {
		keys[i] = key
		rec := ut.Lookup("test_table", 0, intKey(key))
		offsets[i] = rec.Off
		if i == n-1 || rand.IntN(2) == 1 {
			if i == n-1 {
				change = -sum
			} else {
				change = rand.IntN(11) - 5
			}
			sum += change
			changes[i] = change
			count := ToInt(rec.GetVal(1))
			newRecord := int2rec(key, count+change)
			ut.Update(nil, "test_table", offsets[i], newRecord)
		}
	}
	assert.This(sum).Is(0)

	// long transactions to increase outstanding
	// if rand.IntN(1001) == 55 {
	// 	time.Sleep(time.Duration(rand.IntN(10)) * time.Second)
	// }
	if ut.Complete() != "" {
		return false // aborted
	}

	// Transaction committed successfully, update mirror data
	testData.mutex.Lock()
	for i, key := range keys {
		testData.mirrorData[key] += changes[i]
	}
	testData.mutex.Unlock()
	return true
}

func int2rec(key, count int) Record {
	var rb RecordBuilder
	rb.Add(IntVal(key))
	rb.Add(IntVal(count))
	return rb.Build()
}

func intKey(key int) string {
	return Pack(IntVal(key))
}

func verifyData(t *testing.T, db *Database, testData *RandomTestData) {
	// Read all records from database
	rt := db.NewReadTran()

	dbSum := 0
	mirrorSum := 0

	it := index.NewOverIter("test_table", 0)
	for i := 0; i < testData.numRecords; i++ {
		it.Next(rt)
		_, off := it.Cur()
		rec := rt.GetRecord(off)
		count := ToInt(rec.GetVal(1))
		assert.T(t).Msg(i).This(count).Is(testData.mirrorData[i])
		dbSum += count
		mirrorSum += testData.mirrorData[i]
	}
	assert.T(t).Msg("dbSum").This(dbSum).Is(0)
	assert.T(t).Msg("mirrorSum").This(mirrorSum).Is(0)
}

func memory() uint64 {
	sample := make([]metrics.Sample, 1)
	sample[0].Name = "/gc/heap/allocs:bytes"
	metrics.Read(sample)
	return sample[0].Value.Uint64()
}
