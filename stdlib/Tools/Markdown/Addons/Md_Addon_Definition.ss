// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Addon_Base
	{
	MatchInParagraph(line, start, container, raws)
		{
		if false is def = Md_Definition.MatchDefLine(line, start)
			return false

		if raws.Empty?() or (term = raws.Last()).Blank?()
			return false

		raws.PopLast()
		container.AddItem(new Md_Definition(term, def))
		return true
		}

	ParseInline(item)
		{
		if not item.Base?(Md_Definition)
			return
		item.ParseInline()
		}

	PreprocessContainer(container)
		{
		start = end = false
		container.ForEachBlockItem()
			{ |item|
			if item.Base?(Md_Definition)
				{
				if start is false
					start = item
				end = item
				}
			// break the current dl
			else if start isnt false and
				(not item.Base?(Md_Paragraph) or not item.Inline.Blank?())
				{
				start.StartDL = end.EndDL = true
				start = end = false
				}
			}
		if start isnt false
			start.StartDL = end.EndDL = true
		}

	ConvertToHtml(writer, item)
		{
		if not item.Base?(Md_Definition)
			return false

		if item.StartDL is true
			writer.AddOpen('dl')
		writer.Add('dt', MarkdownToHtml.ConvertInline(item.ParsedTerm))
		for def in item.ParsedDefs
			{
			dd = MarkdownToHtml.ConvertInline(def)
			if not dd.Blank?()
				writer.Add('dd', dd)
			}
		if item.EndDL is true
			writer.AddClose('dl')
		return true
		}
	}