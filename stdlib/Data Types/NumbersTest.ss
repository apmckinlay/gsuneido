// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// This is to test the methods in Numbers
// The tests for builtin number methods are in NumberTest
// SuJsWebTest
Test
	{
	Test_ToWords()
		{
		Assert(0.ToWords() is: 'ZERO')
		Assert(8.ToWords() is: 'EIGHT')
		Assert(18.ToWords() is: 'EIGHTEEN')
		Assert(20.ToWords() is: 'TWENTY')
		Assert(100.ToWords() is: 'ONE HUNDRED')
		Assert(750.ToWords() is: 'SEVEN HUNDRED AND FIFTY')
		Assert(1000.ToWords() is: 'ONE THOUSAND')
		Assert(1525.ToWords()
			is: 'ONE THOUSAND FIVE HUNDRED AND TWENTY FIVE')
		Assert(9999.ToWords()
			is: 'NINE THOUSAND NINE HUNDRED AND NINETY NINE')
		Assert(10000.ToWords() is: 'TEN THOUSAND')
		Assert(80750.ToWords()
			is: 'EIGHTY THOUSAND SEVEN HUNDRED AND FIFTY')
		Assert(100000.ToWords() is: 'ONE HUNDRED THOUSAND')
		Assert(785694.ToWords()
			is: 'SEVEN HUNDRED AND EIGHTY FIVE THOUSAND SIX HUNDRED AND NINETY FOUR')
		Assert(1000000.ToWords() is: 'ONE MILLION')
		Assert(1555555.ToWords()
			is: 'ONE MILLION FIVE HUNDRED AND FIFTY FIVE THOUSAND FIVE HUNDRED ' $
				'AND FIFTY FIVE')
		}
	Test_ToWordsSimple()
		{
		data = #(
			(0, "zero") (1, "one") (2, "two") (3, "three") (4, "four") (5, "five")
			(6, "six") (7, "seven") (8, "eight") (9, "nine") (10, "one zero")
			(11, "one one") (12, "one two") (13, "one three") (14, "one four")
			(15, "one five") (16, "one six") (17, "one seven") (18, "one eight")
			(19, "one nine") (20, "two zero") (30, "three zero") (40, "four zero")
			(50, "five zero") (60, "six zero") (70, "seven zero") (80, "eight zero")
			(90, "nine zero") (100, "one zero zero") (1000, "one zero zero zero")
			(10000, "one zero zero zero zero") (100000, "one zero zero zero zero zero")
			(1000000, "one zero zero zero zero zero zero")
			(1000000000, "one zero zero zero zero zero zero zero zero zero")
			(1000000000000, "one zero zero zero zero zero zero zero zero zero zero " $
				"zero zero")
			(1000000000000000, "one zero zero zero zero zero zero zero zero zero zero " $
				"zero zero zero zero zero")
			(122, "one two two") (123, "one two three")
			(123456789012345, "one two three four five six seven eight nine zero " $
				"one two three four five")
			(999999999999999, "nine nine nine nine nine nine nine nine nine nine " $
				"nine nine nine nine nine")
			(1.999999999999, "one decimalpoint nine nine nine nine nine nine nine " $
				"nine nine nine nine nine")
			(1.123456789012, "one decimalpoint one two three four five six seven " $
				"eight nine zero one two")
			(99999999999.9999, "nine nine nine nine nine nine nine nine nine nine " $
				"nine decimalpoint nine nine nine nine")
			(12345678901.1234, "one two three four five six seven eight nine zero " $
				"one decimalpoint one two three four")
			(999999999999.999, "nine nine nine nine nine nine nine nine nine nine " $
				"nine nine decimalpoint nine nine nine")
			(9999.99999999999, "nine nine nine nine decimalpoint nine nine nine nine " $
				"nine nine nine nine nine nine nine")
			(1234.12345678901, "one two three four decimalpoint one two three four " $
				"five six seven eight nine zero one")
			)
		minus = TranslateLanguage("minus") $ ' '
		xlatwords = function (s)
			{
			return s.Split(' ').Map!(TranslateLanguage).Join(' ')
			}
		for x in data
			{
			// translation of the test data
			xtl = xlatwords(x[1])
			// positive number
			Assert((x[0]).ToWordsSimple() is: "* " $ xtl $ " *")
			// negative number
			if x[0] isnt 0
				Assert((-x[0]).ToWordsSimple() is: "* " $ minus $ xtl $ " *")
			}
		}
	Test_Pad()
		{
		Assert(1.Pad(1) is: "1")
		Assert(1.Pad(5) is: "00001")
		Assert(1.Pad(7) is: "0000001")
		Assert(5672.Pad(5) is: "05672")
		Assert(5672.Pad(4) is: "5672")
		Assert(5672.Pad(2) is: "5672")
		Assert(5672.Pad(0) is: "5672")
		Assert(5672.Pad(10) is: "0000005672")
		Assert(5672.Pad(5, "x") is: "x5672")
		Assert((-12).Pad(4) is: "0012")
		}
	Test_Factorial()
		{
		results = #(1, 1, 2, 6, 24, 120)
		for (i in results.Members())
			Assert(i.Factorial() is: results[i])
		}
	Test_EvenOdd()
		{
		Assert(0.Even?() is: true)
		Assert(0.Odd?() is: false)
		Assert(1.Odd?() is: true)
		Assert(1.Even?() is: false)
		Assert(123.Odd?() is: true)
		Assert(123.Even?() is: false)
		Assert(1234.Even?() is: true)
		Assert(1234.Odd?() is: false)
		Assert((-1).Odd?() is: true)
		Assert((-2).Odd?() is: false)
		Assert((-2).Even?() is: true)
		Assert((-3).Even?() is: false)
		}
	Test_Ceiling()
		{
		Assert(0.Ceiling() is: 0)
		Assert(1.Ceiling() is: 1)
		Assert((-1).Ceiling() is: -1)
		Assert(1.5.Ceiling() is: 2)
		Assert((-1.5).Ceiling() is: -1)

		Assert(1.01.Ceiling() is: 2)
		Assert(2.01.Ceiling() is: 3)
		Assert((-1.99).Ceiling() is: -1)
		Assert((-2.99).Ceiling() is: -2)
		Assert(9.999999999.Ceiling() is: 10)
		Assert((-9.999999999).Ceiling() is: -9)
		}
	Test_Floor()
		{
		Assert(0.Floor() is: 0)
		Assert(1.Floor() is: 1)
		Assert((-1).Floor() is: -1)
		Assert(1.5.Floor() is: 1)
		Assert((-1.5).Floor() is: -2)

		Assert(1.01.Floor() is: 1)
		Assert(2.01.Floor() is: 2)
		Assert((-1.99).Floor() is: -2)
		Assert((-2.99).Floor() is: -3)
		Assert(9.999999999.Floor() is: 9)
		Assert((-9.999999999).Floor() is: -10)
		}
	Test_RoundToNearest()
		{
		nearest = -1
		Assert(0.RoundToNearest(nearest) is: 0)
		Assert(1.RoundToNearest(nearest) is: 1)
		Assert(5.RoundToNearest(nearest) is: 5)
		Assert(111111.RoundToNearest(nearest) is: 111111)

		nearest = 5
		Assert(0.RoundToNearest(nearest) is: 0)
		Assert(1.RoundToNearest(nearest) is: 0)
		Assert(3.RoundToNearest(nearest) is: 5)
		Assert(5.RoundToNearest(nearest) is: 5)
		Assert(7.RoundToNearest(nearest) is: 5)
		Assert(8.RoundToNearest(nearest) is: 10)
		Assert(11.RoundToNearest(nearest) is: 10)
		Assert(14.RoundToNearest(nearest) is: 15)

		nearest = 12
		Assert(0.RoundToNearest(nearest) is: 0)
		Assert(1.RoundToNearest(nearest) is: 0)
		Assert(6.RoundToNearest(nearest) is: 12)
		Assert(13.RoundToNearest(nearest) is: 12)
		Assert(17.RoundToNearest(nearest) is: 12)
		Assert(18.RoundToNearest(nearest) is: 24)

		nearest = 50
		Assert(0.RoundToNearest(nearest) is: 0)
		Assert(1.RoundToNearest(nearest) is: 0)
		Assert(49.RoundToNearest(nearest) is: 50)
		Assert(51.RoundToNearest(nearest) is: 50)
		Assert(74.RoundToNearest(nearest) is: 50)
		Assert(75.RoundToNearest(nearest) is: 100)
		}
	Test_FracDigits()
		{
		Assert(12.34.FracDigits() is: 2)
		Assert((-12.34).FracDigits() is: 2)
		Assert(.34.FracDigits() is: 2)
		Assert(12.FracDigits() is: 0)
		Assert(0.FracDigits() is: 0)
		}
	Test_Int?()
		{
		Assert(1.Int?())
		Assert(0.Int?())
		Assert((-1).Int?())
		Assert(0x80000000.Int?())
		Assert(0x7f000000.Int?())
		Assert((0.5).Int().Int?())
		Assert(not (0.1).Int?())
		Assert(not (-2.3).Int?())
		Assert(100000000000.Int?())
		Assert(not (-100.000000001).Int?())
		}
	Test_IntDigits()
		{
		Assert(12.34.IntDigits() is: 2)
		Assert((-12.34).IntDigits() is: 2)
		Assert(.34.IntDigits() is: 0)
		Assert(12.IntDigits() is: 2)
		Assert(0.IntDigits() is: 0)
		}
	Test_Sign()
		{
		Assert(123.Sign() is: 1)
		Assert(0.Sign() is: 1)
		Assert(-123.Sign() is: -1)
		}
	Test_RoundToPrecision()
		{
		Assert(0.RoundToPrecision(2) is: 0)
		Assert(1.RoundToPrecision(2) is: 1)
		Assert(12.RoundToPrecision(2) is: 12)
		Assert(123.RoundToPrecision(2) is: 120)
		Assert(125.RoundToPrecision(2) is: 130)
		Assert(1.234.RoundToPrecision(2) is: 1.2)
		Assert(.1234.RoundToPrecision(2) is: .12)
		Assert(.01234.RoundToPrecision(2) is: .012)
		Assert(.001234.RoundToPrecision(2) is: .0012)

		Assert(.2222.RoundToPrecision(2) is: .22)
		Assert(.8888.RoundToPrecision(2) is: .89)
		Assert(2222.RoundToPrecision(2) is: 2200)
		Assert(8888.RoundToPrecision(2) is: 8900)
		Assert(-.2222.RoundToPrecision(2) is: -.22)
		Assert(-.8888.RoundToPrecision(2) is: -.89)
		Assert(-2222.RoundToPrecision(2) is: -2200)
		Assert(-8888.RoundToPrecision(2) is: -8900)
		}
	Test_ToRGB()
		{
		verify = function (@args)
			{
			Assert(RGB(@args).ToRGB() is: args)
			}
		verify(0, 0, 0)
		verify(1, 2, 3)
		verify(0xff, 0xff, 0xff)

		verify = function (rgb)
			{
			Assert(RGB(@rgb.ToRGB()) is: rgb)
			}
		verify(0x000000)
		verify(0xffffff)
		verify(0x112233)
		}
	Test_Dollar()
		{
		Assert(0.DollarFormat('###.##') is: '$.00')
		Assert(100.DollarFormat('###') is: '$100')
		Assert(100.456.DollarFormat('###.##') is: '$100.46')
		Assert((-100).DollarFormat('-###') is: '$(100)')
		Assert((-100789.456).DollarFormat('-###,###.##') is: '$(100,789.46)')
		}

	Test_PercentToDecimal()
		{
		Assert(85.PercentToDecimal() is: .85)
		Assert(85.55.PercentToDecimal() is: .8555)
		}

	Test_DecimalToPercent()
		{
		Assert(.5678.DecimalToPercent() is: 57)
		Assert(.56789.DecimalToPercent(2) is: 56.79)
		}

	Test_HoursInSeconds()
		{
		Assert(1.HoursInSeconds() is: 3600)
		Assert(.5.HoursInSeconds() is: 1800)
		Assert(.1.HoursInSeconds() is: 360)
		Assert(2.5.HoursInSeconds() is: 9000)
		}
	}