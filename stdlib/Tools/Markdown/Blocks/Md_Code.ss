// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	New(line = false, info = '')
		{
		.Codes = []
		.Info = Md_Helper.Escape(info)
		if line isnt false
			.Codes.Add(line)
		}

	Add(line, start)
		{
		.Codes.Add(line[start..])
		}
	}