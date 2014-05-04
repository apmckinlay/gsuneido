package interp

import . "github.com/apmckinlay/gsuneido/value"

/*
Function is a compiled function, method, or block.

Parameters at the start of names may be prefixed:
'@' for each, '_' for dynamic, '.' for member, or '=' for default.

There can only be one '@' parameter and if present it must be last.

Parameters with default values ('=' or '_')
must come after parameters without defaults.

'.' parameter names for public members are capitalized.
*/
type Function struct {
	Code []byte
	// nparams is the number of values required on the stack
	Nparams int
	// nlocals is the number of parameters and local variables
	Nlocals int
	// strings starts with the parameters, then the locals,
	// and then any argument or member names used in the code,
	// and any argument specs
	Strings []string
	Values  []Value
}

// TODO: Function needs to be a Value
