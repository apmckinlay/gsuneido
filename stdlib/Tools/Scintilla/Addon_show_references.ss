// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{ #("Find References\tF11") }
	On_Find_References()
		{
		FindReferencesControl(.GetCurrentWord())
		}
	}
