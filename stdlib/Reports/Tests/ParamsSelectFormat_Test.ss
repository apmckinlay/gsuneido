// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_layoutRange()
		{
		fn = ParamsSelectFormat.ParamsSelectFormat_layout
		data = #(operation: "not in range", value: #20180101, value2: #20180405)
		result = fn(data, 'date', false)
		Assert(result isSize: 5)
		Assert(result[0][0] is: 'Text')
		Assert(result[0][1] is: 'Excluding From ')

		Assert(result[1][0] is: 'ShortDate')
		Assert(result[1].data is: #20180101)

		Assert(result[2][0] is: 'Text')
		Assert(result[2][1] is: ' To ')

		Assert(result[3][0] is: 'ShortDate')
		Assert(result[3].data is: #20180405)
		}

	Test_layoutString()
		{
		fn = ParamsSelectFormat.ParamsSelectFormat_layout
		data = #(operation: "empty", value: "", value2: "")
		result = fn(data, 'string', false)

		Assert(result[0][0] is: 'Text')
		Assert(result[0][1] is: 'empty')

		data = #(operation: "equals", value: "Test", value2: "")
		result = fn(data, 'string', false)
		Assert(result[0][0] is: 'Text')
		Assert(result[0].data is: 'Test')
		}
	}