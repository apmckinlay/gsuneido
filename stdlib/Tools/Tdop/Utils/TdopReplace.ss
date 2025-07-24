// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(astWriter, pattern, replacement, from = 0, to = false, count = false)
		{
		target = astWriter.GetRoot()
		source = astWriter.GetSource()
		replacementOb = .parseReplacement(replacement)
		replaceCounter = 0
		pos = from
		forever
			{
			if count isnt false and replaceCounter >= count
				break
			match = TdopSearch(target, pattern, pos)
			if match is false or (to isnt false and match[0][0] + match[0][1] > to)
				break
			astWriter.Replace(match[0].node,
				.generateReplacement(match, replacementOb, source))
			pos = match[0][0] + match[0][1]
			replaceCounter++
			}
		s = astWriter.ToString()
		return to is false
			? s[astWriter.GetCurPos(from)..]
			: s[astWriter.GetCurPos(from)..astWriter.GetCurPos(to)]
		}

	parseReplacement(replacement)
		{
		if replacement.Prefix?('\=')
			return Object(replacement.RemovePrefix('\='))

		replacementOb = Object()

		s = replacement
		forever
			{
			match = s.Match(`\\[[:digit:]]`)
			if match is false
				{
				replacementOb.Add(s)
				break
				}
			replacementOb.Add(s[..match[0][0]], Number(s[match[0][0] + 1]))
			s = s[match[0][0] + match[0][1]..]
			}
		return replacementOb
		}

	generateReplacement(match, replacementOb, source)
		{
		replacement = ''
		for item in replacementOb
			replacement $= String?(item)
				? item
				: source[match[item][0]::match[item][1]]
		return replacement
		}
	}