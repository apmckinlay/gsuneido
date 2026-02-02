// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Addon_Base
	{
	New()
		{
		.literals = Object()
		.headings = Object()
		}

	litmark: '\xb2'
	PreParse(s)
		{
		Assert(String?(s))
		s = .extractLiterals(s)
		return s
		}

	AfterHtml(html)
		{
		return .restoreLiterals(
			Opt('<small><a name="0" class="noPrint"><em>Headings:</em></a> ',
				.headings.Join('\n'), '</small>\n') $ html)
		}

	ConvertToHtml(writer, item)
		{
		if item.Base?(Md_ATXheadings)
			{
			n = .headings.Size() + 1
			.headings.Add('<a href="#' $ n $ '" class="noPrint">' $ item.Inline $
				'</a>&nbsp;&nbsp;')
			writer.AddWithBlock('a', attr: Object(name: n))
				{
				writer.Add('h' $ item.Level,
					MarkdownToHtml.ConvertInline(item.ParsedInline))
				}
			return true
			}
		return false
		}

	PreInline(item)
		{
		if item[0] isnt #text
			return false

		item[1] = .extractLinks(item[1])

		link_rx = `\<[A-Z][a-z0-9]+[A-Z][a-zA-Z0-9]*\>`
		// activate page links
		item[1] = item[1].Replace(link_rx)
			{ |s|
			.addLiteral(.InternalLink(s))
			}

		// special substitutions
		item[1] = item[1].Replace("^\[find\]$")
			{ |unused|
			.addLiteral(.searchForm)
			}
		item[1] = item[1].Replace("^\[recent\]$")
			{ |unused|
			.addLiteral(WikiRecent())
			}
		item[1] = item[1].Replace("^\[recent ?[0-9]+\]$")
			{ |unused|
			.addLiteral(WikiRecent(Number(item[1].Extract("^\[recent ?([0-9]+)\]"))))
			}

		return false
		}

	extractLiterals(src)
		{
		src = .replace(src, '[sic]', '[/sic]')
			{ |s|
			.addLiteral(MarkdownToHtml.Encode(s[5 .. -6])) /* = remove tag */
			}
		src = .replace(src, '[esc]', '[/esc]')
			{ |s|
			.addLiteral(MarkdownToHtml.Encode(s[5 .. -6])) /* = remove tag */
			}

		// handle embedded Suneido code
		// if the result contains html tags it's treated as literal
		src = .replace(src, '[$', '$]')
			{ |s|
			try
				{
				code = s[2 .. -2].Trim()
				.assertCodeOkayToEval(code)
				result = String(code.Eval()) // Eval should be safe due to whitelist
				}
			catch (x)
				result = "[$ " $ x $ " $]"
			if result =~ "<[a-zA-Z]+>"
				result = .addLiteral(result)
			result
			}
		return src
		}

	replace(s, start, end, block)
		{
		pos = 0
		replacements = Object()
		while s.Size() isnt startPos = s.Find(start, pos)
			{
			if s.Size() is endPos = s.Find(end, startPos + start.Size())
				break
			endPos += end.Size()
			replacements.Add(
				Object(startPos, endPos - startPos, block(s[startPos..endPos])))
			pos = endPos
			}
		for (i = replacements.Size() - 1; i >= 0; i--)
			s = s.ReplaceSubstr(@replacements[i])
		return s
		}

	assertCodeOkayToEval(code)
		{
		// code must be a single function call that is whitelisted
		if code.Blank?()
			return
		whiteList = GetContributions("WikiWhitelistFunctions")
		if code.AfterFirst("(").Has?("(")
			throw code $ " - can not have multiple calls"
		fn = code.BeforeFirst("(")
		if not whiteList.Has?(fn)
			throw code $ " - function is not in the Wiki function whitelist"
		}

	restoreLiterals(dst)
		{
		dst = dst.Replace(.litmark $ "(\d+)" $ .litmark)
			{ |s|
			.restoreLiterals(.literals[Number(s[1 .. -1])])
			}
		return dst
		}

	protocol_rx: "(https?|ftp|mailto|file|gopher|telnet|news)"
	file_rx: "http://File(\?|/)"
	img_rx: "(?i)\.(gif|jpeg|jpg|png)$"
	extractLinks(s)
		{
		validUrlChars = "[^][ \t\r\n<>\"'()]*[^][ \t\r\n<>\"'(),.?]"
		url_rx = "\<" $ .protocol_rx $ ":"  $ validUrlChars

		s = s.Replace(.file_rx $ validUrlChars)
			{ |s|
			.addLiteral(.InternalFile(s))
			}

		s = s.Replace(url_rx, .extractUrl)
		return s
		}
	extractUrl(url)
		{
		if url =~ `http://appserver(\.axon)?:8080/Wiki`
			return url // don't allow external links to wiki
		return .addLiteral(.externalLink(url))
		}

	addLiteral(literal)
		{
		.literals.Add(literal)
		return .litmark $ (.literals.Size() - 1) $ .litmark
		}

	externalLink(url)
		{
		return url =~ .img_rx
			? '<img src="' $ url $ '">'
			: '<a href="' $ url $ '">' $ url $ '</a>'
		}

	InternalLink(name)
		{
		return .nameNotExist?(name)
			? '<a href="Wiki?edit=' $ name $ '"><b>?</b></a>' $ name
			: '<a href="Wiki?' $ name $ '">' $ name $ '</a>'
		}

	nameNotExist?(name)
		{
		return QueryEmpty?('wiki', :name)
		}

	InternalFile(name)
		{
		name = name.Replace(.file_rx, '')
		return name =~ .img_rx
			? '<img src="File/' $ name $ '">'
			: '<a href="File/' $ name $ '">' $ Paths.Basename(name) $ '</a>'
		}

	searchForm:	'<form method="get" action=Wiki>
				<input type="text" size=40 name="find" spellcheck="true">
				<input type="submit" value="Search">
				</form>'
	}