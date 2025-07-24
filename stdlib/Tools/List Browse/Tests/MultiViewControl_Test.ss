// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_needToApplySelects?()
		{
		fn = MultiViewControl.MultiViewControl_needToApplySelects?
		Assert(fn(fromNew?:) is: false)

		m = Object()
		m.MultiViewControl_prevSelVal = []
		m.MultiViewControl_newRecordsSinceFlip = false
		m.MultiViewControl_access = #(Select_vals: ())
		Assert(m.Eval(fn, fromNew?: false) is: false)

		m.MultiViewControl_newRecordsSinceFlip = true
		Assert(m.Eval(fn, fromNew?: false) is: true)

		m.MultiViewControl_newRecordsSinceFlip = false
		m.MultiViewControl_prevSelVal = [[something:]]
		Assert(m.Eval(fn, fromNew?: false) is: true)
		}

	Test_getGoToField()
		{
		test = Object()
		fn = MultiViewControl.MultiViewControl_getGoToField
		baseQuery = 'test_query rename testquery_a to testquery_a_renamed'

		test.MultiViewControl_query = baseQuery
		test.MultiViewControl_access = class
			{
			GetKeys()
				{
				return #(testquery_num)
				}
			}
		Assert(test.Eval(fn) is: 'testquery_num')

		test.MultiViewControl_access = class
			{
			GetKeys()
				{
				return #(testquery_name)
				}
			}
		Assert(test.Eval(fn) is: 'testquery_name')

		test.MultiViewControl_access = class
			{
			GetKeys()
				{
				return #(testquery_num, testquery_name)
				}
			}
		Assert(test.Eval(fn) is: 'testquery_num')

		sort = ' sort testquery_num'
		test.MultiViewControl_query = baseQuery $ sort
		Assert(test.Eval(fn) is: 'testquery_num')

		sort = ' sort testquery_name'
		test.MultiViewControl_query = baseQuery $ sort
		Assert(test.Eval(fn) is: 'testquery_name')

		sort = ' sort testquery_a_renamed'
		test.MultiViewControl_query = baseQuery $ sort
		Assert(test.Eval(fn) is: 'testquery_num')

		sort = ' sort testquery_b, testquery_c'
		Assert(test.Eval(fn) is: 'testquery_num')

		sort = ' sort testquery_b, testquery_name'
		Assert(test.Eval(fn) is: 'testquery_num')

		sort = ' sort testquery_b, testquery_num'
		Assert(test.Eval(fn) is: 'testquery_num')
		}
	}