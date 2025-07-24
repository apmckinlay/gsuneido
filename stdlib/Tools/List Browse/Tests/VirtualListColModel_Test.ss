// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.TearDownIfTablesNotExist('usercolumns', 'userselects')
		.columnsSaveName = .TempName()

		UserColumns.EnsureTable()
		AccessSelectMgr.Ensure()
		}

	Test_HeaderChanged?()
		{
		colModel = VirtualListColModel(#(a, b), columnsSaveName: .columnsSaveName)
		Assert(QueryEmpty?('usercolumns', usercolumns_title: .columnsSaveName),
			msg: 'not empty')
		colModel.Destroy()
		colModel.SetHeaderChanged()

		colModel.Destroy()
		Assert(not QueryEmpty?('usercolumns', usercolumns_title: .columnsSaveName),
			msg: 'empty')
		}

	Test_InitiateCustomKey()
		{
		colModel = VirtualListColModel(#(a, b))
		Assert(colModel.VirtualListColModel_customKey is: false)
		Assert(colModel.VirtualListColModel_columnsSaveName is: false)

		colModel = VirtualListColModel(#(a, b), query: 'stdlib',
			customKey: 'access title | ctrl')
		Assert(colModel.VirtualListColModel_customKey is: 'access title | ctrl')
		Assert(colModel.VirtualListColModel_columnsSaveName is: 'access title | ctrl')

		colModel = VirtualListColModel(#(a, b), columnsSaveName: 'something')
		Assert(colModel.VirtualListColModel_customKey is: false)
		Assert(colModel.VirtualListColModel_columnsSaveName is: 'something')

		colModel = VirtualListColModel(#(a, b), columnsSaveName: 'something',
			query: 'stdlib', customKey: 'access title | ctrl')
		Assert(colModel.VirtualListColModel_customKey is: 'access title | ctrl')
		Assert(colModel.VirtualListColModel_columnsSaveName is: 'something')
		}

	Test_setCustomKey()
		{
		colModel = VirtualListColModel(#(a, b))
		colModel.VirtualListColModel_setCustomKey('test_custom_key', 'test_query', #())
		Assert(colModel.VirtualListColModel_customKey is: 'test_custom_key')
		Assert(colModel.VirtualListColModel_columnsSaveName is: 'test_custom_key')

		colModel = VirtualListColModel(#(a, b), defaultColumns: #(a), query: 'test_query')
		Assert(colModel.VirtualListColModel_customKey is: false)
		Assert(colModel.VirtualListColModel_columnsSaveName is: false)

		colModel = VirtualListColModel(#(a, b), defaultColumns: #(a),
			customKey: 'test_custom_key', query: 'test_query')
		Assert(colModel.VirtualListColModel_customKey is: 'test_custom_key')
		Assert(colModel.VirtualListColModel_columnsSaveName is: 'test_custom_key')
		defaultCols = Query1('usercolumns', usercolumns_title: "test_custom_key",
			usercolumns_user: "")
		Assert(defaultCols.usercolumns_order is: #(a, b))
		Assert(defaultCols.usercolumns_sizes is: #(100, 0))
		}

	Test_GetMandatoryFields()
		{
		.MakeLibraryRecord([name: "Field_test_mandatory_col", text:
			`Field_string
				{
				Control: #(Field, mandatory:)
				}`])
		.MakeLibraryRecord([name: "Field_test_regular_col", text:
			`Field_string
				{
				}`])
		colModel = VirtualListColModel(
			#(test_mandatory_col, test_regular_col, forced_mandatory),
			mandatoryFields: #(forced_mandatory))
		Assert(colModel.GetMandatoryFields() is: #(test_mandatory_col, forced_mandatory))

		colModel = VirtualListColModel(
			#(test_mandatory_col, test_regular_col, forced_mandatory),
			mandatoryFields: #(forced_mandatory))

		colModel.SetColumns(Object(#test_regular_col))
		colModel.SetColWidth(0, 200)
		colModel.AddMissingMandatoryCols()
		Assert(colModel.GetColumns()
			is: #(test_regular_col, test_mandatory_col, forced_mandatory))
		Assert(colModel.GetColWidth(1) is: 100)
		Assert(colModel.GetColWidth(2) is: 100)
		Assert(colModel.HeaderChanged?())

		fmts = colModel.VirtualListColModel_formatting.ListFormatting_dispFormats
		Assert(fmts.test_mandatory_col isnt: false)
		}

	Test_CapFieldPrompt()
		{
		table = .MakeTable('(num, date, Time) key (num)')
		c = VirtualListColModel(query: table)
		Assert(c.CapFieldPrompt('num') is: false)
		Assert(c.CapFieldPrompt('date') is: false)
		Assert(c.CapFieldPrompt('time') is: false)

		c = VirtualListColModel(query: table, headerSelectPrompt: 'no_prompts')
		Assert(c.CapFieldPrompt('num') is: false)
		Assert(c.CapFieldPrompt('date') is: false)
		Assert(c.CapFieldPrompt('time') is: 'Time')
		}

	Test_GetStretchCol()
		{
		c = new VirtualListColModel
		c.StretchColumn = 'a'
		c.VirtualListColModel_columns = #(a, b, c, e, f, g)
		c.VirtualListColModel_widths = #(10, 20, 10, 14, 13, 0)
		Assert(c.GetStretchCol() is: 0)

		// stretching columns is hidden
		c.StretchColumn = 'c'
		Assert(c.GetStretchCol() is: 2)

		// stretching columns is missing
		c.StretchColumn = 'd'
		Assert(c.GetStretchCol() is: 4)

		// resizing columns before the default stretch col
		c.StretchColumn = 'c'
		Assert(c.GetStretchCol(1 /* b */) is: 2)

		// resizing columns after the default stretch col
		c.StretchColumn = 'c'
		Assert(c.GetStretchCol(3 /* e */) is: 4)

		// resizing stretch column
		c.StretchColumn = 'c'
		Assert(c.GetStretchCol(2 /* c */) is: 4)
		}

	Test_GetSelectWhere()
		{
		_select = AccessSelectMgr(#(
			#('name', '=', 'Hello'),
			#('city', '=', 'Saskatoon')))
		_fields = #('city')
		m = VirtualListColModel
			{
			GetSelectMgr()    { return _select }
			GetSelectFields() { return SelectFields(_fields) }
			}
		colsFn = function (@unused) { #(name, city) }
		Assert(m.GetSelectWhere('', false, colsFn)
			is: ' where name is "Hello" where city is "Saskatoon"')
		specs = m.GetWhereSpecs(_select.Select_vals(), colsFn)
		Assert(specs.fields is: #(name, city))

		colsFn = function (@unused) { #(city) }
		Assert(m.GetSelectWhere('', false, colsFn)
			is: ' extend name  where name is "Hello" where city is "Saskatoon"')

		colsFn = function (@unused) { #() }
		Assert(m.GetSelectWhere('', false, colsFn)
			is: ' extend name,city  where name is "Hello" where city is "Saskatoon"')

		_select = AccessSelectMgr()
		colsFn = function (@unused) { #(name, city) }
		Assert(m.GetSelectWhere('', false, colsFn) is: '')

		id = .MakeIdField()
		_select = AccessSelectMgr(Object(
			Object(id.name, '=', 'Hello'),
			Object('city', 	'=', 'Saskatoon'),
			Object('fax', 	'=', '123')))
		_fields = Object(id.num, 'city')
		colsFn = {|unused| Object(id.num, 'city') }
		Assert(m.GetSelectWhere('', false, colsFn)
			is: ' extend fax ' $
				' where ' $ id.num $ ' < ""' $
				' where city is "Saskatoon"' $
				' where fax is "123"')

		outputRec = Record()
		outputRec[id.num] = tempNum = .TempName()
		outputRec[id.name] = "Hello"
		outputRec[id.abbrev] = .TempName()
		QueryOutput(id.table, outputRec)
		Assert(m.GetSelectWhere('', false, colsFn)
			is: ' extend fax ' $
				' where ' $ id.num $ ' in ("' $ tempNum $ '")' $
				' where city is "Saskatoon"' $
				' where fax is "123"')

		.SpyOn(GetForeignNumsFromNameAbbrevFilter).Return(false)
		colsFn = {|unused| Object(id.num, id.name, 'city') }
		Assert(m.GetSelectWhere('', false, colsFn)
			is: ' extend fax ' $
				id.leftjoin $
				'  where ' $ id.name $ ' is "Hello"' $
				' where city is "Saskatoon"' $
				' where fax is "123"')

		m = VirtualListColModel
			{
			GetSelectMgr()    { return _select }
			GetSelectFields() { return SelectFields(_fields) }
			VirtualListColModel_suppressed: true
			}
		Assert(m.GetSelectWhere('', false, colsFn)
			is: ' extend fax ' $
				id.leftjoin $
				'  where ' $ id.name $ ' is "Hello"' $
				' where city is "Saskatoon"' $
				' where fax is "123"' $
				' /*SLOWQUERY SUPPRESS*/ ')
		}

	Test_GetSelectFields()
		{
		colModel = VirtualListColModel(#(city, name))
		Assert(colModel.GetSelectFields().Prompts() equalsSet: #(City, Name))

		colModel = VirtualListColModel(#(city, name))
		Assert(colModel.GetSelectFields(#(email)).Prompts()
			equalsSet: #(City, Email, Name))
		Assert(colModel.GetSelectFields().Prompts() equalsSet: #(City, Email, Name))
		}

	Test_HasSelectedVals?()
		{
		m = VirtualListColModel.HasSelectedVals?
		ob = Object()
		vals = []
		ob.VirtualListColModel_selectMgr = FakeObject(Select_vals: vals)
		Assert(ob.Eval(m) is: false)

		vals = [[city: #(value: "", value2: "", operation: ""),
			check:, condition_field: "city"]]
		ob.VirtualListColModel_selectMgr = FakeObject(Select_vals: vals)
		Assert(ob.Eval(m) is: false)

		vals = [[city: #(value: "Saskatoon", value2: "", operation: "equals"),
			check:, condition_field: "city"]]
		ob.VirtualListColModel_selectMgr = FakeObject(Select_vals: vals)
		Assert(ob.Eval(m))

		vals = [[city: #(value: "Saskatoon", value2: "", operation: "equals"),
			check: false, condition_field: "city"]]
		ob.VirtualListColModel_selectMgr = FakeObject(Select_vals: vals)
		Assert(ob.Eval(m) is: false)

		vals = [[city: #(value: "Saskatoon", value2: "", operation: "equals"),
			check: false, condition_field: "city"],
			[date: #(value: #20200305, value2: "", operation: "equals"),
			check: false, condition_field: "date"]]
		ob.VirtualListColModel_selectMgr = FakeObject(Select_vals: vals)
		Assert(ob.Eval(m) is: false)

		vals = [[city: #(value: "Saskatoon", value2: "", operation: "equals"),
			check: false, condition_field: "city"],
			[date: #(value: #20200305, value2: "", operation: "equals"),
			check: true, condition_field: "date"]]
		ob.VirtualListColModel_selectMgr = FakeObject(Select_vals: vals)
		Assert(ob.Eval(m))
		}

	Test_defaultCols()
		{
		fn = VirtualListColModel.VirtualListColModel_defaultCols
		t = .MakeTable('(num, num2, date, date2, Time, Time2) key (num) key (num2)')
		cols = Object()
		result = cols.Eval(fn, t)
		Assert(result[..2] equalsSet: #(num, num2))
		Assert(result[-2..] equalsSet: #(time, time2))
		Assert(result[2..4] equalsSet: #(date, date2))

		t = .MakeTable('(num, date) key ()')
		cols = Object()
		result = cols.Eval(fn, t)
		Assert(result equalsSet: #(num, date))

		t = .MakeTable('(num, num2, date) key (num, num2)')
		cols = Object()
		result = cols.Eval(fn, t)
		Assert(result[0] is: 'num')
		Assert(result[1..] equalsSet: #(num2, date))

		t = .MakeTable('(date, num, num2) key (num, num2) key (num, date)')
		cols = Object()
		result = cols.Eval(fn, t)
		Assert(result[0] is: 'num')
		Assert(result[1..] equalsSet: #(num2, date))

		t = .MakeTable('(date, num2, num) key (num, num2) key (num2, date)')
		cols = Object()
		result = cols.Eval(fn, t)
		Assert(result[..2] equalsSet: #('num', 'num2'))
		Assert(result[2] is: 'date')
		}

	Teardown()
		{
		QueryDo('delete usercolumns
			where usercolumns_title is ' $ Display(.columnsSaveName))
		QueryDo('delete userselects
			where userselect_title is ' $ Display(.columnsSaveName))
		QueryDo('delete usercolumns
			where usercolumns_title is "test_custom_key"')
		super.Teardown()
		}
	}
