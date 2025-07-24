// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// This is the test for the Number function
// See also: NumberTest and NumbersTest
// SuJsWebTest
Test
	{
	Test_one()
		{
		test = function (x, expected)
			{
			Assert(Number(x) is: expected)
			}
		test(0, 0)
		test(123, 123)
		test(-123, -123)
		test('1,234,567', 1234567)
		test('', 0)
		test('   ', 0)
		test('0', 0)
		test('00', 0)
		test('01', 1)
		test('123', 123)
		test(' 123 ', 123)
		test('0xff', 255) // hex
		test('0777', 777) // NOT octal
		test('-0777', -777) // NOT octal
		test('+123', 123)
		test('123.', 123)
		test(false, 0)

		cant = function (x)
			{
			Assert({ Number(x) } throws: "can't convert")
			}
		cant('.')
		cant('0x')
		cant('0xt')
		cant('foo')
		cant('false')
		cant(#20170718)
		}
	}