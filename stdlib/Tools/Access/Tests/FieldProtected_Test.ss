// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(FieldProtected?('test_field', [], 'test_protect') is: false)
		Assert(FieldProtected?('test_field',
			[test_field__protect: true], 'test_protect'))
		Assert(FieldProtected?('test_field',
			[test_field__protect: false],'test_protect') is: false)
		Assert(FieldProtected?('test_field', []) is: false)
		Assert(FieldProtected?('test_field',
			[test_protect: Object(test_field:)],'test_protect'))
		Assert(FieldProtected?('test_field',
			[test_protect: Object('allbut' ,test_field:)], 'test_protect') is: false)
		Assert(FieldProtected?('test_field',
			[test_protect: 'protected'],'test_protect'))
		}
	}