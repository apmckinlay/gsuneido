// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_output()
		{
		Assert({ QueryOutput('tables', []) } throws:)
		}

	Test_output_query()
		{
		.test_query('insert { } into tables')
		}

	Test_update_query()
		{
		.test_query("update columns set table = 0")
		}

	Test_delete_query()
		{
		.test_query("delete columns")
		}

	test_query(query)
		{
		Assert({ QueryDo(query) } throws:)
		}

	Test_update()
		{
		.test({|x| x.Update() })
		}

	Test_delete()
		{
		.test({|x| x.Delete() })
		}

	test(block)
		{
		Transaction(update:)
			{|t|
			x = t.QueryFirst('columns sort column')
			Assert({ block(x) } throws:)
			}
		}
	}
