// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// WARNING - getConvertedParamsDateCodes test may fail if test runs over midnight
	//		This should rarely be the case. To handle these cases we'd have to override
	//		the date being used for the conversions somehow in DateControl.convertDate
	Test_getConvertedParamsDateCodes()
		{
		fn = GetConvertedParamsDateCodes
		dateField = .TempName()
		.MakeLibraryRecord([name: "Field_" $ dateField,
			text: `Field_date
				{
				}`])

		date = Date().Replace(hour: 1, minute: 30) // force to have a time
		.SpyOn(DateControl.DateControl_today).Return(date)

		Assert(fn(false) is: false)

		params = Record(
			string_param: #(operation: "equals", value: "hello", value2: ""))
		params[dateField] = #(operation: "equals", value: "t", value2: "")
		paramsResult = fn(params)
		Assert(paramsResult[dateField].value is: date.NoTime())

		params[dateField] = #(operation: 'equals', value: 'bobogoobo', value2: '')
		paramsResult = fn(params)
		Assert(paramsResult[dateField].value is: false)
		Assert(paramsResult[dateField].value2 is: '')

		params[dateField] = #(operation: 'equals', value: #20120214, value2: '')
		paramsResult = fn(params)
		Assert(paramsResult[dateField].value is: #20120214)
		Assert(paramsResult[dateField].value2 is: '')

		params[dateField] = #(operation: 'range', value: #20120214, value2: 't')
		paramsResult = fn(params)
		Assert(paramsResult[dateField].value is: #20120214)
		Assert(paramsResult[dateField].value2 is: date.NoTime())

		// ensure the converted result is a record, so it works with standard report
		Assert(paramsResult isType: 'Record')
		Assert(paramsResult.non_existing_member is: '')
		Assert(paramsResult.custfield_num_new isDate: )
		}

	Test_getConvertedParamsDateCodes_dateTime()
		{
		fn = GetConvertedParamsDateCodes
		dateTimeField = .TempName()
		.MakeLibraryRecord([name: "Field_" $ dateTimeField,
			text: `Field_date_time { }`])

		date = Date().Replace(hour: 1, minute: 30) // force to have a time
		.SpyOn(DateControl.DateControl_today).Return(date)

		params = Record()
		params[dateTimeField] = #(operation: "equals", value: "t")
		paramsResult = fn(params)
		Assert(paramsResult[dateTimeField].value is: date)

		params = Record()
		params[dateTimeField] = #(operation: "equals", value: "h+")
		paramsResult = fn(params)
		val = paramsResult[dateTimeField].value
		Assert(val.NoTime() is: date.EndOfMonth().Plus(days: 1).NoTime())
		Assert(val.Format('HHmm') is: '0130')

		params = Record()
		params[dateTimeField] = #(operation: "range", value: "h-", value2: 'h+')
		paramsResult = fn(params)
		val = paramsResult[dateTimeField].value
		val2 = paramsResult[dateTimeField].value2
		Assert(val.Format('HHmm') is: '0130')
		Assert(val2.Format('HHmm') is: '0130')
		Assert(val is: date.EndOfMonth().Minus(days: 1))
		Assert(val2 is: date.EndOfMonth().Plus(days: 1))

		params = Record()
		params[dateTimeField] = #(operation: "equals", value: "t150", value2: "")
		paramsResult = fn(params)
		Assert(paramsResult[dateTimeField].value
			is: date.NoTime().Plus(hours: 1, minutes: 50))

		params = Record()
		params[dateTimeField] = #(operation: "equals", value: "t--2045", value2: "")
		paramsResult = fn(params)
		Assert(paramsResult[dateTimeField].value
			is: date.NoTime().Plus(days: -2, hours: 20, minutes: 45))
		}

	Test_getConvertedParamsDateCodes_nested()
		{
		fn = GetConvertedParamsDateCodes

		dateField = .TempName()
		.MakeLibraryRecord([name: "Field_" $ dateField, text: `Field_date { }`])

		dateTimeField = .TempName()
		.MakeLibraryRecord([name: "Field_" $ dateTimeField, text: `Field_date_time { }`])

		date = Date().Replace(hour: 1, minute: 30) // force to have a time
		.SpyOn(DateControl.DateControl_today).Return(date)

		params = Record(string_param: #(operation: "equals", value: "hello", value2: ""))
		testFilters = Object()
		testFilters[0] = [condition_field: dateField]
		testFilters[0][dateField] = #(value2: "", value: "m", operation: "greater than")
		testFilters[1] = [condition_field: dateTimeField]
		testFilters[1][dateTimeField] = #(value2: "", value: "t--30", operation: "equals")
		testFilters.nested = testFilters.Copy()
		params.testFilters = testFilters

		newFilters = fn(params).testFilters
		Assert(params.string_param is:
			#(operation: "equals", value: "hello", value2: ""))

		Assert(newFilters[0][dateField].value is: date.NoTime().Replace(day: 1))
		Assert(newFilters[1][dateTimeField].value
			is: date.NoTime().Minus(days: 2).Plus(minutes: 30))

		Assert(newFilters.nested[0][dateField].value is: date.NoTime().Replace(day: 1))
		Assert(newFilters.nested[1][dateTimeField].value
			is: date.NoTime().Minus(days: 2).Plus(minutes: 30))
		}

	Test_getConvertedParamsDateCodes_plainField()
		{
		fn = GetConvertedParamsDateCodes

		dateField = .TempName()
		.MakeLibraryRecord([name: "Field_" $ dateField, text: `Field_date { }`])

		dateTimeField = .TempName()
		.MakeLibraryRecord([name: "Field_" $ dateTimeField, text: `Field_date_time { }`])

		date = Date().Replace(hour: 1, minute: 30) // force to have a time
		.SpyOn(DateControl.DateControl_today).Return(date)

		params = Record(
			string_param: #(operation: "equals", value: "hello", value2: ""))
		params[dateField] = "t++"
		params[dateTimeField] = "t--1800"
		paramsResult = fn(params)
		Assert(paramsResult[dateField] is: date.NoTime().Plus(days: 2))
		Assert(paramsResult[dateTimeField]
			is: date.NoTime().Minus(days: 2).Plus(hours: 18))
		Assert(paramsResult.string_param is:
			#(operation: "equals", value: "hello", value2: ""))
		}

	Test_getConvertedParamsDateCodes_inList()
		{
		fn = GetConvertedParamsDateCodes

		params = #([
			date: #(
				value: #(#20181024, #20181025),
				value2: "",
				operation: "in list"),
			condition_field: "date"
			])
		paramsResult = fn(params)
		Assert(paramsResult is: params)
		}
	}