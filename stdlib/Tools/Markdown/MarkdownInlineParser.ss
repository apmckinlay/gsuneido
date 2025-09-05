// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(inline)
		{
		parser = new this(inline)
		return parser.Result()
		}

	i: 0
	base: 0
	New(.inline)
		{
		.results = Object()
		while false isnt c = .cur()
			{
			if c is '`'
				.matchCode()
			else if c is '<'
				.matchAutoLinkOrHTML()
			else
				.advance()
			}
		.addText(.i, last?:)
		}

	Result()
		{
		return .results
		}

	matchCode()
		{
		start = .i
		length = 0
		while .cur() is '`'
			{
			length++
			.advance()
			}

		if false is end = .findCloseCode(start, length)
			return

		.addText(start)
		.results.Add(Object(#code, .normalizeCodeSpan(.inline[start+length..end])))
		.set(end+length)
		}

	findCloseCode(start, length)
		{
		pos = start + length
		while false isnt match = .inline.Match('`+', pos)
			{
			if match[0][1] is length
				return match[0][0]
			pos = match[0][0] + match[0][1]
			}
		return false
		}

	normalizeCodeSpan(code)
		{
		code = code.Replace('\n', ' ')
		if code[::1] is ' ' and code[-1::1] is ' ' and
			code.Find1of('^ \t') isnt code.Size()
			code = code[1..-1]
		return code
		}

	matchAutoLinkOrHTML()
		{
		if .matchHTMLTag?() or
			.matchHTMLComment?() or
			.matchProcessingInstruction?() or
			.matchDeclaration?() or
			.matchCDATA?()
			return

		if .matchURIAutoLink?() or
			.matchEmailAutoLink?()
			return

		.advance()
		}

	matchHTMLComment?()
		{
		return .matchHTMLByRegex('\A<\!-->|<\!--->|<\!--(.|\s)*?-->')
		}

	matchProcessingInstruction?()
		{
		return .matchHTMLByRegex('\A<\?(.|\s)*?\?>')
		}

	matchDeclaration?()
		{
		return .matchHTMLByRegex('\A<![[:alpha:]](.|\s)*?>')
		}

	matchCDATA?()
		{
		return .matchHTMLByRegex('\A<!\[CDATA\[(.|\s)*?\]\]>')
		}

	matchHTMLByRegex(regex)
		{
		if false isnt match = .inline[.i..].Match(regex)
			{
			.addText(.i)
			.results.Add(Object(#html, .inline[.i::match[0][1]]))
			.set(.i + match[0][1])
			return true
			}
		return false
		}

	matchHTMLTag?()
		{
		if false is length = Md_Helper.MatchHTMLTag(.inline[.i..])
			return false

		.addText(.i)
		.results.Add(Object(#html, .inline[.i::length]))
		.set(.i + length)
		return true
		}

	matchURIAutoLink?()
		{
		if false isnt match = .inline[.i..].Match(
			'\A<([[:alpha:]][-.+[:alnum:]]+:[^[:cntrl:] <>]*)>')
			{
			uri = .inline[.i + match[1][0]::match[1][1]]
			.addText(.i)
			.results.Add(Object(#link, uri, href: uri))
			.set(.i + match[0][1])
			return true
			}
		return false
		}

	matchEmailAutoLink?()
		{
		if false isnt match = .inline[.i..].Match(
			'\A<([-.!#$%&\'*+/=?^_`{|}~[:alnum:]]+' $
				'@[[:alnum:]]([-[:alnum:]]*[[:alnum:]])?' $
				'(\.[[:alnum:]]([-[:alnum:]]*[[:alnum:]])?)*)>')
			{
			email = .inline[.i + match[1][0]::match[1][1]]
			.addText(.i)
			.results.Add(Object(#link, email, href: 'mailto:' $ email))
			.set(.i + match[0][1])
			}
		return false
		}

	addText(to, last? = false)
		{
		if to > .base
			for item in .processText(.inline[.base..to], :last?)
				.results.Add(String?(item) ? Object(#text, item) : item)
		}

	first?: true
	processText(text, last?)
		{
		results = Object()
		lines = text.Lines()
		if text.Suffix?('\n')
			lines.Add('')

		for row in lines.Members()
			{
			if .first? is true or row isnt 0
				lines[row] = lines[row].LeftTrim()
			if last? is true or row isnt lines.Size() - 1
				lines[row] = lines[row].RightTrim()
			}

		if .first? is true
			.first? = false

		results.Add(lines.Join('\n'))
		return results
		}

	set(.i)
		{
		.base = .i
		}

	cur()
		{
		if .i >= .inline.Size()
			return false
		return .inline[.i]
		}

	advance()
		{
		.i++
		}
	}