// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_wr()
		{
		.wr(#((s tag ()), (e tag)), "<tag />")
		.wr(#((s tag ()), (c hello) (e tag)), "<tag>hello</tag>")
		.wr(#((c 'bad & chars < in " here')), "bad &amp; chars &lt; in &quot; here")
		.wr(#((s tag ()), (a tag2, 'word', ()), (e tag)), "<tag><tag2>word</tag2></tag>")
		}
	wr(list, string)
		{
		xw = new XmlWriter
		for x in list
			switch x[0]
				{
			case 's' : xw.StartElement(x[1], x[2])
			case 'c' : xw.Characters(x[1])
			case 'e' : xw.EndElement(x[1])
			case 'a' : xw.AddElement(x[1], x[2], x[3])
			default : throw 'bad data'
				}
		Assert(xw.GetText() is: string)
		}
	Test_rw()
		{
		.rw('')
		.rw('stuff')
		.rw('bad &amp; chars &lt; in &quot; here')
		.rw('<tag>')
		.rw('stuff<tag>')
		.rw('<tag>stuff')
		.rw('<tag />')
		.rw('<tag>stuff</tag>')
		.rw('<tag color="red">')
		}
	rw(text)
		{
		xr = new XmlReader
		xr.SetContentHandler(xw = new XmlWriter)
		xr.Parse(text)
		Assert(xw.GetText() is: text)
		}
	}