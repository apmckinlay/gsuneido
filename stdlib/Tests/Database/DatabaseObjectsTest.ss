// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(name, value) key(name)')
		for i in .data.Members()
			QueryOutput(table, Object(name: i, value: .data[i]))

		x = Object()
		x[0] = x
		try
			{
			QueryOutput(table, Object(name: 99, value: x))
			throw "didn't catch self nested object"
			}

		WithQuery(table)
			{ |q|
			for x in .data
				Assert(q.Next().value is: x)
			}
		}
	data:
		(
		()
		(1, 2, 3)
		(a: 1, b: 2, c: 3)
		(1, 2, 3, a: 4, b: 5, c: 6)
		(1, "hello", (1, "two"), a: 1, b: "hello", c: (1, "two"))
		)
	}