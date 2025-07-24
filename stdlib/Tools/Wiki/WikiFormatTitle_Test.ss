// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		test = function(s, expected)
			{
			Assert(WikiFormatTitle(s) is: expected)
			}
		test("Foo", "Foo")
		test("CamelCase", "Camel Case")
		test("April25Meeting2016", "April 25 Meeting 2016")
		test("EDIStuff", "EDI Stuff")
		test("AboutEDI", "About EDI")
		test("ETA", "ETA")
		}
	}