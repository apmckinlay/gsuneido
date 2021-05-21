// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"sync"
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var timestamp SuDate
var tsLock sync.Mutex

func StartTimestamps() {
	timestamp = Now()
	go ticker()
}

func ticker() {
	for {
		time.Sleep(1 * time.Second)
		t := Now()
		tsLock.Lock()
		if t.Compare(timestamp) > 0 {
			timestamp = t
		}
		tsLock.Unlock()
	}
}

func Timestamp() SuDate {
	tsLock.Lock()
	defer tsLock.Unlock()
	timestamp = timestamp.Increment()
	return timestamp
}
