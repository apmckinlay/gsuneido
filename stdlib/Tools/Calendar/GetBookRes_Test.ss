// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		cl = GetBookRes
			{
			ReturnWithDelay()
				{
				return -1
				}
			}
		Assert(cl(#(queryvalues: (name: 1))) is: -1)
		Assert(cl(#(queryvalues: (name: 'black.bmp', book: 'invalid'))) is: -1)

		x = cl(#(queryvalues: (name: 'foo')))
		Assert(x is: -1)

		x = cl(#(queryvalues: (name: 'black.bmp')))
		y = cl(#(queryvalues: (book: 'imagebook', name: 'black.bmp')))

		y[1].Expires = x[1].Expires
		Assert(x is: y)
		Assert(x[0] is: 'OK')
		Assert(x[1].Content_Type is: 'image/bmp')
		Assert(x[2] startsWith: "BM")
		}
	}
