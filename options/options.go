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
