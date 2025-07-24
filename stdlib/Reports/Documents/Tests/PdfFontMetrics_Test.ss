// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_GetFontReference()
		{
		f = PdfFontMetrics.GetFontReference

		Assert(f(font: Object()) is: '')
		Assert(f(font: Object(italic:)) is: 'I')

		weight = 550
		Assert(f(font: Object(:weight)) is: '')
		Assert(f(font: Object(:weight, italic: false)) is: '')
		Assert(f(font: Object(:weight, italic: true))  is: 'I')

		weight = 'normal' // FW.NORMAL = 400
		Assert(f(font: Object(:weight, italic: false)) is: '')
		Assert(f(font: Object(:weight, italic: true))  is: 'I')

		weight = 'medium' // FW.MEDIUM = 500
		Assert(f(font: Object(:weight, italic: false)) is: '')
		Assert(f(font: Object(:weight, italic: true))  is: 'I')

		weight = 551
		Assert(f(font: Object(:weight, italic: false)) is: 'B')
		Assert(f(font: Object(:weight, italic: true))  is: 'BI')

		weight = 'semibold' // FW.SEMIBOLD = 600
		Assert(f(font: Object(:weight, italic: false)) is: 'B')
		Assert(f(font: Object(:weight, italic: true))  is: 'BI')

		weight = 'bold' // FW.BOLD = 700
		Assert(f(font: Object(:weight, italic: false)) is: 'B')
		Assert(f(font: Object(:weight, italic: true))  is: 'BI')
		}
	}