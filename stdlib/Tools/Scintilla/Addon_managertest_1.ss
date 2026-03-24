// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	For: Test
	Save(s)
		{
		.Parent.S = s
		}

	Get()
		{
		return 'Parent: ' $ .Parent.S
		}
	}