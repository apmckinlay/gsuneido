// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(QueryGetSuppressions("tables") is: #())
		Assert(QueryGetSuppressions("/*CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/tables " $
			"join by (table) columns /*CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/ " $
			"sort column").Copy().Sort!()
			is: #("/* CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/",
				"/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/").Copy().Sort!())
		}
	}