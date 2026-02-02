// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Addon_Base
	{
	New(.attr = #())
		{
		}

	MatchInParagraph(line, start, container, raws)
		{
		if false is cols = .isTableDelimiter(line, start)
			return false

		return .addTable(cols, container, raws) is true
		}

	isTableDelimiter(line, start)
		{
		line = line[start..].Trim()
		cols = Object()
		hasBar? = false

		if line[::1] is '|'
			{
			line = line[1..]
			hasBar? = true
			}

		while false isnt match = line.Match('^\s*(:?)-+(:?)\s*')
			{
			align = match[2][1] is 1
				? match[1][1] is 1
					? #center
					: #right
				: #left
			cols.Add(align)
			next = match[0][1]
			if line[next::1] is '|'
				{
				next++
				hasBar? = true
				}
			line = line[next..]
			}
		if not hasBar? or cols.Empty?() or line isnt ''
			return false

		return cols
		}

	addTable(cols, container, raws)
		{
		if raws.NotEmpty?() and false isnt titles = .matchTableTitle(raws.Last(), cols)
			{
			raws.PopLast()
			Assert(container isnt: false)
			container.AddItem(new Md_Table(cols, titles))
			return true
			}
		return false
		}

	matchTableTitle(titleLine, cols)
		{
		titles = Md_Table.ParseTableLine(titleLine)
		if titles.Size() isnt cols.Size()
			return false
		return titles
		}

	ParseInline(item)
		{
		if not item.Base?(Md_Table)
			return
		item.ParseInline()
		}

	ConvertToHtml(writer, item)
		{
		if not item.Base?(Md_Table)
			return false

		writer.AddWithBlock('table', attr: .attr)
			{
			if item.ParsedTitles.Any?(#NotEmpty?)
				writer.AddWithBlock('tr')
					{
					for (i = 0; i < item.ParsedTitles.Size(); i++)
						writer.Add('th',
							MarkdownToHtml.ConvertInline(item.ParsedTitles[i]),
							extra: 'style="text-align: ' $ item.Aligns[i] $ ';"')
					}
			for row in item.ParsedRows
				writer.AddWithBlock('tr')
					{
					for (i = 0; i < row.Size() and i < item.Aligns.Size(); i++)
						writer.Add('td', MarkdownToHtml.ConvertInline(row[i]),
							extra: 'style="text-align: ' $ item.Aligns[i] $ ';"')
					}
			}
		return true
		}
	}