// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
XmlContentHandler
	{
	StartElement(qname, atts)
		{ Print('START', qname, atts) }
	EndElement(qname)
		{ Print('END', qname) }
	Characters(s)
		{ Print(s) }
	}