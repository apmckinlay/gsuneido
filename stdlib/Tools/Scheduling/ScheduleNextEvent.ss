// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// NOTE: use Scheduler for simpler implementation
class
	{
	destroyed: false
	maxTimeLapse: 900000			// 1000 * 60 * 15: 15 minutes in milliseconds
	CallClass(taskTable = 'tasks', logTasks = false,
		updateEventsFn = false) // used by SchedulerControl
		{
		scheduledEventInstance = new this(taskTable, logTasks, updateEventsFn)
		scheduledEventInstance.NextEvent()
		return scheduledEventInstance
		}

	noRunningEvents: true
	New(.taskTable = 'tasks', .logTasks = false, .updateEventsFn = false)
		{
		.EnsureTaskTable(taskTable)
		}
	EnsureTaskTable(taskTable)
		{
		Database('ensure ' $ taskTable $
			" (uid, taskname, task, runinterval, rundaily, runweekly, runmonthly,
			rundate, interval, interval_units, daily_time, weekly_time, weekly_day,
			monthly_time, monthly_day, date, suspended, threaded, prev_event,
			next_event, prev_interval, time_span, time_start, time_end)
			key (uid)")
		}

	// only for UI - uses Delay instead of Sleep so it doesn't block
	NextEvent()
		{
		Assert(Sys.Win32?())
		if .destroyed is true
			return
		delay = .nextEvent()
		Delay(delay, .NextEvent)
		}

	// runs forever
	EventLoop()
		{
		while not .destroyed
			{
			delay = .nextEvent()
			Thread.Sleep(delay)
			}
		}

	// returns the delay (in ms) till the next event (possibly 0)
	nextEvent()
		{
		now = Date()
		events = Object() // list of event objects of form #(uID, next)
		execEvents = Object()
		// find tasks due to run (execEvents) and closest next event after that
		closestEvent = .getClosestEvent(events, now, execEvents)

		if .noRunningEvents is true
			{
			if .updateEventsFn isnt false
				(.updateEventsFn)(events)
			.noRunningEvents = false
			return closestEvent
			}
		// convert to a time since running tasks may take some of this time
		next = Date().Plus(milliseconds: closestEvent)
		// execute any current tasks
		print = Suneido.Member?('Print') ? Suneido.Print : false
		for (exec in execEvents.Sort!(By(#interval)))
			{
			.log_task(exec.taskname,
				'SchedulerLastProcessStarted', .taskTable, .logTasks)
			try
				{
				exec.task.Eval() // needs Eval
				.log_task(exec.taskname $ ' Completed',
					'SchedulerLastProcessCompleted', .taskTable, .logTasks)
				}
			catch (error)
				{
				SuneidoLog('ERROR: (CAUGHT) ScheduleNextEvent - ' $ exec.taskname $
					' - ' $ error, caughtMsg: 'unattended')
				.log_task(exec.taskname $ ' ERROR (' $ error $ ')',
					'SchedulerLastProcessCompleted', .taskTable, .logTasks)
				}
			}
		if (print isnt false)
			Suneido.Print = print
		// update the ScheduleWindowControl, if it is open
		if .updateEventsFn isnt false
			(.updateEventsFn)(events)

		closestEvent = .milliSecondsBetweenDates(Date(), next)
		return Max(0, closestEvent)
		}

	milliSecondsBetweenDates(fromDate, toDate)
		{
		// this method should not always try to use the milliseconds
		// between the dates because the numbers can overflow to negative
		// values if the dates are too far apart (e.g. you only have a monthly task)
		if fromDate.Plus(milliseconds: .maxTimeLapse) < toDate
			return .maxTimeLapse
		return (toDate.MinusSeconds(fromDate) * 1000).Floor() /*= 1000 milliseconds */
		}

	getClosestEvent(events, now, execEvents)
		{
		idx = 0
		closestEvent = .maxTimeLapse
		RetryTransaction()
			{ |t|
			t.QueryApply(.taskTable $
				' where uid > ' $ Display(ScheduleTask.MinUID() - 1) $ ' sort uid')
				{ |task|
				idx++
				recClosestEvent = closestEvent
				events.Add(Object(uID: task.uid, taskname: task.taskname,
					task: task.task, next: false, threaded: task.threaded,
					interval: .calcInterval(task)))
				// a.  Examine suspended field
				if task.suspended is true
					continue
				// b.  Examine prev_interval and interval fields
				if ((task.runinterval is true or task.time_span is true) and
					task.interval isnt '' and task.interval_units isnt '')
					recClosestEvent = .getNextIntervalTask(now, task, execEvents, events,
						idx, recClosestEvent)
				// c.  Examine daily fields
				if task.rundaily is true
					recClosestEvent = .getNextDaily(now, task, execEvents, events,
						idx, recClosestEvent)
				// d.  Examine weekly fields
				if task.runweekly is true
					recClosestEvent = .getNextWeekly(now, task, execEvents, events,
						idx, recClosestEvent)
				// e.  Examine monthly fields
				if task.runmonthly is true
					recClosestEvent = .getNextMonthly(now, task, execEvents, events,
						idx, recClosestEvent)
				// f.  Examine date field
				if task.rundate is true
					recClosestEvent = .getNextDateTask(now, task, execEvents, events,
						idx, recClosestEvent)
				// f.  Update closest event, if necessary
				if (recClosestEvent < closestEvent)
					closestEvent = recClosestEvent
				}
			}
		return closestEvent
		}

	calcInterval(task)
		{
		if task.runinterval is true
			return .getNextInterval(Date.Begin(), task.interval, task.interval_units)
		else
			return Date.End()
		}

	getNextIntervalTask(now, task, execEvents, events, idx, recClosestEvent)
		{
		result = .evalInterval(now, execEvents, events[idx - 1], task)
		return .evalResult(result, events[idx - 1],
			recClosestEvent, false, false, false)
		}

	getNextDaily(now, task, execEvents, events, idx, recClosestEvent)
		{
		testDate = now.NoTime().Plus(
			hours:	.hours(task.daily_time),
			minutes: .minutes(task.daily_time)
			)
		result = .evalDate(now, testDate, execEvents, events[idx - 1],
			task.rundaily, task)
		return .evalResult(result, events[idx - 1],
			recClosestEvent, testDate, task.rundaily, testDate.Plus(days: 1))
		}

	hours(time)
		{
		return (time / 100).Floor() /*= 2 digits in the front */
		}

	minutes(time)
		{
		return time - .hours(time) * 100 /*= 2 digits in the front */
		}

	getNextWeekly(now, task, execEvents, events, idx, recClosestEvent)
		{
		dayMod = (now.WeekDay() <= task.weekly_day)
			? task.weekly_day - now.WeekDay()
			: 7 - now.WeekDay() + task.weekly_day
		testDate = now.NoTime().Plus(
			hours:	.hours(task.weekly_time),
			minutes: .minutes(task.weekly_time),
			days:	dayMod
			)
		result = .evalDate(now, testDate, execEvents, events[idx - 1],
			task.runweekly, task)
		return .evalResult(result, events[idx - 1],
			recClosestEvent, testDate, task.runweekly, testDate.Plus(days: 7))
		}

	getNextMonthly(now, task, execEvents, events, idx, recClosestEvent)
		{
		day = .GetMonthlyDay(now, task.monthly_day)
		// time needed to create date
		time = String(task.monthly_time).LeftFill(.timeDigits, '0')
		testDate = Date(now.Format("yyyyMM") $ day.Pad(2) $	"." $ time,
			"yyyyMMdd.t")
		result = .evalDate(now, testDate, execEvents, events[idx - 1],
			task.runmonthly, task)
		return .evalResult(result, events[idx - 1],
			recClosestEvent, testDate, task.runmonthly, testDate.Plus(months: 1))
		}

	GetMonthlyDay(date, monthlyDay)
		{
		switch (monthlyDay)
			{
		case 0: day = 1
		case 1: day = date.EndOfMonthDay()
		case 2: day = (date.EndOfMonthDay() / 2).Ceiling()
		default : day = monthlyDay - 1
			}
		return day
		}

	getNextDateTask(now, task, execEvents, events, idx, recClosestEvent)
		{
		testDate = task.date
		result = .evalDate(now, testDate, execEvents, events[idx - 1],
			task.rundate, task)
		return .evalResult(result, events[idx - 1],
			recClosestEvent, testDate, false, false)
		}

	evalInterval(now, execEvents, event, task)
		{
		if task.time_span is true
			{
			end = Date(Display(now.NoTime().Plus(days: -1)) $ '.' $
				Display(task.time_end).LeftFill(.timeDigits, '0'))
			if task.prev_interval > end
				testDate = Date(Display(now.NoTime()) $ '.' $
					Display(task.time_start).LeftFill(.timeDigits, '0'))
			else
				testDate = .getNextInterval(task.prev_interval, task.interval,
					task.interval_units)
			}
		else
			testDate = .getNextInterval(task.prev_interval, task.interval,
				task.interval_units)
		// get correct date from units and interval number
		if (now >= testDate)
			{
			task.prev_interval = now
			task.Update()
			execEvents.Add(event)
			if task.time_span is true
				{
				start = Date(Display(now.NoTime()) $ '.' $
					Display(task.time_start).LeftFill(.timeDigits, '0'))
				end = Date(Display(now.NoTime()) $ '.' $
					Display(task.time_end).LeftFill(.timeDigits, '0'))
				if now < start
					event.next = start
				else if now > end
					event.next = start.Plus(days: 1)
				else
					event.next = .getNextInterval(now, task.interval, task.interval_units)
				}
			else
				event.next = .getNextInterval(now, task.interval, task.interval_units)
			}
		else
			event.next = testDate

		return .milliSecondsBetweenDates(now, event.next)
		}

	getNextInterval(date, interval, units)
		{
		switch (units)
			{
		case 0:	return date.Plus(seconds: interval)
		case 1:	return date.Plus(minutes: interval)
		case 2:	return date.Plus(hours: interval)
		case 3:	return date.Plus(days: interval) /*= days */
		case 4:	return date.Plus(days: (7 * interval)) /*= weeks */
			}
		}

	evalResult(result, event, recClosestEvent, testDate, next, nextDate)
		{
		if Number?(result)
			{
			if (testDate isnt false)
				event.next = .getNext(event.next, testDate)
			if (result < recClosestEvent)
				return result
			}
		else if ((next) and (nextDate isnt false))
			event.next = .getNext(event.next, nextDate)
		return recClosestEvent
		}

	getNext(prevNext, date)
		{
		return prevNext is false or date < prevNext ? date : prevNext
		}

	evalDate(now, testDate, execEvents, event, evaluate?, task)
		{
		if not evaluate?
			return false
		if .simplify(now) <= testDate
			return .current_event(testDate, now, task, execEvents, event)
		else if .simplify(now) > task.next_event // if event skipped for some reason...
			return .skipped_event(task, now, execEvents, event)
		return false
		}

	simplify(date)
		{ return date.NoTime().Plus(hours: date.Hour(), minutes: date.Minute()) }

	current_event(testDate, now, task, execEvents, event)
		{
		origTask = task.Copy()
		nowWithoutSeconds = .simplify(now)
		diff = .milliSecondsBetweenDates(nowWithoutSeconds, testDate)
		if diff is 0		// execute event
			{
			task.next_event = Date.End()
			if nowWithoutSeconds > task.prev_event
				{
				execEvents.Add(event)
				task.prev_event = nowWithoutSeconds
				}
			result = true
			}
		else
			{
			task.next_event = testDate
			result = .milliSecondsBetweenDates(now, testDate)
			}
		if task isnt origTask
			task.Update()
		return result
		}

	skipped_event(task, now, execEvents, event)
		{
		execEvents.Add(event)
		nextRunDate = .next_run_date(task, now)
		if task.next_event isnt nextRunDate
			{
			task.next_event = nextRunDate
			task.Update()
			}
		return true
		}

	next_run_date(task, now)
		{
		if task.suspended is true or task.rundate is true
			return task.next_event
		if task.runinterval is true
			return false
		if not Date?(task.prev_event) or task.prev_event is Date.Begin()
			task.prev_event = .set_prev_event(task)
		if task.rundaily is true
			{
			date = task.prev_event.Replace(year: now.Year(), month: now.Month(),
				day: now.Day())
			if date < now
				date = date.Plus(days: 1)
			return date
			}
		if task.runweekly is true
			{
			date = task.prev_event.Plus(days: 7)
			while date < now
				date = date.Plus(days: 7)
			Assert(date.WeekDay() is: task.prev_event.WeekDay())
			return date
			}
		if task.runmonthly is true
			{
			date = task.prev_event.Replace(month: now.Month())
			if date < now
				date = date.Plus(months: 1)
			return date
			}
		}

	timeDigits: 4
	set_prev_event(task)
		{
		if task.rundaily is true
			return Date(Display(Date.Begin()) $ '.' $
				Display(task.daily_time).LeftFill(.timeDigits, '0'))
		if task.runmonthly is true
			return Date(Display(Date.Begin().Replace(day: task.monthly_day)) $ '.' $
				Display(task.monthly_time).LeftFill(.timeDigits, '0'))
		if task.runweekly is true
			{
			date = Date.Begin()
			while date.WeekDay() isnt task.weekly_day
				date = date.Plus(days: 1)
			return Date(Display(date) $ '.' $
				Display(task.weekly_time).LeftFill(.timeDigits, '0'))
			}
		return Date.Begin()
		}

	OverrideLogTasks: false
	log_task(msg, mem, taskTable, logTasks)
		{
		msg $= " at " $ Display(Date())
		ma = ServerEval('MemoryArena')
		msg $= ' ' $ (ma / 1.Mb()).Round(0) $ 'mb '
		if logTasks is true or .OverrideLogTasks is true
			AddFile('Schedule_Tasks_Log', msg $ '\r\n')

		ob = ServerSuneido.Get(mem, Object())
		ob[taskTable] = msg
		ServerSuneido.Set(mem, ob)
		}
	LogTask(msg, mem, taskTable, logTasks = false)
		{
		.log_task(msg, mem, taskTable, logTasks)
		}

	Destroy()
		{
		.destroyed = true
		}
	}
