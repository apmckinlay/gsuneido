// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(s)
		{
		result = Object()
		scan = .scan(s)
		i = 0
		spaceItem = false
		while i < scan.Size()
			{
			if .match?(scan, i, type: #WHITESPACE) or
				.match?(scan, i, type: #NEWLINE)
				{
				if spaceItem is false
					{
					spaceItem = scan[i]
					spaceItem.type = #WHITESPACE
					}
				else
					{
					spaceItem.end = scan[i].end
					spaceItem.value $= scan[i].value
					spaceItem.text $= scan[i].text
					}
				i++
				continue
				}
			else if spaceItem isnt false
				{
				result.Add(spaceItem)
				spaceItem = false
				}

			if false isnt advance = .matchHole(scan, i, result)
				i += advance
			else
				result.Add(scan[i++])
			}

		if spaceItem isnt false
			result.Add(spaceItem)

		return result
		}

	scan(s)
		{
		scan = Scanner(s)
		list = Object()
		start = 0
		while scan isnt type = scan.Next2()
			{
			item = Object(:start, end: scan.Position(), :type, value: scan.Value(),
				text: scan.Text())
			if type is #IDENTIFIER
				item.keyword? = scan.Keyword?()
			list.Add(item)
			start = scan.Position()
			}
		return list
		}

	matchHole(scan, i, result)
		{
		if .match?(scan, i, ':') and .match?(scan, i + 1, '[') and
			.match?(scan, i + 2, type: #IDENTIFIER) and
			.match?(scan, i + 3 /*=pos*/, ']')
			{
			result.Add(Object(start: scan[i].start, end: scan[i + 3 /*=pos*/].end,
				type: #HOLE,
				value: scan[i + 2].value, text: ':[' $ scan[i + 2].text $ ']'))
			return 4 /*=offset*/
			}

		return false
		}

	match?(scan, i, text = false, type = '')
		{
		return i < scan.Size() and scan[i].type is type and
			(text is false or scan[i].text is text)
		}
	}
