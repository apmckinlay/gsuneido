// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		// NOTE: need to use .Func to bypass caching
		tbl = .MakeTable('(num, name, path, text) key(num)')
		Assert(GetBookText.Func('Foo', tbl) is: false)
		QueryOutput(tbl, [num: 123, name: 'Foo', path: '/res', text: 'foofoo'])
		Assert(GetBookText.Func('Foo', tbl) is: 'foofoo')
		QueryOutput(tbl, [num: 456, name: 'Bar', path: '/res/icon', text: 'barbar'])
		Assert(GetBookText.Func('icon/Bar', tbl) is: 'barbar')
		}
	}