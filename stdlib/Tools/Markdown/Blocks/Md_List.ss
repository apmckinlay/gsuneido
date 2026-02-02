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

	Continue(line, start)
		{
		line2, start = .matchListMarker(line, start)
		if false isnt line2
			return line2, start

		if false isnt listItem = .NextOpenBlock()
			{
			result, unused = listItem.Continue(line, start)
			if result isnt false
				return line, start
			}

		return false, start
		}

	matchListMarker(line, start)
		{
		if false is n = .IgnoreLeadingSpaces(line, start)
			return false, start

		if .Type is 'ul' and .matchMarker?(line, start+n) and
			Md_ThematicBreak.Match(line, start, container: this) is false
			return line, start

		if .Type is 'ol'
			{
			c = .CountLeadingChar(line, start+n, '0-9')
			if c >= 1 and c <= 9/*=max length*/ and .matchMarker?(line, start+n+c)
				return line, start
			}

		return false, start
		}

	matchMarker?(line, start)
		{
		return line[start::1] is .marker and line[start+1::1] in ('', ' ', '\t')
		}

	HasEndingBlankLine?()
		{
		return .Children.Last().HasEndingBlankLine?()
		}

	Loose?: false
	Close()
		{
		super.Close()
		prevItemHasEndingBlankLine? = false
		.ForEachBlockItem()
			{ |listItem|
			if listItem.Loose? is true or prevItemHasEndingBlankLine? is true
				{
				.Loose? = true
				break
				}
			prevItemHasEndingBlankLine? = listItem.HasEndingBlankLine?()
			}
		}
	}