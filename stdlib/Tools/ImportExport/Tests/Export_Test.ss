// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_export_format_date()
		{
		// test date with DateTime
		test_class = Export
			{
			SetTable(name) { .table_name = name }
			Datadict(unused)
				{ return Object(Format: Object('DateTime')) }
			}
		date = Date()
		x = Record(text_val: 'text value', boolean_val: 'true', date_val: date,
			number_val: 10)
		fields = #('text_val', 'boolean_val', 'date_val', 'number_val')
		result = test_class.Export_format_dates(x, fields)
		Assert(result.text_val is: x.text_val)
		Assert(result.boolean_val is: x.boolean_val)
		Assert(result.date_val is: date.ShortDateTime())
		Assert(result.number_val is: x.number_val)

		// test date with ShortDate
		test_class = Export
			{
			SetTable(name) { .table_name = name }
			Datadict(unused)
				{ return Object(Format: Object('ShortDate')) }
			}
		date = Date()
		x = Record(text_val: 'text value', boolean_val: 'true', date_val: date,
			number_val: 10)
		fields = #('text_val', 'boolean_val', 'date_val', 'number_val')
		result = test_class.Export_format_dates(x, fields)
		Assert(result.text_val is: x.text_val)
		Assert(result.boolean_val is: x.boolean_val)
		Assert(result.date_val is: date.ShortDate())
		Assert(result.number_val is: x.number_val)

		// test not date
		date = 'test_string'
		x = Record(text_val: 'text value', boolean_val: 'true', date_val: date,
			number_val: 10)
		fields = #('text_val', 'boolean_val', 'date_val', 'number_val')
		result = Export.Export_format_dates(x, fields)
		Assert(result.text_val is: x.text_val)
		Assert(result.boolean_val is: x.boolean_val)
		Assert(result.date_val is: date)
		Assert(result.number_val is: x.number_val)
		}
	Test_convertNumFieldsToName()
		{
		.MakeFile() // to ensure teardown
		table = .MakeTable(
			'(testtable_num, testtable_name, testtable_date) key(testtable_num)')
		test_class = Export
			{
			HeaderType: 'Fields'
			table_name: ''
			OpenFile() { }
			SetTable(name) { .table_name = name }
			Datadict(field)
				{
				if not #(testtable_num testtable_num_2).Has?(field)
					return Object(Control: Object('Field' width: 15))
				return Object(Control: Object('Id', .table_name, field: 'testtable_num'))
				}
			}
		c = new test_class(table)
		c.SetTable('/* tableHint: */ ' $ table)

		r = Record(a: 1, b: 2, c: 3)

		// test empty field set, nothing should change
		c.Fields = #()
		c.Export_convertNumFieldsToName(r)
		Assert(r is: Record(a: 1, b: 2, c: 3))

		// test field set with no _nums, nothing should change
		c.Fields = #(a b c)
		c.Export_convertNumFieldsToName(r)
		Assert(r is: Record(a: 1, b: 2, c: 3))

		// test converting the num to name
		c.Fields = #(a b c testtable_num)
		QueryOutput(table, Record(testtable_num: 1, testtable_name: 'test1',
			testtable_date: Date()))
		r.testtable_num = 1
		c.Export_convertNumFieldsToName(r)
		Assert(r.testtable_num is: 'test1')

		// Header fields should now have _name not _num fields
		c.GetHead(#())
		Assert(c.Head has: 'testtable_name')
		Assert(c.Head hasnt: 'testtable_num')

		// test field with suffix after "_num"
		c.Fields = #(a b c testtable_num testtable_num_2)
		QueryOutput(table, Record(testtable_num: 2, testtable_name: 'test2',
			testtable_date: Date()))
		r.testtable_num = 1
		r.testtable_num_2 = 2
		c.Export_convertNumFieldsToName(r)
		Assert(r.testtable_num is: 'test1')
		Assert(r.testtable_num_2 is: 'test2')

		// test num field without related name field, should export the original num value
		c.Fields = #(a b c testtable_num testtable2_num)
		r.testtable_num = 1
		r.testtable2_num = 3
		c.Export_convertNumFieldsToName(r)
		Assert(r.testtable_num is: 'test1')
		Assert(r.testtable2_num is: 3)

		// Header fields should have _num not _name fields
		c.GetHead(#())
		Assert(c.Head has: 'testtable2_num')
		Assert(c.Head hasnt: 'testtable2_name')
		}
	Test_format_Encrypt()
		{
		test_class = Export
			{
			Fields: #('test_SSNSIN')
			Datadict(field /*unused*/)
				{
				return Object(Control: Object('SSNSIN'))
				}
			}
		r = Record(test_SSNSIN: '888-88-8888'.Xor(EncryptControlKey()))
		test_class.Export_format_Encrypt(r)
		Assert(r.test_SSNSIN is: '888-88-8888')

		test_class = Export
			{
			Fields: #('test_encrypt')
			Datadict(field /*unused*/)
				{
				return Object(Control: Object('Encrypt'))
				}
			}
		r = Record(test_encrypt: '123456'.Xor(EncryptControlKey()))
		test_class.Export_format_Encrypt(r)
		Assert(r.test_encrypt is: '123456')
		}
	}
