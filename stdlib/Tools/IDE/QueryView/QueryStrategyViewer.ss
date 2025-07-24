// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Top: 0
	Xstretch: 1
	Ystretch: 1
	New(.query, .q_or_c)
		{
		super([ScintillaIDEControl { Addon_suneido_style: (query:) },
			set: .strategy(.query, .q_or_c)
			xmin: 100, ymin: 100 readonly:])
		}
	strategy(query, q_or_c)
		{
		try
			text = QueryStrategyAndWarnings(query, q_or_c is "Cursor")
		catch (e)
			text = e
		return text
		}
	}