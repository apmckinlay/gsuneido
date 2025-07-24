// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(PrintParams(Object(date: '')) is: #(WrapItems))
		Assert(PrintParams(Object(ReportDestination: 'preview')) is: #(WrapItems))
		Assert(PrintParams(Object(date: #(operation: '', value: ''))) is: #(WrapItems))

		result = #("WrapItems",
			#("Horz", #("Text", "Date: "),
			#("ShortDate", data: #20101201),
			font: #(size: 8, name: "Arial")))
		Assert(PrintParams(Object(date: #20101201)) is: result)

		result = #("WrapItems",
			#("Horz", #("Text", "Date: "),
			#("ShortDate", data: #(value: #20101201, opteration: "equals")),
			font: #(size: 8, name: "Arial")))
		Assert(PrintParams(Object(date: Object(opteration: 'equals', value: #20101201)))
			is: result)
		}

	Test_ctrlFormat()
		{
		fn = PrintParams.PrintParams_ctrlFormat
		.MakeLibraryRecord([name: "Field_fakeRule",
			text: `Field_string { Prompt: 'Fake Rule'; ParamsNoSave: }`])
		.MakeLibraryRecord([name: "Rule_fakeRule",
			text: `function () { return 'hello world' }`])
		.MakeLibraryRecord([name: "Field_fakeRule2",
			text: `Field_string { Prompt: 'Fake Rule' }`])
		.MakeLibraryRecord([name: "Rule_fakeRule2",
			text: `function () { return 'return value' }`])
		.MakeLibraryRecord([name: "Field_fakeField",
			text: `Field_string { Prompt: 'Fake Field'; ParamsNoSave: }`])
		fmt = Object()
		expected = #(("Horz",
			("Text", "Fake Rule: "),
			("Text", data: "hello world"),
			font: (size: 8, name: "Arial")))
		paramData = Record(number: 5)
		fn('fakeRule', paramData, fmt)
		Assert(fmt is: expected)
		Assert(paramData is: Record(number: 5, fakeRule: 'hello world'))

		paramData = Record(number: 1)
		fmt = Object()
		fn('name', paramData, fmt)
		Assert(fmt is: #())
		Assert(paramData is: Record(number: 1))

		paramData = Record(name: 'hello', number: 33)
		fmt = Object()
		expected = #(("Horz",
			("Text", "Name: "),
			("Text", data: "hello"), font: (size: 8, name: "Arial")))
		fn('name', paramData, fmt)
		Assert(fmt is: expected)
		Assert(paramData is: Record(name: 'hello', number: 33))

		fmt = Object()
		fn('doesNotExist', paramData, fmt)
		Assert(fmt is: #())
		Assert(paramData is: Record(name: 'hello', number: 33))

		fmt = Object()
		fn('fakeField', paramData, fmt)
		Assert(fmt is: #())
		Assert(paramData is: Record(name: 'hello', number: 33))

		fmt = Object()
		fn('fakeRule2', paramData, fmt)
		Assert(fmt is: #())
		Assert(paramData is: Record(name: 'hello', number: 33))

		param = Object(paramPrompt: 'Test', paramFormat: 'Number', paramField: 'number')
		paramData = Record(number: 10)
		fmt = Object()
		expectedFmt = #(("Horz",
			("Text", "Test: "),
			("Number", data: 10), font: (size: 8, name: "Arial")))
		fn(param, paramData, fmt)
		Assert(fmt is: expectedFmt)
		}
	}