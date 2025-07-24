// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Tablename: 'testConcurrency'
	Duration: 30 // seconds
	Nrecords: 1000
	actionsPerCycle: 100
	CallClass()
		{
		(new this).Run()
		}
	Run()
		{
		.setup()
		.print('Running...')
		n = 0
		startTime = Date()
		do
			{
			.actionsPerCycle.Times(.doAction)
			n += .actionsPerCycle
			}
			while Date().MinusSeconds(startTime) < .Duration
		.print('Completed', n, 'actions')
		}
	setup()
		{
		if TableExists?(.Tablename)
			{
			.Nrecords = QueryCount(.Tablename).Round(-2)
			.print('Using existing ' $ .Tablename, 'with', .Nrecords, 'records')
			}
		else
			.create()
		}
	create()
		{
		.print('Creating', .Tablename, 'with', .Nrecords, 'records')
		Database("create " $ .Tablename $
			' (a, b, c, d, e, f, g)
			key(a) index(b, c)')
		outputsPerIteration = 100
		for (i = 0; i < .Nrecords / outputsPerIteration; ++i)
			Transaction(update:)
				{|t|
				outputsPerIteration.Times { t.QueryOutput(.Tablename, .record()) }
				}
		}
	doAction()
		{
		// using string case values to avoid magic number warnings
		try
			switch (String(Random(5 /* = number of possible actions */)))
				{
			case '0': .lookup()
			case '1': .range()
			case '2': .append()
			case '3': .update()
			case '4': .erase()
				}
		catch (e)
			.print('ERROR', e)
		}
	lookup()
		{
		QueryApply(.Tablename, b: n = Random(.Nrecords))
			{|x| Assert(x.b is: n) }
		}
	range()
		{
		maxRangeSize = 500
		from = Random(.Nrecords)
		to = from + Random(maxRangeSize)
		QueryApply(.Tablename $ ' where b >= ' $ from $ ' and b < ' $ to)
			{|x| Assert(from <= x.b and x.b < to) }
		}
	append()
		{
		r = .record()
		.retryTransaction('append')
			{|t|
			t.QueryOutput(.Tablename, r)
			}
		}
	update()
		{
		.retryTransaction('update')
			{|t|
			t.QueryDo('update ' $ .Tablename $ ' where b = ' $ Random(.Nrecords) $
				'set c = ' $ Random(.Nrecords))
			}
		}
	erase()
		{
		.retryTransaction('erase')
			{|t|
			t.QueryDo('delete ' $ .Tablename $ ' where b = ' $ Random(.Nrecords))
			}
		}
	retryTransaction(which, block)
		{
		try
			RetryTransaction(block)
		catch (e, 'Retry')
			.print('ERROR', which, e)
		}
	record()
		{
		maxNumberValue = 10000
		return [a: Timestamp(), b: Random(.Nrecords),
			c: .randomString(), d: Random(maxNumberValue),
			e: .randomString(), f: Random(maxNumberValue), g: .randomString()]
		}
	randomString()
		{
		maxNumberPart = 10
		strings = #(hello, world, now, the, time, for, all, foo, bar, foobar)
		return strings[Random(strings.Size())].Repeat(Random(maxNumberPart))
		}
	print(@args)
		{
		//Print(@args)
		ServerEval('Print', args.Map!(Display).Join(','))
		}
	}