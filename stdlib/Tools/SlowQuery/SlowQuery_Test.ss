// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.TearDownIfTablesNotExist('slow_queries')
		}

	Test_main()
		{
		query = .TempName()
		Assert(SlowQuery.SlowQuery_slow?(query) is: false)

		SlowQuery.SlowQuery_log(query, 10)

		Assert(SlowQuery.SlowQuery_slow?(query))
		QueryApply1('slow_queries where sq_query is ' $ Display(query))
			{
			Assert(it.sq_last_used isnt: '')
			it.Delete()
			}

		Assert(SlowQuery.SlowQuery_slow?(query) is: false)
		}

	Test_getIndexes()
		{
		f = SlowQuery.SlowQuery_getIndexes.Func
		Assert(f(.TempName()) is: #())
		t = .MakeTable(' (num_internal, date, city, number, boolean,
			end_date, stateprov, start_date)
			key (num_internal) index (date) index (city) index (boolean)
			index (stateprov) index (end_date) index (date, city) index (start_date)')
		spy = .SpyOn(SelectPrompt)
		expected = #(date_new, end_date, start_date, boolean_new, city_default, stateprov)
		for col in expected
			spy.Return(col, when: function(col){ return {|field| field is col } }(col))
		Assert(f(t $ ' rename num_internal to num_internal_new
			rename date to date_new
			rename city to city_default,
				boolean to boolean_new')
			is: expected)
		}

	Test_LogIfTooSlow()
		{
		.WatchTable('slow_queries')
		mock = Mock(SlowQuery)
		mock.When.LogIfTooSlow([anyArgs:]).CallThrough()
		mock.LogIfTooSlow(1, 'query')
		mock.Verify.Never().log([anyArgs:])

		mock.LogIfTooSlow(10, 'query where abc =~ "hello"')
		mock.Verify.Never().log([anyArgs:])

		mock.LogIfTooSlow(10, 'query /*SLOWQUERY SUPPRESS*/')
		mock.Verify.Never().log([anyArgs:])

		hash = ''
		mock.LogIfTooSlow(10, 'query')
			{|h|
			hash = h
			}
		mock.Verify.log([anyArgs:])
		Assert(hash isnt: '')

		mock = Mock(SlowQuery)
		mock.When.LogIfTooSlow([anyArgs:]).CallThrough()
		_slowQueryLog = Object()
		hash = ''
		mock.LogIfTooSlow(10, 'query')
			{|h|
			hash = h
			}
		mock.Verify.log([anyArgs:])
		Assert(hash isnt: '')
		Assert(_slowQueryLog.logged)

		mock = Mock(SlowQuery)
		mock.When.LogIfTooSlow([anyArgs:]).CallThrough()
		_slowQueryLog = Object(logged:)
		mock.LogIfTooSlow(10, 'query')
		mock.Verify.Never().log([anyArgs:])

		_slowQueryLog = Object(suppressed:)
		mock.LogIfTooSlow(10, 'query')
		mock.Verify.Never().log([anyArgs:])
		}
	}
