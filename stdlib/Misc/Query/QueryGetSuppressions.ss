// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (query)
	{
	qsups = Object()
	for sup in #(
		"CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE",
		"CHECKQUERY SUPPRESS: UNION NOT DISJOINT",
		"CHECKQUERY SUPPRESS: JOIN MANY TO MANY")
		if query.Has?(sup)
			qsups.Add("/* " $ sup $ "*/")
	return qsups
	}