// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	NamePath(lib, name)
		{
		name = name.RemovePrefix('_')
		x = .getByName(lib, name)
		return LibRecGetPath(x, lib)
		}
	getByName(lib, name)
		{
		return Query1(lib, group: -1, :name)
		}
	}