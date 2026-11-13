// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// Test TableExists? and basic call through
		table = .MakeTable("(num) key(num)")
		Assert(TableAuthorized?(table))
		Assert(TableAuthorized?(.TempTableName()) is: false)

		// Test unauthorized table
		mock = Mock(TableAuthorized?)
		mock.When.CallClass([anyArgs:]).CallThrough()
		mock.When.query(table).Throw("not authorized (from server)")
		Assert(mock.CallClass(table) is: false)
		mock.When.query(table).Throw("unexpected error")
		Assert({ mock.CallClass(table) } throws: "unexpected error")
		}
	}
