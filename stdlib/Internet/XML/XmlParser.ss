// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// uses SAX to parse XML and builds a tree of XmlNode's
XmlContentHandler
	{
	CallClass(s)
		{
		if s isnt '' and not s.Prefix?('<')
			throw 'Invalid xml format'
		return (new this).Parse(s)
		}
	Parse(text)
		{
		.stack = Stack()
		xr = new XmlReader
		xr.SetContentHandler(this)
		.node = false
		xr.Parse(text)
		return .node
		}
	StartElement(qname, atts)
		{
		.stack.Push(XmlNode(qname, attributes: atts))
		}
	Characters(text)
		{
		.stack.Top().AddChild(XmlNode(:text))
		}
	EndElement(qname)
		{
		if .stack.Top().Name() isnt qname
			throw "unmatched tag: " $ qname
		.node = .stack.Pop()
		if .stack.Count() > 0
			.stack.Top().AddChild(.node)
		}
	}