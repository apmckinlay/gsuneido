// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (query, cursor = false)
	{
	f = (cursor ? Cursor : WithQuery)
	warnings = ""
	try
		f(query)
			{|q|
			strategy = QueryStrategy(q, formatted:)
			}
	catch (e, "PROJECT NOT UNIQUE|UNION NOT DISJOINT|JOIN MANY TO MANY")
		{
		warnings = e.Replace(',', ', ')
		query = QuerySuppress(query)
		f(query)
			{|q| strategy = QueryStrategy(q, formatted:) }
		}
	tempindex = strategy.Has?("TEMPINDEX") ? "TEMPINDEX" : ""
	return Opt("WARNING: ", Join(", ", warnings, tempindex), "\n\n") $ strategy
	}