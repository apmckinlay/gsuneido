// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	HeadingLevel: false
	New(line)
		{
		.raws = [line]
		}

	Continue(line, start, checkingContinuationText? = false)
		{
		if .BlankLine?(line, start)
			return false, start
		if Md_ContainerBlock.MatchParagraphInteruptableBlockItem(
			line, start, :checkingContinuationText?) isnt false
			return false, start
		return line, start
		}

	IsContinuationText?(line, start)
		{
		result, start = .Continue(line, start, checkingContinuationText?:)
		return false isnt result
		}

	Add(line, start, lazyContinuation? = false, container = false, _mdAddons = #())
		{
		if not lazyContinuation?
			{
			if .IsSetextHeadingUnderline?(line, start)
				{
				s = .raws.Join('\n')
				while false isnt def = .matchLinkReferenceDefinition(s)
					s = s[def.end..]

				if not s.Blank?()
					{
					.HeadingLevel = line.Has?('=') ? 1 : 2
					.Close()
					return
					}
				}
			// addon could modify .raws
			else if mdAddons.Any?({ it.MatchInParagraph(line, start, container, .raws) })
				return
			}

		.raws.Add(line[start..])
		}

	IsSetextHeadingUnderline?(line, start)
		{
		return false isnt .IgnoreLeadingSpaces(line, start) and
			line[start..].Trim() =~ '^(=+|-+)$'
		}

	matchLinkReferenceDefinition(s)
		{
		// start should be false with limit: false
		start = .IgnoreLeadingSpaces(s, 0, limit: false)
		if false is result = Md_Helper.MatchLinkLabel(s, start)
			return false

		label = result.s
		p = result.end
		if s[p::1] isnt ':'
			return false

		if false is destOb = .matchLinkDest(s, p+1)
			return false
		dest = destOb.s

		altEnd = .altEnd(s, destOb)
		end = titleStart = Md_Helper.MatchSpaces(s, destOb.end, optional?:)
		title = ''
		if false isnt titleOb = Md_Helper.MatchLinkTitle(s, titleStart)
			{
			title = titleOb.s
			end = titleOb.end
			}

		if title isnt '' and destOb.end is titleStart // no spaces between dest and title
			return false

		return .buildDef(s, end, altEnd, label, dest, title)
		}

	matchLinkDest(s, start)
		{
		p = Md_Helper.MatchSpaces(s, start, optional?:)
		if false is destOb = Md_Helper.MatchLinkDestination(s, p)
			return false
		if destOb.s is '' and destOb.GetDefault(#inBracket?, false) isnt true
			return false
		return destOb
		}

	altEnd(s, destOb)
		{
		for (i = destOb.end; i < s.Size(); i++)
			{
			if s[i] in (' ', '\t')
				continue
			if s[i] is '\n'
				return i+1
			break
			}
		return false
		}

	buildDef(s, end, altEnd, label, dest, title = '')
		{
		newline = s.Find('\n', end)
		prevNewline = s.FindLast('\n', end)

		if s[end..newline].Blank?()
			return Object(:label, :dest, :title, end: newline+1)
		else if prevNewline isnt false and s[prevNewline..end].Blank?()
			return Object(:label, :dest, :title, end: prevNewline+1)
		else if altEnd isnt false
			return Object(:label, :dest, title: '', end: altEnd)
		else
			return false
		}

	Close(_document)
		{
		s = .raws.Join('\n')
		while false isnt def = .matchLinkReferenceDefinition(s)
			{
			document.AddLinkDefs(@def)
			s = s[def.end..]
			}

		.Inline = s
		super.Close()
		}
	}
