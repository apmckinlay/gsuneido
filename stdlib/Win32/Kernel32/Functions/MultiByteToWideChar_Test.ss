// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		test = function (mb, wc)
			{
			Assert(MultiByteToWideChar(mb) is: wc)
			Assert(WideCharToMultiByte(wc) is: mb)
			}
		test("", "\x00\x00")
		test("###", "#\x00#\x00#\x00\x00\x00")
		// TODO more tests!
		}
	}