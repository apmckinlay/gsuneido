// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (s, pat)
	{
	lines = Object()
	s.ForEachMatch(pat)
		{
		lines.Add(Object(s.LineFromPosition(it[0][0])))
		}
	return lines
	}
