// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		fred = [name: 'Fred', age: 23]
		andy = [name: 'Andy', age: 45]
		Assert(By(#age)(fred, andy))
		Assert(By(#age)(andy, fred) is: false)
		Assert(By(#age)(andy, andy) is: false)
		Assert(By(#name)(andy, fred))
		Assert(By(#name)(fred, andy) is: false)

		list = [fred, andy]
		list.Sort!(By(#age))
		Assert(list is: [fred, andy])
		list.Sort!(By(#name))
		Assert(list is: [andy, fred])
		}
	Test_multi()
		{
		x31 = #(a: 3, b: 1)
		x12 = #(a: 1, b: 2)
		x21 = #(a: 2, b: 1)
		x13 = #(a: 1, b: 3)
		list = [x31, x12, x21, x13]
		list.Sort!(By(#a, #b))
		Assert(list is: [x12, x13, x21, x31])
		list.Sort!(By(#b, #a))
		Assert(list is: [x21, x31, x12, x13])
		}
	}