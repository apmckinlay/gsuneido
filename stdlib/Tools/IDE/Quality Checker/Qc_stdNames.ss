// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(fromStdlib?)
		{
		// returns the list as members because .Member?() is faster than .Has?()
		stdNames = BuiltinNames().ListToMembers()
		if not fromStdlib?
			for name in
				QueryList('stdlib where group = -1 and not name.Suffix?("Test")', #name)
				stdNames[name] = true
		return stdNames.Set_readonly() // readonly so Memoize doesn't copy it
		}
	}