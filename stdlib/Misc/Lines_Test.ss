// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Lines('') is: [])
		Assert(Lines('hello') is: ['hello'])
		Assert(Lines('hello\nworld') is: ['hello', 'world'])
		Assert(Lines('hello\r\nworld') is: ['hello', 'world'])
		Assert(Lines('\nhello\n\nworld\n') is: ['', 'hello', '', 'world'])
		}
	}