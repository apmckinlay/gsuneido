package interp

import . "github.com/apmckinlay/gsuneido/base"

/*
Frame is the context for a function/method/block invocation.
*/
type Frame struct {
	// fn is the Function being executed
	fn *SuFunc
	// ip is the current index into the Function's code
	ip int
	// locals is used to reference the local variables
	// needs to be a slice instead of just an index
	// to handle persisted block locals which are no longer on the stack
	locals []Value
	self   Value //TODO
}
