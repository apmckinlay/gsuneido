// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		spyAlert = .SpyOn('KeyException.KeyException_alert').Return(true)
		table = .MakeTable('(a,b) key(a)',
			.r1 = [a: 1, b: 'one'],
			.r2 = [a: 2, b: 'two'],
			.r3 = [a: 3, b: 'three'])
		fktbl = .MakeTable('(c,a) key(c) index (a) in ' $ table)

		Assert(KeyExceptionTransaction()
			{ |t|
			t.QueryDo('delete ' $ table $ ' where a is 1')
			t.QueryOutput(table, Record(a: 1, b: 'test 1'))
			t.QueryDo('delete ' $ fktbl $ ' where a is 1')
			})
		callLog = spyAlert.CallLogs()
		Assert(callLog.Empty?())

		Assert({ KeyExceptionTransaction()
			{ |t|
			t.QueryOutput(table, Record(a: 2, b: 'test 2'))
			} }
			throws: "duplicate key")
		callLog = spyAlert.CallLogs()
		Assert(callLog isSize: 1)
		Assert(callLog[0].msg has: "Duplicate value in field a")
		}
	}