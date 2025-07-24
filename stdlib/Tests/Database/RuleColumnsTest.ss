// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		suppress = " /*CHECKQUERY SUPPRESS: UNION NOT DISJOINT*/"
		a = .MakeTable("(a,test_simple_default) key(a)")
		b = .MakeTable("(a) key(a)")

		query = a $ " union " $ a $ suppress
		Assert(QueryRuleColumns(query) is: #())

		query = a $ " union (" $ b $ " extend test_simple_default)" $ suppress
		Assert(QueryRuleColumns(query) is: #())

		query = "(" $ b $ " extend test_simple_default)
			union (" $ b $ " extend test_simple_default)" $ suppress
		Assert(QueryRuleColumns(query) is: #(test_simple_default))
		}
	}