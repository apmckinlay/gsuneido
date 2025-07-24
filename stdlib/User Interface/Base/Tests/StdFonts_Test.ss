// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_font()
		{
		Assert(StdFonts.Font('abc') is: 'abc')
		Assert(StdFonts.Font('@ui') is: StdFonts.Ui())
		Assert(StdFonts.Font('@mono') is: StdFonts.Mono())
		Assert(StdFonts.Font('@serif') is: StdFonts.Serif())
		Assert(StdFonts.Font('@sans') is: StdFonts.Sans())
		}

	Test_weight()
		{
		Assert(StdFonts.Weight(123) is: 123)
		Assert(StdFonts.Weight('normal') is: FW.NORMAL)
		Assert(StdFonts.Weight('bold') is: FW.BOLD)
		}

	Test_size()
		{
		Assert(StdFonts.FontSize(11) is: 11)
		Assert(StdFonts.FontSize(11, 9) is: 11)
		Assert(StdFonts.FontSize('+2', 10) is: 12)
		Assert(StdFonts.FontSize('-2', 10) is: 8)
		}

	Test_SciSize()
		{
		Assert(StdFonts.SciSize(11) is: 11)
		Assert(StdFonts.SciSize(11, 9) is: 11)
		Assert(StdFonts.SciSize('+2', 10) is: 12)
		Assert(StdFonts.SciSize('+2', 9)  is: 12)
		Assert(StdFonts.SciSize('-2', 10) is: 9)
		}
	}