// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
SvcDisabledTables
	{
	Tables()
		{
		return LibraryTables().Remove(@SvcControl.SvcExcludeLibraries)
		}

	ResetCache()
		{
		LibraryTables.ResetCache()
		super.ResetCache()
		}
	}