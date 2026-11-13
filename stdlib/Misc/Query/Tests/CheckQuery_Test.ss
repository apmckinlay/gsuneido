// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(FormatQuery("columns project column")
			has: "NOT UNIQUE")
		Assert(CheckQuery("columns project column
			/*CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/")
			is: "")
		Assert(CheckQuery("columns project table, column")
			is: "")

		Assert(FormatQuery("tables union columns")
			has: "NOT DISJOINT")
		Assert(CheckQuery("tables union columns
			/*CHECKQUERY SUPPRESS: UNION NOT DISJOINT*/")
			is: "")
		Assert(CheckQuery("tables extend x=1 union tables")
			is: "")

		query = "(views extend x=123) join (tables extend x=123)"
		Assert(FormatQuery(query)
			has: "MANY TO MANY")
		Assert(CheckQuery(query $ " /*CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/")
			is: "")
		Assert(CheckQuery("tables join tables")
			is: "")
		}
	}