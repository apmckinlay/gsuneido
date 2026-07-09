// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
SvcDisabledTables
	{
	Tables()
		{
		return BookTables()
		}

	ResetCache()
		{
		BookTables.ResetCache()
		super.ResetCache()
		}
	}
