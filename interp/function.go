package interp

import . "github.com/apmckinlay/gsuneido/core/value"

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
	code []byte
	// nparams is the number of values required on the stack
	nparams int
	// nlocals is the number of parameters and local variables
	nlocals int
	// strings starts with the parameters, then the locals,
	// and then any argument or member names used in the code,
	// and any argument specs
	strings []string
	values  []Value
}

// TODO: Function needs to be a Value
