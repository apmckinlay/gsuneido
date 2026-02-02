// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
_BookModel
	{
	GetAllChildren(table)
		{
		// Used for ServerEval: the JS client executes this code on the server.
		// It references the same server-side object, so a copy is required
		// to avoid conflicts with access from the EXE client.
		return super.GetAllChildren(table).Copy()
		}
	}