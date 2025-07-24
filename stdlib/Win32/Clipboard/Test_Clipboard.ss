// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// TAGS: win32
Test
	{
	Setup()
		{
		super.Setup()
		.save = ClipboardReadString()
		}
	Test_data()
		{
		test = function(s)
			{
			Thread.Sleep(5)
			Assert(ClipboardWriteData(s, CF.BITMAP) isnt: 0)
			Thread.Sleep(5)
			Assert(ClipboardReadData(CF.BITMAP) is: s)
			}
		test("")
		test("hello world")
		test("hello\x00world")
		test("hello world\x00".Repeat(100))
		}
	Test_string()
		{
		test = function(s)
			{
			Thread.Sleep(5)
			Assert(ClipboardWriteString(s) isnt: 0)
			Thread.Sleep(5)
			Assert(ClipboardReadString() is: s.BeforeFirst('\x00'))
			}
		test("")
		test("hello world")
		test("hello\x00world")
		test("hello world".Repeat(100))
		}
	Teardown()
		{
		if String?(.save)
			ClipboardWriteString(.save)
		super.Teardown()
		}
	}