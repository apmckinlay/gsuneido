// Package options contains configuration options
// including command line flags
package options

// command line flags
var (
	BuiltDate string
	Client    bool
	Repl      bool
	NetAddr   string
)

// Args are the remaining command line arguments
var Args []string

// debugging options
var (
	HeapDebug             = true
	ThreadDisabled        = false
	TimersDisabled        = false
	ClearCallbackDisabled = false
)

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
