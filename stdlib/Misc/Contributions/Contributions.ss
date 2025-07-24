// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(name)
		{
		result = Object()
		for lib in Libraries()
			{
			recname = lib.Capitalize() $ "_" $ name
			try
				{
				result.Add(Global(recname))
				if QueryEmpty?(lib, name: recname, group: -1)
					ProgrammerError("Contribution library mismatch: " $
						recname $ ' not in ' $ lib)
				}
			catch (unused, "can't find|not authorized")
				{ }
			}
		return result
		}

	Init()
		{
		LibUnload.AddObserver(#Contributions, {|name|
			if IsContribution?(name)
				.ResetCache()
			})
		}
	}
