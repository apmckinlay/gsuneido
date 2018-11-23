package runtime

// Frame is the context for a function/method/block invocation.
type Frame struct {
	// fn is the Function being executed
	fn *SuFunc
	// ip is the current index into the Function's code
	ip int
	// bp is the base pointer to the locals
	bp   int
	self Value //TODO self
}
