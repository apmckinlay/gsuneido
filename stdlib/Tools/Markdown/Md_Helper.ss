// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// Only detab the starting substring that only contains ' ' or \t
	Detab(line)
		{
		pos = line.FindRx('[^ \t]')
		return line[..pos].Detab() $ line[pos..]
		}

	MatchHTMLTag(s)
		{
		start = 0
		if s[start::1] isnt '<'
			return false

		start++
		openTag? = true
		if s[start::1] is '/'
			{
			start++
			openTag? = false
			}

		if false is pos = .matchTagName(s, start)
			return false

		if openTag?
			{
			while false isnt end = .matchAttribute(s, pos)
				pos = end
			}

		pos = .matchSpaces(s, pos, optional?:)
		if openTag? is true and s[pos::1] is '/'
			pos++

		if s[pos::1] isnt '>'
			return false

		return pos + 1
		}

	matchTagName(s, start)
		{
		if false isnt match = s[start..].Match('\A[[:alpha:]]([[:alnum:]]|-)*')
			return start + match[0][1]
		return false
		}

	matchAttribute(s, start)
		{
		if false is nameStart = .matchSpaces(s, start)
			return false

		if false is match = s[nameStart..].Match('\A[[:alpha:]_:][-_.:[:alnum:]]*')
			return false

		nameEnd = nameStart + match[0][1]
		valueStart = .matchSpaces(s, nameEnd, optional?:)
		if s[valueStart::1] isnt '='
			return nameEnd

		valueStart = .matchSpaces(s, valueStart + 1, optional?:)
		if false is valueEnd = .matchAttributeValue(s, valueStart)
			return false
		return valueEnd
		}

	matchAttributeValue(s, start)
		{
		if s[start::1] in (`'`, `"`)
			{
			if s.Size() is end = s.Find(s[start], pos: start + 1)
				return false
			return end + 1
			}

		end = s.Find1of(' \t\n"\'=<>`', pos: start)
		if start is end
			return false
		return end
		}

	matchSpaces(s, start, optional? = false)
		{
		end = s.Find1of('^ \t\n', pos: start)

		if optional? is false and start is end
			return false
		return end
		}
	}