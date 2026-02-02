// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Top: 0
	Xstretch: 1
	Ystretch: 1
	New(.query, .q_or_c, one = false)
		{
		super([ScintillaIDEControl { Addon_suneido_style: (query:) },
			set: .strategy(.query, .q_or_c, one)
			xmin: 100, ymin: 100 readonly:])
		}
	strategy(query, q_or_c, one)
		{
		try
			{
			text = one
				? .strategy1(query)
				: QueryStrategyAndWarnings(query, q_or_c is "Cursor")
			}
		catch (e)
			text = e
		return text
		}
	strategy1(query)
		{
		named = ""
		args = QueryToNamed(query)
		if args.Size() > 1
			named = Display(args)[1..] $ "\n\n"
		return named $ Query.Strategy1(@args)
		}
	}