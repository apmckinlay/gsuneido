// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(parent, options)
		{
		.parent = parent
		.setOptions(options)
		}
	setOptions(options)
		{
		if not Object?(options)
			return
		for m in options.Members()
			this[m.Capitalize()] = options[m]
		}
	Getter_Parent()
		{ // provide readonly .Parent
		return .parent
		}
	Default(@args)
		{
		.parent[args[0]](@+1args)
		}
	}