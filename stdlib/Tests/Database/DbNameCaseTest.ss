// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// test that mixed case names work
Test
	{
	Test_main()
		{
		schema = "(aName, aPhone) key (aName) index (aPhone)"
		table = .MakeTable(schema)
		Assert(Schema(table).Tr('\n') like: table $ ' ' $ schema)
		Database("drop " $ table)
		Database("ensure " $ table $ ' ' $ schema)
		Assert(Schema(table).Tr('\n') like: table $ ' ' $ schema)
		}
	}
