// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(WikiTitleValid?('') is: `ERROR: invalid page name: `)
		Assert(WikiTitleValid?(123) is: `ERROR: invalid page name: 123`)
		Assert(WikiTitleValid?('123Title') is: `ERROR: invalid page name: 123Title`)
		Assert(WikiTitleValid?('Title123') is: `ERROR: invalid page name: Title123`)
		Assert(WikiTitleValid?('title') is: `ERROR: invalid page name: title`)
		Assert(WikiTitleValid?('Title') is: `ERROR: invalid page name: Title`)
		Assert(WikiTitleValid?('titleTitle') is: `ERROR: invalid page name: titleTitle`)
		Assert(WikiTitleValid?('Title123Title') is: ``)
		}
	}