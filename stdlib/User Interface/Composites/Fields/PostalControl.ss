// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
PatternControl
	{
	Name: "Postal"
	New(mandatory = false)
		{
		super(.Pattern(), width: 10, :mandatory, status: "A postal code e.g. S7K 5G8")
		}

	Pattern()
		{
		return 'A#A #A#'
		}
	}