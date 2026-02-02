// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(QueryStripSort("tables") is: "tables")
		Assert(QueryStripSort("tables sort totalsize") is: "tables")
		Assert(QueryStripSort("tables where table is 123 sort totalsize")
			is: "tables where table is 123")
		Assert(QueryStripSort("tables extend sort_field = 1")
			is: "tables extend sort_field = 1")
		Assert(QueryStripSort(
			"/* tableHint: tables */ tables join by(table) columns sort column")
			is: "/* tableHint: tables */ tables join by(table) columns")
		Assert(QueryStripSort("/*CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/tables " $
			"join by (table) columns /*CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/ " $
			"sort column")
			is: "/*CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/" $
				"tables join by (table) columns " $
				"/*CHECKQUERY SUPPRESS: JOIN MANY TO MANY*/" )
		}
	}