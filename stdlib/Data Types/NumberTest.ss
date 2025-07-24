// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// This is the test for the built-in number methods
// See also: NumbersTest and Number_Test
// SuJsWebTest
Test
	{
	Test_Format()
		{
		Assert(0.Format('###') is: '0')
		Assert(0.Format('#.##') is: '.00')
		Assert(0.Format('###.') is: '0.')
		Assert(7689924.9449999.Format('################.##'),
			is: "7689924.94")
		}
	Test_mathfunc()
		{
		pi = 3.1415926535
		.asserteq((pi / 2).Sin(), 1)
		.asserteq((pi / 2).Cos(), 0)
		.asserteq((pi / 4).Tan(), 1)
		.asserteq(.32696.ASin(), .333085)
		.asserteq(.32696.ACos(), 1.237711)
		.asserteq(-862.42.ATan(), -1.569637)
		.asserteq(144.Sqrt(), 12)
		.asserteq(3.Pow(3), 27)
		.asserteq(10.Log().Exp(), 10)
		.asserteq(1000.Log10(), 3)
		}
	asserteq(x, y)
		{
		Assert((x - y).Abs() lessThan: 1e-6, msg: x $ " isnt " $ y)
		}
	Test_Round()
		{
		Assert(0.Round(0) is: 0)
		Assert(0.Round(2) is: 0)
		Assert(0.Round(-2) is: 0)
		Assert(1.Round(0) is: 1)
		Assert(12.34.Round(0) is: 12)
		Assert(12.56.Round(0) is: 13)
		Assert(12.34.Round(1) is: 12.3)
		Assert(12.56.Round(1) is: 12.6)
		Assert(12.34.Round(2) is: 12.34)
		Assert(12.56.Round(2) is: 12.56)
		Assert(12.34.Round(3) is: 12.34)
		Assert(12.56.Round(3) is: 12.56)
		Assert(1234.Round(-1) is: 1230)
		Assert(1256.Round(-1) is: 1260)
		Assert(1234.Round(-2) is: 1200)
		Assert(1256.Round(-2) is: 1300)

		Assert((1 / 3).Round(2) is: .33)
		Assert((2 / 3).Round(2) is: .67)
		Assert((-1 / 3).Round(2) is: -.33)
		Assert((-2 / 3).Round(2) is: -.67)

		Assert(.5.Round(0) is: 1)
		Assert(9999.Round(-4) is: 10000)

		Assert(123456789012.34.Round(2) is: 123456789012.34)
		Assert(123456789012345.Round(-2) is: 123456789012300)

		Assert(Display(0.Round(2)) is: "0")
		}
	Test_RoundDown()
		{
		Assert((1 / 3).RoundDown(2) is: .33)
		Assert((2 / 3).RoundDown(2) is: .66)

		Assert((-1 / 3).RoundDown(2) is: -.33)
		Assert((-2 / 3).RoundDown(2) is: -.66)
		Assert((-23.00).RoundDown(0) is: -23)
		Assert((-23.35).RoundDown(0) is: -23)
		Assert((-23.50).RoundDown(0) is: -23)
		Assert((-23.67).RoundDown(0) is: -23)
		Assert((-999.9999999).RoundDown(0) is: -999)
		Assert(123456789012.99.RoundDown(1) is: 123456789012.9)
		Assert(123456789012399.RoundDown(-2) is: 123456789012300)

		Assert(23.67.RoundDown(0) is: 23)
		Assert(23.50.RoundDown(0) is: 23)
		Assert(23.35.RoundDown(0) is: 23)
		Assert(23.00.RoundDown(0) is: 23)
		Assert(0.RoundDown(0) is: 0)
		Assert(23.50.RoundDown(1) is: 23.5)
		Assert(23.35.RoundDown(1) is: 23.3)
		Assert(23.00.RoundDown(1) is: 23.0)
		Assert(23.50.RoundDown(-1) is: 20)
		Assert(23.35.RoundDown(-1) is: 20)
		Assert(23.00.RoundDown(-1) is: 20)
		Assert(999.9999999.RoundDown(0) is: 999)
		}
	Test_RoundUp()
		{
		test = function (x, r, expected)
			{
			result = x.RoundUp(r)
			if result isnt expected
				throw x $ ".RoundUp(" $ r $ ")\nexpected " $ expected $ " got " $ result
			}

		test(1 / 3, 2, .34)
		test(2 / 3, 2, .67)

		test(-1 / 3, 2, -.34)
		test(-2 / 3, 2, -.67)
		test(-23.00, 0, -23)
		test(-23.35, 0, -24)
		test(-23.50, 0, -24)
		test(-23.67, 0, -24)
		test(-999.9999999, 0, -1000)
		test(123456789012.11, 1, 123456789012.2)
		test(123456789012311, -2, 123456789012400)

		test(23.67, 0, 24)
		test(23.50, 0, 24)
		test(23.35, 0, 24)
		test(23.00, 0, 23)
		test(0, 0, 0)
		test(23.50, 1, 23.5)
		test(23.35, 1, 23.4)
		test(23.00, 1, 23.0)
		test(23.50, -1, 30)
		test(23.35, -1, 30)
		test(23.00, -1, 30)
		test(999.9999999, 0, 1000)
		}
	Test_Hex()
		{
		// SuJsWebTest Excluded
		test = function (n, expected)
			{
			Assert(n.Hex() is: expected)
			}
		test(0, "0")
		test(0.0, "0")
		test(1, "1")
		test(1.0, "1")
		test(0xffffffff, "ffffffff")
		test(0x7fffffff, "7fffffff")
		test(2147483647, "7fffffff")
		test(2147483647.0, "7fffffff")
		test(0x80000000, "80000000")
		if BuiltDate() < #20250430
			return
		test(-1, "ffffffffffffffff")
		test(-1.0, "ffffffffffffffff")
		}
	}