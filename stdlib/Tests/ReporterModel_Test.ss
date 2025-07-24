// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		Customizable_Test.RequiresConfiglib(this)
		.TearDownIfTablesNotExist('customizable')
		if not TableExists?('customizable')
			Customizable.EnsureSaveInTable()
		}
	Test_getCalcField()
		{
		data = Record()
		prompt = ''
		Assert(ReporterModel.ReporterModel_getCalcField(prompt, data) is: false)

		data = Record(formulas: #([calc: "Test One"]))
		prompt = 'Test One'
		Assert(ReporterModel.ReporterModel_getCalcField(prompt, data)
			is: [calc: "Test One"])

		data = Record(formulas: #([calc: "Test One"]), summarize_field0: "Test One",
			summarize_func0: "maximum", summarize_func_cols: #("Maximum Test One"))
		prompt = "Maximum Hours Reading"
		Assert(ReporterModel.ReporterModel_getCalcField(prompt, data) is: false)

		data = Record(formulas: #([calc: "Test One"]), summarize_field0: "Test One",
			summarize_func0: "maximum", summarize_func_cols: #("Maximum Test One"))
		prompt = "Maximum Test One"
		Assert(ReporterModel.ReporterModel_getCalcField(prompt, data)
			is: [calc: "Test One"])

		data = Record(formulas: #([calc: "Hours Reading"], [calc: "Test One"]),
			summarize_field1: "Hours Reading", summarize_field0: "Test One",
			summarize_func1: "total", summarize_func0: "maximum",
			summarize_func_cols: #("Maximum Test One", "Total Hours Reading"))
		prompt = "Maximum Test One"
		Assert(ReporterModel.ReporterModel_getCalcField(prompt, data)
			is: [calc: "Test One"])
		}

	Test_get_valid_cols()
		{
		mock = Mock()
		mock.ReporterModel_design_cols = #('aaa')
		mock.ReporterModel_data = #(columns: #(#(text: 'aaa'), #(text: 'bbb'),
			#(text: 'ccc')))
		cols = mock.Eval(ReporterModel.ReporterModel_valid_cols)
		Assert(cols is: #(#(text: 'aaa')))

		mock.ReporterModel_design_cols = #('aaa', 'ddd')
		mock.ReporterModel_data = #(columns: #(#(text: 'aaa'), #(text: 'bbb'),
			#(text: 'ccc')))
		cols = mock.Eval(ReporterModel.ReporterModel_valid_cols)
		Assert(cols is: #(#(text: 'aaa')))

		mock.ReporterModel_design_cols = #('aaa', 'bbb')
		mock.ReporterModel_data = #(columns: #(#(text: 'aaa'), #(text: 'bbb'),
			#(text: 'ccc')))
		cols = mock.Eval(ReporterModel.ReporterModel_valid_cols)
		Assert(cols is: #(#(text: 'aaa'), #(text: 'bbb')))
		}

	Test_design_cols_with_summarize()
		{
		rpt = ReporterModel
			{
			ReporterModel_summarize_cols: #(col1, col2, col3)
			ForEachCalc(block)
				{
				for col in #(col2, col4)
					block(0, col)
				}
			}
		Assert(rpt.ReporterModel_design_cols_with_summarize(),
			is: #(col1, col2, col3, col4))
		}

	Test_printParams()
		{
		Assert(.printParams(false, Record()) is: #())
		fake_sf = FakeObject(PromptToField: function (prompt)
				{ return #(One: one, Two: two)[prompt] })
		data = Record(checkbox0: false, fieldlist0: 'One', oplist0: 'equals',
			val0: 1, print0: true)
		Assert(.printParams(false, data) is: #())
		data.checkbox0 = true
		Assert(.printParams(fake_sf, data) is: #(one))
		data.checkbox1 = true
		data.fieldlist1 = 'Two'
		data.oplist1 = 'equals'
		data.val1 = 2
		data.print1 = false
		Assert(.printParams(fake_sf, data) is: #(one))
		data.print1 = true
		Assert(.printParams(fake_sf, data) equalsSet: #(one, two))
		}

	printParams(sf, data)
		{
		return ReporterModel.ReporterModel_printParams(
			ReporterModel.ReporterModel_params(sf, data))
		}

	Test_format_params()
		{
		format_params = ReporterModel.ReporterModel_format_params
		params = Object()
		field = 'field'
		format_params(params, field, 'greater than or equal to', 1, 1)
		Assert(params isSize: 1)
		Assert(params hasMember: field)

		format_params(params, field, 'less than or equal to', 10, 2)
		Assert(params isSize: 2)
		Assert(params hasMember: field)
		Assert(params hasMember: field $ '_' $ 2)

		format_params(params, field, 'contains', '5', 3)
		Assert(params isSize: 3)
		Assert(params hasMember: field)
		Assert(params hasMember: field $ '_' $ 2)
		Assert(params hasMember: field $ '_' $ 3)
		}

	Test_calc_dd()
		{
		rpt = ReporterModel
			{
			ReporterModel_data: #(formulas: #(
				[type: 'Checkmark', key: '0', calc: 'check', formula: 'true'],
				[type: 'Date and Time', key: '1', calc: 'date', formula: '2018-03-01'],
				[type: 'Number, 3 decimals', key: '2', calc: 'number', formula: '15'],
				[type: 'Text, multi line', key: '3', calc: 'text',
					formula: 'Hello World']))
			}
		Assert(rpt.ReporterModel_calc_dd()
			is: Object(
				calc0: Field_boolean_checkmark,
				calc1: Field_date,
				calc2: Field_number))
		}

	menu_params_tests: #(
		#(params: #(
			paramselects:
				[PrintLines: false,
				printParams: #(
					#(paramField: "calc20190328144357356",
						paramPrompt: "Menu Date",
						paramFormat: 'Field_date { Format: #("ShortDate") }')),
				ReportDestination: "preview",
				calc20190328144357356: #(value: #20190301, value2: "",
					operation: "less than or equal to")],
			extends: "",
			summarize_fields: #("total_etaorder_total_amount",
				"max_calc20190328144357356"))
		result1: ' where calc20190328144357356 <= #20190301',
		result2: ''),

		#(params: #(
			paramselects:
				[PrintLines: false,
				printParams: #(
					#(paramField: "calc20190328141522786",
						paramPrompt: "T1",
						paramFormat: 'Field_date { Format: #("ShortDate") }')),
				ReportDestination: "preview",
				calc20190328141522786: #(value: #20190401, value2: "",
					operation: "greater than or equal to")],
			extends: "",
			summarize_fields: #("total_etaorder_total_amount"))
		result1: ' where calc20190328141522786 >= #20190401',
		result2: ''),

		#(params: #(
			paramselects:
				[PrintLines: false,
				printParams: #(
					#(paramField: "calc20190328142051336",
						paramPrompt: "Bob",
						paramFormat:
							'Field_number { Format: #("Number", mask: "-###,###,###") }')
					),
				ReportDestination: "preview",
				calc20190328142051336: #(value: 3, value2: "",
					operation: "less than or equal to")],
			extends: "extend calc20190328142051336 = Reporter000020()",
			summarize_fields: #())
		result1: '',
		result2: ' where calc20190328142051336 <= 3'),
		)
	Test_menu_params_calc_where()
		{
		fn1 = ReporterModel.ReporterModel_menu_params_calc_where_before_summarize
		fn2 = ReporterModel.ReporterModel_menu_params_calc_where

		for test in .menu_params_tests
			{
			params = test.params
			Assert(fn1(params.paramselects, params.extends, params.summarize_fields)
				is: test.result1)
			Assert(fn2(params.paramselects, params.extends, params.summarize_fields)
				is: test.result2)
			}
		}

	Test_convertParamSelectsField()
		{
		_testMenuParams = Object()
		repModel = ReporterModel
			{
			Menu_params_fields()
				{
				return _testMenuParams
				}
			}
		fn = repModel.ReporterModel_convertParamSelectsField

		paramselects = Record()
		fn('anything', paramselects)
		Assert(paramselects isSize: 0)

		paramselects.testBoolField_param = #(value:, value2: "", operation: "equals")
		fn('testBoolField', paramselects)
		Assert(paramselects
			is: Record(testBoolField_param: #(value:, value2: "", operation: "equals")))

		_testMenuParams.Add("testBoolField?")
		fn('testBoolField?', paramselects)
		Assert(paramselects
			is: Record(
					testBoolField_param: #(value:, value2: "", operation: "equals"),
					"testBoolField?_param": #(value:, value2: "", operation: "equals")))

		paramselects = Record(testField: #(value:, value2: "", operation: "equals"))
		fn('testField', paramselects)
		Assert(paramselects
			is: Record(testField: #(value:, value2: "", operation: "equals")))
		}

	Test_preProcessMenuParams()
		{
		repModel = ReporterModel
			{
			Menu_params_fields()
				{
				return Object('test_name', 'test2_name')
				}
			ReporterModel_sf: class
				{
				GetJoinNumField(field) { return field.Replace('name', 'num') }
				}
			}
		foreigNumsSpy = .SpyOn(GetForeignNumsFromNameAbbrevFilter)
		fn = repModel.ReporterModel_preProcessMenuParams

		paramselects = Record()
		result = fn(paramselects)
		Assert(paramselects isSize: 0)
		Assert(result.paramselects is: paramselects)
		Assert(result.options is: #())
		Assert(result.fields is: #('test_name', 'test2_name'))
		Assert(foreigNumsSpy.CallLogs() isSize: 0)

		foreigNumsSpy.Return(false)
		paramselects = Record(
			test_name_param: [operation: 'contains', value: 'testing', value2: ''],
			test2_name_param: [operation: '', value: '', value2: ''])
		result = fn(paramselects)

		Assert(foreigNumsSpy.CallLogs() isSize: 1)
		Assert(result.paramselects is: paramselects)
		Assert(result.options is: #('test_name'))
		Assert(result.fields is: #('test_name', 'test2_name'))

		foreigNumsSpy.ClearAndReturn(#(numField: 'test_num', nums: ()))
		result = fn(paramselects)
		Assert(foreigNumsSpy.CallLogs() isSize: 2)
		expected = Record(
			test_name_param: [operation: 'contains', value: 'testing', value2: ''],
			test2_name_param: [operation: '', value: '', value2: ''],
			test_num: [operation: 'less than', value: '', value2: ''])
		Assert(result.paramselects is: expected)
		Assert(result.options is: #('test_num'))
		Assert(result.fields is: #('test_num', 'test2_name'))

		foreigNumsSpy.ClearAndReturn(
			#(numField: 'test_num', nums: ()), #(numField: 'test2_num', nums: (1,2,3)))
		paramselects = Record(
			test_name_param: [operation: 'contains', value: 'testing', value2: ''],
			test2_name_param: [operation: 'matches', value: 'test', value2: ''])
		result = fn(paramselects)
		Assert(foreigNumsSpy.CallLogs() isSize: 4)
		expected = Record(
			test_name_param: [operation: 'contains', value: 'testing', value2: ''],
			test_num: [operation: 'less than', value: '', value2: ''],
			test2_name_param: [operation: 'matches', value: 'test', value2: '']
			test2_num: [operation: 'in list', value: #(1,2,3) value2: ''])
		Assert(result.paramselects is: expected)
		Assert(result.options is: #('test_num', 'test2_num'))
		Assert(result.fields is: #('test_num', 'test2_num'))

		repModel = ReporterModel
			{
			Menu_params_fields()
				{
				return Object('test_name', 'test_num')
				}
			ReporterModel_sf: class
				{
				GetJoinNumField(field) { return field.Replace('name', 'num') }
				}
			}
		fn = repModel.ReporterModel_preProcessMenuParams
		paramselects = Record(
			test_name_param: [operation: 'contains', value: 'testing', value2: ''],
			test_num_param: [operation: 'equals to', value: 1, value2: ''])
		result = fn(paramselects)
		Assert(foreigNumsSpy.CallLogs() isSize: 4)
		Assert(result.paramselects is: paramselects)
		Assert(result.options is: #('test_name', 'test_num'))
		Assert(result.fields is: #('test_name', 'test_num'))
		}

	Test_make_query()
		{
		tables = GetForeignNumsFromNameAbbrevFilter_Test.SetupForeignTable(this)
		table1 = tables.table1
		table2 = tables.table2
		if false isnt cache = Suneido.GetDefault(#ForeignKeyTables, false)
			cache.Reset()

		sourceName = .TempName()
		source = Object('Reporter', 'queries',
			name: sourceName,
			tables: Object(table2)
			query: table2 $ ' rename table1_num to table1_num_test')
		.SpyOn(GetCustomReportsSource).Return(source)

		origRpt = [Source: sourceName, report_name: .TempName(),
			columns: Object(
				#(width: 22, text: "Name 2")),
			select: [
				menu_option0:, print0:, checkbox0:, fieldlist0: "Name 2"
				menu_option1:, print1:, checkbox1:, fieldlist1: "Num 1 test Name"],
			coloptions: Object(
				"Name 2": [heading: "Name 2"]),
			heading1: "Sample Trucking Co.",
			nonsummarized_fields: Object(),
			summarize_func_cols: Object()]

		rpt = origRpt.DeepCopy()
		model = ReporterModel(rpt)
		sort = model.ReporterModel_make_sort()
		cols = model.ReporterModel_make_cols()
		// no select
		query = model.ReporterModel_make_query(cols, sort, #())
		Assert(query.Tr("\r\n", "")
			is: '/* tableHint: ' $ table2 $
				' */ (' $ table2 $ ' rename table1_num to table1_num_test)')

		// no name to num optimization
		paramselects = [
			table2_name_param: [operation: 'greater than', value: 'aaa'],
			table1_name_test_param: []]
		query = model.ReporterModel_make_query(cols, sort, paramselects)
		Assert(query.Tr("\r\n", "")
			is: '/* tableHint: ' $ table2 $ ' */ (' $
				table2 $ ' rename table1_num to table1_num_test)' $
				' where table2_name > "aaa"')

		// name to num optimization - multiple matches
		paramselects = [
			table1_name_test_param: [operation: 'greater than', value: 'AAA'],
			table2_name_param: []]
		query = model.ReporterModel_make_query(cols, sort, paramselects)
		Assert(query.Tr("\r\n", "")
			is: '/* tableHint: ' $ table2 $ ' */ (' $
				table2 $ ' rename table1_num to table1_num_test)' $
				' where table1_num_test in (2, 3)')

		// name to num optimization - no match
		paramselects = [
			table1_name_test_param: [operation: 'greater than', value: 'DDD'],
			table2_name_param: []]
		query = model.ReporterModel_make_query(cols, sort, paramselects)
		Assert(query.Tr("\r\n", "")
			is: '/* tableHint: ' $ table2 $
				' */ (' $ table2 $ ' rename table1_num to table1_num_test)' $
				' where table1_num_test < ""')

		// columns: "Num 1 test Abbrev", "Name 2"
		// select by 'Num 1 test Name'
		// should still have leftjoin by() after name to num optimization
		rpt = origRpt.DeepCopy()
		rpt.columns.Add(#(width: 22, text: "Num 1 test Abbrev"))
		rpt.coloptions.Add([heading: "Num 1 test Abbrev"], at: "Num 1 test Abbrev")
		model = ReporterModel(rpt)
		sort = model.ReporterModel_make_sort()
		cols = model.ReporterModel_make_cols()
		query = model.ReporterModel_make_query(cols, sort, paramselects)
		Assert(query.Tr("\r\n", "")
			is: '/* tableHint: ' $ table2 $ ' */ (' $ table2 $
				' rename table1_num to table1_num_test' $
				' leftjoin by(table1_num_test) (' $ table1 $ ' ' $
					'project table1_num, table1_name, table1_abbrev ' $
					'rename table1_num to table1_num_test ' $
					'rename table1_name to table1_name_test ' $
					'rename table1_abbrev to table1_abbrev_test) )' $
				' where table1_num_test < ""')

		// summarize by 'Num 1 test Abbrev'
		// columns: 'Count'
		// select by 'Num 1 test Name'
		// should still have leftjoin by() after name to num optimization
		paramselects = [
			table1_name_test_param: [operation: 'greater than', value: 'AAA'],
			table2_name_param: []]
		rpt = origRpt.DeepCopy()
		rpt.summarize_by = 'Num 1 test Abbrev'
		rpt.summarize_func_cols = #('Count')
		rpt.summarize_func0 = 'count'
		rpt.columns = [#(width: 22, text: "Count")]
		rpt.coloptions = ['Count': [heading: "Count"]]
		model = ReporterModel(rpt)
		sort = model.ReporterModel_make_sort()
		cols = model.ReporterModel_make_cols()
		query = model.ReporterModel_make_query(cols, sort, paramselects)
		Assert(query.Tr("\r\n", "")
			is:  '/* tableHint: ' $ table2 $ ' */ (' $ table2 $
				' rename table1_num to table1_num_test' $
				' leftjoin by(table1_num_test) (' $ table1 $ ' ' $
					'project table1_num, table1_name, table1_abbrev ' $
					'rename table1_num to table1_num_test ' $
					'rename table1_name to table1_name_test ' $
					'rename table1_abbrev to table1_abbrev_test) )' $
				' where table1_num_test in (2, 3)' $
				'summarize table1_abbrev_test,count')

		// summarize by 'Num 1 test Abbrev'
		// no summarize function
		// columns: 'Num 1 test Abbrev'
		// select by 'Num 1 test Name'
		// should have project with the suppress comment
		paramselects = [
			table1_name_test_param: [operation: 'greater than', value: 'AAA'],
			table2_name_param: []]
		rpt = origRpt.DeepCopy()
		rpt.summarize_by = 'Num 1 test Abbrev'
		rpt.columns = [#(width: 22, text: "Num 1 test Abbrev")]
		rpt.coloptions = ['Num 1 test Abbrev': [heading: "Num 1 test Abbrev"]]
		model = ReporterModel(rpt)
		sort = model.ReporterModel_make_sort()
		cols = model.ReporterModel_make_cols()
		query = model.ReporterModel_make_query(cols, sort, paramselects)
		Assert(query.Tr("\r\n", "")
			is:  '/* tableHint: ' $ table2 $ ' */ (' $ table2 $
				' rename table1_num to table1_num_test' $
				' leftjoin by(table1_num_test) (' $ table1 $ ' ' $
					'project table1_num, table1_name, table1_abbrev ' $
					'rename table1_num to table1_num_test ' $
					'rename table1_name to table1_name_test ' $
					'rename table1_abbrev to table1_abbrev_test) )' $
				' where table1_num_test in (2, 3)' $
				'project /*CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE*/ table1_abbrev_test')
		}


	Test_getParamField()
		{
		cl = ReporterModel
			{
			ReporterModel_sf: class
				{
				GetJoinNumField(field)
					{
					return _joins.GetDefault(field, false)
					}
				}
			ReporterModel_buildParam(field, promptInfo)
				{
				_params.Add([:field, :promptInfo])
				}
			}
		fn = cl.ReporterModel_getParamField

		_joins = Object()
		_params = Object()
		// normal field
		field = .MakeDatadict(Prompt: 'Test Prompt', SelectPrompt: 'Test Select Prompt')
		fn(field, field, 'Test Select Prompt')
		Assert(_params[0] is: [:field,
			promptInfo: [prefix: '', suffix: '', promptMethod: 'SelectPrompt',
				baseField: field]])

		// summarized field
		fn('min_' $ field, 'min_' $ field, 'Minimum Test Select Prompt')
		Assert(_params[1] is: [field: 'min_' $ field,
			promptInfo: [prefix: 'Minimum', suffix: '', promptMethod: 'SelectPrompt',
				baseField: field]])

		// _name field
		_joins[field $ '_name'] = field $ '_num'
		fn(field $ '_name', field $ '_name', 'Test Select Prompt Name')
		Assert(_params[2] is: [field: field $ '_name',
			promptInfo: [prefix: '', suffix: 'Name', promptMethod: 'SelectPrompt',
				baseField: field $ '_num']])

		// _abbrev field, which should be converted to _num field by .abbrev_to_num_field
		_joins[field $ '_abbrev'] = field $ '_num'
		fn(field $ '_num', field $ '_abbrev', 'Test Select Prompt Abbrev')
		Assert(_params[3] is: [field: field $ '_num',
			promptInfo: [prefix: '', suffix: 'Abbrev', promptMethod: 'SelectPrompt',
				baseField: field $ '_num']])

		// summarized info field
		infoField = 'reporter_info_' $ (base = .TempName()) $ '_Email'
		fn('max_' $ infoField, 'max_' $ infoField, 'Maximum Test Email')
		Assert(_params[4] is: [field: 'max_' $ infoField,
			promptInfo: [prefix: 'Maximum', suffix: 'Email',
				promptMethod: 'Reporter_extend_info.PrefixPrompt', baseField: base]])

		// custom field using prompt
		customField = .MakeDatadict(fieldName: 'custom_' $ .TempName(),
			Prompt: 'Test Custom Prompt', SelectPrompt: 'Test Custom Prompt ~ Test Table')
		fn('total_' $ customField, 'total_' $ customField, 'Total Test Custom Prompt')
		Assert(_params[5] is: [field: 'total_' $ customField,
			promptInfo: [prefix: 'Total', suffix: '', promptMethod: 'Prompt',
				baseField: customField]])

		// custom field using selectprompt
		fn('total_' $ customField, 'total_' $ customField,
			'Total Test Custom Prompt ~ Test Table')
		Assert(_params[6] is: [field: 'total_' $ customField,
			promptInfo: [prefix: 'Total', suffix: '', promptMethod: 'SelectPrompt',
				baseField: customField]])
		}
	}
