// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name:	"ScheduleAddEdit"
	Title:	"Task Information"
	CallClass(hwnd, table, task = false, options = 'all', interval_options = 'all')
		{
		return OkCancel(Object(this, table, task, options,
			interval_options), .Title, hwnd)
		}
	New(taskTable, task = false, options = 'all', interval_options = 'all')
		{
		super(.layout(options, interval_options))
		.task = task
		.taskTable = taskTable
		.setFieldControls()
		if (task isnt false)
			.setControls(task)
		}
	layout(options, interval_options)
		{
		.options = options
		controls = Object('Vert')
		if .hasOption?('Name')
			controls.Add(#(taskName, name: 'taskName'))
		if .hasOption?('Code')
			controls.Add(#(taskCode, name: "taskCode"))
		if .hasOption?('Name') or .hasOption?('Code')
			controls.Add(#Skip)

		controls.Add(.radioOptions(interval_options))
		controls.Add(#Skip)

		botCtrls = Object('Horz')
		if .hasOption?('Suspend')
			botCtrls.Add(#(Highlight (CheckBox "Disable Task" name: checkSuspend)))
		controls.Add(botCtrls)

		return Object('Record' controls)
		}

	hasOption?(option)
		{
		return .options is 'all' or .options.Has?(option)
		}

	radioOptions(interval_options)
		{
		radioGroups = Object('RadioGroups')
		if interval_options is 'all'
			interval_options = #(seconds minutes hours days weeks)

		if .hasOption?('Interval')
			radioGroups.Add(Object('Horz'
				#(Number, name: taskInterval, width: 6, mask: '##',
					rangefrom: 1, rangeto: 59)
				Object('ChooseList', interval_options,
					name: "taskIntervalUnits" xmin: 120), label: 'Intervals'))

		if .hasOption?('TimeSpan')
			radioGroups.Add(#(Horz (Time name: 'timeStart') Skip (Static 'to') Skip
				(Time name: 'timeEnd') Skip (Static 'every') Skip
				(Number, name: timeSpanInterval, width: 6, mask: '##',
					rangefrom: 0, rangeto: 59)
				(ChooseList, #(seconds minutes hours days weeks),
					name: "timeSpanUnits" xmin: 120), label: 'Time Span'))

		if .hasOption?('Daily')
			radioGroups.Add(#(Horz (Time, name: "taskDaily") Skip
				(CheckBox "Skip Weekends" name: 'taskDailySkip'), label: 'Daily At'))

		if .hasOption?('Weekly')
			radioGroups.Add(#(Horz (Time name: "taskWeekly") Skip (Static "on") Skip
				(ChooseList, #(Sunday Monday Tuesday Wednesday Thursday Friday Saturday),
					name: "taskWeeklyOn" xmin: 120),
				label: 'Weekly At'))

		if .hasOption?('Monthly')
			radioGroups.Add(#(Horz (Time name: "taskMonthly") Skip (Static "on") Skip
				(ChooseList, listField: 'scheduler_monthly',
					name: "taskMonthlyOn" xmin: 120),
				label: 'Monthly At'))

		if .hasOption?('DateTime')
			radioGroups.Add(#(Horz (Time name: "taskDateTime") Skip (Static "on") Skip
				(ChooseDate name: "taskDateOn"),
				label: 'Date At'))

		return radioGroups
		}

	dummyDataCtrl: class { Default(@args /*unused*/) { return '' } }

	setFieldControls()
		{
		.taskName = .setDataCtrlReference('taskName')
		.taskCode = .setDataCtrlReference('taskCode')
		.checkSuspend = .setDataCtrlReference('checkSuspend')
		.radioGroups = .setDataCtrlReference('RadioGroups')
		.taskInterval = .setDataCtrlReference('taskInterval')
		.taskIntervalUnits = .setDataCtrlReference('taskIntervalUnits')
		.timeStart = .setDataCtrlReference('timeStart')
		.timeEnd = .setDataCtrlReference('timeEnd')
		.timeInterval = .setDataCtrlReference('timeSpanInterval')
		.timeIntervalUnits = .setDataCtrlReference('timeSpanUnits')
		.taskDaily = .setDataCtrlReference('taskDaily')
		.taskDailySkip = .setDataCtrlReference('taskDailySkip')
		.taskWeekly = .setDataCtrlReference('taskWeekly')
		.taskWeeklyOn = .setDataCtrlReference('taskWeeklyOn')
		.taskMonthly = .setDataCtrlReference('taskMonthly')
		.taskMonthlyOn = .setDataCtrlReference('taskMonthlyOn')
		.taskDateTime = .setDataCtrlReference('taskDateTime')
		.taskDateOn = .setDataCtrlReference('taskDateOn')
		}

	setDataCtrlReference(name)
		{
		ctrl = .FindControl(name)
		return ctrl isnt false ? ctrl : .dummyDataCtrl
		}

	IntervalUnitsMap: #('seconds', 'minutes', 'hours', 'days', 'weeks')
	setControls(task)
		{
		.taskName.Set(task.taskname)
		.taskCode.Set(task.task)

		picked = .picked(task)

		.radioGroups.Picked(picked)
		if picked is 'Intervals'
			.taskInterval.Set(task.interval)
		else if picked is 'Time Span'
			.timeInterval.Set(task.interval)
		if picked is 'Intervals'
			.taskIntervalUnits.Set(.IntervalUnitsMap[task.interval_units])
		else if picked is 'Time Span'
			.timeIntervalUnits.Set(.IntervalUnitsMap[task.interval_units])
		.timeStart.Set(task.time_start)
		.timeEnd.Set(task.time_end)
		.taskDaily.Set(String(task.daily_time).BeforeFirst(' skip weekends'))
		.taskDailySkip.Set(String(task.daily_time).Has?('skip weekends'))
		.taskWeekly.Set(task.weekly_time)
		.taskWeeklyOn.SelectItem(task.weekly_day)
		.taskMonthly.Set(task.monthly_time)
		.taskMonthlyOn.SelectItem(task.monthly_day)
		.taskDateTime.Set(task.date.Hour() * 100 + task.date.Minute()) /*= two digit */
		.taskDateOn.Set(task.date.ShortDate())
		.checkSuspend.Set(task.suspended)
		}
	picked(task)
		{
		return task.runinterval is true
			? 'Intervals'
			: task.time_span is true
				? 'Time Span'
				: task.rundaily is true
					? 'Daily At'
					: task.runmonthly is true
						? 'Monthly At'
						: task.runweekly is true
							? "Weekly At"
							: "Date At"
		}
	getControls()
		{
		task = .initTask()
		task.taskname = .taskName.Get()
		task.task = .taskCode.Get()
		task.suspended = .checkSuspend.Get()
		switch (.radioGroups.Get())
			{
		case "Intervals":
			.setIntervals(task)
		case "Time Span":
			.setTimeSpan(task)
		case "Daily At":
			.setDaily(task)
		case "Weekly At":
			.setWeek(task)
		case "Monthly At":
			.setMonth(task)
		case "Date At":
			.setDateAt(task)
			}
		return new ScheduleTask(task)
		}
	initTask()
		{
		return Record(taskTable: .taskTable, threaded: false,
			prev_event: (.task isnt false) ? .task.prev_event : Date.Begin(),
			next_event: (.task isnt false) ? .task.next_event : Date.End(),
			prev_interval: (.task isnt false) ? .task.prev_interval : Date())
		}
	setIntervals(task)
		{
			task.runinterval = true
			task.interval = .taskInterval.Get()
			task.interval_units = .IntervalUnitsMap.Find(.taskIntervalUnits.Get())
		}
	setTimeSpan(task)
		{
			task.time_span = true
			task.time_start = .timeStart.Get()
			task.time_end = .timeEnd.Get()
			task.interval = .timeInterval.Get()
			task.interval_units = .IntervalUnitsMap.Find(.timeIntervalUnits.Get())
		}
	setDaily(task)
		{
			task.rundaily = true
			task.daily_time = .taskDaily.Get() $ (.taskDailySkip.Get() is true
				? ' skip weekends'
				: '')
			// need this to make sure that skip Weekends does not get
			// set on monthly/weekly as they default to daily_time if not set.
			task.monthly_time = .taskMonthly.Get()
			task.weekly_time = .taskWeekly.Get()
		}
	setWeek(task)
		{
			task.runweekly = true
			task.weekly_time = .taskWeekly.Get()
			task.weekly_day = .taskWeeklyOn.SelectedItem()
		}
	setMonth(task)
		{
			task.runmonthly = true
			task.monthly_time = .taskMonthly.Get()
			task.monthly_day = .taskMonthlyOn.SelectedItem()
		}
	setDateAt(task)
		{
		task.rundate = true
		time = String(.taskDateTime.Get()).LeftFill(4, '0') /*= 4 digit time,
			time needed to create date */
		date = .taskDateOn.Get() // date needed to create date
		if String?(date) and date is ''	// if date is "", get today's date
			date = Date()
		task.date = Date(date.Format("yyyyMMdd") $ "." $ time, "yyyyMMdd.t")
		}

	inputValid?()
		{
		if .Data.Valid() isnt true
			return false
		switch (.radioGroups.Get())
			{
		case "Intervals":
			return not .emptyInterval?()
		case "Time Span":
			return not .emptyTime?()
		case "Daily At":
			return .taskDaily.Get() isnt ''
		case "Weekly At":
			return not .emptyWeek?()
		case "Monthly At":
			return not .emptyMonth?()
		case "Date At":
			return not .emptyDate?()
			}
		}
	emptyInterval?()
		{
		return .taskInterval.Get() is '' or .taskInterval.Get() is 0 or
			.taskIntervalUnits.Get() is ''
		}
	emptyTime?()
		{
		return .timeStart.Get() is '' or .timeEnd.Get() is '' or
			.timeInterval.Get() is '' or .timeIntervalUnits.Get() is ''
		}
	emptyWeek?()
		{
		return .taskWeekly.Get() is '' or .taskWeeklyOn.Get() is ''
		}
	emptyMonth?()
		{
		return .taskMonthly.Get() is '' or .taskMonthlyOn.Get() is ''
		}
	emptyDate?()
		{
		return .taskDateTime.Get() is '' or .taskDateOn.Get() is ''
		}
	OK()
		{
		if (.inputValid?())
			return .getControls()
		if not YesNo("One or more fields contains incomplete or " $
			"incorrect information.\n\n" $
			"Do you wish to correct it?  " $
			"(If you choose No, the information will be discarded).",
			"Task Information Incomplete",
			.Window.Hwnd, MB.ICONERROR)
			.Send('On_Cancel')
		return false
		}
	}
