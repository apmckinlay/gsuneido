// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	range: 100_000
	steps: 20
	Setup()
		{
		.tbl = .MakeTable("(k, a, b) key(k) index(b) index(a)")
		for k in .. .range
			QueryOutput(.tbl, [:k, a: Random(.range), b: Random(.range)])
		}

	Test_two_indexes()
		{
		for (a = 0; a <= .range; a += .range / .steps)
			for (b = 0; b <= .range; b += .range / .steps)
				{
				query = .tbl $ " where a > " $ a $ " and b > " $ b
				strat = QueryStrategy(query)
				//if strat.Has?("+")
				Print(strat)
				}
		}

	XTest_with_sort()
		{
		for (a = 0; a <= .range; a += .range / .steps)
			{
			query = .tbl $ " where a > " $ a $ " sort b"
			strat = QueryStrategy(query)
			Print(strat)
			}
		}
	}
