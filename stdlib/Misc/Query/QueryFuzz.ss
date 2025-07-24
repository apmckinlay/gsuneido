// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
// Random queries for fuzz testing
class
	{
	tables: (
		(cus, '(ck, c1, c2, c3, c4) key(ck)'),
		(ivc, '(ik, i1, i2, i3, i4, ck) key(ik) index(ck) in cus'),
		(aln, '(ik, ak, a1, a2, a3, a4) key(ik, ak)')
		(bln, '(ik, bk, b1, b2, b3, b4) key(ik, bk)')
		)
	sizes: (c: 10, i: 100, a: 500, b: 1000)

	TestCorpus()
		{
		// TODO multi-thread
		for query in QueryFuzzCorpus
			Assert(QueryHash(query, details:) is: QueryAltHash(query, details:))
		}

	maxErrors: 20
	Multi(secs) // multithreaded, comparing to QueryAlt
		{
		results = Object(count: 0, errors: Object())
		nthreads = Suneido.GoMetric("/sched/gomaxprocs:threads")
		wg = WaitGroup()
		for .. nthreads
			wg.Thread({ .multi(secs, results) })
		for query in QueryFuzzCorpus
			if not .multi1(query.Tr(' \t\r\n', ' ') $ " /*corpus*/", results)
				break
		wg.Wait(10 /*= time for QueryFuzzCorpus */ +
			(secs * 1.75).Int())/*= allow for Timer running over */
		return results
		}
	multi(secs, results)
		{
		start = Date()
		while (Date().MinusSeconds(start) < secs)
			{
			for ..10 /*= loop for time */
				if not .multi1(.MakeQuery(), results)
					return
			}
		}

	multi1(base, results)
		{
		err = {|e|
			if results.errors.Size() < .maxErrors
				results.errors.Add(e)
			else
				return false
			}
		for query in [base, .vary(base), .vary(base)]
			try
				if QueryHash(query) is QueryAltHash(query)
					++results.count
				else
					err(query)
			catch (e)
				err(query $ "\n=> " $ e)
		return true
		}

	Run(secs)
		{
		File("qfuzz.log", "w")
			{|f|
			return .run(secs, { f.Writeline(it) })
			}
		}
	RunPrint(secs)
		{
		return .run(secs, Print)
		}
	run(secs, write)
		{
		count = 0
		for query in QueryFuzzCorpus
			{
			count += .run1vary(query.Tr(' \t\r\n', ' ') $ " /*corpus*/", write)
			}
		QueryApply('views')
			{|x|
			count += .run1(x.view_name $ " /*view*/", write)
			}
		Plugins().ForeachContribution('Reporter', 'queries')
			{ |x|
			query = x.query
			if query.BeforeFirst('.').GlobalName?() and
				not Uninit?(query.BeforeFirst('.'))
				query = Global(query)()
			if query.Size() < 3900 /*= max line length */
				count += .run1(query.Tr(' \t\r\n', ' ') $ " /*reporter " $ x.name $ "*/",
					write)
			}
		Timer(:secs)
			{
			count += .run1vary(.MakeQuery(), write)
			}
		return count
		}
	run1vary(query, write)
		{
		n = .run1(query, write)
		for .. .nvariations
			n += .run1(.vary(query), write)
		return n
		}
	nvariations: 2
	vary(query)
		{
		return query.Replace(`"\d*"`)
			{|unused|
			#(`""`, `"3"`, `"12"`).RandVal()
			}
		}
	run1(query, write)
		{
		try
			{
			hash = QueryHash(query)
			write(hash $ '\t' $ query)
			return 1
			}
		catch (e)
			{
			if e.Has?("can't find")
				return 0 // skip/ignore view that references code not in use
			Print(query)
			Print("    =>", e)
			AddFile("qfuzz.err", query $ '\r\n=> ' $ e $ '\r\n')
			}
		return 0
		}
	Checker(secs)
		{
		if not FileExists?(gs = './gsport.exe')
			gs = '../gsport.exe'
		count = 0
		RunPiped(gs $ ' "Print(QueryFuzz.RunPrint(' $ secs $ ')); Exit()"')
			{|rp|
			while false isnt line = rp.Readline()
				if line =~ "^\d+\t[(a-z]"
					{
					++count
					.check(line)
					}
				else
					Print(line)
			rp.CloseWrite()
			}
		Print(count)
		Exit()
		}
	Check(errProcess = false)
		{
		count = 0
		File("qfuzz.log", "r")
			{|f|
			while false isnt line = f.Readline()
				{
				++count
				.check(line, errProcess)
				}
			}
		return count
		}
	check(line, errProcess = false)
		{
		if errProcess is false
			errProcess = .printQueryDiff
		expected = Number(line.BeforeFirst('\t'))
		query = line.AfterFirst('\t')
		try
			{
			hash = QueryHash(query)
			if hash isnt expected
				errProcess(query)
			}
		catch(e)
			errProcess(e $ '\r\n' $ query)
		}
	printQueryDiff(query)
		{
		Print("DIFF", query)
		}

	MakeTables(seed = false)
		{
		if seed is false
			seed = Random(9999999999) /*= large range */
		Random.Seed(seed)
		.DropTables()
		for x in .tables
			{
			Database('create ' $ x[0] $ ' ' $ x[1])
			for .. .sizes[x[0][0]]
				.MakeRec(x[0])
			}
		return seed
		}
	MakeRec(tbl)
		{
		random = 100
		for .. 20 /*= retries for duplicate keys */
			{
			rec = Object()
			for col in QueryColumns(tbl)
				if col in ('ak', 'bk')
					rec[col] = String(Random(random)) /*= line key range */
				else
					rec[col] = Random(random) < 10 /*= min */
						? "" : String(Random(2 * .sizes[col[0]]))
			try
				{
				QueryOutput(tbl, rec)
				return
				}
			}
		}
	DropTables()
		{
		for x in .tables.Copy().Reverse!()
			if TableExists?(x[0])
				Database('drop ' $ x[0])
		}
	start: (
		cus,
		(cus join ivc)
		(ivc join aln)
		(cus join (ivc join aln))
		((cus join ivc) join aln)
		)
	MakeQuery()
		{
		q = ""
		try
			{
			q = [.copy(.start.RandVal())]
			.add_unions(q)
			.leftjoin(q)
			.aln_to_bln(q)
			.add_query1(q)
			q = Nested.FlatStr(q)
			q = .add_sort(q)
			return q
			}
		catch (e)
			{
			AddFile("qfuzz.err", Display(q) $ '\r\n=> ' $ e $ '\r\n')
			SuneidoLog("ERROR in MakeQuery " $ e)
			return "cus"
			}
		}
	leftjoin(q)
		{
		Nested.Visit(q)
			{|x,j/*unused*/,ob/*unused*/|
			if Object?(x) and 'join' is x.GetDefault(1, false) and
				Random(4) is 3 /*= 1/4 leftjoin */
				{
				x[1] = 'leftjoin'
				if Random(2) is 1 /*= 1/2 the time */
					x.Reverse!()
				}
			}
		}
	add_unions(q)
		{
		for i in ..4 /*= frequency of unions */
			{
			r = Nested.Random(q)
			ob = r[0]
			i = r[1]
			if ob[i] not in ('join', 'leftjoin', 'union')
				ob[i] = [ob[i], 'union', .copy(ob[i])]
			}
		}
	aln_to_bln(q)
		{
		Nested.Visit(q)
			{|x,j,ob|
			if x is 'aln' and Random(101) < 50 /*= 1/2 bln */
				ob[j] = 'bln'
			}
		}
	add_query1(q)
		{
		for n in ..3 /*= number of query1 */
			{
			r = Nested.Random(q)
			ob = r[0]
			i = r[1]
			if ob[i] not in ('join', 'union')
				ob[i] = .add_1query1(ob[i], n)
			}
		}
	add_1query1(x, n)
		{
		allcols = #(ck, c1, c2, c3, c4, ik, i1, i2, i3, i4,
			ak, a1, a2, a3, a4, bk, b1, b2, b3, b4)
		q = Nested.FlatStr(x)
		cols = QueryColumns(q)
		if cols is #()
			return x
		if not q.Identifier?()
			q = '(' $ q $ ')'
		randomStr = 100
		switch Random(8) /*= random selection */
			{
		case 0,1,2: /*= where */
			col = cols.RandVal()
			op = #('is', 'is', 'is', 'is',
				'isnt', '<', '<=', '>', '>=').RandVal()
			val = Random(randomStr + 1) < 25 /*= min */
				? "" : String(Random(randomStr))
			return '(' $ q $ ' where ' $ col $ ' ' $ op $ ' "' $
				val $ '")' /*= range */
		case 3: /*= extend constant */
			return '(' $ q $ ' extend x' $ n $ ' = "' $ String(n) $ '")'
		case 4: /*= extend expression */
			e = allcols.RandVal()
			if cols.Has?(e)
				e = 'x' $ n
			col = cols.RandVal()
			return '(' $ q $ ' extend ' $ e $ ' = ' $ col $ ')'
		case 5: /*= extend <rule> */
			return '(' $ q $ ' extend r' $ n $ ')'
		case 6: /*= rename */
			nonkey = cols.Filter({ it[1] isnt 'k' })
			if nonkey is #()
				return x
			col = nonkey.RandVal()
			return '(' $ q $ ' rename ' $ col $ ' to y' $ n $ ')'
		case 7: /*= remove (project) */
			nonkey = cols.Filter({ it[1] isnt 'k' })
			if nonkey is #()
				return x
			rem = nonkey.Shuffle!()[..Random(2)+1]
			return q $ ' remove ' $ rem.Join(',') // = project
			}
		}
	add_sort(q)
		{
		cols = QueryColumns(q)
		if cols is #() or Random(2) is 1
			return q // no sort
		cols = cols.Copy().Shuffle!()[..Random(3)+1] /*= number of sort fields */
		sort = ' sort '
		if Random(3) /*=frequency*/ is 1
			sort $= 'reverse '
		return q $ sort $ cols.Join(',')
		}
	copy(x)
		{
		if Object?(x)
			x = x.DeepCopy()
		return x
		}
	}