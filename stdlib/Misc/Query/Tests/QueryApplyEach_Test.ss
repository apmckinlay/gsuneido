// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Suneido.forceRetryTransaction = true
		table = .MakeTable('(num, value) key(num, value)')
		for (i = 0; i < 100; ++i)
			QueryOutput(table, Record(num: i, value: String(i).Repeat(50)))

		n = 0
		QueryApplyEach(table $ ' where value =~ "40|70"')
			{ |t /*unused*/, x|
			++n
			Assert(x.value has: '40') // because we update 70 below

			x.value = 12
			x.Update()
			count = QueryDo('update ' $ table $ ' where num is 2
				set value = "x" $ value')
			Assert(count is: 1)

			QueryDo('update ' $ table $ ' where num is 70 set value = ""')
			}
		Assert(n is: 2) // one record but also one retry
		}
	Teardown()
		{
		Suneido.Delete("forceRetryTransaction")
		super.Teardown()
		}
	}