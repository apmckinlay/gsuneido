// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
TestObserverGui
	{
	New()
		{
		.Data = Object()
		}
	Values()
		{
		return Object(Data: .Data, Totals: .Totals)
		}
	SetValues(from)
		{
		for item in #(Data, Totals)
			for m in from[item].Members()
				this[item][m] = from[item][m]
		}
	}