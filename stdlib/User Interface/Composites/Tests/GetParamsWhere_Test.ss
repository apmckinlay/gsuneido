// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(GetParamsWhere('abc' data: Record()) is: '')
		Assert(GetParamsWhere('abc' data: Record(abc: '')) is: '')
		// used to throw error if value2 was not in the param ob, but only range needs it
		params = Record(abc: #(operation: 'equals', value: 'def'))
		Assert(.standardize(GetParamsWhere('abc' data: params)) is: 'whereabcisdef')

		params = Record(abc: #(operation: 'range', value: 1, value2: 2))
		Assert(.standardize(GetParamsWhere('abc' data: params))
			is: 'whereabc>=1andabc<=2')
		}

	// this method is just to remove the possibility of differences in spacing and
	// quotes used for displaying values in the resulting 'where' clause
	standardize(result)
		{
		return result.Tr('^A-Za-z0-9<>=')
		}

	Test_getParamField()
		{
		data = Record()
		field = ''
		Assert(GetParamsWhere.GetParamsWhere_getParamField(data, field) is: '')

		data = Record(arivc_date: '', arivc_date_paidoff: '')
		field = 'arivc_date'
		Assert(GetParamsWhere.GetParamsWhere_getParamField(data, field)
			is: 'arivc_date')

		data = Record(arivc_date_param: '', arivc_date_paidoff_param: '')
		field = 'arivc_date'
		Assert(GetParamsWhere.GetParamsWhere_getParamField(data, field)
			is: 'arivc_date_param')
		}
	Test_isTimestamp?()
		{
		result = GetParamsWhere.GetParamsWhere_isTimestamp?('')
		Assert(result.timestamp? is: false, msg: 'empty field')
		Assert(result.field is: '')

		result = GetParamsWhere.GetParamsWhere_isTimestamp?(#(arivc_num, timestamp))
		Assert(result.timestamp?, msg: 'arivc_num timestamp')
		Assert(result.field is: 'arivc_num')

		result = GetParamsWhere.GetParamsWhere_isTimestamp?(#(arivc_num, string))
		Assert(result.timestamp? is: false, msg: 'arivc_num string')
		Assert(result.field is: #(arivc_num, string))
		}
	Test_addRestrictions()
		{
		build_callable = false
		ob = #(operation: "range")
		restrictions = Object()
		field = 'arivc_num'
		val = ''
		val2 = ''
		GetParamsWhere.GetParamsWhere_addRestrictions(ob, restrictions, field,
			val, val2, build_callable)
		Assert(restrictions
			is: Object(
				Object(field, ">=", val)
				Object(field, "<=", val2)))

		ob = #(operation: "not in range")
		restrictions = Object()
		GetParamsWhere.GetParamsWhere_addRestrictions(ob, restrictions, field,
			val, val2, build_callable)
		Assert(restrictions is: Object(Object(field, "not in range", val, val2)))

		ob = #(operation: "contains")
		val = 'test'
		restrictions = Object()
		GetParamsWhere.GetParamsWhere_addRestrictions(ob, restrictions, field,
			val, val2, build_callable)
		Assert(restrictions is: Object(Object(field, '=~', '(?i)(?q)test')))

		ob = #(operation: "not empty")
		val = ""
		restrictions = Object()
		GetParamsWhere.GetParamsWhere_addRestrictions(ob, restrictions, field,
			val, val2, build_callable)
		Assert(restrictions is: Object(Object(field, 'isnt', "")), msg: 'not empty')
		ob = #(operation: "not empty")
		val = ""
		restrictions = Object()
		GetParamsWhere.GetParamsWhere_addRestrictions(ob, restrictions, field,
			val, val2, build_callable: true)
		Assert(restrictions is: Object(Object(field, 'isnt', "")),
			msg: 'not empty callable')
		}

	Test_handleIdAllowOther()
		{
		ob = Object(operation: 'all', value: '', value2: "")
		method = GetParamsWhere.GetParamsWhere_handleIdAllowOther
		method(ob)
		Assert(Object(operation: 'all', value: '', value2: "") is: ob)

		ob = Object(operation: 'empty', value: '', value2: "")
		method(ob)
		Assert(Object(operation: 'empty', value: '', value2: "") is: ob)

		ob = Object(operation: 'not empty', value: '', value2: "")
		method(ob)
		Assert(Object(operation: 'not empty', value: '', value2: "") is: ob)

		ob = Object(operation: 'in list',
			value: #(#20110309.101353107), value2: "")
		method(ob)
		Assert(Object(operation: 'in list', value: #(#20110309.101353107), value2: "")
			is: ob)

		ob = Object(operation: 'equals', value: 'Barney', value2: "")
		method(ob)
		Assert(#(operation: 'equals', value: 'Barney', value2: "") is: ob)

		ts = Timestamp()
		ob = Object(operation: 'equals', value: ts, value2: "", idtext: 'Fred Flinstone')
		method(ob)
		Assert(Object(ts, 'Fred Flinstone') is: ob.value)

		ts = Timestamp()
		ob = Object(operation: 'not equal to', value: ts, value2: "",
			idtext: 'Fred Flinstone')
		method(ob)
		Assert(Object(ts, 'Fred Flinstone') is: ob.value)
		}

	Test_buildRestrictions()
		{
		// buildRestrictions(fields, data, no_encode)
		m = GetParamsWhere.GetParamsWhere_buildRestrictions
		build_callable = false

		// test empty
		Assert(m(#(), [], false, build_callable) is: Object())

		//test simple
		data = [date_time: #(operation: 'equals', value: #20170727)]
		Assert(m(#('date_time'), data, false, build_callable)
			is: Object(Object('date_time', 'is', #20170727)))

		// test timestamp option
		data = [date_time: #(operation: 'equals', value: #20170727)]
		Assert(m(#(('date_time', 'timestamp')), data, false, build_callable)
			is: Object(
				Object('date_time', '>=', #20170727.StartOfDay())
				Object('date_time', "<=", #20170727.EndOfDay()))
			)
		data = [date_time: #(operation: 'range', value: #20170727, value2: #20170728)]
		Assert(m(#(('date_time', 'timestamp')), data, false, build_callable)
			is: Object(
				Object('date_time', '>=', #20170727.StartOfDay())
				Object('date_time', "<=", #20170728.EndOfDay()))
			)
		data = [date_time: #(operation: 'not empty', value: #20170727, value2: #20170728)]
		Assert(m(#(('date_time', 'timestamp')), data, false, build_callable)
			is: Object(Object('date_time', 'isnt', '')))
		data = [date_time: #(operation: 'empty', value: #20170727, value2: #20170728)]
		Assert(m(#(('date_time', 'timestamp')), data, false, build_callable)
			is: Object(Object('date_time', 'is', '')))
		data = [date_time: #(operation: 'less than or equal to',
			value: #20170727, value2: #20170728)]
		Assert(m(#(('date_time', 'timestamp')), data, false, build_callable)
			is: Object(Object('date_time', '<=', #20170727.EndOfDay())))
		data = [date_time: #(operation: 'greater than or equal to',
			value: #20170727, value2: #20170728)]
		Assert(m(#(('date_time', 'timestamp')), data, false, build_callable)
			is: Object(Object('date_time', '>=', #20170727)))
		}

	Test_buildRestrictions_boolean()
		{
		// buildRestrictions(fields, data, no_encode)
		m = GetParamsWhere.GetParamsWhere_buildRestrictions
		build_callable = false

		data = [boolean: #(operation: 'equals', value: true)]
		Assert(m(#('boolean'), data, false, build_callable)
			is: Object(Object('boolean', 'is', true)))

		data = [boolean: #(operation: 'not equal to', value: true)]
		Assert(m(#('boolean'), data, false, build_callable)
			is: Object(Object('boolean', 'isnt', true)))

		data = [boolean: #(operation: 'equals', value: false)]
		Assert(m(#('boolean'), data, false, build_callable)
			is: Object(Object('boolean', 'isnt', true)))

		data = [boolean: #(operation: 'equals', value: '')]
		Assert(m(#('boolean'), data, false, build_callable)
			is: Object(Object('boolean', 'isnt', true)))

		// not equal to
		data = [boolean: #(operation: 'not equal to', value: true)]
		Assert(m(#('boolean'), data, false, build_callable)
			is: Object(Object('boolean', 'isnt', true)))

		data = [boolean: #(operation: 'not equal to', value: "")]
		Assert(m(#('boolean'), data, false, build_callable)
			is: Object(Object('boolean', 'is', true)))

		data = [boolean: #(operation: 'not equal to', value: false)]
		Assert(m(#('boolean'), data, false, build_callable)
			is: Object(Object('boolean', 'is', true)))
		}

	Test_buildRestrictions_invalid()
		{
		// sometimes restriction fields get data in them other than the correct
		// ParamsSelect object. We haven't been able to track down how the data gets
		// this way (maybe rules?), the purpose of this test is to ensure it is handled
		// gracefully when it does happen
		m = GetParamsWhere.GetParamsWhere_buildRestrictions
		build_callable = false
		watch = .WatchTable('suneidolog')

		data = [some_field: 0]
		Assert(m(#('some_field'), data, false, build_callable) is: Object())
		calls = .GetWatchTable(watch)
		Assert(calls isSize: 1)
		Assert(calls[0].sulog_message has: 'invalid param data')
		}

	Test_callable()
		{
		restrictions = [date: Object(operation: 'equals', value: #20100101, value2: '')]
		fn = GetParamsWhere('date', data: restrictions, build_callable:)
		Assert(fn([date: #20100101]), msg: 'equals #20100101')
		Assert(fn([date: #20190101]) is: false, msg: 'equals #20190101')

		restrictions = [date: Object(
			operation: 'not equal to', value: #20100101, value2: '')]
		fn = GetParamsWhere('date', data: restrictions, build_callable:)
		Assert(fn([date: #20100101]) is: false, msg: 'not equal to #20100101')
		Assert(fn([date: #20190101]), msg: 'not equal to #20190101')

		restrictions = [comment: Object(
			operation: 'contains', value: 'hello', value2: '')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']), msg: 'hello world contains hello')
		Assert(fn([comment: 'world']) is: false, msg: 'world contains hello')

		restrictions = [comment: Object(
			operation: 'does not contain', value: 'hello', value2: '')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']) is: false,
			msg: 'hello world does not contain hello')
		Assert(fn([comment: 'world']), msg: 'world does not contain hello')

		restrictions = [comment: Object(
			operation: 'matches', value: '.*world', value2: '')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']), msg: '.*world matches hello world')
		Assert(fn([comment: 'hello']) is: false, msg: '.*world matches hello')

		restrictions = [comment: Object(
			operation: 'does not match', value: '.*world', value2: '')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']) is: false,
			msg: '.*world does not matche hello world')
		Assert(fn([comment: 'hello']), msg: '.*world does not matche hello')

		restrictions = [comment: Object(
			operation: 'starts with', value: 'hello', value2: '')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']), msg: 'hello world starts with hello')
		Assert(fn([comment: 'world']) is: false, msg: 'world starts with hello')

		restrictions = [comment: Object(
			operation: 'ends with', value: 'world', value2: '')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']), msg: 'hello world ends with world')
		Assert(fn([comment: 'hello']) is: false, msg: 'hello ends with world')

		restrictions = [comment: Object(
			operation: 'range', value: 'a', value2: 'i')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']), msg: 'hello world range')
		Assert(fn([comment: 'world']) is: false, msg: 'world range')

		restrictions = [comment: Object(
			operation: 'not in range', value: 'a', value2: 'i')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 'hello world']) is: false, msg: 'hello world not in range')
		Assert(fn([comment: 'world']), msg: 'world not in range')

		restrictions = [comment: Object(operation: 'not empty', value: '', value2: '')]
		fn = GetParamsWhere('comment', data: restrictions, build_callable:)
		Assert(fn([comment: 65]), msg: 'number not empty')
		}

	// GetParamsWhere used to use GetDefault to obtain the filter data for each field,
	// but now GetDefault is invoking rules when the filter data is not set which causes
	// the field data to be the rule value instead of the filter object that it is
	// supposed to be.
	Test_paramWhenFieldIsRule()
		{
		r = Record()
		r.AttachRule('biztest_getparamswhere_test', function () { return 'test' })
		Assert(GetParamsWhere('biztest_getparamswhere_test', data: r) is: '')

		r.biztest_getparamswhere_test = Object(operation: 'equals', value: 'test')
		Assert(GetParamsWhere('biztest_getparamswhere_test', data: r)
			is: ' where biztest_getparamswhere_test is "test"')
		}

	Test_handleEmptyValues()
		{
		cl = GetParamsWhere { GetParamsWhere_getSelectEmptyVals(unused) { return false } }
		m = cl.GetParamsWhere_handleEmptyValues

		ob = #()
		m('field', ob)
		Assert(ob is: #())

		ob = #(field: 'field', operation: 'equals', value: 'something')
		m('field', ob)
		Assert(ob is: #(field: 'field', operation: 'equals', value: 'something'))

		cl = GetParamsWhere { GetParamsWhere_getSelectEmptyVals(unused)
			{ return Object('', 'None') } }
		m = cl.GetParamsWhere_handleEmptyValues

		ob = #(field: 'field', operation: 'contains', value: 'something')
		m('field', ob)
		Assert(ob is: #(field: 'field', operation: 'contains', value: 'something'))

		ob = #(field: 'field', operation: 'equals', value: 'something')
		m('field', ob)
		Assert(ob is: #(field: 'field', operation: 'equals', value: 'something'))

		ob = Object(field: 'field', operation: 'equals', value: '')
		m('field', ob)
		Assert(ob is: #(field: 'field', operation: 'in list', value: #('', 'None')))

		ob = Object(field: 'field', operation: 'not empty', value: '')
		m('field', ob)
		Assert(ob is: #(field: 'field', operation: 'not in list', value: #('', 'None')))

		ob = Object(field: 'field', operation: 'in list', value: Object('some value'))
		m('field', ob)
		Assert(ob is: #(field: 'field', operation: 'in list', value: #('some value')))

		ob = Object(field: 'field', operation: 'in list', value: Object('None'))
		m('field', ob)
		Assert(ob.operation is: 'in list')
		Assert(ob.value.Size() is: 2)
		Assert(ob.value has: '')
		Assert(ob.value has: 'None')

		ob = Object(field: 'field', operation: 'in list', value: Object('', 'v1', 'v2'))
		m('field', ob)
		Assert(ob.operation is: 'in list')
		Assert(ob.value.Size() is: 4)
		Assert(ob.value has: '')
		Assert(ob.value has: 'None')
		Assert(ob.value has: 'v1')
		Assert(ob.value has: 'v2')
		}
	}