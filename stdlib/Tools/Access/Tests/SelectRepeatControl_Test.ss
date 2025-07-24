// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_BuildWhere()
		{
		GetForeignNumsFromNameAbbrevFilter_Test.SetupForeignTable(this)
		sf = SelectFields(#(table2_num, table2_name, table1_num_test))

		fn = SelectRepeatControl.BuildWhere
		Assert(fn(sf, []) is: #(where: '', errs: '', joinflds: (), fields: #()))

		condition1 = #(condition_field: table2_num, check:,
			table2_num: [operation: 'equals', value: 'test'])
		Assert(fn(sf, [condition1])
			is: #(where: ' where table2_num is "test"', errs: '', joinflds: (table2_num),
				fields: #(table2_num)))

		condition2 = #(condition_field: table1_abbrev_test, check:,
			table1_abbrev_test: [operation: 'range', value: 'aaa', value2: 'bbb'])
		Assert(fn(sf, [condition2])
			is: #(where: ' where table1_num_test in (1, 2)', errs: '',
				joinflds: (table1_num_test), fields: #(table1_num_test)))

		condition3 = #(condition_field: table1_name_test, check:,
			table1_name_test: [operation: 'not in range', value: 'BBB', value2: 'ZZZ'])
		Assert(fn(sf, [condition2, condition3])
			is: #(
				where: ' where table1_num_test in (1, 2) ' $
					'where table1_num_test in (1, "")',
				errs: '', joinflds: (table1_num_test, table1_num_test),
				fields: #(table1_num_test)))

		condition4 = #(condition_field: table1_num_test, check:,
			table1_num_test: [operation: 'not equal to', value: 1])
		Assert(fn(sf, [condition1, condition2, condition3, condition4])
			is: #(
				where: ' where table2_num is "test" ' $
					'where table1_num_test in (1, 2) ' $
					'where table1_num_test in (1, "") ' $
					'where table1_num_test isnt 1',
				errs: '',
				joinflds: (table2_num, table1_num_test, table1_num_test, table1_num_test),
				fields: #(table2_num, table1_num_test))
			)

		condition5 = #(condition_field: table1_abbrev_test, check:,
			table1_abbrev_test: [operation: 'not empty'])
		Assert(fn(sf, [condition1, condition2, condition3, condition4, condition5])
			is: #(
				where: ' where table2_num is "test" ' $
					'where table1_num_test in (1, 2) ' $
					'where table1_num_test in (1, "") ' $
					'where table1_num_test isnt 1 ' $
					'where table1_num_test isnt ""',
				errs: '',
				joinflds: (table2_num, table1_num_test, table1_num_test, table1_num_test,
					table1_num_test)
				fields: #(table2_num, table1_num_test))

			)
		}

	Test_ProcessPresets()
		{
		cl = SelectRepeatControl
			{
			Data: class
				{
				Get() { return _curConditions }
				}
			Set(conditions)
				{
				_set.Add(conditions)
				}
			SelectRepeatControl_valid?()
				{
				return true
				}
			}

		Assert(cl.ProcessPresets(#()) is true)
		Assert(cl.ProcessPresets('') is true)

		_set = Object()
		_curConditions = #(conditions: #())
		newConditions = #(conditions: #(#(condition_field: 'date',
			check:, date: [operation: 'not empty'])))
		Assert(cl.ProcessPresets(newConditions) is true)
		Assert(_set.Last() is:
			#(conditions: #(#(condition_field: 'date',
				check:, date: [operation: 'not empty']))))

		// same condition, different operators
		_set = Object()
		_curConditions = [conditions: Object(Object(condition_field: 'date',
			check:, date: [operation: 'not empty']))]
		newConditions = [conditions: Object(Object(condition_field: 'date',
			check:, date: [operation: 'empty']))]
		Assert(cl.ProcessPresets(newConditions) is true)
		Assert(_set.Last() is:
			#(conditions: #(#(condition_field: 'date',
				check:, date: [operation: 'empty']))))

		// same condition, same operator
		_set = Object()
		_curConditions = [conditions: Object(Object(condition_field: 'date',
			check:, date: [operation: 'empty']))]
		newConditions = [conditions: Object(Object(condition_field: 'date',
			check:, date: [operation: 'empty']))]
		Assert(cl.ProcessPresets(newConditions) is true)
		Assert(_set.Last() is:
			#(conditions: #(#(condition_field: 'date',
				check:, date: [operation: 'empty']))))

		// different conditions
		_set = Object()
		_curConditions = [conditions: Object(Object(condition_field: 'date',
			check:, date: [operation: 'empty']))]
		newConditions = [conditions: Object(Object(condition_field: 'text',
			check:, date: [operation: 'empty']))]
		Assert(cl.ProcessPresets(newConditions) is true)
		Assert(_set.Last() is:
			#(conditions: #(
				#(condition_field: 'text',
					check:, date: [operation: 'empty']),
				#(condition_field: 'date',
					check: false, date: [operation: 'empty']))))

		// remove duplicates
		_set = Object()
		_curConditions = [conditions: Object(
			Object(condition_field: 'date', check:, date: [operation: 'not empty']))]
		for .. 20
			{
			condition = Object(condition_field: 'text',
				check:, text: [operation: 'not empty'])
			_curConditions.conditions.Add(condition)
			}
		newConditions = [conditions: Object(Object(condition_field: 'text',
			check:, text: [operation: 'empty']))]
		Assert(cl.ProcessPresets(newConditions) is true)
		Assert(_set.Last() is:
			#(conditions: (
				(condition_field: 'text', check:, text: [operation: 'empty']),
				(condition_field: 'date', check: false, date: [operation: 'not empty']))))


		// max number of conditions
		_set = Object()
		_curConditions = [conditions: Object()]
		for i in .. 30
			{
			condition = Object(condition_field: 'text' $ i, check:)
			condition['text' $ i] = [operation: 'not empty']
			_curConditions.conditions.Add(condition)
			}
		newConditions = [conditions: Object()]
		for i in .. 30
			{
			condition = Object(condition_field: 'date' $ i, check:)
			condition['date' $ i] = [operation: 'not empty']
			newConditions.conditions.Add(condition)
			}
		Assert(cl.ProcessPresets(newConditions) is true)
		Assert(_set.Last().conditions.Size() is: 50)
		Assert(_set.Last().conditions[..30].CountIf({
			it.condition_field.Has?('date') and it.check is true })
			is: 30)
		Assert(_set.Last().conditions[20..].CountIf({
			it.condition_field.Has?('text') and it.check is false })
			is: 20)

		// new conditions with more than 50
		_set = Object()
		_curConditions = [conditions: Object()]
		for i in .. 20
			{
			condition = Object(condition_field: 'text' $ i, check:)
			condition['text' $ i] = [operation: 'not empty']
			_curConditions.conditions.Add(condition)
			}
		newConditions = [conditions: Object()]
		for i in .. 50
			{
			condition = Object(condition_field: 'date' $ i, check:)
			condition['date' $ i] = [operation: 'not empty']
			newConditions.conditions.Add(condition)
			}
		Assert(cl.ProcessPresets(newConditions) is true)
		Assert(_set.Last().conditions.Size() is: 50)
		Assert(_set.Last().conditions.CountIf({
			it.condition_field.Has?('date') and it.check is true })
			is: 50)
		}
	}