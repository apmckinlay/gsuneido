// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(QueryAddMore("tables", "")
			is: "tables")
		Assert(QueryAddMore("tables", "join columns")
			is: "tables\njoin columns")
		Assert(QueryAddMore("tables sort nrows", "join columns")
			is: "tables\njoin columns\nsort nrows")
		Assert(QueryAddMore("/* tableHint: tables */ tables sort nrows", "join columns")
			is: "/* tableHint: tables */ tables\njoin columns\nsort nrows")
		Assert(QueryAddMore("tables /* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/ " $
			"/* CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/ sort nrows", "join columns")
			is: "tables /* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/ " $
				"/* CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/\n" $
				"join columns\nsort nrows")
		}
	}