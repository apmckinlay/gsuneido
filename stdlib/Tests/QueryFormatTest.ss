// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_accum_totals1()
		{
		mock = Mock()
		mock.QueryFormat_totalFields = #('string', 'number', 'boolean', 'date')
		mock.QueryFormat_totals = Object().Set_default(Object().Set_default(0))
		rec1 = [string: '3', number: 1, boolean: false, date: '']
		rec2 = [string: '', number: '', boolean: '', date: #20170101]
		rec3 = [string: 'hello', number: 2, boolean: true, date: #20170201]
		rec4 = [string: '2', number: 3, boolean: true, date: #20170301]
		records = Object(rec1, rec2, rec3, rec4)
		for mem in records.Members()
			mock.Eval(QueryFormat.QueryFormat_accum_totals1, records[mem], 0)
		rec1.string = 1
		for mem in records.Members()
			mock.Eval(QueryFormat.QueryFormat_accum_totals1, records[mem], 999)
		Assert(mock.QueryFormat_totals[0] is: #(5, 6))
		Assert(mock.QueryFormat_totals[999] is: #(3, 6))
		}
	Test_summarizeMsgWarning()
		{
		func = QueryFormat.QueryFormat_getSummarizeMsgWarning
		Assert(func(true) is: 'Unable to preform action.\r\nToo many records to ' $
			'summarize.  Please use the "Select..." button to reduce the number of ' $
			'records')
		Assert(func(false) is: 'Unable to preform action.\r\nToo many records to ' $
			'summarize.  Please use filters to reduce the number of ' $
			'records')
		}
	}
