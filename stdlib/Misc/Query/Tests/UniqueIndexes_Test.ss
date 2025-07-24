// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(a,b,c) key(a)')
		Assert(UniqueIndexes(table) is: [])

		table = .MakeTable('(a,b) key(a) index unique(b)')
		Assert(UniqueIndexes(table) is: #('b'))

		table = .MakeTable('(a,b,c) key(a) index unique(b, c)')
		Assert(UniqueIndexes(table) is: #('b,c'))

		table = .MakeTable('(a,b,c,d,e) key(a)
			index unique(d, e)
			index unique(b)
			index unique(c)')
		Assert(UniqueIndexes(table) equalsSet: #('b','c', 'd,e'))
		}
	}