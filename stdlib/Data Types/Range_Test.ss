// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// SuJsWebTest
Test
	{
	// interface:
	Setup()
		{
		.range = new Range(5, 10)
		}
	Test_Main()
		{
		Assert(.range.GetLow() is: 5)
		Assert(.range.GetHigh() is: 10)
		}
	Test_Contains()
		{
		Assert(.range.Contains?(7))
		Assert(.range.Contains?(5))
		Assert(.range.Contains?(10))
		Assert(.range.Contains?(-5) is: false)
		Assert(.range.Contains?(11) is: false)
		Assert(.range.Contains?(4) is: false)
		}
	Test_Overlaps()
		{
		Assert(.range.Overlaps?(Range(1, 2)) is: false)
		Assert(.range.Overlaps?(Range(1, 6)))
		Assert(.range.Overlaps?(Range(5, 10)))
		Assert(.range.Overlaps?(Range(6, 9)))
		Assert(.range.Overlaps?(Range(6, 12)))
		Assert(.range.Overlaps?(Range(1000, 1230)) is: false)
		Assert(.range.Overlaps?(Range(-12, 56)))
		}
	Test_Plus()
		{
		Assert(Range(10, 20).Plus(Range(15, 25)) is: Range(10, 25))
		Assert(Range(5, 10).Plus(Range(10,25)) is: Range(5, 25))
		Assert(Range(5, 10).Plus(Range(100, 200)) is: Range(5, 200))
		Assert(Range(5, 10).Plus(Range(-4, -2)) is: Range(-4, 10))
		Assert(Range(5, 10).Plus(Range(2, 7)) is: Range(2, 10))
		Assert(Range(5, 10).Plus(Range(8, 19)) is: Range(5, 19))
		}
	Test_Minus()
		{
		Assert(Range(10, 20).Minus(Range(15, 25)) is: Range(10, 15))
		Assert(Range(5, 10).Minus(Range(10,25)) is: Range(5, 10))
		Assert(Range(5, 10).Minus(Range(100, 200)) is: Range(5, 10))
		Assert(Range(5, 10).Minus(Range(-4, -2)) is: Range(5, 10))
		Assert(Range(5, 10).Minus(Range(2, 7)) is: Range(7, 10))
		Assert(Range(5, 10).Minus(Range(8, 19)) is: Range(5, 8))
		Assert(Range(5, 10).Minus(Range(8, 10)) is: Range(5, 8))
		Assert(Range(5, 10).Minus(Range(0, 20)).Empty?())
		Assert(Range(#20010718, #20010802).Minus(Range(#20010718, #20010719))
			is: Range(#20010719, #20010802))
		}
	Test_Separate()
		{
		Assert(
			Range(0, 10).Separate(Range(0, 10)) is: Object(Range(0, 0), Range(10, 10)))
		Assert(
			Range(0, 10).Separate(Range(5, 10)) is: Object(Range(0, 5), Range(10, 10)))
		Assert(
			Range(0, 10).Separate(Range(0, 7)) is: Object(Range(0, 0),	Range(7, 10)))
		Assert(
			Range(0, 10).Separate(Range(3, 4.2)) is: Object(Range(0, 3), Range(4.2, 10)))
		}
	}
