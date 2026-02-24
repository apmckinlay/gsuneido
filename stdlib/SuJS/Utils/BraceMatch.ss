// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	braces: ('(', ')', '[', ']', '{', '}', '<', '>')
	CallClass(s, pos)
		{
		c = s[pos]
		if false is i = .braces.Find(c)
			return -1
		return i % 2 is 0 // open
			? .searchBrace(s, c, .braces[i + 1],
				start: pos + 1, end: s.Size())
			: .searchBrace(s, c, .braces[i - 1],
				start: pos - 1, end: -1, step: -1)
		}

	searchBrace(s, self, target, start, end, step = 1)
		{
		c = 0
		for (i = start; step is 1 ? i < end : i > end; i += step)
			{
			if s[i] is self
				c++
			else if s[i] is target
				{
				if c > 0
					c--
				else
					return i
				}
			}
		return -1
		}
	}