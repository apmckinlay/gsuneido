// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	New(line = false, .Info = '')
		{
		.Codes = []
		if line isnt false
			.Codes.Add(line)
		}

	Add(line)
		{
		.Codes.Add(line)
		}
	}