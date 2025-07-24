// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		XmlCheckNest("")
		XmlCheckNest("hello world")
		XmlCheckNest("<p>paragraph</p>")
		XmlCheckNest("<p>one <b>two</b> three</p>")
		XmlCheckNest("<p>one<br />two</p>")

		Assert({ XmlCheckNest("<p>") } throws: "unclosed tags: p @ 1")
		Assert({ XmlCheckNest("</p>") } throws: "unmatched closing tag: p @ 1")
		Assert({ XmlCheckNest("<a>\n<b>") } throws: "unclosed tags: a @ 1, b @ 2")
		Assert({ XmlCheckNest("<p><i></p></i>") }
			throws: "unmatched closing tag: p @ 1 expecting: i @ 1")
		}
	}