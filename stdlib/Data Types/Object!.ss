// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function(@args)
	{
	x = Function?(args.GetDefault(0, false)) ? (args[0])(@+1args) : args[0]
	if not Object?(x)
		throw "expected an object but it was " $ Type(x)
	return x
	}
