// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// Only detab the starting substring that only contains ' ' or \t
	Detab(line, from = 0)
		{
		pos = line.Find1of('^ \t', from)
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

		pos = .MatchSpaces(s, pos, optional?:)
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
		if false is nameStart = .MatchSpaces(s, start)
			return false

		if false is match = s[nameStart..].Match('\A[[:alpha:]_:][-_.:[:alnum:]]*')
			return false

		nameEnd = nameStart + match[0][1]
		valueStart = .MatchSpaces(s, nameEnd, optional?:)
		if s[valueStart::1] isnt '='
			return nameEnd

		valueStart = .MatchSpaces(s, valueStart + 1, optional?:)
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

	MatchSpaces(s, start, optional? = false)
		{
		end = s.Find1of('^ \t\n', pos: start)

		if optional? is false and start is end
			return false
		return end
		}

	MatchLinkDestination(s, start)
		{
		if s[start::1] is '<'
			{
			for (p = start + 1; p < s.Size(); p++)
				{
				if s[p] is '\\'
					p++
				else if s[p] is '>'
					return Object(end: p + 1, s: .Escape(s[start+1..p]), inBracket?:)
				else if s[p] is '\n'
					return false
				}
			return false
			}
		else
			{
			parentheses = 0
			for (p = start; p < s.Size(); p++)
				{
				if s[p] =~ '[[:cntrl:] ]'
					return Object(end: p, s: .Escape(s[start..p]))
				else if s[p] is '\\'
					p++
				else if s[p] is '('
					parentheses++
				else if s[p] is ')'
					{
					if parentheses > 0
						parentheses--
					else
						return Object(end: p, s: .Escape(s[start..p]))
					}
				}
			return Object(end: s.Size(), s: .Escape(s[start..]))
			}
		}

	MatchLinkTitle(s, start)
		{
		close = false
		if s[start::1] in (`"`, `'`)
			close = s[start]
		else if s[start::1] is '('
			close = '()'
		if close is false
			return Object(end: start, s: '')

		for (p = start + 1; p < s.Size(); p++)
			{
			if s[p] is '\\'
				p++
			else if close.Has?(s[p])
				return Object(end: p + 1, s: .Escape(s[start+1..p]))
			}
		return false
		}

	MatchLinkLabel(s, start, allowBlank? = false)
		{
		if s[start::1] isnt '['
			return false

		c = 0
		end = false
		for (i = start+1; i < s.Size(); i++, c++)
			{
			if c > 999 /*=max*/
				return false

			if s[i] is '\\'
				i++
			else if s[i] is '['
				return false
			else if s[i] is ']'
				{
				end = i
				break
				}
			}

		if end is false
			return false

		label = s[start+1..end]
		if label.Blank?() and not allowBlank?
			return false

		return Object(end: end + 1, s: .NormalizeLinkLabel(label))
		}

	NormalizeLinkLabel(s)
		{
		//MISSING: Unicode case fold
		return s.Trim().Tr(' \t\n', ' ').Lower()
		}

	Escape(s)
		{
		return s.Replace(`\\([[:punct:]])`, `\1`)
		}
	}
