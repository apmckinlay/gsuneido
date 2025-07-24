// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		book = .MakeTable("(name, path, num, text) key(path, name) key(num)",
			[num: 1, name: 'One'],
			[num: 2, name: 'Two'],
			[num: 3, name: 'Three']
			)
		helpbook = book $ "Help"
		Database("ensure " $
			helpbook $ " (name, path, num, text) key(num)")
		index = helpbook $ "HowToIndex"

		BookHowToIndex(helpbook)
		Assert(QueryCount(index) is: 0)

		QueryOutput(helpbook,
			[num: 1, name: 'Two'])
		QueryOutput(helpbook,
			[num: 2, path: '/Two', name: 'How Do I'])
		QueryOutput(helpbook,
			[num: 3, path: '/Two/How Do I', name: 'How do I frizzle?',
				text: '<!-- option: Two -->'])
		bad = BookHowToIndex(helpbook)
		Assert(bad is: #())
		Assert(QueryCount(index) is: 1)
		Assert(QueryFirst(index $ ' sort name')
			is: [name: 'Two', howtos: #('/Two/How Do I/How do I frizzle?')])

		QueryOutput(helpbook,
			[num: 4, path: '/Two/How Do I', name: 'How do I cogitate?',
				text: '<!-- option: Two -->'])
		bad = BookHowToIndex(helpbook)
		Assert(bad is: #())
		Assert(QueryCount(index) is: 1)
		Assert(QueryFirst(index $ ' sort name')
			is: [name: 'Two', howtos: #('/Two/How Do I/How do I frizzle?',
			'/Two/How Do I/How do I cogitate?')])

		QueryOutput(helpbook,
			[num: 5, path: '/Two/How Do I', name: 'How do I munge?',
				text: '<!-- option: Three --> <!-- option: One -->'])
		bad = BookHowToIndex(helpbook)
		Assert(bad is: #())
		Assert(QueryAll(index $ ' sort name')
			is: [
			[name: 'One', howtos: #('/Two/How Do I/How do I munge?')],
			[name: 'Three', howtos: #('/Two/How Do I/How do I munge?')],
			[name: 'Two', howtos: #('/Two/How Do I/How do I frizzle?',
				'/Two/How Do I/How do I cogitate?')]
			])

		QueryOutput(helpbook,
			[num: 6, path: '/Two/How Do I', name: 'How do I munge?',
				text: '<!-- option: Thre -->'])
		bad = BookHowToIndex(helpbook)
		Assert(bad is: #('Thre'))
		}
	Test_Duplicates()
		{
		book = .MakeTable("(name, path, num, text) key(path, name) key(num)",
			[num: 1, name: 'One'],
			[num: 2, name: 'Two'],
			[num: 3, name: 'Three']
			)
		helpbook = book $ "Help"
		Database("ensure " $ helpbook $ " (name, path, num, text) key(num)")
		index = helpbook $ "HowToIndex"

		QueryOutput(book, [num: 4, path: '/Two', name: 'Quacking'])
		QueryOutput(book, [num: 5, path: '/Three', name: 'Quacking'])
		QueryOutput(helpbook, [num: 1, name: 'Two'])
		QueryOutput(helpbook, [num: 2, path: '/Two', name: 'How Do I'])
		QueryOutput(helpbook, [num: 3, path: '/Two/How Do I', name: 'How do I quack two?',
			text: '<!-- option: /Two/Quacking -->'])
		QueryOutput(helpbook, [num: 4, name: 'Three'])
		QueryOutput(helpbook, [num: 5, path: '/Three', name: 'How Do I'])
		QueryOutput(helpbook, [num: 6, path: '/Three/How Do I',
			name: 'How do I quack three?', text: '<!-- option: /Three/Quacking -->'])

		bad = BookHowToIndex(helpbook)
		Assert(bad is: #())
		Assert(QueryAll(index $ ' sort name')
			is: [
			[name: '/Three/Quacking', howtos: #('/Three/How Do I/How do I quack three?')],
			[name: '/Two/Quacking', howtos: #('/Two/How Do I/How do I quack two?')]
			])
		}
	}