// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(nrecords = 100, nthreads = 10, seconds = 10)
		{
		(new this(:nrecords, :nthreads, :seconds)).Run()
		}
	New(.nrecords = 100, .nthreads = 10, .seconds = 10)
		{
		}
	Run()
		{
		.create()
		for .. .nthreads
			Thread(.threadfn)
		}
	create()
		{
		try Database("drop testupdates")
		Database("create testupdates (num, count) key(num)")
		for (i = 0; i < .nrecords; ++i)
			QueryOutput('testupdates', [num: i, count: 0])
		}
	threadfn()
		{
		try
			{
			n = 0
			conflicts = 0
			end = Date().Plus(seconds: .seconds)
			forever
				{
				for .. 10
					conflicts += .delete_output2()
				n += 100
				if Date() > end
					break
				}
			Print('finished:', n, '-', (conflicts * 100 / n).Round(0) $ '%')
			}
		catch (e)
			Print(crashed: e)
		}
//	update()
//		{
//		i = Random(.nrecords / 2) + Random(.nrecords / 2)
//		try
//			Transaction(update:)
//				{|t|
//				x = t.Query1('testupdates', num: i)
//				++x.count
//				Thread.Sleep(Random(5))
//				x.Update()
//				}
//		catch (e, 'Transaction: block commit failed: transaction conflict')
//			{
//			.log.Writeline(e)
//			return 1
//			}
//		return 0
//		}
//	delete_output()
//		{
//		i = Random(.nrecords / 2) + Random(.nrecords / 2)
//		try
//			Transaction(update:)
//				{|t|
//				x = t.Query1('testupdates', num: i)
//				x.Delete()
//				++x.count
//				t.QueryOutput('testupdates', x)
//				}
//		catch (e, 'Transaction: block commit failed: transaction conflict')
//			{
//			if not e.Has?('read conflict')
//				.log.Writeline(e)
//			return 1
//			}
//		return 0
//		}
	delete_output2()
		{
		x = false
		i = Random(.nrecords / 2) + Random(.nrecords / 2)
		try
			Transaction(update:)
				{|t|
				if false is x = t.Query1('testupdates', num: i)
					return 0
				x.Delete()
				}
		catch (unused, 'Transaction: block commit failed: transaction conflict')
			return 1
		try
			Transaction(update:)
				{|t|
				++x.count
				t.QueryOutput('testupdates', x)
				}
		catch (unused, 'Transaction: block commit failed: transaction conflict')
			return 1
		return 0
		}
	}