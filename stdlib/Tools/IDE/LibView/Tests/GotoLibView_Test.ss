// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_split()
		{
		split = GotoLibView.GotoLibView_split
		Assert(split('LibLocateControl') is: #(name: LibLocateControl, method: ''))
		Assert(split('Config.Invalidate') is: #(name: Config, method: Invalidate))
		Assert(split('abc.css') is: #(name: "abc.css", method: ''))
		Assert(split('abc.js') is: #(name: "abc.js", method: ''))
		Assert(split('12345') is: false)
		}
	}