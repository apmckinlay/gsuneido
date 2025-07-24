// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_GetFontSize()
		{
		Assert(ReportBase.GetFontSize(#()) is: 10)
		Assert(ReportBase.GetFontSize(#(size: 12)) is: 12)

		rpt = new ListFormatting(false, false)
		Assert(rpt.GetFontSize(#(size: '+3')) is: 13)
		Assert(rpt.GetFontSize(#(size: '-3')) is: 7)

		rpt.ListFormatting_curFont = #(size: 11)
		Assert(rpt.GetFontSize(#(size: '+3')) is: 14)
		Assert(rpt.GetFontSize(#(size: '-3')) is: 8)

		rpt = new .rptClass
		Assert(rpt.Eval(Report.GetFontSize, #(size: '+3')) is: 13)
		Assert(rpt.Eval(Report.GetFontSize, #(size: '-3')) is: 7)

		rpt.SetFont(#(size: 11))
		Assert(rpt.Eval(Report.GetFontSize, #(size: '+3')) is: 14)
		Assert(rpt.Eval(Report.GetFontSize, #(size: '-3')) is: 8)
		}

	rptClass: ReportBase
		{
		curFont: false
		GetFont()
			{
			return .curFont
			}
		SetFont(.curFont)
			{
			}
		}
	}