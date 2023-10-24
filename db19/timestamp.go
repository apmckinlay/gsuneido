// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"sync"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

var timestamp SuDate
var tsLock sync.Mutex
var timeError = false

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
		d := t.MinusMs(timestamp)
		if d > 0 {
			if d > 5000 {
				log.Println("ERROR: time skip from", timestamp, "to", t,
					"=", time.Duration(d) * time.Millisecond)
			}
			timestamp = t // normal case
			timeError = false
		} else if d < 0 && !timeError {
			log.Println("ERROR: time went backwards from", timestamp, "to", t)
			timeError = true
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
