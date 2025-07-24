// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// NOTE: This is a temporary way to bypass the library check.
// NOTE: Onec suneido.js (sujslib) is used by everyone, this can be removed
class
	{
	CallGlobal(@args)
		{
		if not Sys.SuneidoJs?()
			return false

		Global(args[0])(@+1args)
		}

	CallOnRenderBackend(@args)
		{
		if not Sys.SuneidoJs?()
			return false

		(Global('SuRenderBackend')()[args[0]])(@+1args)
		}
	}