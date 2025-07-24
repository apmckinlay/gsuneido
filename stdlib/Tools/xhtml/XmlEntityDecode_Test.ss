// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	decodeValues: #(
		#('&amp;', '&'),
		#('&lt;', '<'),
		#('&gt;', '>'),
		#('&quot;', '"'),
		#('&apos;', "'"),
		// check escaped text
		#('&amp;amp;', '&amp;'),
		#('&amp;lt;', '&lt;'),
		#('&amp;gt;', '&gt;'),
		#('&amp;quot;', '&quot;'),
		#('&amp;apos;', '&apos;')
		)
	Test_main()
		{
		for value in .decodeValues
			Assert(XmlEntityDecode(value[0]) is: value[1])
		}
	}