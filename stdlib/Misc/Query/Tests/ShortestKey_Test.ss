// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_unit()
		{
		Assert(ShortestKey(#('a')) is: 'a')
		Assert(ShortestKey(#('a,b')) is: 'a,b')
		Assert(ShortestKey(#('a,b', 'c')) is: 'c')
		Assert(ShortestKey(#('a,b', 'c', 'd')) is: 'c') // first
		Assert(ShortestKey(#('a,b', 'c', 'c_num')) is: 'c_num') // prefer num
		Assert(ShortestKey(#('a,b', 'c_name', 'c_num')) is: 'c_num')
		Assert(ShortestKey(#('a,b', 'c_name', 'c')) is: 'c') // avoid name
		}
	Test_main()
		{
		table = .MakeTable('(aaaa,b,c) key(aaaa) key(b,c)')
		Assert(ShortestKey(table) is: 'aaaa')
		Assert(ShortestKey(table $ ' sort b') is: 'aaaa')
		Assert(ShortestKey(QueryKeys(table)) is: 'aaaa')
		WithQuery(table)
			{ |q|
			Assert(ShortestKey(q) is: 'aaaa')
			}
		Cursor(table)
			{ |c|
			Assert(ShortestKey(c) is: 'aaaa')
			}
		}
	}