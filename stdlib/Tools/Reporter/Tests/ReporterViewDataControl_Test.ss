// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.TearDownIfTablesNotExist('customizable')
		if not TableExists?('customizable')
			Customizable.EnsureSaveInTable()
		}
	Test_sort()
		{
		func = ReporterViewDataControl.ReporterViewDataControl_getSort

		query = 'fake_table sort field1, field2, field3'
		renamed = #(field1_name: field1_name_testSuffix,
			field1_abbrev: field1_abbrev_testSuffix)
		Assert(func(query, renamed) is: 'field1, field2, field3')

		query = 'fake_table sort reverse field1, field2, field3'
		renamed = #(field1_name: field1_name_testSuffix,
			field1_abbrev: field1_abbrev_testSuffix)
		Assert(func(query, renamed) is: 'reverse field1, field2, field3')

		query = 'fake_table sort field1_name, field1_abbrev'
		renamed = #(field1_name: field1_name_testSuffix,
			field1_abbrev: field1_abbrev_testSuffix)
		Assert(func(query, renamed) is: 'field1_name_testSuffix, ' $
			'field1_abbrev_testSuffix')

		query = 'fake_table sort reverse field1_name, field1_abbrev'
		renamed = #(field1_name: field1_name_testSuffix,
			field1_abbrev: field1_abbrev_testSuffix)
		Assert(func(query, renamed) is: 'reverse field1_name_testSuffix, ' $
			'field1_abbrev_testSuffix')

		query = 'fake_table sort field1, field2, field3_name, field3_abbrev'
		renamed = #(field3_name: field3_name_testSuffix,
			field3_abbrev: field3_abbrev_testSuffix)
		Assert(func(query, renamed) is: 'field1, field2, ' $
			'field3_name_testSuffix, field3_abbrev_testSuffix')

		query = 'fake_table sort reverse field1, field2, field3_name, field3_abbrev'
		renamed = #(field3_name: field3_name_testSuffix,
			field3_abbrev: field3_abbrev_testSuffix)
		Assert(func(query, renamed) is: 'reverse field1, field2, ' $
			'field3_name_testSuffix, field3_abbrev_testSuffix')
		}

	Test_one()
		{
		if not TableExists?('configlib')
			return

		last = QueryLast('configlib sort num')
		if last is false
			.AddTeardown({ QueryDo('delete configlib') })
		else
			.AddTeardown({ QueryDo('delete configlib where num > ' $
				Display(last['num'])) })

		cl = ReporterViewDataControl
			{
			ReporterViewDataControl_suffix()
				{
				return '_test_suffix'
				}
			ReporterViewDataControl_showData(@unused)
				{
				}
			}
		table = .MakeTable('(a_num, a_name, b?, c?) key(a_num)')
		.MakeLibraryRecord([name: "Field_a_name", text:
			`Field_string
				{
				Prompt: 'Test A Name'
				}`])
		.MakeLibraryRecord([name: "Field_b?", text:
			`Field_string
				{
				Prompt: 'Test B?'
				}`])
		customOb = .MakeCustomField(table, 'Number')

		q = table $
			' summarize a_name, total b?, total c?, max ' $ customOb.field $
			' sort total_b?, total_c?'
		sf = SelectFields(#('a_name'))
		sf.AddField('total_b?', 'Total B')
		sf.AddField('total_c?', 'Total C')
		sf.AddField('max_' $ customOb.field, 'Total ' $ customOb.prompt)
		result = cl(q, #(), sf,)
		Assert(result.query is: QueryStripSort(q) $
				' rename total_c? to total_c_test_suffix' $
				' sort total_b?, total_c_test_suffix')
		Assert(result.cols is: Object('max_' $ customOb.field, 'a_name',
			'total_c_test_suffix', 'total_b?'))
		}
	}