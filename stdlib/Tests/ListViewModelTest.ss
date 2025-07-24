// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.table = .MakeTable('(name, age) key(name)')
		for (x in .data)
			QueryOutput(.table, x)
		.test(ListViewModel(.table $ " where false is name.Has?('x')"))
		.test(ListViewModelCached(ListViewModel(.table $
			" where false is name.Has?('x')")))
		}

	test(vlv)
		{
		.testReadingPrevious(vlv)

		Assert(vlv.Getnumrows() is: 24)
		vlv.Close()
		vlv = ListViewModel(.table $ " where false is name.Has?('x')")

		Assert(vlv.Getnumrows() is: 24)

		.loopThroughListTwice(vlv)

		.getRandomLinesFromList(vlv)

		Assert(vlv.Getnumrows() is: 24)
		vlv.Close()
		}

	// test reading previous (and the reset of the size)
	testReadingPrevious(vlv)
		{
		for (i = 0; i < 3; ++i)
			{
			Assert(numrows = vlv.Getnumrows() is: 24)
			for (j = 0; j < 25; ++j)
				{
				--numrows
				if (j < 24)
					Assert(vlv.Getrecord(numrows) is: .data[23 - j])
				else
					Assert(vlv.Getrecord(numrows) is: Object())
				}
			}
		}

	loopThroughListTwice(vlv)
		{
		for (i = 0; i < 24; ++i)
			Assert(vlv.Getitem(i, "name") is: .data[i].name)
		Assert(vlv.Getitem(i, "name") is: "", msg: "2 loop test")
		for (i = 0; i < 24; ++i)
			Assert(vlv.Getitem(i, "name") is: .data[i].name)
		}

	getRandomLinesFromList(vlv)
		{
		// for loop that generated random numbers and calls GetRecord with it...
		for (i = 0; i < 100; ++i)
			{
			// determine the field to get
			field_name = (i % 2 is 0) ? "name" : "age"
			rand = Random(50)
			result = vlv.Getrecord(rand)
			item_result = vlv.Getitem(rand, field_name)
			if (rand <= 23 and rand >= 0)
				{
				Assert(result is: .data[rand])
				Assert(item_result is: .data[rand][field_name])
				}
			else if (rand > 23)
				{
				Assert(result is: Object())
				Assert(item_result is: "")
				}
			}
		}

	data:
		(
		(name: "aaron", age: 27)
		(name: "anthony", age: 28)
		(name: "bill", age: 29)
		(name: "bruce", age: 30)
		(name: "cam", age: 32)
		(name: "carter", age: 31)
		(name: "craig", age: 33)
		(name: "fred", age: 23)
		(name: "geoff", age: 24)
		(name: "george", age: 23)
		(name: "grant", age: 26)
		(name: "greg", age: 25)
		(name: "jared", age: 19)
		(name: "jason", age: 100)
		(name: "jeff", age: 22)
		(name: "jim", age: 21)
		(name: "joe", age: 35)
		(name: "john", age: 50)
		(name: "josh", age: 20)
		(name: "kent", age: 34)
		(name: "kevin", age: 36)
		(name: "kirk", age: 35)
		(name: "pete", age: 67)
		(name: "steve", age: 18)
		(name: "xavier", age: 24)
		(name: "xavier2", age: 25)
		)
	}
