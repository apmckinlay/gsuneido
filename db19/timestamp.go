// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"sync"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/generic/ints"
)

var timestamp SuDate
var tsLock sync.Mutex

// StartTimestamp is called by gsuneido.go openDbms
func StartTimestamps() {
	// start timestamp ahead one second
	// to avoid overlap if it restarts in less than a second.
	timestamp = Now().WithoutMs().Plus(0, 0, 0, 0, 0, 0, 990)
	go ticker()
}

// ticker runs on the dbms i.e. server or standalone, not client
func ticker() {
	prev := Now().WithoutMs()
	for {
		time.Sleep(1 * time.Second)
		t := Now().WithoutMs()
		if d := t.MinusMs(prev); ints.Abs(d) > 5000 {
			log.Println("ERROR: time skip from", prev, "to", t,
				"=", time.Duration(d)*time.Millisecond)
		}
		prev = t
		tsLock.Lock()
		if t.Compare(timestamp) > 0 {
			// only update timestamp forwards
			timestamp = t
		}
		tsLock.Unlock()
	}
}

// Timestamp is the backend. See also Thread.Timestamp
func Timestamp() SuDate {
	tsLock.Lock()
	defer tsLock.Unlock()
	ts := timestamp
	if ts.Millisecond() < TsThreshold {
		timestamp = timestamp.AddMs(TsInitialBatch)
	} else {
		timestamp = timestamp.AddMs(1)
	}
	return ts
}
