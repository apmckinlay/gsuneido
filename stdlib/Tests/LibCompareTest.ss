// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.t1 = .MakeTable('(num, parent, group, name, text) key(num)')
		.t2 = .MakeTable('(num, parent, group, name, text) key(num)')
		.t3 = .MakeTable('(num, parent, group, name, text) key(num)')
		.t4 = .MakeTable('(num, parent, group, name, text) key(num)')

		for (i = 0; i < 100; ++i)
			{
			QueryOutput(.t1,
				Object(num: i, parent: 0, group: -1, name: i, text: "record" $ i))
			if i > 2
				{
				QueryOutput(.t2,
					Object(num: i, parent: 0, group: -1, name: i, text: "record" $ i))
				QueryOutput(.t3,
					Object(num: i, parent: 0, group: -1, name: i, text: "record" $ i))
				}
			else
				{
				QueryOutput(.t3,
					Object(num: i, parent: 0, group: -1, name: i, text: "record" $ i $ i))
				QueryOutput(.t4,
					Object(num: i, parent: 0, group: -1, name: i + 1000,
						text: "record" $ i $ i))
				}
			if i > 4 and i < 98
				QueryOutput(.t4,
					Object(num: i, parent: 0, group: -1, name: i, text: "record" $ i))
			if i >=98
				QueryOutput(.t4,
					Object(num: i, parent: 0, group: -1, name: i, text: "record" $ i $ i))
			}
		}
	Test_same()
		{
		result = LibCompare(.t1, .t1)
		Assert(result is: Object())  //no differences
		}
	Test_removed()
		{
		result = LibCompare(.t1, .t2)
		Assert(result is: Object(
			#("-", 0, 0)
			#("-", 1, 1)
			#("-", 2, 2)))
		}
	Test_added()
		{
		result = LibCompare(.t2, .t1)
		Assert(result is: Object(
			#("+", 0, 0)
			#("+", 1, 1)
			#("+", 2, 2)))
		}
	Test_changed()
		{
		result = LibCompare(.t1, .t3)
		Assert(result is: Object(
			#("#", 0, 0)
			#("#", 1, 1)
			#("#", 2, 2)))
		}
	Test_mixed()
		{
		result = LibCompare(.t1, .t4)
		Assert(result is: Object(
			#("-", 0, 0)
			#("-", 1, 1)
			#("-", 2, 2)
			#("-", 3, 3)
			#("-", 4, 4)
			#("#", 98, 98)
			#("#", 99, 99)
			#("+", 0, 1000)
			#("+", 1, 1001)
			#("+", 2, 1002)))
		}
	}
