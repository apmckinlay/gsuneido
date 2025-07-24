// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_NeedImport?()
		{
		mock = Mock(MessageHistoryControl)
		mock.Table = table = .MakeTable('(a) key(a)')
		mock.NumField = 'a'
		mock.When.NeedImport?().CallThrough()
		Assert(mock.NeedImport?() is: false)

		QueryOutput(table, [a: ''])
		Assert(mock.NeedImport?() is: false)

		QueryDo('delete ' $ table)
		QueryOutput(table, [a: Timestamp()])
		Assert(mock.NeedImport?() is: false)

		QueryDo('delete ' $ table)
		QueryOutput(table, [a: Timestamp().Plus(minutes: 1)])
		Assert(mock.NeedImport?() is: false)

		QueryDo('delete ' $ table)
		QueryOutput(table, [a: Timestamp().Minus(minutes: 10)])
		Assert(mock.NeedImport?())

		QueryDo('delete ' $ table)
		QueryOutput(table, [a: Timestamp().Minus(minutes: 11)])
		Assert(mock.NeedImport?())
		}
	}