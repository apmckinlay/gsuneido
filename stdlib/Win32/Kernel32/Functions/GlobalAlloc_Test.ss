// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// TAGS: win32
Test
	{
	Test_data()
		{
		test = function (s)
			{
			hm = GlobalAllocData(s)
			Assert(GlobalSize(hm) is: s.Size())
			Assert(GlobalData(hm) is: s)
			}
		test("")
		test('x')
		test('\x00')
		test("hello world")
		test("hello world\x00".Repeat(100))
		}
	Test_string()
		{
		test = function (s)
			{
			hm = GlobalAllocString(s)
			s = s.BeforeFirst('\x00')
			Assert(GlobalSize(hm) is: s.Size() + 1) // +1 for nul terminator
			Assert(GlobalString(hm) is: s)
			}
		test("")
		test('x')
		test('\x00')
		test("hello\x00world")
		test("hello world".Repeat(100))
		}
	}
