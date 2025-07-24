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

		QueryOutput(helpbook,
			[num: 1, name: 'Two'])
		QueryOutput(helpbook,
			[num: 5, path: '/Two', name: 'Three', text: '<html><p>Testing</p></html>'])
		QueryOutput(helpbook,
			[num: 2, path: '/Two', name: 'How Do I'])
		QueryOutput(helpbook,
			[num: 3, path: '/Two/How Do I', name: 'How do I frizzle?',
				text: '<!-- option: /Two/Three -->'])
		BookHowToIndex(helpbook)
		text = BookHowToLinks('/Two/Three', helpbook)
		Assert(text
			matches: '(?q)<a href="/' $ helpbook $ '/Two/How Do I/How do I frizzle?">')

		QueryOutput(helpbook,
			[num: 4, path: '/Two/How Do I ...?', name: 'How do I spizzle?',
				text: '<!-- option: /Two/Three -->'])
		BookHowToIndex(helpbook)
		text = BookHowToLinks('/Two/Three', helpbook)
		Assert(text
			matches: '(?q)<a href="/' $ helpbook $ '/Two/How Do I/How do I frizzle?">')
		Assert(text
			matches: '(?q)<a href="/' $ helpbook $
				'/Two/How Do I ...?/How do I spizzle?">')

		QueryOutput(helpbook,
			[num: 6, path: '/main/sub', name: 'Order Tender History',
				text: '<html><p>Testing</p></html>'])
		QueryOutput(helpbook,
			[num: 7, path: '/main/How Do I', name: 'How do I test?',
				text: "<!-- option: /main/sub/Order Tender History -->

<h2>How do test?</h2>

<$GetBookPage('HDIUseEdi')$>"])
		BookHowToIndex(helpbook)
		text = BookHowToLinks('/main/sub/Order Tender History', helpbook)
		Assert(text
			matches: '(?q)<a href="/' $ helpbook $ '/main/How Do I/How do I test?">')
		}
	}