// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	New(.Aligns, .titles)
		{
		.rows = Object()
		.ParsedTitles = Object()
		.ParsedRows = Object()
		}

	Continue(line, start)
		{
		if .BlankLine?(line, start)
			return false, start
		if Md_ContainerBlock.MatchParagraphInteruptableBlockItem(line, start) isnt false
			return false, start
		return line, start
		}

	Add(line, start)
		{
		.rows.Add(.ParseTableLine(line[start..]))
		}

	ParseTableLine(line)
		{
		line = line.Trim()
		values = Object()
		pos = 0
		if line[pos::1] is '|' // skip the start '|'
			pos++

		value = ''
		escape? = false
		while pos < line.Size()
			{
			if line[pos] is '|' and not escape?
				{
				values.Add(value.Trim())
				value = ''
				}
			else
				{
				value $= line[pos]
				escape? = line[pos] is '\\'
				}
			pos++
			}

		if not value.Blank?()
			values.Add(value.Trim())

		return values
		}

	ParseInline()
		{
		for title in .titles
			.ParsedTitles.Add(MarkdownInlineParser2(title))
		for row in .rows
			{
			.ParsedRows.Add(parsed = Object())
			for value in row
				parsed.Add(MarkdownInlineParser2(value))
			}
		}
	}