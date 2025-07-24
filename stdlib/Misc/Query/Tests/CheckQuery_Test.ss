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

		Assert(FormatQuery("tables rename table to t join tables")
			has: "MANY TO MANY")
		Assert(CheckQuery("tables rename table to t join tables
			/*CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/")
			is: "")
		Assert(CheckQuery("tables join tables")
			is: "")
		}
	}