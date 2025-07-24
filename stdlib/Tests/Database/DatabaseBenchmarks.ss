// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		results = Object().Set_default(Object().Set_default(Object()))
		for keysize in #(10, 20, 30, 40 50, 60, 70, 80, 90, 100)
			.benchmarks(results, 'keysize', :keysize, recsize: 200)
		for recsize in #(100, 200, 300, 400 500, 600, 700, 800, 900, 1000)
			.benchmarks(results, 'recsize', :recsize)
		for nrecs in #(10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000,
			90000, 100000)
			.benchmarks(results, 'nrecs', :nrecs)
		.print(results)
		}
	print(results)
		{
		for b in .Members().Filter({|m| m.Has?('_b_') }).Sort!()
			{
			name = b.AfterFirst('_b_')
			Print(name)
			Print('-'.Repeat(name.Size()))
			Print('keysize')
			for x in results[name]['keysize']
				Print(x[0].Pad(3 /*= keysize padding*/, ' '), x[1])
			Print()
			Print('recsize')
			for x in results[name]['recsize']
				Print(x[0].Pad(4 /*= recsize padding*/, ' '), x[1])
			Print()
			Print('nrecs')
			for x in results[name]['nrecs']
				Print(x[0].Pad(6 /*= nrecs padding*/, ' '), x[1])
			Print()
			}
		}
	benchmarks(results, which, keysize = 10, recsize = 100, nrecs = 10000)
		{
		if TableExists?('test')
			Database('drop test')
		TestData.Make('test', :keysize, :recsize, :nrecs)
		for b in .Members().Filter({|m| m.Has?('_b_') })
			{
			name = b.AfterFirst('_b_')
			t = Timer(reps: 10)
				{
				this[b]()
				}
			results[name][which].Add([.varying(which, keysize, recsize, nrecs), t])
			}
		}
	varying(which, keysize, recsize, nrecs)
		{
		switch which
			{
		case 'keysize': return keysize
		case 'recsize': return recsize
		case 'nrecs': return nrecs
			}
		}
	b_Read()
		{
		QueryApply('test') {|unused| }
		}
	b_Count()
		{
		Query1('test summarize count')
		}
	b_Select()
		{
		QueryApply('test where fld = 0') {|unused| }
		}
	b_Sort()
		{
		QueryFirst('test extend sortkey = BenchmarkSortKey(key.Size()) sort sortkey')
		}
	b_Join()
		{
		QueryApply('test join by (key) (test rename fld to fld2)') {|unused| }
		}
	}