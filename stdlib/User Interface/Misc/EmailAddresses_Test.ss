// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.TearDownIfTablesNotExist('email_addresses')
		EmailAddresses.Ensure()
		ea = new EmailAddresses
			{
			EmailAddresses_add1(word, email, t/*unused*/)
				{ .List.Add([word, email]) }
			}
		ea.List = []
		ea.OutputAddr('a@b.c')
		Assert(ea.List is: [['a@b.c', 'a@b.c']])

		ea.List = []
		ea.OutputAddr('joe smith <a@b.c>')
		Assert(ea.List is: [
			['joe smith', 'joe smith <a@b.c>'],
			['smith', 'joe smith <a@b.c>'],
			['a@b.c', 'joe smith <a@b.c>']])
		}

	Test_stripInternalDesc()
		{
		m = EmailAddresses.EmailAddresses_stripInternalDesc

		Assert(m('') is: '')
		Assert(m('test@test.com') is: 'test@test.com')
		Assert(m('<test@test.com>') is: '<test@test.com>')
		Assert(m('(abbrev)<test@test.com>') is: '(abbrev)<test@test.com>')
		Assert(m('Sample Company (abbrev)<test@test.com>') is:
			'Sample Company (abbrev)<test@test.com>')
		Assert(m('Airways Ltd.*Center(1) (airwaycenter) <info@testemail.com>') is:
			'Airways Ltd. (airwaycenter) <info@testemail.com>')
		Assert(m('Airways Ltd.*Center (airwaycenter) <info@testemail.com>') is:
			'Airways Ltd. (airwaycenter) <info@testemail.com>')
		Assert(m('Airways Ltd.*Center (Test <info@testemail.com>') is:
			'Airways Ltd. (Test <info@testemail.com>')
		}

	Test_splitMultipleAddresses()
		{
		m = EmailAddresses.EmailAddresses_splitMultipleAddresses

		list = Object("")
		expected = Object("")
		Assert(m(list) is: expected)

		list = Object("name<first@test.com;second@test.com>")
		expected = Object("name<first@test.com>", "name<second@test.com>")
		Assert(m(list) is: expected)
		list = Object("name<first@test.com;second@test.com,third@test.com>")
		expected = Object("name<first@test.com>", "name<second@test.com>",
			"name<third@test.com>")
		Assert(m(list) is: expected)
		list = Object("name<first@test.com>")
		expected = Object("name<first@test.com>")
		Assert(m(list) is: expected)
		list = Object("first@test.com")
		expected = Object("first@test.com")
		Assert(m(list) is: expected)
		list = Object("<first@test.com;second@test.com>")
		expected = Object("<first@test.com>", "<second@test.com>")
		Assert(m(list) is: expected)

		list = Object("test@a.com", "name<first@test.com;second@test.com>")
		expected = Object("test@a.com", "name<first@test.com>", "name<second@test.com>")
		Assert(m(list) is: expected)

		list = Object("name<first@test.com;second@test.com>", "test@a.com")
		expected = Object("name<first@test.com>", "name<second@test.com>","test@a.com")
		Assert(m(list) is: expected)

		list = Object("a@b.com", "name<first@test.com;second@test.com>", "test@a.com")
		expected = Object("a@b.com", "name<first@test.com>", "name<second@test.com>",
			"test@a.com")
		Assert(m(list) is: expected)

		list = Object("last, name<first@test.com;second@test.com>")
		expected = Object("last, name<first@test.com>", "last, name<second@test.com>")
		Assert(m(list) is: expected)

		list = Object("name*test<first@test.com;second@test.com>")
		expected = Object("name*test<first@test.com>", "name*test<second@test.com>")
		Assert(m(list) is: expected)

		list = Object("name<first@test.com;second@test.com>","a@b.com",
			"name<third@test.com;fourth@test.com>", "test@a.com")
		expected = Object("name<first@test.com>", "name<second@test.com>", "a@b.com",
			"name<third@test.com>", "name<fourth@test.com>", "test@a.com")
		Assert(m(list) is: expected)

		// angle brackets in description
		list = Object("Name < Test <stopokay@test.com,gonnabreak@okay>")
		expected = Object("Name < Test <stopokay@test.com>",
			"Name < Test <gonnabreak@okay>")
		Assert(m(list) is: expected)

		list = Object("Name << Test <stopokay@test.com,gonnabreak@okay>")
		expected = Object("Name << Test <stopokay@test.com>",
			"Name << Test <gonnabreak@okay>")
		Assert(m(list) is: expected)

		list = Object("Name > Test <first@test.com>")
		expected = Object("Name > Test <first@test.com>")
		Assert(m(list) is: expected)

		list = Object("Name <A> Test <first@test.com>")
		expected = Object("Name <A> Test <first@test.com>")
		Assert(m(list) is: expected)

		// long description
		description = 'a'.Repeat(110)
		list = Object(description $ "<tester@test.com,test2@test.com>")
		expected = Object(description $ "<tester@test.com>",
			description $ "<test2@test.com>")
		Assert(m(list) is: expected)

		// long email addresses
		address1 = 'a'.Repeat(100) $ '@test.com'
		list = Object('Description ' $ "<" $ address1 $ ", test@test.com>")
		expected = Object('Description ' $ "<" $ address1 $ ">",
			'Description <test@test.com>')
		Assert(m(list) is: expected)

		// long description and long addresses
		description = 'a'.Repeat(110)
		address1 = 'a'.Repeat(100) $ '@test.com'
		list = Object(description $ "<" $ address1 $ ",test2@test.com>")
		expected = Object(description $ "<" $ address1 $ ">",
			description $ "<test2@test.com>")
		Assert(m(list) is: expected)
		}

	Test_preparePrefix()
		{
		fn = EmailAddresses.EmailAddresses_preparePrefix

		prefix = ''
		Assert(fn(prefix) is: '')

		prefix = 'abc'
		Assert(fn(prefix) is: 'abc')

		prefix = 'ABC'
		Assert(fn(prefix) is: 'abc')

		prefix = 'ABCdefGHIjkl'
		Assert(fn(prefix) is: 'abcdefghijkl')

		// longer than prefix limit, should truncate at limit
		prefix = 'A'.Repeat(600)
		Assert(fn(prefix) is: 'a'.Repeat(500))

		// test invalid prefix data type
		prefix = 123
		Assert({ fn(prefix) }
			throws: 'unexpected data type for EmailAddresses.preparePrefix')
		}
	}
