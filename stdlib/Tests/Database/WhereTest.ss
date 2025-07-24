// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_constant_folding()
		{
		test = function(query, expected)
			{
			Assert(QueryStrategy(query)
				is: expected)
			}
		test('tables where false or false', "nothing")
		test('tables extend x = 5 where x = 6', "nothing")
		test('tables where table = 5 and false', "nothing")
		test('tables where true and true', "tables where true")
		test('tables where table = 5 or true', "tables where true")
		}
	}