// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function(@args)
	{
	x = Function?(args.GetDefault(0, false)) ? (args[0])(@+1args) : args[0]
	if not Date?(x)
		throw "expected a date but it was " $ Type(x)
	return x
	}
