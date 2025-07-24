// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = GetCodeRes
			{
			ReturnWithDelay()
				{
				return -1
				}
			}
		Assert(cl(#(queryvalues: (name: 1))) is: -1)
		Assert(cl(#(queryvalues: (name: 'black.bmp'))) is: -1)

		Assert(cl(#(queryvalues: (name: 'calendar1_7.css', lib: 'invalid'))) is: -1)

		x = cl(#(queryvalues: (name: 'calendar1_7.css')))
		y = cl(#(queryvalues: (lib: 'stdlib', name: 'calendar1_7.css')))

		y[1].Expires = x[1].Expires
		Assert(x is: y)
		Assert(x[0] is: 'OK')
		Assert(x[1].Content_Type is: 'text/css')
		Assert(x[2] startsWith: "html")
		}
	}