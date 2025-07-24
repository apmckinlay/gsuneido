// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_buildFontPDFSize()
		{
		m = EnsurePDFFont.EnsurePDFFont_buildFontPDFSize
		Assert(m([], []) is: [])
		Assert(m([size: '+1'], [size: 9]) is: [size: 10])
		Assert(m([], [size: 9]) is: [size: 9])
		Assert(m([size: 12], [size: 9]) is: [size: 12])
		}
	}