// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.oldBook = Suneido.GetDefault(#CurrentBook, false)
		}

	Test_main()
		{
		Suneido.CurrentBook = 'fakeBook'
		Assert(FindCurrentBook() is: 'fakeBook')

		Suneido.Delete(#CurrentBook)
		.SpyOn(AccessPermissions.GetDefaultBook).Return('defaultBook')
		Assert(FindCurrentBook() is: 'defaultBook')
		}

	Teardown()
		{
		if .oldBook isnt false
			Suneido.CurrentBook = .oldBook
		}
	}
