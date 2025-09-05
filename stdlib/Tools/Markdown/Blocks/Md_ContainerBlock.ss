// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	New()
		{
		.children = Object()
		}

	NextOpenBlock()
		{
		if .children.Empty?()
			return false
		last = .children.Last()
		if last.Closed?
			return false
		return last
		}

	Add(line)
		{
		if false isnt next = .NextOpenBlock()
			next.Close()

		if false isnt blockItem = .build(line)
			.children.Add(blockItem)
		}

	build(line)
		{
		if .BlankLine?(line)
			return false
		if false isnt item = Md_IndentedCode.Match(line)
			return item
		if false isnt item = .MatchParagraphInteruptableBlockItem(line, container: this)
			return item
		return Md_Paragraph(line)
		}

	MatchParagraphInteruptableBlockItem(line, checkingContinuationText? = false,
		container = false)
		{
		if false isnt item = Md_ThematicBreak.Match(line, :container,
			:checkingContinuationText?)
			return item
		if false isnt item = Md_ATXheadings.Match(line)
			return item
		if false isnt item = Md_FencedCode.Match(line)
			return item
		if false isnt item = Md_BlockQuote.Match(line)
			return item
		if false isnt item = Md_ListItem.Match(line, :checkingContinuationText?,
			:container)
			return item
		if false isnt item = Md_Html.Match(line, :container)
			return item
		return false
		}

	AddItem(item)
		{
		if false isnt next = .NextOpenBlock()
			next.Close()

		.children.Add(item)
		}

	Close()
		{
		if false isnt next = .NextOpenBlock()
			next.Close()
		super.Close()
		}

	ForEachBlockItem(block)
		{
		.children.Each(block)
		}
	}