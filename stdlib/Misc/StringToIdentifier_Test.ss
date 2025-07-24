// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(StringToIdentifier(1234) is: '1234')
		Assert(StringToIdentifier('1234') is: '1234')
		Assert(StringToIdentifier('string to identifier test')
			is: 'string20to20identifier20test')
		Assert(StringToIdentifier('string_to_identifier_test')
			is: 'string_to_identifier_test')
		Assert(StringToIdentifier("string's test") is: 'string27s20test')
		Assert(StringToIdentifier("string #test2") is: 'string2023test2')
		Assert(StringToIdentifier("String Test2") is: 'String20Test2')
		Assert(StringToIdentifier("string test2") is: 'string20test2')
		}
	}