// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		for case in .cases
			Assert(LibRecordType(case[0]) is: case[1])
		}
	cases:
		(
		('dll', 'dll')
		('class {', 'class')
		('class : Base {', 'class')
		('class Base {', 'class')
		('Base {', 'class')
		('_Base {', 'class')
		('123', 'number')
		('-123', 'number')
		('"hello"', 'string')
		('function', 'function')
		('#()', 'object')
		('#{}', 'record')
		('[]', 'record')
		)
	}