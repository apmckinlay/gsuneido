// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_fieldObject?()
		{
		fn = Report.Report_fieldObject?
		Assert(fn('field') is: false)
		Assert(fn(#(field)))
		Assert(fn(#(Field)) is: false)
		Assert(fn(Object(TextFormat)) is: false)
		}

	Test_applyProperies()
		{
		fn = Report.Report_applyProperies
		fn(Object(x: 1, y: 2, heading: 'hello', xstretch: 1), fmt = Object())
		Assert(fmt.X is: 1440)
		Assert(fmt.Y is: 2880)
		Assert(fmt.Heading is: 'hello')
		Assert(fmt.Xstretch is: 1)
		}

	Test_buildBookLogParams()
		{
		build = Report.Report_buildBookLogParams
		Assert(build(#()) is: #())
		Assert(build(#(printParams: #('test'), sort: 'test')) is: #(sort: 'test'))
		Assert(build(#(test: #(operation: 'equals', value: 't')))
			is: #(test: #(operation: 'equals', value: 't')))

		params = #(Filters:
			([number: #(operation: 'equals', value: 1), condition_field: number]))
		Assert(build(params)
			is: #(filterWheres: #(Filters: " where number is 1")))

		params = #(Filters:	([condition_field: 'invalid'], ['invalid'], []))
		Assert(build(params) is: #(filterWheres: #(Filters: "")))

		params = #(Filters:
			([number: #(operation: 'equals', value: 2), condition_field: 'number'])
			a_Filters:
			([date: #(operation: 'equals', value: #20240101), condition_field: 'date'])
			)
		Assert(build(params)
			is: #(filterWheres: #(Filters: " where number is 2",
				a_Filters: " where date is #20240101")))
		}
	}
