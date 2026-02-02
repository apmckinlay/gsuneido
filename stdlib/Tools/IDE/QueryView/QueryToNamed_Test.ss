// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		test = function(query)
			{
			return {|@expected|
				if expected is [false]
					expected = [query]
				Assert(QueryToNamed(query) is: expected)
				}
			}
		test("stdlib")("stdlib")
		test("stdlib where name = 'Max'")("stdlib", name: "Max")
		test("stdlib where 123 = num")("stdlib", num: 123)
		test("stdlib where name is 'Max' and num = 123")(
			"stdlib", name: "Max", num: 123)

		test("tables join columns")(false)
		test("tables project table")(false)
		test("tables where nrows = totalsize")(false)
		}
	}