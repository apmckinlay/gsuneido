// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package options contains configuration options
// including command line flags
package options

var BuiltDate string

// command line flags
var (
	Mode       string
	Action     string
	Error      string
	Arg        string
	Port       string
	Unattended bool
	NoRelaunch bool
)

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
// Should be accessed atomically. Zero means disabled.
var Coverage int64

var Trace = 0

const (
	TraceFunctions = 1 << iota
	TraceStatements
	TraceOpcodes
	TraceRecords
	TraceLibraries
	TraceSlowQuery
	TraceQuery
	TraceSymbol
	TraceAllIndex
	TraceTable
	TraceSelect
	TraceTempIndex
	TraceQueryOpt

	TraceConsole
	TraceLogFile

	TraceClientServer
	TraceExceptions
	TraceGlobals

	TraceJoinOpt
)
