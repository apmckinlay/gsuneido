// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// for global and dynamic variables
function (name)
	{
	Assert(name.GlobalName?() or name.DynamicName?())
	try
		name.Eval() // needs to be Eval for dynamic names
	catch (unused, "can't find|uninitialized|error loading")
		return true
	return false
	}
