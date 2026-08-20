// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
/*
CombyReplace is
	- token-based structural search and replace
	- delegates matching to CombyMatch
	- replaces matched regions with a replace pattern that can also
		contain holes (:[name]), which get filled from the corresponding
		captures in the search pattern
*/
class
	{
	CallClass(source, search, replace, from = 0, to = false)
		{
		replaceItems = CombyTemplate(replace)

		result = source[..from]
		prevEnd = from
		while false isnt m = CombyMatch(search, source, pos: prevEnd)
			{
			if to isnt false and m.end > to
				break
			result $= source[prevEnd..m.pos]
			result $= .buildReplacement(replaceItems, m.holes)
			prevEnd = m.end
			}
		result $= source[prevEnd..]
		return result
		}

	buildReplacement(items, holes)
		{
		result = ''
		for item in items
			{
			if item.type is #HOLE
				result $= holes.Member?(item.value)
					? holes[item.value] : item.text
			else
				result $= item.text
			}
		return result
		}
	}
