// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

// rogsChan is used by other threads to Run code On the Go Side UI thread
// Need buffer so we can send to channel and then notifyCside
var rogsChan = make(chan func(), 1)

// runOnGoSide is called by interp via runtime.RunOnGoSide
// and cside via goc.RunOnGoSide
func runOnGoSide() {
	if InRunUI {
		return // don't want to reenter and run recursively
	}
	InRunUI = true
	defer func() { InRunUI = false }()
	for range 8 { // process available messages, but not forever
		select {
		case fn := <-rogsChan:
			fn()
		default: // non-blocking
			return
		}
	}
}
