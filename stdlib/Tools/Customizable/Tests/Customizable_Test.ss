// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.RequiresConfiglib(this)
		.TearDownIfTablesNotExist('customizable_fields')
		.TearDownIfTablesNotExist('customizable')
		if not TableExists?('customizable')
			Customizable.EnsureSaveInTable()
		if not TableExists?('customizable_fields')
			Customizable.EnsureTable()
		.fields = Object()
		_customTableNameVal = ''
		}

	RequiresConfiglib(test)
		{
		if not Libraries().Has?('configlib')
			test.AddTeardown(function () { ServerEval("Unuse", "configlib") })
		if not TableExists?('configlib')
			test.AddTeardown(function () { Database("destroy configlib") })
		}

	getTestCustomizableClass(table = 'test_table', name = false, user = '')
		{
		.AddTeardown(function() { try Database('destroy test_customizable') })
		c = Customizable
			{
			Customizable_savein: 'test_customizable'
			Customizable_logError(@unused) {}
			Customizable_resetCustomTableDataSources() { }
			Customizable_customTableName() { return _customTableNameVal isnt ''
				? _customTableNameVal : 'custom_table_' $ Display(Timestamp()).Tr('#.') }
			}
		testCl = c(table, name, :user)
		if user is ""
			QueryDo('delete test_customizable')	// ensure no leftovers from previous tests
		return testCl
		}

	getSelectFields(c)
		{
		return SelectFields(QueryColumns(c.GetTable()))
		}

	Test_CustomFields()
		{
		table = .MakeTable('(k, custom, my_custom, customer) key(k)')
		c = .getTestCustomizableClass(table)
		Assert(c.CustomFields() is: #())
		table = .MakeTable('(k, custom, my_custom, customer, custom_1, custom_123)
			key(k)')
		c = .getTestCustomizableClass(table)
		Assert(c.CustomFields() is: #(custom_1, custom_123))
		}

	Test_CreateField()
		{
		table = .MakeTable('(k) key(k)')
		c = .getTestCustomizableClass(table)
		fld = c.CreateField('My Field', 'Text, single line', SelectFields())
		.fields.Add(fld)
		Assert(Prompt(fld) is: 'My Field')
		}

	Test_Layout()
		{
		c = .getTestCustomizableClass('test_layout', 'test_layout')
		sf = .getSelectFields(c)
		Assert(c.Layout(sf) is: '')
		s = 'hello world'
		c.SaveLayout(s, sf)
		Assert(c.Layout(sf) is: s)

		user = .TempName()
		c2 = .getTestCustomizableClass('test_layout', 'test_layout', :user)
		Assert(c2.Layout(sf) is: s)

		s2 = 'hello world 2'
		c2.SaveLayout(s2, sf)
		Assert(c2.Layout(sf) is: s2)

		c3 = .getTestCustomizableClass('test_layout', 'test_layout', user: .TempName())
		Assert(c3.Layout(sf) is: s)
		}

	Test_Layout_tab()
		{
		c = .getTestCustomizableClass('test_layout', 'test_layout')
		sf = .getSelectFields(c)
		Assert(c.Layout(sf, 'test_tab') is: '')
		s = 'hello tab'
		c.SaveLayout(s, sf, tab: 'test_tab')
		Assert(c.Layout(sf, 'test_tab') is: s)
		}

	Test_build_field_name()
		{
		table = .MakeTable('(group, lib_committed, lib_modified, name,
			num, parent, text)
			key (name,group)
			key (num)
			index (parent,name)')
		fields = #(
			'Field_custom_000001': 'My Test Field'
			'Field_custom_000002': 'My Test Field2'
			'Field_custom_000002_param': 'My Test Param Field2')
		for field in fields.Members()
			QueryOutput(table, Record(name: field,
				text: 'Field_string\n' $
					'\t{\n' $
					'\tPrompt: ' $ Display(fields[field]) $ '\n' $
					'\t}'
				num: QueryMax(table, 'num', 0) + 1
				group: -1, parent: 0))

		Assert(Customizable.Customizable_build_field_name(table)
			greaterThanOrEqualTo: 'custom_000003')
		}

	Test_promptToField()
		{
		table = .MakeTable('(key, field) key(key)')
		.MakeLibraryRecord([name: "Table_" $ table,text: `class { Name: ` $ table $ ` }`])
		c = .getTestCustomizableClass(table)
		Assert(c.PromptToField("Not Here") is: false)
		fld2 = c.CreateField('My Field', 'Text, single line', SelectFields())
		.fields.Add(fld2)
		Assert(c.PromptToField("My Field") is: fld2)
		}

	Test_SetRecordDefaultValues()
		{
		field1 = .TempName().Lower()
		table = .MakeTable('(k) key(k)')

		rule1 = .TempName().Lower()
		.MakeLibraryRecord([name: 'Rule_' $ rule1, text: 'function()
			{
			.field_a
			.field_b
			return "RULE1 KICKED IN"
			}'])

		rule2 = .TempName().Lower()
		.MakeLibraryRecord([name: 'Rule_' $ rule2, text: 'function()
			{
			.' $ rule1 $ '\r\n' $
			'.' $ field1 $ '
			return "RULE2 KICKED IN"
			}'])

		.AddTeardown(
			{ QueryDo('delete customizable_fields
			where custfield_field in (' $ Display(field1) $ ', ' $
				Display(rule1) $ ', ' $ Display(rule2) $ ')') })
		QueryOutput('customizable_fields', [custfield_num: Timestamp(),
			custfield_name: table,
			custfield_field: rule1,
			custfield_default_value: "TEST1"])
		QueryOutput('customizable_fields', [custfield_num: Timestamp(),
			custfield_name: table,
			custfield_field: rule2,
			custfield_default_value: "TEST2"])
		QueryOutput('customizable_fields', [custfield_num: Timestamp(),
			custfield_name: table,
			custfield_field: field1,
			custfield_default_value: "TEST3"])
		rec = Record()
		Customizable.SetRecordDefaultValues(table, rec)
		Assert(rec[rule1] is: "TEST1")
		Assert(rec[rule2] is: "TEST2")
		Assert(rec[field1] is: "TEST3")
		rec.field_a = 'HELLO'
		Assert(rec[rule2] is: "RULE2 KICKED IN")
		Assert(rec[rule1] is: "RULE1 KICKED IN")
		}

	Test_handleProtectField()
		{
		m = Customizable.Customizable_handleProtectField
		rec = []
		protectField = 'testProtectField'
		rec[protectField] = 0
		values = #(testField: 'testDefaultValue')
		orig = rec.Copy()
		msg = "ERROR: invalid return type from protect rule"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		m(rec, protectField, values)
		Assert(orig is: rec)

		rec[protectField] = true
		orig = rec.Copy()
		m(rec, protectField, values)
		Assert(orig is: rec)

		rec[protectField] = "test protection message"
		orig = rec.Copy()
		m(rec, protectField, values)
		Assert(orig is: rec)

		rec[protectField] = ""
		m(rec, protectField, values)
		Assert(rec.testField is: 'testDefaultValue')

		rec[protectField] = #(testField:, testField2:)
		m(rec, protectField, values)
		Assert(rec.testField is: '')

		rec[protectField] = #('allbut', testField:, testField2:)
		m(rec, protectField, values)
		Assert(rec.testField is: 'testDefaultValue')

		rec[protectField] = ""
		rec['testField__protect'] = true
		m(rec, protectField, values)
		Assert(rec.testField is: '')
		}

	Test_DeleteField_setInternal()
		{
		fn = Customizable.Customizable_setInternal

		text = "Field_num { ExcludeSelect: true }"
		Assert(fn(text) is: 'Field_num {\n\tInternal: true\n\tExcludeSelect: true }')

		text = "Field_num\n\t{\n\tExcludeSelect: true\n\t}"
		Assert(fn(text) is:
			'Field_num\n\t{\n\tInternal: true\n\tExcludeSelect: true\n\t}')

		text = "Field_num\n\t{\n\tInternal: true\n\tPrompt: 'Test'\n\t}"
		Assert(fn(text) is: false)
		}

	Test_fieldControlAndFormat()
		{
		build = Customizable.Customizable_fieldControlAndFormat
		Assert(build(#(control: #(mask: '###.##')))
			is: '\tControl_mask: "###.##"\n')
		Assert(build(#(format: #(mask: '###.##')))
			is: '\tFormat_mask: "###.##"\n')
		str = build(#(format: #(mask: '###.##'), control: #(mask: '###.##')))
		Assert(str has: '\tFormat_mask: "###.##"\n')
		Assert(str has: '\tControl_mask: "###.##"\n')
		}

	Test_TabCustom?()
		{
		cl = .getTestCustomizableClass()
		Assert(cl.TabCustom?('barney') is: false)
		QueryOutput('test_customizable',
			[name: 'test_table', tab: 'barney', custom?: true])
		Assert(cl.TabCustom?('barney'))
		}

	Test_SaveLayout()
		{
		cl = .getTestCustomizableClass()
		sf = SelectFields()
		sf.AddField('test_0001', 'Test')
		sf.AddField('test_0002', 'Test2')
		sf.AddField('test_0003', 'Test3')
		sf.AddField('test_0004', 'Test4')
		cl.SaveLayout('Test\r\nTest2', sf, 'general_fred')
		layoutRec = Query1('test_customizable', name: 'test_table', tab: 'general_fred')
		Assert(layoutRec.layout is: #(Form, (test_0001, group: 0), nl,
			(test_0002, group: 0)))
		Assert(layoutRec.custom? is: false)
		Assert(layoutRec.hidden? is: false)

		cl.SaveLayout('Test\r\nTest3', sf, 'custom_fred', custom:)
		layoutRec = Query1('test_customizable', name: 'test_table', tab: 'custom_fred')
		Assert(layoutRec.layout is: #(Form, (test_0001, group: 0), nl,
			(test_0003, group: 0)))
		Assert(layoutRec.custom?)
		Assert(layoutRec.hidden? is: false)

		// call without the custom argument, should retain setting from table
		cl.SaveLayout('Test\r\nTest4', sf, 'custom_fred')
		layoutRec = Query1('test_customizable', name: 'test_table', tab: 'custom_fred')
		Assert(layoutRec.layout is: #(Form, (test_0001, group: 0), nl,
			(test_0004, group: 0)))
		Assert(layoutRec.custom?)
		Assert(layoutRec.hidden? is: false)
		}

	Test_ListHiddenLayouts()
		{
		cl = .getTestCustomizableClass()
		QueryOutput('test_customizable',
			[name: 'test_table', tab: 'barney', hidden?: true])
		QueryOutput('test_customizable',
			[name: 'test_table', tab: 'fred', hidden?: false])
		QueryOutput('test_customizable',
			[name: 'test_table', tab: 'wilma', hidden?: true])
		QueryOutput('test_customizable',
			[name: 'test_table', tab: 'betty', hidden?: true])
		Assert(cl.ListHiddenLayouts() is: #(barney betty wilma))
		}

	Test_GetNonPermissableFields()
		{
		fn = Customizable.GetNonPermissableFields
		Assert(fn(false) is: #())
		Assert(fn('invalid_table') is: #())
		Assert(fn('stdlib') is: #())
		}

	Test_GetEditableCustomFields()
		{
		table = .MakeTable('(k) key(k)')
		.MakeLibraryRecord([name: "Table_" $ table,text: `class { Name: ` $ table $ ` }`])
		testTable = Test.TempTableName()
		Assert(Customizable.GetEditableCustomFields(table, testTable) is: #())
		custfield = .makeField(table)
		Assert(Customizable.GetEditableCustomFields(table, testTable) is:
			Object(custfield))

		hidden = .makeField(table)
		.MakeCustomizeField(testTable, hidden, extrafields: #(custfield_hidden: true))
		Assert(Customizable.GetEditableCustomFields(table, testTable) is:
			Object(custfield))

		readonly = .makeField(table)
		.MakeCustomizeField(testTable, readonly, extrafields: #(custfield_readonly: true))
		Assert(Customizable.GetEditableCustomFields(table, testTable) is:
			Object(custfield))

		invisible = .makeField(table)
		.MakeCustomizeField(testTable, invisible, extrafields: #(custfield_readonly: true,
			custfield_hidden: true))
		Assert(Customizable.GetEditableCustomFields(table, testTable) is:
			Object(custfield))

		mandatory = .makeField(table)
		.MakeCustomizeField(testTable, mandatory,
			extrafields: #(custfield_mandatory: true))
		Assert(Customizable.GetEditableCustomFields(table, testTable)
			is: Object(custfield, mandatory))
		}

	makeField(table)
		{
		c = .getTestCustomizableClass(table)
		fld = c.CreateField('My Field', 'Text, single line', SelectFields())
		.fields.Add(fld)
		return fld
		}

	Test_formatDeletedCustomFields()
		{
		mock = Mock(Customizable)
		mock.When.formatDeletedCustomFields([anyArgs:]).CallThrough()
		mock.When.DeletedField?([anyArgs:]).Return(false)
		mock.When.DeletedField?('d1').Return(true)
		form = #(Form, (f1, group: 0), (f2, group: 1), nl
			(d1, group: 0), (d2, group: 1)).DeepCopy()
		mock.formatDeletedCustomFields(form)
		Assert(form is: #(Form, (f1, group: 0), (f2, group: 1), nl,
			(Static, 'd1', group: 0), (d2, group: 1)))
		}

	Test_handleRemovedFields()
		{
		mock = Mock(Customizable)
		mock.When.handleRemovedFields([anyArgs:]).CallThrough()
		mock.When.programmerError([anyArgs:]).Return('')
		mock.Customizable_availableFields = false
		form = #(Form, (f1, group: 0), (Static, 'test'), nl, (f2, group: 0))
		f = form.DeepCopy()
		mock.handleRemovedFields(f)
		Assert(f is: form)

		mock.Customizable_availableFields = Object('f1', 'f2')
		mock.handleRemovedFields(f)
		Assert(f is: form)
		mock.Customizable_availableFields = Object('f2')
		mock.handleRemovedFields(f)
		Assert(f
			is: #(Form, (Static, '???', group: 0), (Static, 'test'), nl, (f2, group: 0)))

		tbl = .MakeTable('(a,b) key(a)')
		custOb = .MakeCustomField(tbl, 'Text, single line')
		form = Object('Form', Object(custOb.field))
		f = form.DeepCopy()
		mock.handleRemovedFields(f)
		Assert(f is: Object('Form', Object('Static', custOb.prompt)))
		}

	Test_removeEmptyStatic()
		{
		fn = Customizable.Customizable_removeEmptyStatic
		form = #(Form, (Static,  ' '), (f1, group: 0), (Static, '  '), (f2, group: 1), nl,
			(Static, ' '), (f3, group: 0), (Static, 'text'), (f4, group: 1)).DeepCopy()
		fn(form)
		Assert(form is: #(Form, (f1, group: 0), (f2, group: 1), nl,
			(f3, group: 0), (Static, 'text'), (f4, group: 1)))
		}

	Test_CustomTableTabs()
		{
		custTable = .TempName()
		_customTableNameVal = custTable
		tab = 'FakeTab_' $ custTable
		name = .MakeTable('(key_num, custom, my_custom, customer) key(key_num)')
		c = .getTestCustomizableClass(name)

		// Tab exists, not a Custom Table Tab
		Assert(c.CustomTab(tab) is: false)
		Assert(c.CustomTableTab?(tab) is: false)
		Assert(c.CustomTableName(name) is: name)

		c.SaveTab(tab, true)
		noteField = QueryColumns(custTable).FindOne({ it.Prefix?('custom_') })
		ctrlSuffix = c.Customizable_customTableLinkField('key_num', custTable)
		.AddTeardown({ .customTableTabTeardown(ctrlSuffix, noteField, custTable) })

		Assert(noteField isnt: false)
		Assert(Query1('configlib', name: 'Rule_' $ ctrlSuffix) isnt: false)
		Assert(Query1('configlib', name: 'Field_' $ ctrlSuffix) isnt: false)
		Assert(c.CustomTab(tab) isnt: false)
		Assert(c.CustomTableTab?(tab))
		Assert(c.CustomTableName(name) is: name)

		ctrl = c.Table(tab)
		Assert(ctrl[0] is: #LineItemControl)
		Assert(ctrl.name is: ctrlSuffix)
		Assert(ctrl.linkField is: #custtable_FK)
		Assert(ctrl.headerFields is: #(key_num_new))
		Assert(ctrl.protectField is: #custtable_protect)
		Assert(ctrl.columns equalsSet: [#bizuser_user_cur, #custtable_num_new, noteField])
		Assert(ctrl.query
			is: custTable $ ' rename custtable_num to custtable_num_new,' $
				'\r\n\t\t\t\tbizuser_user to bizuser_user_cur\r\n\t\t\t\t' $
				'sort custtable_num_new')
		_customTableNameVal = ''
		}

	customTableTabTeardown(ctrlSuffix, noteField, custTable)
		{
		QueryDo('delete configlib where name.Suffix?("' $ ctrlSuffix $ '")')
		QueryDo('delete configlib where name is "Field_' $ noteField $ '"')
		Database('drop ' $ custTable)
		}

	Test_getModifiedForm()
		{
		fn = Customizable.Customizable_getModifiedForm
		field = 'fake_field'
		prompt = 'aaa'
		layout = ''
		Assert(fn(layout, field, prompt) is: '')

		layout = #(Form, (Static, aaa))
		Assert(fn(layout, field, prompt) is: #(Form, (fake_field)))

		layout = #(Form, (Static, 'aaa aaa aaa\r\naaa'))
		Assert(fn(layout, field, prompt) is: #(Form, (fake_field),
			(Static, ' aaa aaa\r\naaa')))

		layout = #(Form, (Static, 'text before the field aaa text after the field'))
		Assert(fn(layout, field, prompt) is: #(Form, (Static, 'text before the field '),
			(fake_field), (Static, ' text after the field')))

		layout = #(Form, (field1, group: 0), (Static ' '), (field2, group: 1), nl
			(field3, group: 0), (Static ' aaa '), (field4, group: 1))
		Assert(fn(layout, field, prompt) is:
			#(Form, (field1, group: 0), (Static ' '), (field2, group: 1), nl
			(field3, group: 0), (Static ' '), (fake_field), (Static ' '),
				(field4, group: 1)))

		layout = #(Form, (field1, group: 0), (Static ' '), (field2, group: 1), nl
			(field3, group: 0), (Static ' aaaaaa '), (field4, group: 1))
		Assert(fn(layout, field, prompt) is:
			#(Form, (field1, group: 0), (Static ' '), (field2, group: 1), nl
			(field3, group: 0), (Static ' '), (fake_field), (Static 'aaa '),
				(field4, group: 1)))
		}

	Teardown()
		{
		for fld in .fields
			{
			QueryDo("delete configlib where name = " $ Display('Field_' $ fld))
			Unload('Field_' $ fld)
			}

		super.Teardown()
		}
	}
