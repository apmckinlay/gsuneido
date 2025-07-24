// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		table = .MakeTable('(num, name, abbrev, type, type2) key(num)',
			[num: 1, name: 'One', abbrev: 'one', type: 'a', type2: 'a'],
			[num: 2, name: 'Two', abbrev: 'two', type: 'b', type2: 'b'],
			[num: 3, name: 'Three', abbrev: 'three', type: 'a', type2: 'a'],
			[num: 4, name: 'Four', abbrev: 'four', type: 'b', type2: 'b'],
			[num: 5, name: 'Five', abbrev: 'five', type: 'a', type2: 'a'],
			[num: 6, name: 'Six', abbrev: 'six', type: 'b', type2: 'b'],
			)

		rec = QueryAny1(table $ ' where type is "c"', 'type2')
		Assert(rec is: false)

		rec = QueryAny1(table $ ' where name is "One"', 'type', 'type2')
		Assert(rec.type is: 'a')
		Assert(rec.type2 is: 'a')

		rec = QueryAny1(table $ ' where type is "a"', 'type2')
		Assert(rec.type2 is: 'a')

		rec = QueryAny1(table $ ' where type is "b"', 'type2')
		Assert(rec.type2 is: 'b')

		orig = Suneido.GetDefault('ValidateQueryAny1?', false)
		.AddTeardown({ Suneido.ValidateQueryAny1? = orig })
		Suneido.ValidateQueryAny1? = true
		err = ''
		try QueryAny1(table, 'type2')
		catch (err) { }
		Assert(err is: "QueryAny1: all recs do not match")
		}
	}