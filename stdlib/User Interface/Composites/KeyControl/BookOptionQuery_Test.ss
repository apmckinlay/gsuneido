// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		book = .MakeTable('(path, name, text) key(path, name)',
			[path: '/A',		name: 'x', text: 'one'],
			[path: '/A/sub',	name: 'y', text: 'two()'],
			[path: '/B',		name: 'z', text: 'three'],
			[path: '/res',		name: '1', text: 'four'],
			[path: '/res/sub',	name: '2', text: 'five'])
		Assert(Query1(BookOptionQuery(book, 'abc')) is: false)
		Assert(Query1(BookOptionQuery(book, 'one')) isnt: false)
		Assert(Query1(BookOptionQuery(book, 'two')) isnt: false)
		Assert(Query1(BookOptionQuery(book, 'four')) is: false)
		Assert(Query1(BookOptionQuery(book, 'five')) is: false)
		}
	}