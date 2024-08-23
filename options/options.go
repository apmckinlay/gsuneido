// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package options contains configuration options
// including command line flags
package options

import (
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/util/generic/atomics"
	"github.com/apmckinlay/gsuneido/util/regex"
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
	IgnoreVersion  bool
	WebServer      bool
	WebPort        string
	TimeoutMinutes = 2 * 60 // 2 hours
	Passphrase     string   // used with -load
)

// StrictCompare determines whether comparisons between different types
// are allowed (old behavior) or throw an exception (new behavior)
// NOT thread safe, but don't want overhead on every compare
var (
	StrictCompare   = false
	StrictCompareDb = false
)

// CmdLine is the remaining command line arguments
var CmdLine string

// debugging options
const (
	ThreadDisabled        = false
	TimersDisabled        = false
	ClearCallbackDisabled = false
)

var (
	AllWarningsThrow = regex.Compile("")
	NoWarningsThrow  = regex.Compile(`\A\Z`)
	WarningsThrow    atomics.Value[regex.Pattern]
)

func init() {
	WarningsThrow.Store(NoWarningsThrow)
}

// Coverage controls whether Cover op codes are added by codegen.
var Coverage atomic.Bool

var Nworkers = func() int {
	return min(8, max(1, runtime.NumCPU()-1)) // ???
}()

// DbStatus should be set to one of:
// - "" (running normally)
// - starting
// - repairing
// - corrupted
var DbStatus atomics.String

func init() {
	DbStatus.Store("starting")
}

func BuiltStr() string {
	goos := strings.Replace(runtime.GOOS, "darwin", "macos", 1)
	return BuiltDate + " (" + runtime.Version() + " " +
		goos + "/" + runtime.GOARCH + BuiltExtra + ")"
}
