// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// This is the second stage for parsing nested emphasis and links
// Based on
// https://spec.commonmark.org/0.31.2/#an-algorithm-for-parsing-nested-emphasis-and-links
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
		.nodes = DoubleLinkedList()
		.delimiters = DoubleLinkedList()

		while false isnt c = .cur()
			{
			if c is '\\'
				{
				if .next() =~ '[[:punct:]]'
					{
					.flushText()
					.addTextNode(.next())
					.set(.i+2)
					}
				else if .next() is '\n'
					{
					.lineBreak(hard?:)
					.set(.i+1)
					}
				else
					.advance()
				}
			else if c is '&'
				.matchEntityOrNumericChar()
			else if c is '\n'
				.lineBreak(hard?: .inline[.i-1::1] is ' ' and .inline[.i-2::1] is ' ')
			else if c is '`'
				.matchCode()
			else if c is '<'
				.matchAutoLinkOrHTML()
			else if c is '['
				{
				.flushText()
				node = .addTextNode('[')
				.delimiters.Append(Object('[', active:, :node, pos: .i))
				.set(.i+1)
				}
			else if c is '!' and .next() is '['
				{
				.flushText()
				node = .addTextNode('![')
				.delimiters.Append(Object('![', active:, :node, pos: .i))
				.set(.i+2)
				}
			else if c is ']'
				{
				.flushText()
				.matchLinkImage()
				}
			else if c in ('*', '_')
				.emphasis()
			else
				.advance()
			}
		.flushText(beforeLineBreak?:)
		.result = .processEmphasis(.nodes)
		}

	Result()
		{
		return .result
		}

	matchEntityOrNumericChar()
		{
		if .matchEntity() or .matchNumericChar()
			return

		.advance()
		}

	maxEntityLen: 33
	matchEntity()
		{
		piece = .inline[.i::.maxEntityLen]
		if piece.Size() is j = piece.Find(';')
			return false
		if not Html5Entities.Member?(name = piece[..j+1])
			return false
		.flushText()
		.addTextNode(name, codepoints: Html5Entities[name])
		.set(.i+j+1)
		return true
		}

	maxNumericCharLen: 10
	matchNumericChar()
		{
		piece = .inline[.i::.maxNumericCharLen]
		if piece.Size() is j = piece.Find(';')
			return false
		if piece[..j+1] !~ '(?i)\A&#([0-9]+)|(x[0-9a-f]+);\Z'
			return false
		.flushText()
		s = piece[..j+1].Lower()
		codepoint = Number((s[2] is 'x' ? '0' : '') $ s[2..-1])
		.addTextNode(piece[..j+1], codepoints: [codepoint])
		.set(.i+j+1)
		return true
		}

	lineBreak(hard? = false)
		{
		.flushText(beforeLineBreak?:)
		.addNode(Object(#linkbreak, :hard?))
		.set(.i+1)
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

		.flushText(start)
		.addNode(Object(#code, .normalizeCodeSpan(.inline[start+length..end])))
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
		if .matchHTML?()
			return

		if .matchURIAutoLink?() or
			.matchEmailAutoLink?()
			return

		.advance()
		}

	matchHTML?(_options= #())
		{
		if options.GetDefault(#turnOffHtml, false) is true
			return false

		return .matchHTMLTag?() or
			.matchHTMLComment?() or
			.matchProcessingInstruction?() or
			.matchDeclaration?() or
			.matchCDATA?()
		}

	matchHTMLTag?()
		{
		if false is length = Md_Helper.MatchHTMLTag(.inline[.i..])
			return false

		.flushText()
		.addNode(Object(#html, .inline[.i::length]))
		.set(.i + length)
		return true
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
			.flushText()
			.addNode(Object(#html, .inline[.i::match[0][1]]))
			.set(.i + match[0][1])
			return true
			}
		return false
		}

	matchURIAutoLink?()
		{
		if false isnt match = .inline[.i..].Match(
			'\A<([[:alpha:]][-.+[:alnum:]]+:[^[:cntrl:] <>]*)>')
			{
			uri = .inline[.i + match[1][0]::match[1][1]]
			.flushText()
			.addNode(Object(#link, Object(Object(#text uri)), href: uri))
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
			.flushText()
			.addNode(Object(#link, Object(Object(#text, email)), href: 'mailto:' $ email))
			.set(.i + match[0][1])
			}
		return false
		}

	matchLinkImage()
		{
		open = false
		for (item = .delimiters.Prev(); item isnt false;
			item = .delimiters.Prev(item))
			{
			if item.value[0] in ('[', '![')
				{
				open = item
				break
				}
			}

		if open is false
			{
			.advance()
			return
			}
		if open.value.active is false
			{
			.delimiters.Del(open)
			.advance()
			return
			}
		if .matchInlineLink?(open) or
			.matchRefLink?(open)
			return

		.delimiters.Del(open)
		.advance()
		}

	matchInlineLink?(open)
		{
		// .i is ']'
		if .next() isnt '('
			return false

		pos = Md_Helper.MatchSpaces(.inline, .i + 2, optional?:)
		if false is destOb = Md_Helper.MatchLinkDestination(.inline, pos)
			return false
		linkDest = destOb.s

		pos = Md_Helper.MatchSpaces(.inline, destOb.end, optional?:)
		if linkDest isnt '' and pos is destOb.end and .inline[pos::1] isnt ')'
			return false

		if false is titleOb = Md_Helper.MatchLinkTitle(.inline, pos)
			return false
		linkTitle = titleOb.s

		pos = Md_Helper.MatchSpaces(.inline, titleOb.end, optional?:)
		if .inline[pos::1] isnt ')'
			return false

		.addLinkImage(open, linkDest, linkTitle)
		.set(pos+1)
		return true
		}

	matchRefLink?(open, _document)
		{
		// .i is ']'
		// shortcut reference link
		label =  .inline[open.value.pos + open.value[0].Size()...i]
		end = .i+1
		if false isnt labelOb = Md_Helper.MatchLinkLabel(.inline, .i+1, allowBlank?:)
			{
			if labelOb.s isnt ''
				label = labelOb.s
			end = labelOb.end
			}

		if false is def = document.GetLinkDefs(Md_Helper.NormalizeLinkLabel(label))
			return false

		.addLinkImage(open, def.dest, def.title)
		.set(end)
		return true
		}

	addLinkImage(open, dest, title)
		{
		type = open.value[0] is '[' ? #link : #image
		linkText = .extractSubList(.nodes, open.value.node)
		.nodes.Del(open.value.node)
		.addNode(Object(type,
			.processEmphasis(linkText, bottom: open), href: dest, :title))
		if type is #link
			for (item = .delimiters.Prev(open);
				item isnt false;
				item = .delimiters.Prev(item))
				{
				if item.value[0] is '['
					item.value.active = false
				}
		.delimiters.Del(open)
		}

	extractSubList(list, after, before = false)
		{
		if Same?(start = list.Next(after), before)
			return DoubleLinkedList()
		end = list.Prev(before)
		list.Extract(start, end)
		return DoubleLinkedList(start, end)
		}

	emphasis()
		{
		.flushText()
		start = .i
		prev = .inline[.i-1::1]
		ch = .cur()
		length = 0
		do
			{
			length++
			.advance()
			}
		while .cur() is ch
		node = .addTextNode(.inline[start...i])

		next = .inline[.i::1]
		left? = .isLeftFlanking?(prev, next)
		right? = .isLeftFlanking?(next, prev)
		open? = ch is '*' ?
			left? :
			left? and (not right? or prev =~ '[[:punct:]]') // for '_'
		close? = ch is '*' ?
			right? :
			right? and (not left? or next =~ '[[:punct:]]')

		.delimiters.Append(Object(ch, active:, :node, pos: start,
			:length, :open?, :close?))
		.set(.i)
		}

	isLeftFlanking?(prev, next)
		{
		if next.Blank?()
			return false
		if next =~ '[[:punct:]]' and not prev.Blank?() and prev !~ '[[:punct:]]'
			return false
		return true
		}

	three: 3 // Rule 9 & 10 of Emphasis and strong emphasis
	processEmphasis(nodeList, bottom = false)
		{
		openers = Object(
			'*': Object(bottom, bottom, bottom),
			'_': Object(bottom, bottom, bottom))
		cur = .delimiters.Next(bottom)
		while cur isnt false
			{
			curDelimiter = cur.value
			if curDelimiter[0] not in ('*', '_') or curDelimiter.close? isnt true
				{
				cur = .delimiters.Next(cur)
				continue
				}

			origLength = curDelimiter.length
			prev = .delimiters.Prev(cur)
			while not Same?(openers[curDelimiter[0]][origLength % .three],
				prev)
				{
				prevDelimiter = prev.value
				if prevDelimiter[0] is curDelimiter[0] and prevDelimiter.open? is true
					{
					if ((curDelimiter.open? and curDelimiter.close? or
						prevDelimiter.open? and prevDelimiter.close?) and
						(curDelimiter.length + prevDelimiter.length) % .three is 0 and
						curDelimiter.length % .three isnt 0)
						{
						prev = .delimiters.Prev(prev)
						continue
						}

					type = curDelimiter.length >= 2 and prevDelimiter.length >= 2
						? #strong
						: #regular
					len = type is #strong ? 2 : 1
					subList = .extractSubList(nodeList,
						prevDelimiter.node, curDelimiter.node)
					nodeList.Insert(
						Object(#emph,
							subList.ToList(.mergeTextNode), strong?: type is #strong),
						after: prevDelimiter.node)

					.extractSubList(.delimiters, prev, cur)
					prevDelimiter.length -= len
					if prevDelimiter.length is 0
						{
						nodeList.Del(prevDelimiter.node)
						temp = .delimiters.Prev(prev)
						.delimiters.Del(prev)
						prev = temp
						}
					else
						prevDelimiter.node.value[1] = prevDelimiter.node.value[1][..-len]

					curDelimiter.length -= len
					if curDelimiter.length is 0
						break
					else
						curDelimiter.node.value[1] = curDelimiter.node.value[1][..-len]
					}
				else
					prev = .delimiters.Prev(prev)
				}
			if curDelimiter.length is 0
				{
				nodeList.Del(curDelimiter.node)
				temp = .delimiters.Next(cur)
				.delimiters.Del(cur)
				cur = temp
				}
			else
				{
				openers[curDelimiter[0]][origLength % 3/*=fixed*/] = .delimiters.Prev(cur)
				old = cur
				cur = .delimiters.Next(cur)
				if curDelimiter.open? is false
					.delimiters.Del(old)
				}
			}

		while false isnt node = .delimiters.Next(bottom)
			.delimiters.Del(node)

		return nodeList.ToList(.mergeTextNode)
		}

	set(.i)
		{
		.base = .i
		}

	flushText(to = false, beforeLineBreak? = false)
		{
		if to is false
			to = .i
		if to <= .base
			return
		s = .inline[.base..to]
		if .nodes.Empty?() or .nodes.Prev().value[0] is #linkbreak
			s = s.LeftTrim()
		if beforeLineBreak? is true
			s = s.RightTrim()
		.addNode(Object(#text, Md_Helper.Escape(s)))
		.base = to
		}

	addTextNode(text, codepoints = false)
		{
		node = Object(#text, text)
		if codepoints isnt false
			node.codepoints = codepoints
		.addNode(node)
		}

	addNode(node)
		{
		return .nodes.Append(node)
		}

	cur()
		{
		if .i >= .inline.Size()
			return false
		return .inline[.i]
		}

	next()
		{
		if .i + 1 >= .inline.Size()
			return ''
		return .inline[.i + 1]
		}

	advance(by = 1)
		{
		.i += by
		}

	mergeTextNode(prev, cur)
		{
		if prev[0] is #text and cur[0] is #text and
			not prev.Member?(#codepoints) and not cur.Member?(#codepoints)
			{
			prev[1] $= cur[1]
			return true
			}
		return false
		}
	}