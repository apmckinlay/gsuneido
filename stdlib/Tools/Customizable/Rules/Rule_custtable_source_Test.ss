// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rec = []
		Assert(rec.custtable_source is '')

		rec.params = ''
		Assert(rec.custtable_source is '')

		rec.params = 'random string'
		Assert(rec.custtable_source is '')

		rec.params = #()
		Assert(rec.custtable_source is '')

		rec = [params: #(m: 'a value')]
		Assert(rec.custtable_source is '')

		rec = [params: #(m: 'a value', Source: 'test source')]
		Assert(rec.custtable_source is 'test source')
		}
	}
