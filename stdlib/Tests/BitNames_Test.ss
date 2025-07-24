// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(BitNames(0, 'WS') is: "0")
		Assert(BitNames(0x40000002, 'WS', 'TTS') is: "WS.CHILD | TTS.NOPREFIX")
		Assert(BitNames(0x0200, 'WM') is: "WM.MOUSEMOVE")
		Assert(BitNames(0x020a, 'WM') is: "WM.MOUSEWHEEL")
		Assert(BitNames(WM.KEYUP, 'WM') is: "WM.KEYUP")
		}
	}