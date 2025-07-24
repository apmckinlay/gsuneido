// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// similar to Matcher_like
// but collapses runs of whitespace (including newlines) to single spaces
Matcher_like
	{
	Tr(s)
		{
		return s.Trim().Tr(' \t\r\n', ' ')
		}
	}