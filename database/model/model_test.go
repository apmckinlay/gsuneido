// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build interactive

package model

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

// model of commit/merge/checkpoint/persist
//
// []client()
//		->commitChan	commitMsg = set of index overlays
// commit()
//		->mergeChan		mergeMsg = which index to merge
//		<-cpStateChan	*dbState = valid checkpoint state
// []merge()
//
//		->checkpointChan *dbState = state to write
// checkpoint()
//		->persistInChan	persistInMsg = which index to persist
//		->persistOutChan persistOutMsg = new index root
// []persist()

func TestModel(*testing.T) {
	start()
}

// trans identifies a client
type trans int32

func (t trans) String() string {
	return "T" + strconv.Itoa(int(t))
}

// value represents an index value
type value int32

func (v value) String() string {
	return "V" + strconv.Itoa(int(v))
}

// index identifies an index
type index int

func (i index) String() string {
	return string('a' + i)
}

// commitMsg is sent by clients to the commit process
type commitMsg map[index]value

// mergeMsg is sent by the commit process to the index mergers
type mergeMsg struct {
	tran  trans
	idx   index
	state *dbState
}

const nIndexes = 13

// number of goroutines
const (
	nClients  = 4
	nMerges   = 4
	nPersists = 2
)

// duration of simulated work
const (
	timeUnits       = time.Millisecond
	testDuration    = 3000 * timeUnits
	commitDelay     = 10
	mergeDelay      = 20
	clientDelay     = 200
	checkpointDelay = 50
)

// buffers is the input channel buffer size for merges and persists
// a buffer size of 4 seems to make a modest (~10%) increase in throughput
const buffers = 4

func start() {
	commitChan := make(chan commitMsg)
	mergeChan := make(chan mergeMsg, buffers)
	cpStateChan := make(chan *dbState, 1)
	checkpointChan := make(chan *dbState)

	for i := range vals {
		vals[i] = int32(i * 100)
	}
	rand.Seed(time.Now().UnixNano()) // BEWARE of test caching
	for i := 0; i < nClients; i++ {
		go client(commitChan)
	}
	for i := 0; i < nMerges; i++ {
		go merge(mergeChan, cpStateChan)
	}
	go checkpoint(checkpointChan)
	commit(commitChan, mergeChan, cpStateChan, checkpointChan)
}

// tran is the next client number
// must be accessed atomically
var tran int32

// vals is the next value number for each index
// must be access atomically
var vals [nIndexes]int32

// client models a database client
// sends commitMsg to commitChan to the commit process
func client(commitChan chan<- commitMsg) {
	for {
		randSleep(clientDelay)
		msg := make(map[index]value)
		ni := 1 + rand.Intn(nIndexes/2)
		for i := 0; i < nIndexes && ni > 0; i++ {
			if rand.Int()%2 == 0 {
				msg[index(i)] = value(atomic.AddInt32(&vals[i], 1))
				ni--
			}
		}
		commitChan <- commitMsg(msg)
	}
}

// commit models the commit process
// - receive commitMsg's from clients
// - add overlays to state
// - send mergeMsg's to merge's, and dbState to checkpoint
func commit(commitChan <-chan commitMsg, mergeChan chan<- mergeMsg,
	cpStateChan <-chan *dbState, checkpointChan chan<- *dbState) {
	var cpState *dbState
	ticker := time.Tick(1000 * timeUnits)
	terminator := time.NewTimer(testDuration)
	for {
		var cm commitMsg
		select {
		case cm = <-commitChan: // commitMsg from clients
			t := trans(atomic.AddInt32(&tran, 1))
			println("commit", t, cm)
			state := updateState(func(state *dbState) {
				for i, val := range cm {
					// add overlay to state
					idx := &state.indexes[i]
					idx.overlays = append(idx.overlays, val)
					randSleep(commitDelay)
				}
			})
			// notify the relevant merges
			for i := range cm {
				atomic.AddInt32(&activeMergeCount, +1)
				mergeChan <- mergeMsg{t, i, state}
			}
		case cpState = <-cpStateChan:
			// just accept new cpState
		case <-ticker:
			// send ticks from here so commit will block if checkpoint slow
			if cpState == nil {
				// haven't had a lull so need to force one
				println(strings.Repeat("-", 70), "force")
				cpState = <-cpStateChan
			}
			checkpointChan <- cpState
			cpState = nil
		case <-terminator.C:
			if atomic.LoadInt32(&activeMergeCount) > 0 {
				cpState = <-cpStateChan
				checkpointChan <- cpState
				//TODO wait for checkpoint to finish
			}
			return
		}
	}
}

