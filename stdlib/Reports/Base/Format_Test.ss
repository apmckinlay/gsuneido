// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_AddAccessReportPoint()
		{
		data = Object()
		Format.AddAccessReportPoint(data, 'test', 'TESTREPORT',
			Object(test_params: 'test'))
		Assert(data['_access_test'] is:
			Object(control: 'ReportGoTo', report: 'TESTREPORT',
				params: Object(test_params: 'test')))
		}

	Test_AddAccessPoint()
		{
		data = Object()
		Format.AddAccessPoint(data, 'test', 'TestAccess', 'test_num', 'test_value')
		Assert(data['_access_test'] is:
			Object(control: 'AccessGoTo', access: 'TestAccess',
				goto_field: 'test_num', goto_value: 'test_value'))

		data = Object()
		func = function (field_value /*unused*/) { return 'test_func' }
		Format.AddAccessPoint(data, 'test', func, 'test_num', 'test_value')
		Assert(data['_access_test'] is:
			Object(control: 'AccessGoTo', access: 'test_func',
				goto_field: 'test_num', goto_value: 'test_value'))
		}

	Test_AddAccessField()
		{
		data = Object(test_value_field: 'test_value2')
		Format.AddAccessField(data, 'test', Object('AccessGoTo',
			access: 'TestAccess'
			goto_field: 'test_num', goto_value: 'test_value_field'))
		Assert(data['_access_test'] is:
			Object(control: 'AccessGoTo', access: 'TestAccess',
				goto_field: 'test_num', goto_value: 'test_value2'))

		text = '
class
	{
	Prompt: "TEST Prompt"
	Control: (Id access: "TestAccess")
	}'
		.MakeLibraryRecord(Record(name: 'Field_test_field_num', :text))
		data = Object(test_field_num: 'test_num_value')
		Format.AddAccessField(data, 'test_field_num', 'test_field_num')
		Assert(data['_access_test_field_num'] is:
			Object(control: 'AccessGoTo', access: 'TestAccess',
				goto_field: 'test_field_num', goto_value: 'test_num_value'))

		}

	Test_CSVExport()
		{
		cl = Format
			{ Export: true }
		Assert('""' is cl.CSVExportString(''))
		Assert('"normal text"' is cl.CSVExportString('normal text'))

		str = 'a, b, c'
		Assert('"a, b, c"' is cl.CSVExportString(str))

		str = 'a, b\r\n, c'
		Assert('"a, b , c"' is cl.CSVExportString(str))

		str = `a, "b, c`
		Assert(`"a, ""b, c"` is cl.CSVExportString(str))

		str = `a, "", c`
		Assert(`"a, """", c"` is cl.CSVExportString(str))

		str = `'a', 'b', 'c'`
		Assert(`"'a', 'b', 'c'"` is cl.CSVExportString(str))

		str = `'a', 'b', 'c'`
		Assert(`"'a', 'b', 'c'"` is cl.CSVExportString(str))

		str = `"a", "b", "c"`
		Assert(`"""a"", ""b"", ""c"""` is cl.CSVExportString(str))

		Assert(`"123"` is cl.CSVExportString(123))

		Assert(`"#20160413"` is cl.CSVExportString(#20160413))

		Assert(`"true"` is cl.CSVExportString(true))
		}
	}
