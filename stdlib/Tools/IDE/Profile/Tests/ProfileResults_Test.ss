// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildDataValues()
		{
		results = Object()
		data = ProfileResults.BuildDataValues(results, 1)
		Assert(data.Empty?(), msg: "no data")

		results = Object(
			[name: "one", self: 2, total: 4, calls: 3],
			[name: "odd", self: 5, total: 10, calls: 65])
		data = ProfileResults.BuildDataValues(results, reps: 3)
		sbe = #(
			one: (self: 20, total: 40, calls: 1),
			odd: (self: 50, total: 100, calls: 22))
		.validateData(data, sbe, 'BuildDataValues')
		}
	validateData(data, sbe, desc)
		{
		for ob in data
			{
			mem = ob.profile_name
			Assert(sbe.Member?(mem), msg: desc $ ' missing name: ' $ mem)
			for i in sbe[mem].Members()
				Assert(ob['profile_' $ i] is: sbe[mem][i],
					msg: desc $ ' ' $ i $ ' ' $ mem)
			}
		}
	}