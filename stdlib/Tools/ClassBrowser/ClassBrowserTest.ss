// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_base()
		{
		f = ClassBrowserModel.ClassBrowserModel_base
		Assert(f('class { }') is: 'class')
		Assert(f('A { }') is: 'A')
		Assert(f('/* Comment */ A { }') is: 'A')
		Assert(f('class B { }') is: 'B')
		Assert(f('function () {}') is: false)
		}
	Test_Main()
		{
		table = .MakeTable('(num, name, text) key(name)')
		cbm = ClassBrowserModel(Object(table))
		Assert(cbm.Children(0) is: #())
		Assert(cbm.Children?(0) is: false)

		QueryOutput(table, #(num: 10, name: A, text: 'class { }'))
		cbm = ClassBrowserModel(Object(table))
		Assert(cbm.Children(0) is: #((num: 10, name: A, group: true)))
		Assert(cbm.Children?(0))

		QueryOutput(table, #(num: 11, name: B, text: '/* Comment */ A { }'))
		cbm = ClassBrowserModel(Object(table))
		Assert(cbm.Children(0) is: #((num: 10, name: A, group: true)))
		Assert(cbm.Children(10) is: #((num: 11, name: B, group: true)))
		Assert(cbm.Children?(0))
		Assert(cbm.Children?(10))

		QueryOutput(table, #(num: 12, name: C, text: 'class B {}'))
		cbm = ClassBrowserModel(Object(table))
		Assert(cbm.Children(11) is: #((num: 12, name: C, group: true)))

		QueryOutput(table, #(num: 13, name: D, text: 'function () {}'))
		cbm = ClassBrowserModel(Object(table))
		Assert(cbm.Children(0) is: #((num: 10, name: A, group: true)))
		}
	}
