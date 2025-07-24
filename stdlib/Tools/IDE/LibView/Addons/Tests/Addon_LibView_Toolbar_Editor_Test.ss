// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_test_method()
		{
		m = Addon_LibView_Toolbar_Editor.Addon_LibView_Toolbar_Editor_test_method
		Assert(m([], 0) is: false)
		ranges = Object(Object(from: 0, to: 50, name: "Fred"))
		Assert(m(ranges, 25) is: false)
		ranges = Object(Object(from: 0, to: 50, name: "Test_Fred"))
		Assert(m(ranges, 25) is: 'Test_Fred')
		ranges = Object(Object(from: 0, to: 50, name: "Test_Fred"))
		Assert(m(ranges, 75) is: false)
		}
	}