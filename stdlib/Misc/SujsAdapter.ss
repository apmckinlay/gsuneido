// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
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

		((SuRenderBackend())[args[0]])(@+1args)
		}
	}