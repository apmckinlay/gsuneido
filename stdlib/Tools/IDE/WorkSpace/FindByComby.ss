// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(name, code, search, hint = false)
		{
		lines = Object()
		if name =~ '\.(js|css)$'
			return lines

		if hint isnt false and not code.Has?(hint)
			return lines

		for match in CombyMatch(search, code)
			{
			from = code.LineFromPosition(match.pos)
			to = code.LineFromPosition(match.end - 1)
			lines.Add(Seq(from, to + 1))
			}
		return lines
		}
	}
