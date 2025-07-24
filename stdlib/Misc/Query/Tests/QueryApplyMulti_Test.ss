// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_cursorCost()
		{
		cursorCost = QueryApplyMulti.QueryApplyMulti_cursorCost
		Assert(cursorCost('stdlib sort name,group').query is: 'stdlib sort name,group')
		Assert(cursorCost('stdlib').query is: 'stdlib sort num')
		Assert(cursorCost('stdlib where name = "Test"').query
			is: 'stdlib where name = "Test" sort name,group')
		tbl = .MakeTable("(a) key()")
		Assert(cursorCost(tbl).query is: tbl)
		}
	Test_which()
		{
		cursorCost = QueryApplyMulti.QueryApplyMulti_cursorCost
		lookupCost = QueryApplyMulti.QueryApplyMulti_lookupCost
		q = "stdlib"
		Assert(cursorCost(q).cost lessThan: lookupCost(q).cost)
		q = "stdlib where parent is 123"
		Assert(cursorCost(q).cost greaterThan: lookupCost(q).cost)
		}
	Test_composite_key()
		{
		tbl = .MakeTable('(a,b) key(a,b)', [a: 1, b: 2], [a: 3, b: 4], [a: 3], [b: 2])
		withLookup = QueryApplyMulti.QueryApplyMulti_withLookup
		withLookup(tbl, {|unused| })
		}
	}