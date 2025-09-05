// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_ContainerBlock
	{
	New(.indent)
		{
		}

	Match(line, checkingContinuationText? = false, container = false)
		{
		if false is n = .IgnoreLeadingSpaces(line)
			return false

		line = line[n..]
		marker = .advanceMarker(line, container, checkingContinuationText?)
		if marker isnt false
			{
			line = line[marker.Size()..]
			spaces = .advanceSpaces(line, container, checkingContinuationText?)
			if spaces isnt false
				{
				line = line[spaces..]
				listItem = new this(n + marker.Size() + spaces)
				listItem.Add(line)
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

	advanceMarker(line, container, checkingContinuationText?)
		{
		marker = false
		if line[::1] in ('-', '+', '*')
			marker = line[::1]
		else
			{
			m = .CountLeadingChar(line, '0-9')
			accept? = .isMatchingParagraphDirectly(container, checkingContinuationText?)
				? m is 1 and line[0] is '1'
				: m >= 1 and m <= 9/*=max length*/
			if accept? and line[m::1] in ('.', ')')
				{
				marker = line[::m+1]
				}
			}
		return marker
		}

	isMatchingParagraphDirectly(container, checkingContinuationText?)
		{
		return container is false and // from Md_Paragraph.Continue
			checkingContinuationText? is false
		}

	advanceSpaces(line, container, checkingContinuationText?)
		{
		spaces = false
		if .BlankLine?(line) // rule #3
			{
			// an empty list item cannot interrupt a paragraph
			if not .isMatchingParagraphDirectly(container, checkingContinuationText?)
				spaces = 1
			}
		else
			{
			line = Md_Helper.Detab(line)
			c = .CountLeadingChar(line, ' ')
			if c >= 1 and c <= 4/*=max length from rule #1*/
				spaces = c
			else if c >= 5/*=length of indented code from rule #2*/
				spaces = 1
			}
		return spaces
		}

	Continue(line)
		{
		if .BlankLine?(line)
			return line

		if .indent <= .CountLeadingChar(line, ' ')
			return line[.indent..]

		return false
		}

	Loose?: false
	hasBlankLine?: false
	checkStartBlankLine: 0
	HasEndBlankLine?: false
	Add(line)
		{
		if .checkStartBlankLine isnt true
			{
			if .BlankLine?(line)
				{
				if ++.checkStartBlankLine >= 2
					{
					.HasEndBlankLine? = true
					.Close()
					}
				}
			else
				.checkStartBlankLine = true
			}
		else
			{
			if .Loose? is false
				{
				if .HasEndBlankLine? = .BlankLine?(line)
					.hasBlankLine? = true
				else if .hasBlankLine?
					.Loose? = true
				}
			}

		super.Add(line)
		}
	}