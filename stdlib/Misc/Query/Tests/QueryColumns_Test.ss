// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		tbl = .MakeTable("(a,b) key(a)")
		Database("alter " $ tbl $ " create (c)")
		Assert(QueryColumns(tbl) equalsSet: #(a,b,c))
		Database("alter " $ tbl $ " create (d)")
		Assert(QueryColumns(tbl) equalsSet: #(a,b,c)) // cached so not changed
		QueryColumns.ResetCache()
		Assert(QueryColumns(tbl) equalsSet: #(a,b,c,d))
		}

	Test_expired()
		{
		Suneido.Delete('Memoize_QueryColumns_Test.Test_expired queryColumnsExpired')
		queryColumnsExpired = QueryColumns
			{
			ExpirySeconds: -3600
			}

		tbl = .MakeTable("(a,b,c) key(a)")
		queryColumnsExpired.ResetCache()
		Assert(queryColumnsExpired(tbl) equalsSet: #(a,b,c))
		Database("alter " $ tbl $ " create (d)")
		Assert(queryColumnsExpired(tbl) equalsSet: #(a,b,c,d))
		queryColumnsExpired.ResetCache()
		}
	}