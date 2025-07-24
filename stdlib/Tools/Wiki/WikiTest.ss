// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_wikilocking()
		{
		Assert(WikiLock('Test_wikilock'))
		Assert(WikiLock('Test_wikilock') isnt: true)
		WikiUnlock('Test_wikilock')
		Assert(WikiLock('Test_wikilock'))
		}
	Test_multiuser()
		{
		table = .MakeTable('(name, text, created, edited) key(name)')
		pg = WikiEdit('edit', 'WikiTest', :table)
		Assert(pg hasnt: 'unlock')
		pg2 = WikiEdit('edit', 'WikiTest', :table)
		Assert(pg2 has: 'unlock')
		pg2 = WikiEdit('edit', 'WikiTest2', :table)
		Assert(pg2 hasnt: 'unlock')
		WikiSave('WikiTest', 'editmode=edit&text=test', :table)
		pg = WikiEdit('edit', 'WikiTest', :table)
		Assert(pg hasnt: 'unlock')
		}
	Teardown()
		{
		Suneido.Delete('WikiLock')
		super.Teardown()
		}
	}