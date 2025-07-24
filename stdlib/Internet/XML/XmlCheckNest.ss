// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
XmlContentHandler
	{
	CallClass(s)
		{
		(new this).Process(s)
		}
	Process(.text)
		{
		.stack = Stack()
		xr = new XmlReader
		xr.SetContentHandler(this)
		xr.Parse(text)
		unclosed = .stack.List().Join(", ")
		if unclosed isnt ""
			throw "unclosed tags: " $ unclosed
		}
	StartElement(qname, atts/*unused*/, pos)
		{
		line = .text.LineFromPosition(pos) + 1
		.stack.Push(qname $ " @ " $ line)
		}
	EndElement(qname, pos)
		{
		line = .text.LineFromPosition(pos) + 1
		if .stack.Count() is 0
			throw "unmatched closing tag: " $ qname $ " @ " $ line
		if .stack.Top().BeforeFirst(' ') isnt qname
			throw "unmatched closing tag: " $ qname $ " @ " $ line $
				" expecting: " $ .stack.Top()
		.stack.Pop()
		}
	}