// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(wilma, fred, barney, betty, num) key(num)')
		.checkColumns(table, false, false)
		.checkColumns('', #(), #())
		.checkColumns('', #(fred, wilma), #(fred, wilma))

		table = .MakeTable('(wilma, fred, barney, betty, num, custom_999999) key(num)')
		.checkColumns(table, #(fred, wilma), #(fred, wilma, custom_999999))

		table = .MakeTable(
			'(wilma, fred, barney, betty, num, custom_999998, custom_999999) key(num)')
		.checkColumns(table, #(fred, wilma),
			#(fred, wilma, custom_999998, custom_999999))
		}

	checkColumns(name, cols, result)
		{
		Assert(ListCustomize.AddCustomColumns(name, cols) is: result)
		}

	Test_ReasonProtected()
		{
		table = .MakeTable('(wilma, fred, barney, betty, num) key(num)')
		mock = Mock(ListCustomize)
		mock.When.noInfo([anyArgs:]).CallThrough()
		mock.Eval(ListCustomize.ReasonProtected, [], protectField: false, hwnd: 0)
		mock.Verify.alert('No Information', 0)

		mock.Eval(ListCustomize.ReasonProtected, [], protectField: 'protect', hwnd: 0)
		mock.Verify.Times(2).alert('No Information', 0)

		.MakeLibraryRecord([name: table.Capitalize() $ "_allow_delete",
			text: `function() { return 'finalized!' }`])
		mock.Eval(ListCustomize.ReasonProtected, [], protectField: 'protect', hwnd: 0
			query: table)
		mock.Verify.alert('This record can not be deleted.\n\nfinalized!', 0)

		.MakeLibraryRecord([name: 'Rule_' $ table $ '_protect',
			text: `function() { return 'record protected' }`])
		mock.Eval(ListCustomize.ReasonProtected, [],
			protectField: table $ '_protect', hwnd: 0, query: table)
		mock.Verify.alert('This record can not be deleted.\n\n' $
			'finalized!\n\nrecord protected', 0)
		}

	Test_BuildCustomKeyFromQueryTitle()
		{
		fn = ListCustomize.BuildCustomKeyFromQueryTitle
		Assert(fn(title: false, query: 'test_query') is: false)
		Assert(fn(title: 'Hello Screen', query: 'test_query') is: 'Hello Screen ~ ')
		Assert(fn(title: 'Hello Screen', query: 'stdlib') is: 'Hello Screen ~ stdlib')
		}

	Test_AddCustomColumns()
		{
		c = ListCustomize
			{
			ListCustomize_getPermissableFields(unused) { return Object('custom_000001') }
			}
		fn = c.AddCustomColumns
		columns = #('name', 'address')
		query = 'people'
		Assert(fn(query, columns) is: Object('name', 'address', 'custom_000001'))

		query = ''
		Assert(fn(query, columns) is: columns)

		columns = #('name', 'address', 'custom_000002')
		query = 'people'
		Assert(fn(query, columns)
			is: Object('name', 'address', 'custom_000002', 'custom_000001'))

		columns = false
		Assert(fn(query, columns) is: false)
		}
	}