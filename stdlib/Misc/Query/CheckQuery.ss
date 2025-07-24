// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(query, queryInstance = false)
		{
		try
			strategy = QueryStrategy(queryInstance isnt false ? queryInstance : query)
		catch (err)
			return err
		return Join(", ",
			.check(query, strategy, "UNION NOT DISJOINT", "(?i)union-(merge|lookup)"),
			.check(query, strategy, "PROJECT NOT UNIQUE", "(?i)project-(map|seq)"),
			.check(query, strategy, "JOIN MANY TO MANY", "(?i)join n:n"))
		}

	check(query, strategy, ref, pat)
		{
		if String?(query) and query.Has?("CHECKQUERY SUPPRESS: " $ ref)
			return ""
		return strategy =~ pat ? ref : ""
		}
	}
