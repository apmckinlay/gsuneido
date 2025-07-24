// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
PatternControl
	{
	Name: "ZipPostal"
	New(mandatory = false, readonly = false, hidden = false, tabover = false)
		{
		super(.Pattern(), width: 10,
			status: "A zip or postal code e.g. 12345, 12345-1234, or S7K 5G8",
			:mandatory, :readonly, :hidden, :tabover)
		}
	Pattern()
		{
		return 'A#A #A#|#####|#####-####'
		}
	}