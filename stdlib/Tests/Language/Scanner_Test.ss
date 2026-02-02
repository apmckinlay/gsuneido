// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_extras()
		{
		Assert(ScannerFind("test sort test1", "sort") is: 5)
		Assert(ScannerHas?("test sort test1", "sort"))
		Assert(ScannerHas?("test where s = 'sort'", "sort") is: false)
		}
	Test_Scanner()
		{
		.test('')
		.test(`hello true "q\n"`,
			[tok: 'hello', val: 'hello', pos: 0, end: 5, type: #IDENTIFIER],
			[tok: 'true', pos: 6, end: 10, key:, type: #IDENTIFIER],
			[tok: `"q\n"`, val: 'q\n', pos: 11, end: 16, type: #STRING],
			)
		.test('and or is isnt not is:',
			[key:, type: ""], [key:, type: ""], [key:, type: ""], [key:, type: ""],
				[key:, type: ""], [key: false], [tok: ':'])
		}
	Test_QueryScanner()
		{
		.qtest('')
		.qtest('tables join by(table) columns where',
			[tok: 'tables'],
			[tok: 'join', key:],
			[tok: 'by'],
			[tok: '('],
			[tok: 'table'],
			[tok: ')'],
			[tok: 'columns'],
			[tok: 'where', key:],
			)
		}
	test(@args)
		{
		sc = Scanner(args[0])
		.match(sc, args)
		}
	qtest(@args)
		{
		sc = QueryScanner(args[0])
		.match(sc, args)
		}
	match(sc, args)
		{
		for (i = 1; false isnt x = .next(sc); ++i)
			for m in args[i].Members()
				Assert(x[m] is: args[i][m])
		Assert(i is: args.Size(), msg: "wrong number of tokens")
		}
	next(sc)
		{
		do
			{
			pos = sc.Position()
			if sc is type = sc.Next2()
				return false
			}
			while sc.Type() is #WHITESPACE
		return [tok: sc.Text(), :pos, end: sc.Position(),
			val: sc.Value(), :type, key: sc.Keyword?()]
		}
	}
