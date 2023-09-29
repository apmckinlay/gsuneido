// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"sync"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

var timestamp SuDate
var tsLock sync.Mutex

// StartTimestamp is called by gsuneido.go openDbms
func StartTimestamps() {
	timestamp = Now().WithoutMs()
	go ticker()
}

// ticker runs on the dbms i.e. server or standalone, not client
func ticker() {
	for {
		time.Sleep(1 * time.Second)
		t := Now().WithoutMs()
		tsLock.Lock()
		if t.Compare(timestamp) > 0 {
			timestamp = t
		}
		tsLock.Unlock()
	}
}

// Timestamp is the backend
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
