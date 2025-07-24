// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		fn = .cl.ThreadTotalControl_runCalc
		_newTaskList = Object(
			#(id: 1, query: 'q1', filters: #('test filter1')),
			#(id: 2, query: 'q2', filters: #('test filter2')),
			#(id: 3, query: 'q3', filters: #('test filter3')),
			#(id: 4, query: 'q4', filters: #('test filter4')))
		_resultList = Object('q1', 'q3', 'q4')
		msg = "ERROR: ThreadTotal calculation failed - exception"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg))
		fn('calcFunc', 'screenTotal')
		}

	cl: ThreadTotalControl
		{
		Destroyed?()
			{
			return _destroyed?
			}
		ThreadTotalControl_doCalc(@args)
			{
			query = args[1]
			if query is 'q2'
				throw 'exception'
			return args[1]
			}
		ThreadTotalControl_fetchNewTask()
			{
			return not _newTaskList.Empty?() ? _newTaskList.PopFirst() : false
			}
		ThreadTotalControl_doneCalc(x, unused)
			{
			Assert(x is: _resultList.PopFirst())
			}
		ThreadTotalControl_sleep()
			{
			}
		}
	}
