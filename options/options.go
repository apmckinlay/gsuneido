// Package options contains configuration options
// including command line flags
package options

// command line flags
var (
	BuiltDate string
	Repl      bool
	Client    string
	Port      = "3147"
	Version   bool
	Help      bool
)

// CmdLine is the remaining command line arguments
var CmdLine string

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
