// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// an XmlContentHandler to write xml text
// TODO: break lines and indent i.e. pretty print
XmlContentHandler
	{
	New()
		{
		.text = ""
		.element = false
		}
	StartElement(qname, atts)
		{
		.flush()
		s = '<' $ qname
		for m in atts.Members()
			s $= ' ' $ m $ '=' $ '"' $ atts[m] $ '"'
		.element = s
		}
	Characters(string)
		{
		.flush()
		.text $= XmlEntityEncode(string)
		}
	EndElement(qname)
		{
		if .element is false
			.text $= '</' $ qname $ '>'
		else
			{
			.text $= .element $ ' />'
			.element = false
			}
		}
	AddElement(tag, value, atts = #())
		{
		.StartElement(tag, atts)
		.Characters(value)
		.EndElement(tag)
		}
	flush()
		{
		if .element is false
			return
		.text $= .element $ '>'
		.element = false
		}
	GetText()
		{
		.flush()
		return .text
		}
	}