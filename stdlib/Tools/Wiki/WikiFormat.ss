// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// convert Wiki text to html
// based on code in The WikiWiki Way by Bo Leuf and Ward Cunningham
class
	{
	CallClass(src)
		{
		try
			return (new this).Format(src)
		catch (e)
			return '<pre style="color:red;"><b>' $ e $ '<br>' $
				FormatCallStack(e.Callstack(), levels: 10, indent:) $
				'</b></pre><pre>' $ XmlEntityEncode(src) $ '<pre>'
		}

	img_rx: "(?i)\.(gif|jpeg|jpg|png)$"


	litmark: '\xb2'
	mark: '\xb3'
	fileMark: '\xb4'
	Format(src)
		{
		src = src.
			Replace('&', '\&amp;').
			Replace('<', '\&lt;').
			Replace('>', '\&gt;')

		// replace literal sections with markers
		literals = Object()
		src = .setHeaderMarkers(src, literals)

		dst = ""

		// replace headings with anchors, add links at top
		nh = 1
		headings = ""
		heading_rx = '(!!+)(.*)'
		src = src.Replace(heading_rx)
			{ |s|
			exclaims = s.Extract(heading_rx, 1)
			heading = s.Extract(heading_rx, 2)
			headings $= '<a href="#' $ nh $ '" class="noPrint">' $ heading $
				'</a>&nbsp;&nbsp;'
			exclaims $ '<a name="' $ nh++ $ '">' $ heading $ '</a>'
			}
		if headings isnt ""
			dst $= '<small><a name="0" class="noPrint"><em>Headings:</em></a> ' $
				headings $ '</small><p>\n'

		.stack = new Stack
		for line in src.Lines()
			dst $= .handleLine(line)

		dst $= .emit('', 0)

		// restore literal sections
		dst = .getHeaderMarkers(dst, literals)

		return dst
		}

	setHeaderMarkers(src, literals)
		{
		src = .replace(src, '[sic]', '[/sic]')
			{ |s|
			literals.Add(s[5 .. -6]) /* = part between tag and end tag */
			.litmark $ (literals.Size() - 1) $ .litmark
			}
		src = .replace(src, '[esc]', '[/esc]')
			{ |s|
			literals.Add(s[5 .. -6]) /* = part between tag and end tag */
			.litmark $ (literals.Size() - 1) $ .litmark
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
				{
				literals.Add(result)
				result = .litmark $ (literals.Size() - 1) $ .litmark
				}
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

	getHeaderMarkers(dst, literals)
		{
		dst = dst.Replace(.litmark $ "(\d+)" $ .litmark)
			{ |s|
			literals[Number(s[1 .. -1])]
			}
		return dst
		}

	handleLine(line)
		{
		dst = ''
		// replace external links with markers
		markers = Object()
		line = .setMarkers(line, markers)

		line = .handleLineStart(line)

		// horizontal rule
		line = line.Replace("^----*", "<hr>")

		// strong
		line = line.Replace("'''(.*?)'''", "<strong>\1</strong>")

		// emphasis
		line = line.Replace("''(.*?)''", "<em>\1</em>")

		// monospaced
		line = line.Replace("==(.*?)==", "<tt>\1</tt>")

		link_rx = `\<[A-Z][a-z0-9]+[A-Z][a-zA-Z0-9]*\>`
		// activate page links
		line = line.Replace(link_rx, .InternalLink)

		// special substitutions
		line = line.Replace("\[find\]", .SearchForm)
		if line =~ "^\[recent\]"
			line = WikiRecent()
		else if line =~ "^\[recent ?[0-9]+\]"
			line = WikiRecent(Number(line.Extract("^\[recent ?([0-9]+)\]")))

		line = .restoreMarkers(line, markers)

		dst $= line $ "\n"

		return dst
		}

	handleLineStart(line)
		{
		dst = ''
		if line =~ "^\s*$" and .stack.Count() > 0 and .stack.Top() is "table"
			dst $= .emit('', 0) $ '<p>\n' // end table
		else if (line =~ "^\s*$" and
			(.stack.Count() is 0 or .stack.Top() isnt "pre"))
			// blank line means new paragraph
			line = "<p>"
		else if line =~ "^\s"
			// line starting with whitespace means preformatted
			dst $= .emit('pre', 1)
		else if line =~ "^\*+"
			// line starting with asterisks means unordered list (bullets)
			{
			dst $= .emit('ul', line.Extract("^\*+").Size())
			line = line.Replace("^\*+", "<li>")
			}
		else if line =~ "^#+"
			// line starting with pound sign means ordered list (numbered)
			{
			dst $= .emit('ol', line.Extract("^#+").Size())
			line = line.Replace("^#+", "<li>")
			}
		else if line =~ '^""'
			{ // line starting with two double quotes ("") is block quote
			dst $= .emit('blockquote', 1)
			line = line.Replace('^""', "")
			}
		else if line =~ '^:+'
			{ // line starting with colon means definition list
			dst $= .emit('dl', line.Extract("^:+").Size())
			line = .buildDefinitionLine(line)
			}
		else if line.Prefix?('|') and line.Suffix?('|')
			{ // |table|
			if .stack.Count() is 0 or .stack.Top() isnt 'table'
				dst $= .emit('table', 1, 'border="1" cellpadding="3"')
			line = '<td>' $ line[1 .. -1] $ '</td>'
			line = line.Replace('\|', '</td><td>')
			line = Xml('tr', line)
			}
		else if line =~ '^!!'
			line = .headings(line)
		else
			dst $= .emit('', 0)
		return dst $ line
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
				line = "<dt>" $ replace $ "<dd>" $ suffix
				}
			else if text is ""
				return "<dt><dd>"
			}
		return line
		}

	protocol_rx: "(https?|ftp|mailto|file|gopher|telnet|news)"
	file_rx: "http://File(\?|/)"
	setMarkers(line, markers)
		{
		files = Object()
		urls = Object()
		validUrlChars = "[^][ \t\r\n<>\"'()]*[^][ \t\r\n<>\"'(),.?]"
		url_rx = "\<" $ .protocol_rx $ ":"  $ validUrlChars

		line = line.Replace(.file_rx $ validUrlChars)
			{ |s|
			files.Add(.InternalFile(s))
			.fileMark $ (files.Size() - 1) $ .fileMark
			}

		line = line.Replace(url_rx, { .addMarker(it, urls) })

		markers.Add(files at: 'files')
		markers.Add(urls at: 'urls')

		return line
		}
	addMarker(s, urls)
		{
		if s =~ `http://appserver(\.axon)?:8080/Wiki`
			return s // don't allow external links to wiki
		urls.Add(s)
		return .mark $ (urls.Size() - 1) $ .mark
		}
	restoreMarkers(line, markers)
		{
		// restore external links
		line = line.Replace(.mark $ "(\d+)" $ .mark)
			{ |s|
			.external_link(markers.urls[Number(s[1 .. -1])])
			}

		// restore internal files
		line = line.Replace(.fileMark $ "(\d+)" $ .fileMark)
			{ |s|
			markers.files[Number(s[1 .. -1])]
			}
		return line
		}

	headings(line)
		{
		dst = ''
		if line =~ "^!!!!"
			{
			dst $= .emit('h4', 1)
			line = line.Replace("^!!!!", "")
			}
		else if line =~ "^!!!"
			{
			dst $= .emit('h3', 1)
			line = line.Replace("^!!!", "")
			}
		else if line =~ "^!!"
			{
			dst $= .emit('h2', 1)
			line = line.Replace("^!!", "")
			}
		return dst $ line
		}

	emit(code, depth, attrs = "")
		{
		if attrs isnt ""
			attrs = ' ' $ attrs
		tags = ''
		while .stack.Count() > depth						// end tags
			tags $= '</' $ .stack.Pop() $ '>\n'
		while .stack.Count() < depth						// start tags
			{
			tags $= '\n<' $ code $ attrs $ '>'
			.stack.Push(code)
			}
		if .stack.Count() > 0 and .stack.Top() isnt code	// change tag
			{
			tags $= '</' $ .stack.Top() $ '>\n<' $ code $ attrs $ '>\n'
			.stack.Pop()
			.stack.Push(code)
			}
		return tags
		}

	InternalLink(name)
		{
		return false is Query1('wiki', :name)
			? '<a href="Wiki?edit=' $ name $ '"><b>?</b></a>' $ name
			: '<a href="Wiki?' $ name $ '">' $ name $ '</a>'
		}
	InternalFile(name)
		{
		name = name.Replace(.file_rx, '')
		return name =~ .img_rx
			? '<img src="File/' $ name $ '">'
			: '<a href="File/' $ name $ '">' $ Paths.Basename(name) $ '</a>'
		}
	external_link(url)
		{
		return url =~ .img_rx
			? '<img src="' $ url $ '">'
			: '<a href="' $ url $ '">' $ url $ '</a>'
		}

	SearchForm:	'<form method="get" action=Wiki>
				<input type="text" size=40 name="find" spellcheck="true">
				<input type="submit" value="Search">
				</form>'
	}
