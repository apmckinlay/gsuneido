// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		table = .MakeTable('(a, b, a_lower!) key (a) key (a_lower!)')
		for i in 'A'.Asc() .. 'Z'.Asc() + 1
			{
			QueryOutput(table, [a: i.Chr(), b: i.Odd?()])
			}
		result = Object()
		QueryApply(table $ ' minus (' $ table $ ' where a is "R")
			where a_lower! >= "r" sort a_lower!')
			{
			result.Add(it.a)
			}
		Assert(result is: #(S, T, U, V, W, X, Y, Z))

		result = Object()
		QueryApply(table $ ' intersect (' $ table $ ' where a is "R")
			where a_lower! >= "r" sort a_lower!')
			{
			result.Add(it.a)
			}
		Assert(result is: #(R))
		}
	}