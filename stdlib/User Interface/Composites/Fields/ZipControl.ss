// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
PatternControl
	{
	Name: "Zip"
	New(mandatory = false)
		{
		super(.Pattern(), width: 10, :mandatory,
			status: "A zip code e.g. 12345 or 12345-1234")
		}

	Pattern()
		{
		return '#####|#####-####'
		}
	}