// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_ContainerBlock
	{
	New(.indent)
		{
		}

	Match(line, start, checkingContinuationText? = false, container = false)
		{
		if false is n = .IgnoreLeadingSpaces(line, start)
			return false

		start += n
		marker = .advanceMarker(line, start, container, checkingContinuationText?)
		if marker isnt false
			{
			start += marker.Size()
			spaces, line = .advanceSpaces(line, start, container,
				checkingContinuationText?)
			if spaces isnt false
				{
				start += spaces
				listItem = new this(n + marker.Size() + spaces)
				listItem.Add(line, start)
				if container isnt false and container.Base?(Md_List)
					return listItem
				else
					{
					list = new Md_List(marker)
					list.AddItem(listItem)
					return list
					}
				}
			}
		return false
		}

	advanceMarker(line, start, container, checkingContinuationText?)
		{
		marker = false
		if line[start::1] in ('-', '+', '*')
			marker = line[start::1]
		else
			{
			m = .CountLeadingChar(line, start, '0-9')
			accept? = .isMatchingParagraphDirectly(container, checkingContinuationText?)
				? m is 1 and line[start] is '1'
				: m >= 1 and m <= 9/*=max length*/
			if accept? and line[start+m::1] in ('.', ')')
				{
				marker = line[start::m+1]
				}
			}
		return marker
		}

	isMatchingParagraphDirectly(container, checkingContinuationText?)
		{
		return container is false and // from Md_Paragraph.Continue
			checkingContinuationText? is false
		}

	advanceSpaces(line, start, container, checkingContinuationText?)
		{
		spaces = false
		if .BlankLine?(line, start) // rule #3
			{
			// an empty list item cannot interrupt a paragraph
			if not .isMatchingParagraphDirectly(container, checkingContinuationText?)
				spaces = 1
			}
		else
			{
			line = Md_Helper.Detab(line, start)
			c = .CountLeadingChar(line, start, ' ')
			if c >= 1 and c <= 4/*=max length from rule #1*/
				spaces = c
			else if c >= 5/*=length of indented code from rule #2*/
				spaces = 1
			}
		return spaces, line
		}

	Continue(line, start)
		{
		if .BlankLine?(line, start)
			return line, start
		if .indent <= .CountLeadingChar(line, start, ' ')
			return line, start+.indent

		return false, start
		}

	endingBlankLines: 0
	Add(line, start)
		{
		if .BlankLine?(line, start) is false
			{
			if .HasEndingBlankLine?()
				.Loose? = true
			.endingBlankLines = 0
			}
		else
			{
			++.endingBlankLines
			if .Children.Empty?() and .endingBlankLines >= 2
				{
				.Loose? = true
				.Close()
				}
			}
		super.Add(line, start)
		}

	HasEndingBlankLine?()
		{
		return .Children.NotEmpty?() and .endingBlankLines > 0 or
			.lastItemHasEndingBlankLine?()
		}

	lastItemHasEndingBlankLine?()
		{
		children = .Children
		return children.Size() > 0 and children.Last().HasEndingBlankLine?()
		}

	Loose?: false
	}
