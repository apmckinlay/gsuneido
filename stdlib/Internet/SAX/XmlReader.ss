// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// a SAX like XML parser
// see XmlContentHandler for interface
// TODO: add SetLexicalHandler
class
	{
	New()
		{
		.contentHandler = XmlContentHandler
		}
	SetContentHandler(contentHandler)
		{
		.contentHandler = contentHandler
		}
	Parse(text)
		{
		.size = text.Size()
		forever
			{
			i = text.Find('<')
			if i > 0
				{
				s = text[.. i]
				if s.White?()
					.contentHandler.IgnorableWhitespace(s)
				else
					.contentHandler.Characters(XmlEntityDecode(s))
				}
			text = text[i + 1 ..]
			if text is ''
				break
			if false isnt t = .comment(text)
				{
				text = t
				continue
				}
			if false isnt t = .cdata(text)
				{
				text = t
				continue
				}
			if false isnt t = .code(text)
				{
				text = t
				continue
				}
			j = ScannerFindIf(text, { it.Prefix?('>') })
			.processElement(text, j)
			text = text[j + 1 ..]
			}
		}
	comment(text)
		{
		if text.Prefix?("!--")
			{ // comment
			commentSuffix = "-->"
			j = text.Find(commentSuffix)
			text = text[j + commentSuffix.Size() ..]
			return text
			}
		return false
		}
	cdata(text)
		{
		dataPrefix = "![CDATA["
		if text.Prefix?(dataPrefix)
			{
			dataSuffix = "]]>"
			j = text.Find(dataSuffix)
			// should this be XmlEntityDecode'd ???
			dataPrefixSize = dataPrefix.Size()
			.contentHandler.Characters(text[dataPrefixSize :: j - dataPrefixSize])
			return text[j + dataSuffix.Size() ..]
			}
		return false
		}
	code(text)
		{
		if text.Prefix?("$")
			{
			j = ScannerFindPrefix(text, "$>")
			if eatnewline = text[j - 1] is '-'
				--j
			.contentHandler.Code(text[1 :: j - 2])
			text = text[j + 2 + (eatnewline ? 1 : 0) ..]
			if eatnewline
				{
				if text[0] is '\r'
					text = text[1 ..]
				if text[0] is '\n'
					text = text[1 ..]
				}
			return text
			}
		return false
		}
	processElement(text, j)
		{
		tag = text[.. j]
		i = tag.Find1of(' \r\n\t')
		qname = tag[..i].Lower()
		atts = .attributes(tag[i+1..])
		if qname.Prefix?('?')
			; // do nothing
		else if qname.Prefix?('/')
			.contentHandler.EndElement(qname[1 ..], pos: .size - text.Size())
		else if tag.Suffix?('/')
			{
			qname = qname.Tr('/')
			.contentHandler.StartElement(qname, atts, pos: .size - text.Size())
			.contentHandler.EndElement(qname, pos: .size - text.Size())
			}
		else
			.contentHandler.StartElement(qname, atts, pos: .size - text.Size())
		}
	attributes(s)
		{
		atts = Object()
		curState = Object(state: 'name', name: false)
		for (scan = Scanner(s); scan isnt (token = scan.Next()); )
			{
			if token in ('/', '?')
				break
			if scan.Type() in (#WHITESPACE, #NEWLINE)
				continue
			switch curState.state
				{
			case 'name' :
				Assert(scan.Type() is: #IDENTIFIER,
					msg: "XmlReader: expecting identifier, got: " $ scan.Text())
				curState.name = token
				curState.state = '='
			case '=' :
				.processEqual(scan.Text(), curState, atts, token)
			case 'value' :
				.valueAttr(scan.Type(), scan.Value(), atts, curState)
				}
			}
		return atts
		}
	processEqual(text, curState, atts, token)
		{
		if text is '='
			curState.state = 'value'
		else
			{
			atts[curState.name.Lower()] = true
			curState.name = token
			}
		}
	valueAttr(type, value, atts, curState)
		{
		if type isnt #STRING and type isnt #NUMBER and type isnt #IDENTIFIER
			throw "XmlReader: expecting string, number, or identifier"
		atts[curState.name.Lower()] = XmlEntityDecode(value)
		curState.state = 'name'
		}
	}