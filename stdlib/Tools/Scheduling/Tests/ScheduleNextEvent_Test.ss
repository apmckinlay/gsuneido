// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	testTable1:	"schedulecontroltesttable1"
	testTable2:	"schedulecontroltesttable2"
	Setup()
		{
		.now = Date()
		.nowPlusOne = .now.Plus(hours: 1)
		// destroy test tables - they shouldn't exist prior to test
		try
			Database("destroy " $ .testTable1)
		catch (x /*unused*/) { }
		try
			Database("destroy " $ .testTable2)
		catch (x /*unused*/) { }
		}

	Test_Basic_Scheduling()
		{
		// test getNext
		getNext = ScheduleNextEvent.ScheduleNextEvent_getNext
		Assert(getNext(false, .now) is: .now, msg: "getNext")
		Assert(getNext(.now, .nowPlusOne) is: .now, msg:"getNext")
		Assert(getNext(.nowPlusOne, .now) is: .now, msg:"getNext")
		// test evalDate
		schedule = new ScheduleNextEvent(taskTable: .testTable2)
		evalDate = schedule.ScheduleNextEvent_evalDate
		simplify = ScheduleNextEvent.ScheduleNextEvent_simplify
		testTask = new ScheduleTask(Object(taskTable: .testTable2,
			taskname: "x", task:"x"))
		testExecEvents = Object()
		testEvent = Object()
		QueryOutput(.testTable2, [].Merge(testTask))

		r = Mock().Merge(Query1(.testTable2))
		Assert(r isnt: false, msg: "evalDate: record doesn't exist")

		result = evalDate(Date.Begin(), Date(), testExecEvents, testEvent, true, r)
		Assert(Number?(result), msg: 'number result one')
		Assert(testExecEvents.Size() is: 0, msg: 'size one')

		result = evalDate(Date.Begin(), Date(), testExecEvents, testEvent, false, r)
		Assert(result is: false, msg: 'result false two')
		Assert(testExecEvents.Size() is: 0, msg: 'size two')

		now = Date()
		result = evalDate(now, simplify(now), testExecEvents, testEvent, true, r)
		Assert(result, msg: 'result true three')
		Assert(testExecEvents.Size() is: 1, msg: 'size three')
		Assert(r.prev_event is: simplify(now))

		result = evalDate(now, simplify(now), testExecEvents, testEvent, false, r)
		Assert(result is: false, msg: 'result false four')
		Assert(testExecEvents.Size() is: 1, msg: 'size four')

		now = simplify(now)
		result = evalDate(now, now.Plus(minutes: 1), testExecEvents, testEvent, true, r)
		Assert(result is: 60 * 1000)
		Assert(testExecEvents.Size() is: 1, msg: 'size five')

		Assert(simplify(.now)
			is: .now.NoTime().Plus(hours: .now.Hour(), minutes: .now.Minute()),
			msg: "simplify1")
		Assert(simplify(simplify(.now)) is: simplify(.now),	msg: "simplify2")
		}

	Test_getNext()
		{
		getNext = ScheduleNextEvent.ScheduleNextEvent_getNext
		Assert(getNext(false, #20050725) is: #20050725)
		Assert(getNext(#20050724, #20050725) is: #20050724)
		Assert(getNext(#20050725, #20050725) is: #20050725)
		Assert(getNext(#20050726, #20050725) is: #20050725)
		}

	Test_getNextInterval()
		{
		date = #20050701.120000
		intervals =  #(1, 2, 5, 7)
		results = #(
			0: #(1: #20050701.120001, 2: #20050701.120002, 5: #20050701.120005,
				7: #20050701.120007)
			1: #(1: #20050701.1201, 2: #20050701.1202, 5: #20050701.1205,
				7: #20050701.1207)
			2: #(1: #20050701.1300, 2: #20050701.1400, 5: #20050701.1700,
				7: #20050701.1900)
			3: #(1: #20050702.1200, 2: #20050703.1200, 5: #20050706.1200,
				7: #20050708.1200)
			4: #(1: #20050708.1200, 2: #20050715.1200, 5: #20050805.1200,
				7: #20050819.1200))
		for i in ..5
			for x in intervals
				{
				getNextInterval = ScheduleNextEvent.ScheduleNextEvent_getNextInterval
				Assert(getNextInterval(date, x, i) is: results[i][x])
				}
		}

	Test_evalDate()
		{
		c = ScheduleNextEvent
		ob2 = class { Update() { } }
		task = ob2()
		task.prev_event = #20050720
		task.next_event = #20050724
		task.taskname = 'test'

		// do nothing
		test = #(now: false, testDate: false, execEvents: false,
			event: false, evaluate?: false, task: false,
			result: false)
		Assert(test.result
			is: c.ScheduleNextEvent_evalDate(test.now, test.testDate, test.execEvents,
				test.event, test.evaluate?, test.task))

		// execute event
		test = Object(now: #20050725, testDate: #20050725, execEvents: Object(),
			event: false, evaluate?: true, :task, result: true)
		Assert(test.result
			is: c.ScheduleNextEvent_evalDate(test.now, test.testDate, test.execEvents,
				test.event, test.evaluate?, test.task))
		Assert(task.next_event is: Date.End())
		Assert(test.execEvents isSize: 1)
		Assert(task.prev_event is: #20050725)

		// do not execute yet - set next_event
		test = Object(now: #20050724.2359, testDate: #20050725, execEvents: Object(),
			event: false, evaluate?: true, :task, result: 60000)
		Assert(test.result
			is: c.ScheduleNextEvent_evalDate(test.now, test.testDate, test.execEvents,
				test.event, test.evaluate?, test.task))
		Assert(task.next_event is: #20050725)
		Assert(test.execEvents.Empty?())

		// do not execute yet - ensure milliseconds calculation does not overflow,
		// should result in maxTimeLapse
		test = Object(now: #20050724.2359, testDate: #20050825, execEvents: Object(),
			event: false, evaluate?: true, :task, result: 900000)
		Assert(test.result
			is: c.ScheduleNextEvent_evalDate(test.now, test.testDate, test.execEvents,
				test.event, test.evaluate?, test.task))
		Assert(task.next_event is: #20050825)
		Assert(test.execEvents.Empty?())
		}

	Test_evalResult()
		{
		Assert(false is: ScheduleNextEvent.ScheduleNextEvent_evalResult(
			result: false, event: false, recClosestEvent: false,
			testDate: false, next: false, nextDate: false))
		// result is a number
		Assert(ScheduleNextEvent.ScheduleNextEvent_evalResult(
				result: 2, event: false, recClosestEvent: 1,
				testDate: false, next: false, nextDate: false)
			is: 1)
		Assert(ScheduleNextEvent.ScheduleNextEvent_evalResult(
				result: 1, event: false, recClosestEvent: 2,
				testDate: false, next: false, nextDate: false)
			is: 1)
		event = Object(next: #20050726)
		Assert(ScheduleNextEvent.ScheduleNextEvent_evalResult(
				result: 1, :event, recClosestEvent: 2,
				testDate: #20050725, next: false, nextDate: false)
			is: 1)
		Assert(event.next is: #20050725)
		// result is not a number - next is true and nextDate isnt false
		event = Object(next: #20050726)
		Assert(ScheduleNextEvent.ScheduleNextEvent_evalResult(
				result: false, :event, recClosestEvent: 'test',
				testDate: false, next: true, nextDate: #20050725)
			is: 'test')
		Assert(event.next is: #20050725)
		}

	Test_evalInterval_TimeSpan()
		{
		execEvents = Object()
		event = Record()
		ob2 = class { Update() { } }
		task = ob2()
		task.time_span = true
		task.time_end = 20
		task.prev_interval = #20050725.1200
		task.time_start = 1200
		evalInterval = ScheduleNextEvent.ScheduleNextEvent_evalInterval

		now = #20050726.1150
		Assert(evalInterval(now, execEvents, event, task) is: 10 * 60 * 1000)

		now = #20050725
		Assert(evalInterval(now, execEvents, event, task) is: 900000)

		now = #20050725.1700
		Assert(evalInterval(now, execEvents, event, task) is: 900000)
		Assert(task.prev_interval is: now)
		Assert(execEvents isSize: 1)
		Assert(event.next is: #20050726.1200)
		}

	Test_evalInterval()
		{
		execEvents = Object()
		event = Record()
		ob2 = class { Update() { } }
		task = ob2()
		task.time_span = false
		task.runinterval = true
		task.interval = 3
		task.interval_units = 1
		task.prev_interval = #20050725.1200
		evalInterval = ScheduleNextEvent.ScheduleNextEvent_evalInterval

		now = #20050725.1201
		Assert(evalInterval(now, execEvents, event, task) is: 2 * 60 * 1000 )

		task.interval = 2
		task.interval_units = 2
		now = #20050725.1350
		Assert(evalInterval(now, execEvents, event, task) is: 10 * 60 * 1000)

		now = #20050725.1300
		Assert(evalInterval(now, execEvents, event, task) is: 900000)
		}

	Test_simplify()
		{
		Assert(ScheduleNextEvent.ScheduleNextEvent_simplify(#20050725.124533)
			is: #20050725.1245 )
		Assert(ScheduleNextEvent.ScheduleNextEvent_simplify(#20050725.1300)
			is: #20050725.1300)
		Assert(ScheduleNextEvent.ScheduleNextEvent_simplify(#20050725.000125)
			is: #20050725.0001)
		Assert(ScheduleNextEvent.ScheduleNextEvent_simplify(#20050725)
			is: #20050725)
		}

	Test_next_run_date()
		{
		nextRunDate = ScheduleNextEvent.ScheduleNextEvent_next_run_date
		Assert(nextRunDate(#(suspended: true, next_event: 'test'), Date) is: 'test')
		Assert(nextRunDate(#(suspended: false, next_event: 'test', rundate: true), Date)
			is: 'test')
		task = Record(runinterval: true)
		Assert(nextRunDate(task, Date) is: false)
		// run daily
		task.Delete('runinterval')
		task.rundaily = true
		task.prev_event = #20050101.1200
		Assert(nextRunDate(task, #20050102.1100) is: #20050102.1200)
		Assert(nextRunDate(task, #20050102.1300) is: #20050103.1200)
		// run monthly
		task.Delete('rundaily')
		task.runmonthly = true
		task.prev_event = #20050301.1200
		Assert(nextRunDate(task, #20050401.1100) is: #20050401.1200)
		Assert(nextRunDate(task, #20050401.1300) is: #20050501.1200)
		// run weekly
		task.Delete('runmonthly')
		task.runweekly = true
		task.prev_event = #20050301.1200
		Assert(nextRunDate(task, #20050305.1100) is: #20050308.1200)
		Assert(nextRunDate(task, #20050310.1300) is: #20050315.1200)

		// run daily
		task = Record(weekly_day: 0, prev_event: #17000101, date: #20050817.1200,
			interval_units: 0, suspended: false, interval: 10, runweekly: false,
			weekly_time: 1200, task: "Alerrt('Bye')", taskname: "test2",
			prev_interval: #20050817.1336, monthly_time: 1200,
			next_event: #20050816.1336, runinterval: false, threaded: false,
			daily_time: 1336, monthly_day: 0, rundaily: true, time_start: 0,
			time_span: false, uid: 1, runmonthly: false, time_end: 0, rundate: false)
		Assert(nextRunDate(task, #20050817.1358) is: #20050818.1336)
		}

	Test_set_prev_event()
		{
		setPrevEvent = ScheduleNextEvent.ScheduleNextEvent_set_prev_event
		Assert(setPrevEvent([rundaily: true, daily_time: 1306]) is: #17000101.1306)
		task = Record(runmonthly: true, monthly_day: 13, monthly_time: 536)
		Assert(setPrevEvent(task) is: #17000113.0536)
		task = Record(runweekly: true, weekly_day: 3, weekly_time: 2239)
		Assert(setPrevEvent(task) is: #17000106.2239)
		}

	Test_skipped_event()
		{
		skipped_event = ScheduleNextEvent.ScheduleNextEvent_skipped_event

		scheduler = Mock()
		task = Mock()
		now = Date().NoTime()
		task.next_event = now
		scheduler.When.ScheduleNextEvent_next_run_date([anyArgs:]).Return(now)
		scheduler.Eval(skipped_event, task, now, Object(), Object())
		task.Verify.Never().Update()

		scheduler.When.ScheduleNextEvent_next_run_date([anyArgs:]).Return(Timestamp())
		scheduler.Eval(skipped_event, task, now, Object(), Object())
		task.Verify.Update()
		}

	Test_current_event()
		{
		current_event = ScheduleNextEvent.ScheduleNextEvent_current_event
		now = Date().NoTime()
		scheduler = Mock()

		task = Mock()
		task.update? = ''
		scheduler.When.ScheduleNextEvent_milliSecondsBetweenDates([anyArgs:]).Return(1)
		scheduler.When.ScheduleNextEvent_simplify([anyArgs:]).Return(now.NoTime())

		task.next_event = now.Plus(days: 1)
		scheduler.Eval(current_event, now.Plus(days: 1), now, task, [], [])
		task.Verify.Never().Update()

		task.next_event = now.Plus(days: 2)
		scheduler.Eval(current_event, now.Plus(days: 1), now, task, [], [])
		task.Verify.Update()

		task = Mock()
		task.update? = ''
		scheduler.When.ScheduleNextEvent_milliSecondsBetweenDates([anyArgs:]).Return(0)
		task.next_event = Date.End()
		task.prev_event = now.NoTime()
		scheduler.Eval(current_event, now.Plus(days: 1), now, task, [], [])
		task.Verify.Never().Update()

		task = Mock()
		task.update? = ''
		task.next_event = Date.End()
		task.prev_event = now.Plus(days: -1)
		scheduler.Eval(current_event, now.Plus(days: 1), now, task, [], [])
		task.Verify.Update()

		task = Mock()
		task.update? = ''
		task.next_event = now
		task.prev_event = now.Plus(days: -1)
		scheduler.Eval(current_event, now.Plus(days: 1), now, task, [], [])
		task.Verify.Update()
		}

	Test_getClosestEvent()
		{
		ScheduleNextEvent.EnsureTaskTable(.testTable1)

		mock = new Mock(ScheduleNextEvent)
		mock.ScheduleNextEvent_maxTimeLapse = 1000
		mock.ScheduleNextEvent_taskTable = .testTable1

		intervalTask = new ScheduleTask(Object(taskTable: .testTable1,
			taskname: "interval", task: "interval", runinterval: true))
		QueryOutput(.testTable1, [].Merge(intervalTask))

		dailyTask = new ScheduleTask(Object(taskTable: .testTable1,
			taskname: "daily", task: "daily", rundaily: true, daily_time: 1124))
		QueryOutput(.testTable1, [].Merge(dailyTask))

		weeklyTask = new ScheduleTask(Object(taskTable: .testTable1,
			taskname: "weekly", task: "weekly", runweekly: true, weekly_time: 1124))
		QueryOutput(.testTable1, [].Merge(weeklyTask))

		monthlyTask = new ScheduleTask(Object(taskTable: .testTable1,
			taskname: "weekly", task: "weekly", runmonthly: true, monthly_time: 1124))
		QueryOutput(.testTable1, [].Merge(monthlyTask))

		dateTask = new ScheduleTask(Object(taskTable: .testTable1,
			taskname: "date", task: "date", rundate: true,
			date: Date().NoTime().Plus(days: 1)))
		QueryOutput(.testTable1, [].Merge(dateTask))

		events = Object()
		execEvents = Object()
		mock.When.calcInterval([anyArgs:]).Return('')
		mock.When.getNextIntervalTask([anyArgs:]).Return(100)
		mock.When.getNextDaily([anyArgs:]).Return(200)
		mock.When.getNextWeekly([anyArgs:]).Return(300)
		mock.When.getNextMonthly([anyArgs:]).Return(400)
		mock.When.getNextDateTask([anyArgs:]).Return(500)
		mock.When.getClosestEvent([anyArgs:]).CallThrough()
		Assert(mock.getClosestEvent(events, Date(), execEvents) is: 100)

		Assert(events isSize: 5)
		mock.Verify.getNextIntervalTask([anyArgs:])
		mock.Verify.getNextDaily([anyArgs:])
		mock.Verify.getNextWeekly([anyArgs:])
		mock.Verify.getNextMonthly([anyArgs:])
		mock.Verify.getNextDateTask([anyArgs:])
		mock.Verify.Times(5).calcInterval([anyArgs:])
		}

	Teardown()
		{
		try Database("destroy " $ .testTable1)
		try Database("destroy " $ .testTable2)
		}
	}
