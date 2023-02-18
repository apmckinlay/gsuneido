// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package options contains configuration options
// including command line flags
package options

import (
	"runtime"
	"sync/atomic"

	myatomic "github.com/apmckinlay/gsuneido/util/generic/atomic"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
)

var BuiltDate string
var BuiltExtra string

// command line flags
var (
	Mode           string
	Action         string
	Error          string
	Arg            string
	Port           string
	Unattended     bool
	IgnoreVersion  bool
	TimeoutMinutes = 2 * 60 // 2 hours
)

// StrictCompare determines whether comparisons between different types
// are allowed (old behavior) or throw an exception (new behavior)
var StrictCompare = false

// CmdLine is the remaining command line arguments
var CmdLine string

// log file names, port is added when client
var (
	Errlog = "error.log"
)

// debugging options
const (
	ThreadDisabled        = false
	TimersDisabled        = false
	ClearCallbackDisabled = false
)

// Coverage controls whether Cover op codes are added by codegen.
var Coverage atomic.Bool

var Nworkers = func() int {
	return ord.Min(8, ord.Max(1, runtime.NumCPU()-1)) // ???
}()

// DbStatus should be set to one of:
// - "" (running normally)
// - starting
// - repairing
// - corrupted
var DbStatus myatomic.String

func init() {
	DbStatus.Store("starting")
}

func BuiltStr() string {
	return BuiltDate +
		" (" + runtime.Version() + " " + runtime.GOARCH + BuiltExtra + ")"
}
