// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Test
	{
	// SuJsWebTest
	Test_one()
		{
		.check([], [], [])
		.check([[[1], [2]]], [[[1], [2]]], [])
		.check([[[[a: 1]]]], [[['FO_#0#']]], ['FO_#0#': [a: 1]])
		.check([[[[1], [2]]]], [[['FO_#0#', 'FO_#1#']]], ['FO_#0#': [1], 'FO_#1#': [2]])
		.check([#(((((((1)))))))], [[['FO_#0#']]],
			['FO_#0#': [['FO_#1#']], 'FO_#1#': [['FO_#2#']], 'FO_#2#': [1]])

		ob1 = Object()
		.check([[ob1], [[ob1]]], [[ob1], [['FO_#0#']]], ['FO_#0#': ob1])

		ob1 = Object()
		ob2 = Object(ob1)
		ob1.Add(ob2)
		Assert({ FlatObject(ob1, 3) } throws: 'found circular object')
		}

	check(ob, expectedOb, expectedExtra)
		{
		origin = ob.DeepCopy()
		res = FlatObject(ob, 3)
		Assert(res.extra is: expectedExtra)
		Assert(res.ob is: expectedOb)
		Assert(FlatObject.Build(res) is: origin)
		}
	}