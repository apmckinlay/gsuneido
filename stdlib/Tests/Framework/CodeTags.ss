// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Matches(code)
		{
		return .MatchTags(.ExtractTags(code))
		}
	ExtractTags(code)
		{
		s = code.Extract("^// ?TAGS: (.*)")
		return s is false ? Object() : s.Trim().Tr(' \t', ' ').Split(' ')
		}
	MatchTags(codetags, systags = false) // systags parameter is for testing
		{
		if systags is false
			systags = Sys.Systags()
		for tag in codetags
			if tag[0] is '!'
				{
				if systags.Has?(tag[1..])
					return false
				}
			else if not systags.Has?(tag)
				return false
		return true
		}
	}