// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.MakeFile() // to ensure file teardown
		recs = Object()
		table = .MakeTable('(a,b,c) key(a)',
			recs[0] = [a: 1, b: 2, c: 'hello'],
			recs[1] = [a: 4, b: 5, c: 'hello\nworld'])
		schema = Schema(table)
		DumpText(table)
		Database('destroy ' $ table)
		Assert(not TableExists?(table))
		LoadText(table)
		Assert(Schema(table) is: schema)
		i = 0
		QueryApply(table)
			{ |x|
			Assert(x is: recs[i++])
			}
		}
	}