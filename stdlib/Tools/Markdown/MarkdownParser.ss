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

		for line in lines
			{
			line = Md_Helper.Detab(line)
			opens = document.GetOpenBlockItems()
			for (i = 0; i < opens.Size(); i++)
				{
				if false is result = opens[i].Continue(line)
					break
				line = result
				}
			lastOpen = opens[i - 1]
			lazyContinuation? = false
			if ((opens.Last()).Base?(Md_Paragraph) and
				i < opens.Size() - 1 and // there is unsatisfied container block
				opens.Last().IsContinuationText?(line))
				{
				lastOpen = opens.Last()
				lazyContinuation? = true
				}
			lastOpen.Add(line, :lazyContinuation?)
			}

		document.Finish()
		return document
		}

	phase2(document)
		{
		.forEachInline(document)
		}

	forEachInline(container)
		{
		container.ForEachBlockItem()
			{
			if it.Base?(Md_ContainerBlock)
				.forEachInline(it)
			else if it.Base?(Md_Paragraph)
				it.ParsedInline = MarkdownInlineParser(it.Inline)
			else if it.Base?(Md_ATXheadings)
				it.ParsedInline = MarkdownInlineParser(it.Inline)
			}
		}
	}