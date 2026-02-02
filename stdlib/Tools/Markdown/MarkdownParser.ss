// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// The parse algorithm is based on CommonMark Spec@0.31.2
// - https://spec.commonmark.org/0.31.2/
class
	{
	CallClass(md)
		{
		document = .phase1(md)
		.phase2(document)
		return document
		}

	// parse into blocks
	phase1(md)
		{
		lines = md.Lines()
		document = Md_Document()
		_document = document

		for line in lines
			.processLine(document, line)

		document.Finish()
		return document
		}

	processLine(document, line)
		{
		line = Md_Helper.Detab(line)
		start = 0
		opens = document.GetOpenBlockItems()
		n = opens.Size()
		for (i = 0; i < n; i++)
			{
			result, start = opens[i].Continue(line, start)
			if false is result
				break
			line = result
			}
		lastOpen = opens[i - 1]
		container = opens.GetDefault(i - 2, false)
		lazyContinuation? = false
		if (opens[n - 1].Base?(Md_Paragraph) and
			i < n - 1 and // there is unsatisfied container block
			opens[n - 1].IsContinuationText?(line, start))
			{
			lastOpen = opens[n - 1]
			container = opens.GetDefault(n - 2, false)
			lazyContinuation? = true
			}
		// container is used by Md_Paragraph.Add for tables
		lastOpen.Add(line, start, :lazyContinuation?, :container)
		}

	phase2(document)
		{
		_document = document
		.forEachInline(document)
		}

	forEachInline(container, _mdAddons = #())
		{
		container.ForEachBlockItem()
			{
			if it.Base?(Md_ContainerBlock)
				.forEachInline(it)
			else if it.Base?(Md_Paragraph)
				it.ParsedInline = MarkdownInlineParser2(it.Inline)
			else if it.Base?(Md_ATXheadings)
				it.ParsedInline = MarkdownInlineParser2(it.Inline)
			else
				{
				for addon in mdAddons
					addon.ParseInline(it)
				}
			}
		}
	}
