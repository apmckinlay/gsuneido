// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.origBookOption = false
		if false isnt option = Suneido.GetDefault("CurrentBookOption", false)
			.origBookOption = option
		}

	Test_main()
		{
		Suneido.CurrentBookOption = ''
		Assert(CurrentModule() is: '')

		Suneido.CurrentBookOption = '/TestModule/Testing Stuff'
		Assert(CurrentModule() is: 'TestModule')

		Suneido.CurrentBookOption = '/Test Test/Testing'
		Assert(CurrentModule() is: 'Test Test')

		Suneido.CurrentBookOption = '/TestApp/SubMenu/Test Page'
		Assert(CurrentModule() is: 'TestApp')
		}

	Teardown()
		{
		if .origBookOption is false
			Suneido.Delete("CurrentBookOption")
		else
			Suneido.CurrentBookOption = .origBookOption
		}
	}