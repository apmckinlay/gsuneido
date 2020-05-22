// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build interactive

// model of commit/merge/persist
//
// []client()
//		->commitChan	commitMsg = set of index overlays
// commit()
//		->mergeChan		mergeMsg
// merge()
//		->persistChan *dbState = writable state
// persist()
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

func TestModel(*testing.T) {
	start()
}

// tranNum is a sequential transaction number
type tranNum int32

func (t tranNum) String() string {
	return "T" + strconv.Itoa(int(t))
}

// value represents an index value
type value int32

func (v value) String() string {
	return "V" + strconv.Itoa(int(v))
}

type indexBase []value

// indexNum identifies an index
type indexNum int

func (i indexNum) String() string {
	return string('a' + i)
}

// commitMsg is sent by clients to the commit process
type commitMsg map[indexNum]value

// mergeMsg is sent by the commit process to the index mergers
type mergeMsg struct {
	state *dbState
	tran  tranNum
	cm    commitMsg
}

type persistMsg struct{}

const nIndexes = 13

const nClients = 4

// duration of simulated work
const (
	testDuration    = 3000 * timeUnits
	mergeDelay      = 8
	clientDelay     = 200 // avg 100 / 4 clients = commit every ~25
	persistDelay    = 50
	persistInterval = 400
)

// buffers is the channel buffer size
const buffers = 4

// terminate is used to signal the clients to end.
// It should be accessed atomically.
var terminate int32

// allDone is closed by persist when it finishes shutdown
var allDone = make(chan bool)

func start() {
	commitChan := make(chan commitMsg, buffers)
	mergeChan := make(chan *mergeMsg, buffers)
	persistChan := make(chan persistMsg, buffers)

	for i := range vals {
		vals[i] = int32(i * 100)
	}
	rand.Seed(time.Now().UnixNano()) // BEWARE of test caching
	var wg sync.WaitGroup
	for i := 0; i < nClients; i++ {
		wg.Add(1)
		go client(commitChan, &wg)
	}
	go merge(mergeChan, persistChan)
	go persist(persistChan)
	go commit(commitChan, mergeChan)
	time.Sleep(testDuration)
	atomic.AddInt32(&terminate, nClients)
	wg.Wait()
	close(commitChan)
	<-allDone
}

// vals is the next value number for each index
// must be accessed atomically
var vals [nIndexes]int32

//-------------------------------------------------------------------

// client models a database client.
// It sends commitMsg to commitChan to the commit process
func client(commitChan chan<- commitMsg, wg *sync.WaitGroup) {
	defer wg.Done()
	for atomic.LoadInt32(&terminate) == 0 {
		randSleep(clientDelay)
		msg := make(map[indexNum]value)
		ni := 1 + rand.Intn(nIndexes/2)
		for i := 0; i < nIndexes && ni > 0; i++ {
			if rand.Int()%2 == 0 {
				msg[indexNum(i)] = value(atomic.AddInt32(&vals[i], 1))
				ni--
			}
		}
		commitChan <- commitMsg(msg)
	}
}

//-------------------------------------------------------------------

// commit completes transactions.
// It is single threaded, serializes transactions.
//
// - receive commitMsg's from clients
//
// - add overlays to state
//
// - send mergeMsg's to merge's
func commit(commitChan <-chan commitMsg, mergeChan chan<- *mergeMsg) {
	var tran tranNum
	// receive commitMsg from clients
	for cm := range commitChan {
		t := tran
		tran++
		println("commit", t, cm, "commitChan", len(commitChan))
		// add overlays to state
		state := updateState(func(state *dbState) {
			for i, val := range cm {
				idx := &state.indexes[i]
				idx.overlays = append(idx.overlays, val)
			}
		})
		// notify merger
		mergeChan <- &mergeMsg{state, t, cm}
	}
	close(mergeChan)
}

//-------------------------------------------------------------------

// merge models an index merger.
// This only changes the representation, not any actual data.
//
// - receives mergeMsg's from commit
//
// - updates database state with results of merges
func merge(mergeChan <-chan *mergeMsg, persistChan chan<- persistMsg) {
	ticker := time.NewTicker(persistInterval * timeUnits)
loop:
	for {
		select {
		case m := <-mergeChan: // receive mergeMsg's from commit
			if m == nil {
				break loop
			}
			println(indent(1)+"merge", m.tran, "mergeChan", len(mergeChan))
			newBases := map[indexNum]indexBase{}
			for idx := range m.cm { // NOTE: could merge indexes in parallel
				newBases[idx] = merge1(m, idx)
			}
			updateState(func(state *dbState) {
				for idx, newBase := range newBases {
					updateIndexState(&state.indexes[idx], newBase)
				}
				state.lastMergedTran = m.tran
			})
		case <-ticker.C:
			// send ticks from here so we get back pressure
			persistChan <- struct{}{}
		}
	}
	persistChan <- struct{}{}
	close(persistChan)
}

func merge1(m *mergeMsg, idx indexNum) indexBase {
	println(indent(2)+"merge", m.tran, idx)
	randSleep(mergeDelay)
	return merge2(&m.state.indexes[idx])
}

// merge2 merges the first overlay into the base for a single index state
// and returns the new indexState
func merge2(iState *indexState) indexBase {
	newBase := make([]value, len(iState.base), len(iState.base)+1)
	return append(newBase, iState.overlays[0])
}

// updateIndexState updates the global database state with the result of an index merge
func updateIndexState(iState *indexState, newBase indexBase) {
	iState.overlays = iState.overlays[1:]
	iState.base = newBase
}

//-------------------------------------------------------------------

// persist models the writing index updates to disk.
// It receives new states from merge and at intervals
// writes the index updates and then a new root state.
func persist(persistChan <-chan persistMsg) {
	var lastTran tranNum
	for range persistChan {
		state := getState()
		if state.lastMergedTran == lastTran {
			continue
		}
		lastTran = state.lastMergedTran
		println(strings.Repeat("^", 40)+"PERSIST up to", state.lastMergedTran,
			"persistChan", len(persistChan))
		randSleep(persistDelay)
		println(strings.Repeat(".", 40) + "PERSIST end")
		// NOTE: could persist indexes in parallel
	}
	close(allDone)
}

//-------------------------------------------------------------------

var begin = time.Now()

const timeUnits = time.Millisecond

func randSleep(ms int) {
	time.Sleep(time.Duration(rand.Intn(ms)) * timeUnits)
}

func println(args ...interface{}) {
	args2 := make([]interface{}, 0, 20)
	args2 = append(args2, fmt.Sprintf("%4d ", time.Since(begin)/timeUnits))
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
	indexes        [nIndexes]indexState
	lastMergedTran tranNum
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
