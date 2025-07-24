// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Html_table(#(((topleft da: 1), topright ra: 2)
			(bottomleft, bottomright), ta: 3))
			is: '<table ta="3">\n' $
				'<tr ra="2"><td da="1">topleft</td><td>topright</td></tr>\n' $
				'<tr><td>bottomleft</td><td>bottomright</td></tr>\n</table>')
		}
	}