// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Getrecord()
		{
		mock = Mock(ExplorerListModel)
		mock.ExplorerListModel_keyfields = #(date, name)
		mock.ExplorerListModel_query = 'tables sort table'
		t = Mock()
		q = Mock()
		q.When.Next().Return([])
		t.When.Query1([anyArgs:]).Do({|call| query = call.Delete(0); q })
		mock.Eval(ExplorerListModel.Getrecord, t, [date: #20170101, name: 'Fred'])
		Assert(query is: #('tables', date: #20170101, name: 'Fred'))
		}
	}