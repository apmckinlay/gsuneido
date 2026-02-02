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

	Add(line, start)
		{
		if false isnt next = .NextOpenBlock()
			next.Close()

		if false isnt blockItem = .build(line, start)
			.children.Add(blockItem)
		}

	build(line, start)
		{
		if .BlankLine?(line, start)
			return false
		if false isnt item = Md_IndentedCode.Match(line, start)
			return item
		if false isnt item = .MatchParagraphInteruptableBlockItem(line, start,
			container: this)
			return item
		return Md_Paragraph(line[start..])
		}

	MatchParagraphInteruptableBlockItem(line, start, checkingContinuationText? = false,
		container = false, _mdAddons = #())
		{
		if false isnt item = Md_ThematicBreak.Match(line, start, :container,
			:checkingContinuationText?)
			return item
		if false isnt item = Md_ATXheadings.Match(line, start)
			return item
		if false isnt item = Md_FencedCode.Match(line, start)
			return item
		if false isnt item = Md_BlockQuote.Match(line, start)
			return item
		if false isnt item = Md_ListItem.Match(line, start, :checkingContinuationText?,
			:container)
			return item
		if false isnt item = Md_Html.Match(line, start, :container)
			return item
		for addon in mdAddons
			if false isnt item = addon.Match(line, start)
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

	Getter_Children()
		{
		return .children
		}
	}
