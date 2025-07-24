// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// This is the first line of PageHeadFormat
// It handles shrinking the font size of the center title so it fits
Format
	{
	New(title)
		{
		datetime = Date().LongDate() $ ' ' $ Date().Time()
		.left = _report.Construct(Object('Text', datetime))

		.center = _report.Construct(Object('Text', title, justify: "center"))

		page =  TranslateLanguage("Page") $ " " $ _report.GetPage()
		.right = _report.Construct(Object('Text', page, justify: "right"))
		}
	GetSize(data/*unused*/='')
		{
		if .fontsize is false
			.setFontSize()
		oldfont = _report.SelectFont(Object(size: .fontsize, weight: 600))
		szc = .center.GetSize()
		szc.w = .pageWidth()
		_report.SelectFont(oldfont)
		return szc
		}
	fontsize: false
	setFontSize()
		{
		oldfont = _report.SelectFont(#(size: 8))
		szl = .left.GetSize()
		szr = .right.GetSize()
		max = .pageWidth() - 2 * Max(szl.w, szr.w) - 360
		szc = false
		for (.fontsize = 13; .fontsize > 8; --.fontsize)
			{
			_report.SelectFont(Object(size: .fontsize, weight: 600))
			szc = .center.GetSize()
			if szc.w < max
				break
			}
		.d = (szc.h - szc.d) - (szl.h - szl.d)
		_report.SelectFont(oldfont)
		}
	Print(x, y, w, h, data/*unused*/='')
		{
		oldfont = _report.SelectFont(#(size: 8))
		.left.Print(x, y + .d, w, h)
		.right.Print(x, y + .d, w, h)
		_report.SelectFont(Object(size: .fontsize, weight: 600))
		.center.Print(x, y, w, h)
		_report.SelectFont(oldfont)
		}
	pageWidth()
		{
		d = _report.GetDimens()
		return d.W
		}
	}