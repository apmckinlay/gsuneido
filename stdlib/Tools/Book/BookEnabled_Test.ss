// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		bookEnabled = BookEnabled
			{
			BookEnabled_options(book /*unused*/)
				{ return Object("-/test/disabled", "+/test", "--/hidden option",
					'+/book', '--/book/duplicate menu', '--/hiddenSection',
					'--Another Hidden Option') }
			}

		.olduser = Suneido.User
		.oldroles = Suneido.user_roles
		Suneido.User = 'test'
		Suneido.user_roles = #('test')

		// not using "real" options; not overriding
		Assert(BookEnabled.Enabled('Book', '/test') is: 'not found')
		Assert(BookEnabled('Book', '/test') is: false)

		Suneido.User = 'default'
		Assert(BookEnabled.Enabled('Book', '/test') is: 'not found')
		Assert(BookEnabled('Book', '/test') is: false)

		// --options should always return 'hidden' regardless of user
		// since option is often hidden because library that defines the option
		// is not being used.
		Suneido.User = .olduser
		Assert(bookEnabled.Enabled('Book', '/hidden option') is: 'hidden')
		Assert(bookEnabled('Book', '/hidden option') is: 'hidden')

		Suneido.User = 'default'
		Assert(bookEnabled.Enabled('Book', '/hidden option') is: 'hidden')
		Assert(bookEnabled('Book', '/hidden option') is: 'hidden')
		Assert(bookEnabled('Book', '/test'), msg: '/test is not enabled')
		Assert(bookEnabled.Enabled('Book', '/test'), msg: '/test is not enabled')
		Assert(bookEnabled('Book', '/test/disabled')) // disabled still true for default
		Assert(bookEnabled.Enabled('Book', '/test/disabled') is: false)

		Suneido.User = 'test'
		Assert(bookEnabled('Book', '/test'), msg: '/test is not enabled')
		Assert(bookEnabled.Enabled('Book', '/test'), msg: '/test is not enabled')
		Assert(bookEnabled('Book', '/test/disabled') is: false)
		Assert(bookEnabled.Enabled('Book', '/test/disabled') is: false)

		Assert(bookEnabled('Book', '/hiddenSection/screen') is: 'hidden'
			msg: 'hiddenSection screen not hidden')

		Assert(bookEnabled('Book', '/book/duplicate menu') is: 'hidden'
			msg: 'duplicate menu not hidden')
		Assert(bookEnabled('Book', '/book/duplicate menu name'),
			msg: 'duplicate menu name is hidden')

		Assert(bookEnabled('Book', '/book/Another Hidden Option') is: 'hidden'
			msg: 'option with no leading "/" not hidden')
		}
	olduser: false
	oldroles: false
	Teardown()
		{
		if .olduser isnt false
			Suneido.User = .olduser
		if .oldroles isnt false
			Suneido.user_roles = .oldroles
		}
	}
