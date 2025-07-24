// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BadContrast?()
		{
		b = RGBColors.BadContrast?
		// contrast against grey
		Assert(b(#(255, 255, 255)) is: false) 	//white
		Assert(b(#(128, 128, 128)) is: false) 	//grey
		Assert(b(#(127, 127, 127)) is: true)
		Assert(b(#(1, 1, 1)) is: true)
		Assert(b(#(0, 0, 0)) is: true)			//black
		Assert(b(#(100, 183, 100)) is: true)
		Assert(b(#(100, 184, 100)) is: false)
		Assert(b(#(0, 229, 0)) is: true)
		Assert(b(#(0, 231, 0)) is: false) 		//bright green
		}
	Test_GetContrast()
		{
		g = RGBColors.GetContrast
		// defaults that are returned by the method
		dark = 0x000000
		light = 0xffffff
		Assert(g(RGB(255, 255, 255)) is: dark)	//white
		Assert(g(RGB(100, 184, 100)) is: dark)
		Assert(g(RGB(100, 183, 100)) is: light)
		Assert(g(RGB(0, 231, 0)) is: dark) 		//bright green
		Assert(g(RGB(0, 229, 0)) is: light)
		Assert(g('0xe500') is: light) 			// test hex val
		Assert(g(RGB(1, 1, 1)) is: light)
		Assert(g(RGB(0, 0, 0)) is: light)			//black

		//ensure passed in dark/light colors are returned
		dark = 0x303030
		light = 0xbbbbbb
		Assert(g(RGB(255, 255, 255), dark, light) is: dark)	//white
		Assert(g(RGB(100, 184, 100), dark, light) is: dark)
		Assert(g(RGB(100, 183, 100), dark, light) is: light)
		Assert(g(RGB(0, 231, 0), dark, light) is: dark) 		//bright green
		Assert(g(RGB(0, 229, 0), dark, light) is: light)
		Assert(g(RGB(1, 1, 1), dark, light) is: light)
		Assert(g(RGB(0, 0, 0), dark, light) is: light)			//black
		}
	}