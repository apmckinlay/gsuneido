// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta"
)

type void = struct{}

const chanBuffers = 4 // ???

// StartConcur starts the database pipeline -
// starting the goroutines and connecting them with channels.
//
// checker -> merger
//
// persist is called by merger every persistInterval
//
// To stop we close the checker channel, and then each following stage
// closes its output channel.
// Finally the merger closes the allDone channel
// so we know the shutdown has finished.
func StartConcur(db *Database, persistInterval time.Duration) {
	mergeChan := make(chan todo, chanBuffers)
	allDone := make(chan void)
	go merger(db, mergeChan, persistInterval, allDone)
	db.ck = StartCheckCo(db, mergeChan, allDone)
}

type mergeT struct {
	db        *Database
	mergeChan chan todo
	merges    *mergeList
	em        *execMulti
}

func merger(db *Database, mergeChan chan todo,
	persistInterval time.Duration, allDone chan void) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			log.Fatalln("FATAL ERROR in merger:", e)
		}
	}()
	em := startMergeWorkers()
	ep := startExecPersistMulti()
	// ep := &execPersistSingle{}
	merges := &mergeList{}
	mt := mergeT{db: db, mergeChan: mergeChan, merges: merges, em: em}
	ticker := time.NewTicker(persistInterval)
	prevState := db.GetState()
loop:
	for {
		select {
		case m := <-mergeChan: // receive todo's from commit
			if m.isZero() { // channel closed
				break loop
			}
			mt.dispatch(m)
		case <-ticker.C:
			if db.GetState() != prevState {
				prevState = db.Persist(ep, false)
			}
		}
	}
	close(em.jobChan)
	if db.GetState() != prevState {
		db.Persist(ep, false)
	}
	close(allDone)
}

type todo struct {
	tables []string
	meta   *meta.Meta
	sync   chan struct{}
}

func (td todo) isZero() bool {
	return td.tables == nil && td.meta == nil && td.sync == nil
}

func (mt *mergeT) dispatch(m todo) {
	for {
		if m.sync != nil {
			m.sync <- struct{}{}
			return
		}
		mt.merges.start(m)
		m = mt.merges.drain(mt.mergeChan)
		mt.db.Merge(mt.merges.meta, mt.em.merge, mt.merges)
		// mt.db.Merge(mergeSingle, merges)
		if m.isZero() {
			return
		}
	}
}

// mergeSingle is a single threaded merge for tran_test
func mergeSingle(metaWas, metaCur *meta.Meta, merges *mergeList) []meta.MergeUpdate {
	var results []meta.MergeUpdate
	for _, tn := range merges.tn {
		result := metaCur.Merge(metaWas, tn.table, tn.nmerge)
		if !result.Skip() {
			results = append(results, result)
		}
	}
	return results
}

const nMergeWorkers = 8 // ???

type execMulti struct {
	jobChan    chan job
	resultChan chan meta.MergeUpdate
}

func startMergeWorkers() *execMulti {
	em := &execMulti{
		jobChan:    make(chan job, 1),
		resultChan: make(chan meta.MergeUpdate, 1),
	}
	for i := 0; i < nMergeWorkers; i++ {
		go em.worker()
	}
	return em
}

func (em *execMulti) merge(metaWas, metaCur *meta.Meta, merges *mergeList) []meta.MergeUpdate {
	// if only one table, just merge it in this thread
	// and avoid overhead of channels and worker
	if len(merges.tn) == 1 {
		m := merges.tn[0]
		result := metaCur.Merge(metaWas, m.table, m.nmerge)
		if !result.Skip() {
			return append(merges.results, result)
		}
	}
	for i := 0; i < len(merges.tn); {
		select {
		case em.jobChan <- job{metaCur: metaCur, metaWas: metaWas,
			table: merges.tn[i].table, nmerge: merges.tn[i].nmerge}:
			i++
		case result := <-em.resultChan:
			merges.results = append(merges.results, result)
		}
	}
	for len(merges.results) < len(merges.tn) {
		result := <-em.resultChan
		merges.results = append(merges.results, result)
	}
	return merges.results
}

type job struct {
	metaWas *meta.Meta
	metaCur *meta.Meta
	table   string
	nmerge  int
}

func (em *execMulti) worker() {
	for j := range em.jobChan {
		result := j.metaCur.Merge(j.metaWas, j.table, j.nmerge)
		if !result.Skip() {
			em.resultChan <- result
		}
	}
}

type mergeList struct {
	meta    *meta.Meta
	tn      []tableCount
	results []meta.MergeUpdate
}

type tableCount struct {
	table  string
	nmerge int
}

func (ml *mergeList) start(m todo) {
	ml.meta = m.meta
	ml.tn = ml.tn[:0]
	ml.results = ml.results[:0]
	ml.add(m.tables)
}

func (ml *mergeList) add(tables []string) {
outer:
	for _, table := range tables {
		for i := range ml.tn {
			if ml.tn[i].table == table {
				ml.tn[i].nmerge++
				continue outer
			}
		}
		ml.tn = append(ml.tn, tableCount{table: table, nmerge: 1})
	}
}

// drain returns the next message that can't be added to the mergeList
// and must be processed separately
func (ml *mergeList) drain(mergeChan chan todo) todo {
	for {
		select {
		case m := <-mergeChan:
			if m.isZero() { // channel closed
				return todo{}
			}
			if m.sync == nil && ml.meta.SameSchemaAs(m.meta) {
				ml.add(m.tables)
			} else {
				return m // not added to merge
			}
		default: // channel empty
			return todo{}
		}
	}
}

// ------------------------------------------------------------------

type execPersist interface {
	Submit(fn func() meta.PersistUpdate)
	Results() []meta.PersistUpdate
}

type execPersistSingle struct {
	results []meta.PersistUpdate
}

func (ep *execPersistSingle) Submit(fn func() meta.PersistUpdate) {
	result := fn()
	ep.results = append(ep.results, result)

}

func (ep *execPersistSingle) Results() []meta.PersistUpdate {
	results := ep.results
	ep.results = ep.results[:0]
	return results
}

//-------------------------------------------------------------------

type execPersistMulti struct {
	count      int
	results    []meta.PersistUpdate
	workChan   chan func() meta.PersistUpdate
	resultChan chan meta.PersistUpdate
}

const nPersistWorkers = 8 // ???

func startExecPersistMulti() *execPersistMulti {
	workChan := make(chan func() meta.PersistUpdate, 1)
	resultChan := make(chan meta.PersistUpdate, 1)
	for i := 0; i < nPersistWorkers; i++ {
		go persistWorker(workChan, resultChan)
	}
	return &execPersistMulti{workChan: workChan, resultChan: resultChan}
}

func (ep *execPersistMulti) Submit(fn func() meta.PersistUpdate) {
	for {
		select {
		case ep.workChan <- fn:
			ep.count++
			return
		case result := <-ep.resultChan:
			ep.count--
			ep.results = append(ep.results, result)
		}
	}
}

func (ep *execPersistMulti) Results() []meta.PersistUpdate {
	for ; ep.count > 0; ep.count-- {
		result := <-ep.resultChan
		ep.results = append(ep.results, result)
	}
	results := ep.results
	ep.results = ep.results[:0]
	return results
}

func persistWorker(
	workChan chan func() meta.PersistUpdate,
	resultChan chan meta.PersistUpdate) {
	for fn := range workChan {
		resultChan <- fn()
	}
}
