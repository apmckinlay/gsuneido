// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Ntables: 100
	Nrecords: 1000
	CreateTables()
		{
		for i in .. .Ntables
			{
			try Database('drop transtress' $ i)
			Database('create transtress' $ i $ " (a,b) key(a)")
			}
		}
	CallClass()
		{
		(new this).Run()
		}
	Run()
		{
		sw = Stopwatch()
		.tran = Object()
		.func = [.startTran, .commitTran,
			.readRec, .outputRec, .updateRec, .deleteRec,
			.readRec, .outputRec, .updateRec, .deleteRec]
		.ntran = 0
		.startFails = 0
		.commitFails = 0
		for ..100000 /*= number of actions */
			.run1()
		for t in .tran
			try
				t.Complete()
		Print(sw(), "ntran", .ntran,
			"startFails", .startFails, "commitFails", .commitFails)
		}
	run1()
		{
		if .tran.Empty?()
			.startTran()
		else
			{
			f = .func.RandVal()
			t = .tran[.rand(.tran.Size())]
			try
				f(t)
			catch (unused, "*transaction aborted")
				.tran.Remove(t)
			}
		}
	startTran(unused = false)
		{
		++.ntran
		try
			.tran.Add(Transaction(update:))
		catch
			++.startFails
		}
	commitTran(t)
		{
		.tran.Remove(t)
		try
			t.Complete()
		catch
			++.commitFails
		}
	readRec(t)
		{
		tbl = .randTable()
		t.Query1(tbl, a: .rand(.Nrecords))
		}
	outputRec(t)
		{
		tbl = .randTable()
		try t.QueryOutput(tbl, [a: .rand(.Nrecords)])
		}
	updateRec(t)
		{
		tbl = .randTable()
		if false isnt x = t.Query1(tbl, a: .rand(.Nrecords))
			{
			x.b = Random(123456) /*= range */
			x.Update()
			}
		}
	deleteRec(t)
		{
		tbl = .randTable()
		t.QueryDo('delete ' $ tbl $ ' where a = ' $ .rand(.Nrecords))
		}
	randTable()
		{
		return 'transtress' $ .rand(.Ntables)
		}
	rand(n) // biased to smaller numbers
		{
		return Random(1+Random(1+Random(n)))
		}
	}