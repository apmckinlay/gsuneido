// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

type profile struct {
	enabled bool
	// calls is the number of times the function is called
	calls map[*SuFunc]int32
	// total is the time (tsc) spent in each function (and builtins it calls)
	total map[*SuFunc]int64
	// self is total minus the time spent in functions called by this function
	self map[*SuFunc]int64
}

func (th *Thread) StartProfile() {
	th.profile.total = make(map[*SuFunc]int64)
	th.profile.self = make(map[*SuFunc]int64)
	th.profile.calls = make(map[*SuFunc]int32)
	th.profile.enabled = true
}

func (th *Thread) StopProfile() (total, self map[*SuFunc]int64, calls map[*SuFunc]int32) {
	if !th.profile.enabled {
		return
	}
	th.profile.enabled = false
	total, self, calls = th.profile.total, th.profile.self, th.profile.calls
	th.profile.total, th.profile.self, th.profile.calls = nil, nil, nil
	return total, self, calls
}
