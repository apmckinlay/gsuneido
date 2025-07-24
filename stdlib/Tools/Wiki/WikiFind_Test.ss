// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_find()
		{
		cl = WikiFind
			{
			Ftsearch(unused)
				{
				return Object(
					#(name: 'TestyTest1'),
					#(name: 'TestyTest2'),
					#(name: 'TestyTest3')
					)
				}
			WikiFind_displayOrphanedMessage(unused)
				{
				return ''
				}
			}
		find = cl.WikiFind_findByFtsearch
		Assert(find('Test') has: 'TestyTest1')
		Assert(find('Test') has: 'TestyTest2')
		Assert(find('Test') has: 'TestyTest3')
		Assert(find('') is: '<p><b>Please enter something to search for</b></p>')
		}
	}