// merge models an index merger
// - receives mergeMsg's from commit
// - updates database state with results of merges
func merge(mergeChan <-chan mergeMsg, cpStateChan chan<- *dbState) {
	// locks prevent two merges on the same index at the same time
	var locks [nIndexes]sync.Mutex

	// receive mergeMsg's from commit
	for m := range mergeChan {
		merge1(m, locks[m.idx], cpStateChan)
	}
}

// activeMergeCount is used to detect when all merges are complete
// to set cpState
var activeMergeCount int32

func merge1(m mergeMsg, lock sync.Mutex, cpStateChan chan<- *dbState) {
	// lock to prevent two merges on the same index at the same time
	lock.Lock()
	defer lock.Unlock()
	randSleep(mergeDelay)
	newBase := merge2(&m.state.indexes[m.idx])
	println(indent(1+int(m.idx))+"merge", m.tran, m.idx)
	updateState(func(state *dbState) {
		updateIndexState(&state.indexes[m.idx], newBase)
		if 0 == atomic.AddInt32(&activeMergeCount, -1) {
			println(strings.Repeat("-", 70), "cpstate")
			cpStateChan <- state
		}
	})
}

// merge2 merges the first overlay into the base for a single index state
// and returns the new indexState
func merge2(iState *indexState) []value {
	newBase := make([]value, len(iState.base), len(iState.base)+1)
	return append(newBase, iState.overlays[0])
}

// updateIndexState updates the global database state with the result of an index merge
func updateIndexState(iState *indexState, newBase []value) {
	iState.overlays = iState.overlays[1:]
	iState.base = newBase
}

// checkpoint models the checkpoint process that writes index updates to disk.
// At intervals it gets the latest merged state (set by merges)
// concurrently persists indexes,
// and then writes a new root state
func checkpoint(checkpointChan <-chan *dbState) {
	for {
		<-checkpointChan
		println(strings.Repeat("=", 70), "checkpoint")
		//TODO get the current persistable state
		//TODO send a message for each index that requires updating
		//TODO wait for all the results
		//TODO write a new root state
	}
}

func randSleep(ms int) {
	time.Sleep(time.Duration(rand.Intn(ms)) * timeUnits)
}

var begin = time.Now()

func println(args ...interface{}) {
	args2 := make([]interface{}, 0, 20)
	args2 = append(args2, fmt.Sprintf("%4d ", time.Now().Sub(begin)/timeUnits))
	args2 = append(args2, args...)
	fmt.Println(args2...)
}

func indent(i int) string {
	return strings.Repeat(" ", 4*i)
}

// state ------------------------------------------------------------

// state is the global in-memory database state, immutable (copy on write).
// Should only be accessed by getState and updateState
var state = unsafe.Pointer(&dbState{})

// stateMutex guards updates to the state
var stateMutex sync.Mutex

// getState returns the current global in-memory database state
func getState() *dbState {
	return (*dbState)(atomic.LoadPointer(&state))
}

// updateState applies the given update function to a copy of the state
// and sets the state to the result.
// Note: the state passed to the update function is a *shallow* copy,
// it is up to the function to make copies of any nested containers.
// state.indexes is an array (not a slice) so it will already be copied.
func updateState(fn func(*dbState)) *dbState {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	newState := *(*dbState)(state) // copy (including indexes array)
	fn(&newState)
	atomic.StorePointer(&state, unsafe.Pointer(&newState))
	return &newState
}

// dbState must be immutable, updated by copy on write
type dbState struct {
	indexes [nIndexes]indexState
}

// indexState is the state of an index
type indexState struct {
	// modified is set to true if the index has been updated
	modified bool
	// overlays are what is added by each transaction
	// for simplicity each overlay is a single value
	overlays []value
	// base represents the primary (persisted) index
	base []value
}
