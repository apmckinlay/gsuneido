// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(a, b) key(a)', Record(a: 1, b: 2))
		Assert({ RetryTransaction()
				{ |t/*unused*/|
				throw 'transaction conflict'
				} },
			throws: "RetryTransaction: too many retries, " $
				"last error: transaction conflict")
		i = 0
		RetryTransaction()
			{ |t|
			rec = t.Query1(table)
			if i++ is 0
				throw 'transaction conflict'
			rec.Update()
			}
		Assert(i is: 2)
		}

	Test_isTransactionConflict?()
		{
		f = RetryTransaction.IsTransactionConflict?
		Assert(f("Transaction: block commit failed"))
		Assert(f("Transaction: preempted by exclusive"))
		Assert(f("Transaction: conflict with exclusive (from server)"))
		Assert(f('Not Conflict') is: false)
		Assert(f('some error occurred') is: false)
		Assert(f('boolean does not support get') is: false)
		}
	}