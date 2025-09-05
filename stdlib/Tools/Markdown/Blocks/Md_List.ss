// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_ContainerBlock
	{
	Start: false
	New(marker)
		{
		if marker in ('-', '+', '*')
			{
			.Type = 'ul'
			.marker = marker
			}
		else
			{
			.Type = 'ol'
			.marker = marker[-1..]
			.Start = Number(marker[..-1])
			}
		}

	Continue(line)
		{
		if false isnt line2 = .matchListMarker(line)
			return line2

		if false isnt listItem = .NextOpenBlock()
			{
			if listItem.Continue(line) isnt false
				return line
			}

		return false
		}

	matchListMarker(line)
		{
		if false is n = .IgnoreLeadingSpaces(line)
			return false

		if .Type is 'ul' and line[n::1] is .marker and
			Md_ThematicBreak.Match(line, container: this) is false
			return line

		if .Type is 'ol'
			{
			c = .CountLeadingChar(line[n..], '0-9')
			if c >= 1 and c <= 9/*=max length*/ and line[n+c::1] is .marker
				return line
			}

		return false
		}

	Loose?: false
	Close()
		{
		prevItemEndWithBlankLine? = false
		.ForEachBlockItem()
			{ |listItem|
			if listItem.Loose? is true or prevItemEndWithBlankLine? is true
				{
				.Loose? = true
				break
				}
			if listItem.HasEndBlankLine? is true
				prevItemEndWithBlankLine? = true
			}
		super.Close()
		}
	}