// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// convert Wiki text to md
class
	{
	CallClass(src)
		{
		return (new this).Convert(src)
		}

	img_rx: "(?i)\.(gif|jpeg|jpg|png)$"


	litmark: '\xb2'
	mark: '\xb3'
	linkmark: '\xb4'
	evalmark: '\xb5'
	Convert(src)
		{
		// replace literal sections with markers
		literals = Object()
		src = .extractLiterals(src, literals)

		dst = ""
		for line in src.Lines()
			dst $= .handleLine(line)

		if .pre? is true
			dst $= '```'

		// restore literal sections
		dst = .restoreLiterals(dst, literals)

		return dst
		}

	extractLiterals(src, literals)
		{
		src = Md_Addon_Wiki.Md_Addon_Wiki_replace(src, '[sic]', '[/sic]')
			{ |s|
			literals.Add(s) /* = part between tag and end tag */
			.litmark $ (literals.Size() - 1) $ .litmark
			}
		src = Md_Addon_Wiki.Md_Addon_Wiki_replace(src, '[esc]', '[/esc]')
			{ |s|
			literals.Add(s) /* = part between tag and end tag */
			.litmark $ (literals.Size() - 1) $ .litmark
			}
		src = Md_Addon_Wiki.Md_Addon_Wiki_replace(src, '[$', '$]')
			{ |s|
			literals.Add(s)
			.litmark $ (literals.Size() - 1) $ .litmark
			}
		return src
		}

	restoreLiterals(dst, literals)
		{
		dst = dst.Replace(.litmark $ "(\d+)" $ .litmark)
			{ |s|
			literals[Number(s[1 .. -1])]
			}
		return dst
		}

	prevBlank?: true
	handleLine(line)
		{
		dst = ''
		links = Object()
		line = .extractLinks(line, links)

		line = .handleLineStart(line)

		// horizontal rule
		if line =~ "^----*" and not .prevBlank?
			line = '\n' $ line // to avoid being treated as a header

		// strong
		line = line.Replace("'''(.*?)'''")
			{ |s|
			"**" $ s[3..-3].Trim() $ '**' /*= remove quotes*/
			}

		// emphasis
		line = line.Replace("''(.*?)''")
			{ |s|
			"*" $ s[2..-2].Trim() $ '*' /*= remove quotes*/
			}

		// monospaced
		line = line.Replace("==(.*?)==")
			{ |s|
			links.Add('`' $ s[2..-2] $ '`')
			.linkmark $ (links.Size() - 1) $ .linkmark
			}

		// escape html tag
		if not .pre?
			line = line.Replace("<([a-zA-Z]+)", `\\<\1`)

		line = .restoreLinks(line, links)

		dst $= line $ "\n"

		.prevBlank? = line.Blank?()

		return dst
		}

	table?: false
	pre?: false
	handleLineStart(line)
		{
		tableLine? = false
		preLine? = false
		// line starting with whitespace means preformatted
		if line =~ "^\s"
			{
			preLine? = true
			if .pre? is false
				{
				.pre? = true
				line = '```\n' $ line
				}
			}
		else if line =~ "^\*+"
			// line starting with asterisks means unordered list (bullets)
			line = '\t'.Repeat(line.Extract("^\*+").Size()-1) $ line.Replace("^\*+", "* ")
		else if line =~ "^#+"
			// line starting with pound sign means ordered list (numbered)
			line = '\t'.Repeat(line.Extract("^\#+").Size()-1) $ line.Replace("^#+", "1. ")
		else if line =~ '^""'
			{ // line starting with two double quotes ("") is block quote
			line = line.Replace('^""', "> ")
			}
		else if line =~ '^:+'
			{ // line starting with colon means definition list
			line = .buildDefinitionLine(line)
			}
		else if line.Prefix?('|') and line.Suffix?('|')
			{
			tableLine? = true
			if .table? is false
				{
				.table? = true
				line $= '\n|' $ '--|'.Repeat(line.Count('|')-1)
				}
			}
		else if line =~ '^!!+'
			line = '#'.Repeat(line.Extract("^!!+").Size()) $ line.Replace("^!!+", ' ')

		if tableLine? is false and .table? is true
			{
			line = '\n' $ line
			.table? = false
			}
		if preLine? is false and .pre? is true
			{
			line = '```\n' $ line
			.pre? = false
			}
		return line
		}

	buildDefinitionLine(line)
		{
		if line.Prefix?(":") and line.Count(':') > 1
			{
			text = line.LeftTrim(":")
			if text.Has?(":")
				{
				replace = text.BeforeLast(":")
				suffix = text.AfterLast(":").LeftTrim()
				line = replace $ "\n: " $ suffix
				}
			else if text is ""
				return ""
			}
		return line
		}

	protocol_rx: "(https?|ftp|mailto|file|gopher|telnet|news)"
	file_rx: "http://File(\?|/)"
	extractLinks(line, links)
		{
		validUrlChars = "[^][ \t\r\n<>\"'()]*[^][ \t\r\n<>\"'(),.?]"
		url_rx = "\<" $ .protocol_rx $ ":"  $ validUrlChars

		line = line.Replace(.file_rx $ validUrlChars)
			{ |s|
			links.Add(s)
			.linkmark $ (links.Size() - 1) $ .linkmark
			}

		line = line.Replace(url_rx)
			{ |s|
			links.Add(s)
			.linkmark $ (links.Size() - 1) $ .linkmark
			}
		return line
		}
	restoreLinks(line, links)
		{
		// restore external links
		line = line.Replace(.linkmark $ "(\d+)" $ .linkmark)
			{ |s|
			links[Number(s[1 .. -1])]
			}
		return line
		}
	}
