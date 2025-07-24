// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// TAGS: win32
Test
	{
	Test_main()
		{
		Assert(LONG(#(x: 0x12345678)).ToHex() is: "78563412")
		}
	}
