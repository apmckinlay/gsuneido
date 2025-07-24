// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.testCsv(#(Text 'test'), '"test"\r\n')

		.testCsv(#(Horz
			#(Text 'test')
			#(Text 'test2')), '"test","test2"\r\n')

		// testing ObjectFormat, QueryFormat, RowFormat, RowHeaderFormat
		data = Object(
			[table: "foo", column: "Table"],
			[table: "bar", column: "Column"])
		data.order = #(table)
		data.columns = #(table, column)
		.testCsv(Object('Object' data), '"Table","Column"\r\n' $
			'"foo","Table"\r\n' $
			'"bar","Column"\r\n')
		}

	testCsv(rpt, csv)
		{
		.f = FakeFile('')
		report = Report(rpt)
		report.ExportCSV('test_file', quiet?:, fileCl: .f)
		Assert(.f.Get() is: csv)
		}
	}