// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
// Note: 'is' can be evaluated raw/packed, =~ must be unpacked
Test
	{
	Test_stored_rule()
		{
		// Note: this behavior seems wrong
		// but it is the way suneido has worked for a long time
		tbl = .MakeTable("(k, test_simple_default) key(k)", #(k: 1)) // object not record
		Assert(QueryEmpty?(tbl, test_simple_default: 'simple default'), msg: 'one')
		Assert(QueryEmpty?(tbl $ " where test_simple_default =~ 'simple default'"),
			msg: 'two')
		Assert(QueryEmpty?(tbl $ " where test_simple_default in ('simple default')"),
			msg: 'three')
		Assert(Query1(tbl $ " summarize min test_simple_default"), is: [])
		}
	Test_extended_rule()
		{
		tbl = .MakeTable("(k) key(k)", #(k: 1))
		q = tbl $ ' extend test_simple_default'
		Assert(Query1(q $ " where test_simple_default is 'simple default'") is: [k: 1])
		Assert(Query1(q $ " where test_simple_default =~ 'simple default'") is: [k: 1])
		Assert(Query1(q $ " where test_simple_default in ('simple default')") is: [k: 1])
		Assert(Query1(q $ " summarize min test_simple_default")
			is: [min_test_simple_default: "simple default"])
		}
	}
