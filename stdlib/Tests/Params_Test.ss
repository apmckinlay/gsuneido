// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_devmode_reportname()
		{
		report = Object(name: 'name')
		Assert(Params.Params_devmode_reportname(report) is: report.name)

		report.devmode_name = 'devmode_name'
		Assert(Params.Params_devmode_reportname(report) is: report.devmode_name)
		}

	Test_checkAndResetParams()
		{
		record = Record()
		method = Params.Params_checkAndResetParams
		method(record)
		Assert(record hasntMember: 'EmailAttachments')

		record.EmailAttachments = #()
		method(record)
		Assert(record hasMember: 'EmailAttachments')
		Assert(record.EmailAttachments.Empty?())

		record.EmailAttachments = #(file1, file2, file3)
		method(record)
		Assert(record hasMember: 'EmailAttachments')
		Assert(record.EmailAttachments.Empty?())
		}

	Test_export?()
		{
		f = Params.Params_export?
		.MakeLibraryRecord([name: "TestTrueFormat", text: "QueryFormat { Export: true }"])

		Assert(f(Object()) is: false)
		Assert(f(Object(name: 'name')) is: false)

		Assert(f(Object(Object(Object()))) is: false)
		Assert(f(Object(Object('ReporterFormat'))))
		Assert(f(Object(Object(ReporterFormat))))

		Assert(f(Object(Object('TestTrueFormat'))))
		Assert(f(Object(Object('TestTrue'))))
		Assert(f(Object(Object(Global('TestTrueFormat')))))

		.MakeLibraryRecord([name: "TestFalseFormat",
			text: "QueryFormat { Export: false }"])
		Assert(f(Object(Object('TestFalseFormat'))) is: false)
		Assert(f(Object(Object('TestFalse'))) is: false)
		Assert(f(Object(Object(Global('TestFalseFormat')))) is: false)
		}

	Test_Get_devmode()
		{
		Assert(Params.Get_devmode([]) is:
			Record(width: 8.5, height: 11, left: .5, right: .5, top: .5, bottom: .5))
		Assert(Params.Get_devmode([default_orientation: false]) is:
			Record(width: 8.5, height: 11, left: .5, right: .5, top: .5, bottom: .5))
		Assert(Params.Get_devmode([default_orientation: 'Portrait']) is:
			Record(width: 8.5, height: 11, left: .5, right: .5, top: .5, bottom: .5))
		Assert(Params.Get_devmode([default_orientation: 'Landscape']) is:
			Record(width: 11, height: 8.5, left: .5, right: .5, top: .5, bottom: .5))
		}

	Test_filename()
		{
		fn = Params.Params_pdfName

		Assert(fn('') is: '.pdf')
		Assert(fn('testFile') is: 'testFile.pdf')
		Assert(fn('testFile.tmp') is: 'testFile.pdf')
		Assert(fn('testFile.txt') is: 'testFile.pdf')
		Assert(fn('testFile.jpeg') is: 'testFile.pdf')
		Assert(fn('testFile.txt.tmp') is: 'testFile.txt.pdf')
		Assert(fn('testFile.bak.tmp') is: 'testFile.bak.pdf')
		Assert(fn('test_tmp_File.tmp') is: 'test_tmp_File.pdf')
		}

	Test_emailPdfSubject()
		{
		fn = Params.Params_emailPdfSubject
		.SpyOn(PageHeadName).Return('Test company')

		Assert(fn("") is: "Test company - ")
		Assert(fn("Print") is: "Test company - Print")
		Assert(fn("Print ") is: "Test company - ")
		Assert(fn("Generate") is: "Test company - Generate")
		Assert(fn("Generate ") is: "Test company - ")

		Assert(fn("generate Invoice") is: 'Test company - generate Invoice')
		Assert(fn("Generate Invoice") is: 'Test company - Invoice')
		Assert(fn("print Orders") is: 'Test company - print Orders')
		Assert(fn("Print Orders") is: 'Test company - Orders')

		Assert(fn("Printing Orders") is: 'Test company - Printing Orders')
		Assert(fn("Generate and Print Orders") is: 'Test company - and Print Orders')
		}

	Test_RemoveIgnoreFields()
		{
		fieldName1 = .TempName().Lower()
		.MakeLibraryRecord(
			[name: "Field_" $ fieldName1, text: `class { ParamsNoSave: }`])

		fieldName2 = .TempName().Lower()
		.MakeLibraryRecord(
			[name:"Field_" $  fieldName2, text: `class { ParamsNoSave: true }`])

		fieldName3 = .TempName().Lower()
		.MakeLibraryRecord(
			[name: "Field_" $ fieldName3, text: `class { ParamsNoSave: false }`])

		fieldName4 = .TempName().Lower()
		.MakeLibraryRecord(
			[name: "Field_" $ fieldName4, text: `class { Prompt: "Hello World" }`])

		rec = []
		rec[fieldName1] = rec[fieldName2] = rec[fieldName3] = rec[fieldName4] = 'Fred'
		Params.RemoveIgnoreFields(rec)
		Assert(rec hasntMember: fieldName1)
		Assert(rec hasntMember: fieldName2)
		Assert(rec hasMember: fieldName3)
		Assert(rec hasMember: fieldName4)
		}

	Test_addFilterIfSlowQuery()
		{
		fn = Params.Params_addFilterIfSlowQuery
		Assert(fn(#()) is: false)
		Assert(fn(#(slowQueryFilter: 'hello')) is: false)
		Assert(fn(#(slowQueryFilter: false)) is: false)
		Assert(fn(#(slowQueryFilter: 'hello', previewWindow:)) is: false)

		p = Params
			{
			FindControl(unused)
				{
				return FakeObject(FindControl: false)
				}
			}
		fn = p.Params_addFilterIfSlowQuery
		Assert(fn(#(slowQueryFilter: 'hello')) is: false)

		p = Params
			{
			FindControl(unused)
				{
				return FakeObject(FindControl: FakeObject(GetFields: #(a, b)))
				}
			}
		fn = p.Params_addFilterIfSlowQuery
		Assert(fn(#(slowQueryFilter: 'hello')) is: false)

		p = Params
			{
			FindControl(unused)
				{
				return FakeObject(FindControl: FakeObject(GetFields: #(a, b)))
				}
			Params_paramsScreen?(unused)
				{
				return true
				}
			}
		fn = p.Params_addFilterIfSlowQuery
		Assert(fn(#(slowQueryFilter: 'hello')) is: false)

		Assert(fn(#("Vert", slowQueryFilter: 'hello')) is: false)

		Assert(fn(#(#(ReporterCanvasFormat), slowQueryFilter: 'hello')) is: false)

		// not a class
		Assert(fn(Object(Internal?, slowQueryFilter: 'hello')) is: false)

		Assert(fn(#(ObjectFormat { }, slowQueryFilter: 'hello')) is: false)

		Assert(fn(#(QueryFormat
			{
			Query()
				{
				return Object('hello_query')
				}
			}, slowQueryFilter: 'hello', paramsdata: [])) is: false)

		Assert(fn(#(QueryFormat
			{
			Query()
				{
				return Test.TempTableName()
				}
			}, slowQueryFilter: 'hello', paramsdata: [])) is: false)

		rpt = Object(.testFmt, slowQueryFilter: 'hello', paramsdata: [])
		Assert(fn(rpt) is: false)
		Assert(rpt hasntMember: 'suppressSlowQuery')

		spy = .SpyOn('SlowQuery.Validate')
		spy.Return(false)
		rpt = Object(.testFmt, slowQueryFilter: 'hello', paramsdata: [])
		Assert(fn(rpt) is: false)
		Assert(rpt.suppressSlowQuery)
		Assert(spy.CallLogs()[0].query endsWith: ' extend some_field sort hello')
		}

	testFmt: QueryFormat
		{
		testSort: 'hello'
		Query()
			{
			return Test.TempTableName() $ .extend $ ' sort ' $ .testSort
			}

		getter_extend()
			{
			return .extend = ' extend some_field'
			}
		}
	}
