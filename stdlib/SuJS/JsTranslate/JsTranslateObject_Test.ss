// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		test = function (src, expected)
			{ Assert(JsTranslate(src) is: expected) }
		test('()', 'su.empty_object')
		test('(1)', 'su.mkObject(1)')
		test('(1, "x")', 'su.mkObject(1, "x")')
		test('(a:)', 'su.mkObject(null, "a", true)')
		test('(a: 1)', 'su.mkObject(null, "a", 1)')
		test('(a: 1, 9: 2)', 'su.mkObject(null, "a", 1, 9, 2)')
		test('(a?: 1, c?:)', 'su.mkObject(null, "a?", 1, "c?", true)')
		test('(1, 2, a: 3, 9: 4)', 'su.mkObject(1, 2, null, "a", 3, 9, 4)')
		test('(1, (2))', 'su.mkObject(1, su.mkObject(2))')
		}
	}
